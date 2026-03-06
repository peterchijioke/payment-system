package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"take-Home-assignment/internal/dto"
	"take-Home-assignment/internal/models"
	"take-Home-assignment/internal/repositories"

	"gorm.io/gorm"
)

type PaymentService struct {
	accountRepo     repositories.AccountRepository
	transactionRepo repositories.TransactionRepository
	ledgerRepo      repositories.LedgerEntryRepository
	fxQuoteRepo     repositories.FXQuoteRepository
	webhookRepo     repositories.WebhookEventRepository
	idempotencyRepo repositories.IdempotencyKeyRepository
	db              *gorm.DB
	webhookSecret   string
}

func NewPaymentService(
	accountRepo repositories.AccountRepository,
	transactionRepo repositories.TransactionRepository,
	ledgerRepo repositories.LedgerEntryRepository,
	fxQuoteRepo repositories.FXQuoteRepository,
	webhookRepo repositories.WebhookEventRepository,
	idempotencyRepo repositories.IdempotencyKeyRepository,
	db *gorm.DB,
	webhookSecret string,
) *PaymentService {
	return &PaymentService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		ledgerRepo:      ledgerRepo,
		fxQuoteRepo:     fxQuoteRepo,
		webhookRepo:     webhookRepo,
		idempotencyRepo: idempotencyRepo,
		db:              db,
		webhookSecret:   webhookSecret,
	}
}

func (s *PaymentService) ProcessPayment(db *gorm.DB, req *dto.PaymentRequest, idempotencyKey string) (*dto.PaymentResponse, error) {
	req.IdempotencyKey = idempotencyKey

	existingResp, err := s.checkIdempotency(db, req)
	if err != nil {
		return nil, err
	}
	if existingResp != nil {
		return existingResp, nil
	}

	tx := db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.validatePaymentRequest(tx, req); err != nil {
		tx.Rollback()
		return nil, err
	}

	quote, err := s.getFXQuote(tx, req.Currency, req.DestinationCurrency)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := s.accountRepo.LockFunds(tx, req.AccountID, req.Currency, req.Amount); err != nil {
		tx.Rollback()
		return nil, errors.New("failed to lock funds - insufficient available balance")
	}

	txRef := s.generateTransactionRef()
	providerRef := s.generateProviderReference()
	transaction := &models.PaymentTransaction{
		TransactionRef: txRef,
		IdempotencyKey: req.IdempotencyKey,
		AccountID:      req.AccountID,
		Type:           models.TransactionTypePayment,
		Status:         models.TransactionStatusInitiated,
		Amount:         req.Amount,
		Currency:       req.Currency,
		FXQuoteID:      &quote.ID,
		FXRate:         &quote.Rate,
		FXAmount:       floatPtr(quote.Rate * req.Amount),
		FXCurrency:     &req.DestinationCurrency,
		Description:    fmt.Sprintf("Payment to %s", req.RecipientName),
		Reference:      req.Reference,
	}

	if err := s.transactionRepo.Create(tx, transaction); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err := s.createLedgerEntries(tx, transaction); err != nil {
		tx.Rollback()
		return nil, err
	}

	err = s.submitToDownstream(transaction, req)
	if err != nil {
		transaction.Status = models.TransactionStatusProcessing
		tx.Save(transaction)
		tx.Commit()
		return &dto.PaymentResponse{
			TransactionID:     transaction.ID,
			TransactionRef:    transaction.TransactionRef,
			ProviderReference: providerRef,
			Status:            string(transaction.Status),
			Amount:            transaction.Amount,
			Currency:          transaction.Currency,
			CreatedAt:         transaction.InitiatedAt,
		}, nil
	}

	now := time.Now().UTC()
	s.fxQuoteRepo.LockQuote(tx, quote.ID, now)

	s.saveIdempotencyKey(tx, req, http.StatusOK, transaction)

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &dto.PaymentResponse{
		TransactionID:       transaction.ID,
		TransactionRef:      transaction.TransactionRef,
		ProviderReference:   providerRef,
		Status:              string(transaction.Status),
		Amount:              transaction.Amount,
		Currency:            transaction.Currency,
		FXRate:              *transaction.FXRate,
		FXAmount:            *transaction.FXAmount,
		DestinationCurrency: req.DestinationCurrency,
		CreatedAt:           transaction.InitiatedAt,
	}, nil
}

func (s *PaymentService) GetTransactionDetails(transactionID string) (*dto.TransactionDetails, error) {
	transaction, err := s.transactionRepo.FindByID(s.db, transactionID)
	if err != nil {
		return nil, err
	}

	entries, err := s.ledgerRepo.FindByTransactionID(s.db, transactionID)
	if err != nil {
		entries = []models.LedgerEntry{}
	}

	timeline := s.buildTimeline(transaction)

	return &dto.TransactionDetails{
		Transaction:   transaction,
		LedgerEntries: entries,
		Timeline:      timeline,
	}, nil
}

func (s *PaymentService) ListPayments(limit, offset int, status, startDate, endDate string) ([]dto.TransactionDetails, int64, error) {
	transactions, count, err := s.transactionRepo.FindAll(s.db, limit, offset, status, startDate, endDate)
	if err != nil {
		return nil, 0, err
	}

	var results []dto.TransactionDetails
	for _, tx := range transactions {
		entries, _ := s.ledgerRepo.FindByTransactionID(s.db, tx.ID)
		timeline := s.buildTimeline(&tx)
		results = append(results, dto.TransactionDetails{
			Transaction:   &tx,
			LedgerEntries: entries,
			Timeline:      timeline,
		})
	}

	return results, count, nil
}

func (s *PaymentService) ProcessWebhook(db *gorm.DB, rawBody io.Reader, signature string) (*dto.WebhookPayload, error) {
	bodyBytes, err := io.ReadAll(rawBody)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	if !s.verifyWebhookSignature(string(bodyBytes), signature) {
		return nil, errors.New("invalid webhook signature")
	}

	var payload dto.WebhookPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	existingEvent, _ := s.webhookRepo.FindByEventIDAndStatus(db, payload.EventID, "processed")
	if existingEvent != nil {
		return &payload, nil
	}

	var payloadStr = string(bodyBytes)
	webhookEvent := &models.WebhookEvent{
		Source:           "downstream_provider",
		EventType:        payload.EventType,
		EventID:          payload.EventID,
		ProcessingStatus: "processing",
		Payload:          &payloadStr,
	}

	if err := s.webhookRepo.Create(db, webhookEvent); err != nil {
		return nil, fmt.Errorf("failed to log webhook: %w", err)
	}

	tx := db.Begin()

	transaction, err := s.transactionRepo.FindByIDForUpdate(tx, payload.TransactionID)
	if err != nil {
		tx.Rollback()
		s.webhookRepo.MarkFailed(db, payload.EventID, "transaction not found")
		return &payload, nil
	}

	if payload.Amount != "" {
		webhookAmount, err := strconv.ParseFloat(payload.Amount, 64)
		if err == nil && webhookAmount != transaction.Amount {
			tx.Rollback()
			s.webhookRepo.MarkFailed(db, payload.EventID, "amount mismatch")
			return nil, errors.New("webhook amount does not match transaction amount - possible fraud")
		}
	}

	if err := s.validateStateTransition(transaction.Status, models.TransactionStatus(payload.Status)); err != nil {
		tx.Rollback()
		s.webhookRepo.MarkFailed(db, payload.EventID, err.Error())
		return nil, fmt.Errorf("invalid state transition: %w", err)
	}

	switch models.TransactionStatus(payload.Status) {
	case models.TransactionStatusSettled, models.TransactionStatusCompleted:
		if err := s.settlePayment(tx, transaction); err != nil {
			tx.Rollback()
			return nil, err
		}
	case models.TransactionStatusReversed:
		if err := s.reversePayment(tx, transaction, payload.FailureReason); err != nil {
			tx.Rollback()
			return nil, err
		}
	case models.TransactionStatusFailed:
		if err := s.reversePayment(tx, transaction, payload.FailureReason); err != nil {
			tx.Rollback()
			return nil, err
		}
	default:
	}

	s.webhookRepo.MarkProcessed(tx, payload.EventID)

	return &payload, tx.Commit().Error
}

func (s *PaymentService) validateStateTransition(currentStatus, newStatus models.TransactionStatus) error {
	validTransitions := map[models.TransactionStatus][]models.TransactionStatus{
		models.TransactionStatusInitiated:  {models.TransactionStatusProcessing, models.TransactionStatusSettled, models.TransactionStatusCompleted, models.TransactionStatusFailed},
		models.TransactionStatusProcessing: {models.TransactionStatusSettled, models.TransactionStatusCompleted, models.TransactionStatusFailed},
		models.TransactionStatusSettled:    {models.TransactionStatusReversed},
		models.TransactionStatusCompleted:  {models.TransactionStatusReversed},
	}

	allowed, exists := validTransitions[currentStatus]
	if !exists {
		if currentStatus == models.TransactionStatusCompleted ||
			currentStatus == models.TransactionStatusFailed ||
			currentStatus == models.TransactionStatusReversed ||
			currentStatus == models.TransactionStatusSettled {
			return nil
		}
		return fmt.Errorf("current status %s has no valid transitions", currentStatus)
	}

	for _, allowedStatus := range allowed {
		if newStatus == allowedStatus {
			return nil
		}
	}

	return fmt.Errorf("invalid transition from %s to %s", currentStatus, newStatus)
}

func (s *PaymentService) verifyWebhookSignature(payload, signature string) bool {
	if s.webhookSecret == "" {
		return true
	}

	mac := hmac.New(sha256.New, []byte(s.webhookSecret))
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}

func (s *PaymentService) settlePayment(tx *gorm.DB, transaction *models.PaymentTransaction) error {
	now := time.Now().UTC()
	transaction.Status = models.TransactionStatusSettled
	transaction.SettledAt = &now
	transaction.CompletedAt = &now

	if err := s.transactionRepo.Update(tx, transaction); err != nil {
		return err
	}

	if err := s.ledgerRepo.UpdateStatus(tx, transaction.ID, models.LedgerEntryStatusPosted, now, nil); err != nil {
		return err
	}

	// Credit the recipient's account in destination currency
	if transaction.CounterpartyID != nil && transaction.FXAmount != nil && transaction.FXCurrency != nil {
		// Create settlement credit ledger entry
		creditEntry := &models.LedgerEntry{
			EntryReference: fmt.Sprintf("LED-%s-SETTLE-CR", transaction.TransactionRef),
			TransactionID:  transaction.ID,
			AccountID:      *transaction.CounterpartyID,
			EntryType:      models.LedgerEntryTypeSettlementCredit,
			Amount:         *transaction.FXAmount,
			Currency:       *transaction.FXCurrency,
			Status:         models.LedgerEntryStatusPosted,
			Description:    fmt.Sprintf("Payment settlement credit for %s", transaction.TransactionRef),
			EffectiveDate:  now,
			PostedAt:       &now,
		}

		if err := s.ledgerRepo.Create(tx, creditEntry); err != nil {
			return fmt.Errorf("failed to create settlement credit entry: %w", err)
		}

		// Credit the recipient's account balance
		if err := s.accountRepo.CreditFunds(tx, *transaction.CounterpartyID, *transaction.FXCurrency, *transaction.FXAmount); err != nil {
			return fmt.Errorf("failed to credit recipient account: %w", err)
		}
	}

	return nil
}

func (s *PaymentService) reversePayment(tx *gorm.DB, transaction *models.PaymentTransaction, reason string) error {
	now := time.Now().UTC()

	if transaction.Status == models.TransactionStatusSettled || transaction.Status == models.TransactionStatusCompleted {
		transaction.Status = models.TransactionStatusReversed
		transaction.ReversedByID = &transaction.ID
		transaction.ReversalReason = reason

		if err := s.transactionRepo.Update(tx, transaction); err != nil {
			return err
		}

		if err := s.createSettlementReversalEntries(tx, transaction); err != nil {
			return err
		}

		return nil
	}

	transaction.Status = models.TransactionStatusFailed
	transaction.FailedAt = &now
	transaction.FailureReason = reason

	if err := s.transactionRepo.Update(tx, transaction); err != nil {
		return err
	}

	if err := s.ledgerRepo.UpdateStatus(tx, transaction.ID, models.LedgerEntryStatusReversed, nil, transaction.ID); err != nil {
		return err
	}

	return s.accountRepo.ReleaseFunds(tx, transaction.AccountID, transaction.Currency, transaction.Amount)
}

func (s *PaymentService) createSettlementReversalEntries(tx *gorm.DB, transaction *models.PaymentTransaction) error {
	entries, err := s.ledgerRepo.FindByTransactionID(tx, transaction.ID)
	if err != nil {
		return fmt.Errorf("failed to find original ledger entries: %w", err)
	}

	now := time.Now().UTC()

	for _, originalEntry := range entries {
		reversalEntryType := models.LedgerEntryTypeDebit
		if originalEntry.EntryType == models.LedgerEntryTypeDebit {
			reversalEntryType = models.LedgerEntryTypeCredit
		} else if originalEntry.EntryType == models.LedgerEntryTypeSettlementDebit {
			reversalEntryType = models.LedgerEntryTypeSettlementReversalCredit
		}

		reversalEntry := &models.LedgerEntry{
			EntryReference:  fmt.Sprintf("LED-%s-REV-%s", transaction.TransactionRef, originalEntry.EntryReference[len(originalEntry.EntryReference)-4:]),
			TransactionID:   transaction.ID,
			AccountID:       originalEntry.AccountID,
			EntryType:       reversalEntryType,
			Amount:          originalEntry.Amount,
			Currency:        originalEntry.Currency,
			OriginalEntryID: &originalEntry.ID,
			Status:          models.LedgerEntryStatusPosted,
			ReversalReason:  transaction.ReversalReason,
			Description:     fmt.Sprintf("Reversal of %s", originalEntry.EntryReference),
			EffectiveDate:   now,
			PostedAt:        &now,
			ReversedByID:    &transaction.ID,
		}

		if err := s.ledgerRepo.Create(tx, reversalEntry); err != nil {
			return fmt.Errorf("failed to create reversal ledger entry: %w", err)
		}
	}

	return nil
}

func (s *PaymentService) validatePaymentRequest(tx *gorm.DB, req *dto.PaymentRequest) error {
	_, err := s.accountRepo.FindByIDAndActive(tx, req.AccountID)
	if err != nil {
		return errors.New("account not found or inactive")
	}

	_, err = s.accountRepo.FindBalance(tx, req.AccountID, req.Currency)
	if err != nil {
		return errors.New("no balance for specified currency")
	}

	account, _ := s.accountRepo.FindByID(tx, req.AccountID)
	balance, _ := s.accountRepo.FindBalance(tx, req.AccountID, req.Currency)
	available := balance.AvailableBalance - balance.ReservedBalance
	if available < req.Amount {
		return errors.New("insufficient funds")
	}

	if account.DailyLimit > 0 {
		since := time.Now().UTC().Truncate(24 * time.Hour)
		total, _ := s.accountRepo.GetDailyTransactionTotal(tx, req.AccountID, req.Currency, since)
		if total+req.Amount > account.DailyLimit {
			return errors.New("would exceed daily limit")
		}
	}

	return nil
}

func (s *PaymentService) getFXQuote(tx *gorm.DB, fromCurrency, toCurrency string) (*models.FXQuote, error) {
	quote, err := s.fxQuoteRepo.FindValidQuote(tx, fromCurrency, toCurrency, time.Now().UTC())
	if err == nil {
		return quote, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	rate := s.getMockFXRate(fromCurrency, toCurrency)
	quote = &models.FXQuote{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         rate,
		ValidUntil:   time.Now().UTC().Add(15 * time.Minute),
		QuoteID:      fmt.Sprintf("QUOTE-%d-%s-%s", time.Now().Unix(), fromCurrency, toCurrency),
	}

	if err := s.fxQuoteRepo.Create(tx, quote); err != nil {
		return nil, err
	}

	return quote, nil
}

func (s *PaymentService) getMockFXRate(from, to string) float64 {
	rates := map[string]float64{
		"NGN-USD": 0.00065,
		"USD-NGN": 1535.0,
		"NGN-EUR": 0.00060,
		"EUR-NGN": 1667.0,
		"NGN-GBP": 0.00052,
		"GBP-NGN": 1923.0,
		"USD-GBP": 0.79,
		"GBP-USD": 1.27,
		"USD-EUR": 0.92,
		"EUR-USD": 1.09,
	}

	if rate, ok := rates[from+"-"+to]; ok {
		return rate
	}

	return 1.0 + (rand.Float64()*0.1 - 0.05)
}

func (s *PaymentService) createLedgerEntries(tx *gorm.DB, transaction *models.PaymentTransaction) error {
	settlement, err := s.accountRepo.GetSettlementAccount(tx)
	if err != nil {
		return fmt.Errorf("settlement account not found: %w", err)
	}

	debitEntry := &models.LedgerEntry{
		EntryReference: fmt.Sprintf("LED-%s-DR", transaction.TransactionRef),
		TransactionID:  transaction.ID,
		AccountID:      transaction.AccountID,
		EntryType:      models.LedgerEntryTypeDebit,
		Amount:         transaction.Amount,
		Currency:       transaction.Currency,
		Description:    fmt.Sprintf("Payment %s - Debit", transaction.TransactionRef),
		EffectiveDate:  time.Now().UTC().Truncate(24 * time.Hour),
	}

	creditEntry := &models.LedgerEntry{
		EntryReference: fmt.Sprintf("LED-%s-CR", transaction.TransactionRef),
		TransactionID:  transaction.ID,
		AccountID:      settlement.ID,
		EntryType:      models.LedgerEntryTypeCredit,
		Amount:         transaction.Amount,
		Currency:       transaction.Currency,
		Description:    fmt.Sprintf("Payment %s - Credit", transaction.TransactionRef),
		EffectiveDate:  time.Now().UTC().Truncate(24 * time.Hour),
	}

	if err := s.ledgerRepo.Create(tx, debitEntry); err != nil {
		return err
	}

	if err := s.ledgerRepo.Create(tx, creditEntry); err != nil {
		return err
	}

	debitEntry.CounterpartEntryID = &creditEntry.ID
	creditEntry.CounterpartEntryID = &debitEntry.ID

	s.ledgerRepo.Update(tx, debitEntry)
	s.ledgerRepo.Update(tx, creditEntry)

	return nil
}

func (s *PaymentService) submitToDownstream(transaction *models.PaymentTransaction, req *dto.PaymentRequest) error {
	transaction.Status = models.TransactionStatusProcessing
	return nil
}

func (s *PaymentService) checkIdempotency(db *gorm.DB, req *dto.PaymentRequest) (*dto.PaymentResponse, error) {
	hash := generateRequestHash(req)

	idempotency, err := s.idempotencyRepo.FindValid(db, req.IdempotencyKey, hash, time.Now().UTC())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var cachedResp dto.PaymentResponse
	if err := json.Unmarshal([]byte(idempotency.ResponseBody), &cachedResp); err != nil {
		return nil, nil
	}

	return &cachedResp, nil
}

func (s *PaymentService) saveIdempotencyKey(tx *gorm.DB, req *dto.PaymentRequest, status int, transaction *models.PaymentTransaction) error {
	hash := generateRequestHash(req)
	resp := dto.PaymentResponse{
		TransactionID:       transaction.ID,
		TransactionRef:      transaction.TransactionRef,
		Status:              string(transaction.Status),
		Amount:              transaction.Amount,
		Currency:            transaction.Currency,
		DestinationCurrency: *transaction.FXCurrency,
		CreatedAt:           transaction.InitiatedAt,
	}
	respBytes, _ := json.Marshal(resp)

	idempotencyKey := &models.IdempotencyKey{
		Key:              req.IdempotencyKey,
		AccountID:        req.AccountID,
		RequestHash:      hash,
		RequestMethod:    "POST",
		RequestPath:      "/payments",
		OriginalAmount:   req.Amount,
		OriginalCurrency: req.Currency,
		ResponseStatus:   &status,
		ResponseBody:     string(respBytes),
		ExpiresAt:        time.Now().UTC().Add(24 * time.Hour),
	}

	return s.idempotencyRepo.Create(tx, idempotencyKey)
}

func (s *PaymentService) buildTimeline(transaction *models.PaymentTransaction) []dto.StatusChange {
	var timeline []dto.StatusChange

	timeline = append(timeline, dto.StatusChange{
		Status:    string(models.TransactionStatusInitiated),
		Timestamp: transaction.InitiatedAt.Format(time.RFC3339),
	})

	if transaction.ProcessedAt != nil {
		timeline = append(timeline, dto.StatusChange{
			Status:    "processing",
			Timestamp: transaction.ProcessedAt.Format(time.RFC3339),
		})
	}

	if transaction.CompletedAt != nil {
		timeline = append(timeline, dto.StatusChange{
			Status:    string(models.TransactionStatusCompleted),
			Timestamp: transaction.CompletedAt.Format(time.RFC3339),
		})
	}

	if transaction.FailedAt != nil {
		timeline = append(timeline, dto.StatusChange{
			Status:    string(models.TransactionStatusFailed),
			Timestamp: transaction.FailedAt.Format(time.RFC3339),
		})
	}

	return timeline
}

func generateRequestHash(req *dto.PaymentRequest) string {
	data := fmt.Sprintf("%s|%s|%f|%s", req.AccountID, req.Currency, req.Amount, req.IdempotencyKey)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (s *PaymentService) generateTransactionRef() string {
	return fmt.Sprintf("TXN-%d-%s", time.Now().UnixNano(), generateRandomString(8))
}

func (s *PaymentService) generateProviderReference() string {
	return fmt.Sprintf("PRV-%d-%s", time.Now().UnixNano(), strings.ToUpper(generateRandomString(12)))
}

func generateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func floatPtr(f float64) *float64 {
	return &f
}

func (s *PaymentService) GetTransaction(transactionID string) (*models.PaymentTransaction, error) {
	return s.transactionRepo.FindByID(s.db, transactionID)
}

func (s *PaymentService) GetDB() *gorm.DB {
	return s.db
}

func (s *PaymentService) GetAccountBalance(accountID, currency string) (*models.AccountBalance, error) {
	return s.accountRepo.FindBalance(s.db, accountID, currency)
}

func (s *PaymentService) ListTransactions(accountID string, limit, offset int) ([]models.PaymentTransaction, int64, error) {
	return s.transactionRepo.FindByAccountID(s.db, accountID, limit, offset)
}

func (s *PaymentService) ListAccounts() ([]models.Account, error) {
	return s.accountRepo.FindAll(s.db)
}
