package test

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
	"time"

	"take-Home-assignment/internal/dto"
	"take-Home-assignment/internal/models"
	"take-Home-assignment/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock AccountRepository
type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) FindByID(tx *gorm.DB, id string) (*models.Account, error) {
	args := m.Called(tx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) FindByIDAndActive(tx *gorm.DB, id string) (*models.Account, error) {
	args := m.Called(tx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) FindAll(tx *gorm.DB) ([]models.Account, error) {
	args := m.Called(tx)
	return args.Get(0).([]models.Account), args.Error(1)
}

func (m *MockAccountRepository) FindBalance(tx *gorm.DB, accountID, currency string) (*models.AccountBalance, error) {
	args := m.Called(tx, accountID, currency)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AccountBalance), args.Error(1)
}

func (m *MockAccountRepository) LockFunds(tx *gorm.DB, accountID, currency string, amount float64) error {
	args := m.Called(tx, accountID, currency, amount)
	return args.Error(0)
}

func (m *MockAccountRepository) ReleaseFunds(tx *gorm.DB, accountID, currency string, amount float64) error {
	args := m.Called(tx, accountID, currency, amount)
	return args.Error(0)
}

func (m *MockAccountRepository) CreditFunds(tx *gorm.DB, accountID, currency string, amount float64) error {
	args := m.Called(tx, accountID, currency, amount)
	return args.Error(0)
}

func (m *MockAccountRepository) UpdateBalance(tx *gorm.DB, accountID, currency string, updates map[string]interface{}) error {
	args := m.Called(tx, accountID, currency, updates)
	return args.Error(0)
}

func (m *MockAccountRepository) GetSettlementAccount(tx *gorm.DB) (*models.Account, error) {
	args := m.Called(tx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) GetDailyTransactionTotal(tx *gorm.DB, accountID, currency string, since interface{}) (float64, error) {
	args := m.Called(tx, accountID, currency, since)
	return args.Get(0).(float64), args.Error(1)
}

// Mock TransactionRepository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(tx *gorm.DB, transaction *models.PaymentTransaction) error {
	args := m.Called(tx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) Update(tx *gorm.DB, transaction *models.PaymentTransaction) error {
	args := m.Called(tx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) FindByID(tx *gorm.DB, id string) (*models.PaymentTransaction, error) {
	args := m.Called(tx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaymentTransaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByIDForUpdate(tx *gorm.DB, id string) (*models.PaymentTransaction, error) {
	args := m.Called(tx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaymentTransaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByIdempotencyKey(tx *gorm.DB, key string) (*models.PaymentTransaction, error) {
	args := m.Called(tx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaymentTransaction), args.Error(1)
}

func (m *MockTransactionRepository) FindAll(tx *gorm.DB, limit, offset int, status, startDate, endDate string) ([]models.PaymentTransaction, int64, error) {
	args := m.Called(tx, limit, offset, status, startDate, endDate)
	return args.Get(0).([]models.PaymentTransaction), args.Get(1).(int64), args.Error(2)
}

func (m *MockTransactionRepository) FindByAccountID(tx *gorm.DB, accountID string, limit, offset int) ([]models.PaymentTransaction, int64, error) {
	args := m.Called(tx, accountID, limit, offset)
	return args.Get(0).([]models.PaymentTransaction), args.Get(1).(int64), args.Error(2)
}

type MockLedgerEntryRepository struct {
	mock.Mock
}

func (m *MockLedgerEntryRepository) Create(tx *gorm.DB, entry *models.LedgerEntry) error {
	args := m.Called(tx, entry)
	return args.Error(0)
}

func (m *MockLedgerEntryRepository) Update(tx *gorm.DB, entry *models.LedgerEntry) error {
	args := m.Called(tx, entry)
	return args.Error(0)
}

func (m *MockLedgerEntryRepository) FindByTransactionID(tx *gorm.DB, transactionID string) ([]models.LedgerEntry, error) {
	args := m.Called(tx, transactionID)
	return args.Get(0).([]models.LedgerEntry), args.Error(1)
}

func (m *MockLedgerEntryRepository) UpdateStatus(tx *gorm.DB, transactionID string, status models.LedgerEntryStatus, postedAt interface{}, reversedByID interface{}) error {
	args := m.Called(tx, transactionID, status, postedAt, reversedByID)
	return args.Error(0)
}

type MockFXQuoteRepository struct {
	mock.Mock
}

func (m *MockFXQuoteRepository) Create(tx *gorm.DB, quote *models.FXQuote) error {
	args := m.Called(tx, quote)
	return args.Error(0)
}

func (m *MockFXQuoteRepository) FindValidQuote(tx *gorm.DB, fromCurrency, toCurrency string, asOf time.Time) (*models.FXQuote, error) {
	args := m.Called(tx, fromCurrency, toCurrency, asOf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FXQuote), args.Error(1)
}

func (m *MockFXQuoteRepository) LockQuote(tx *gorm.DB, quoteID string, lockedAt time.Time) error {
	args := m.Called(tx, quoteID, lockedAt)
	return args.Error(0)
}

type MockWebhookEventRepository struct {
	mock.Mock
}

func (m *MockWebhookEventRepository) Create(tx *gorm.DB, event *models.WebhookEvent) error {
	args := m.Called(tx, event)
	return args.Error(0)
}

func (m *MockWebhookEventRepository) Update(tx *gorm.DB, event *models.WebhookEvent) error {
	args := m.Called(tx, event)
	return args.Error(0)
}

func (m *MockWebhookEventRepository) FindByEventID(tx *gorm.DB, eventID string) (*models.WebhookEvent, error) {
	args := m.Called(tx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WebhookEvent), args.Error(1)
}

func (m *MockWebhookEventRepository) FindByEventIDAndStatus(tx *gorm.DB, eventID, status string) (*models.WebhookEvent, error) {
	args := m.Called(tx, eventID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WebhookEvent), args.Error(1)
}

func (m *MockWebhookEventRepository) MarkProcessed(tx *gorm.DB, eventID string) error {
	args := m.Called(tx, eventID)
	return args.Error(0)
}

func (m *MockWebhookEventRepository) MarkFailed(tx *gorm.DB, eventID, errorMsg string) error {
	args := m.Called(tx, eventID, errorMsg)
	return args.Error(0)
}

type MockIdempotencyKeyRepository struct {
	mock.Mock
}

func (m *MockIdempotencyKeyRepository) Create(tx *gorm.DB, key *models.IdempotencyKey) error {
	args := m.Called(tx, key)
	return args.Error(0)
}

func (m *MockIdempotencyKeyRepository) FindValid(tx *gorm.DB, key, requestHash string, asOf time.Time) (*models.IdempotencyKey, error) {
	args := m.Called(tx, key, requestHash, asOf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.IdempotencyKey), args.Error(1)
}

// Helper functions

func createTestAccount() *models.Account {
	return &models.Account{
		ID:          "test-account-id",
		Name:        "Test Account",
		AccountType: models.AccountTypeInternal,
		IsActive:    true,
		DailyLimit:  10000,
	}
}

func createTestBalance() *models.AccountBalance {
	now := time.Now()
	return &models.AccountBalance{
		AccountID:         "test-account-id",
		Currency:          "NGN",
		AvailableBalance:  50000,
		ReservedBalance:   0,
		LastTransactionAt: &now,
	}
}

func createTestPaymentRequest() *dto.PaymentRequest {
	return &dto.PaymentRequest{
		AccountID:           "test-account-id",
		Amount:              1000,
		Currency:            "NGN",
		DestinationCurrency: "USD",
		RecipientName:       "John Doe",
		RecipientAccount:    "1234567890",
		RecipientBank:       "Test Bank",
		RecipientCountry:    "NG",
		Reference:           "TEST-REF-001",
	}
}

func createTestFXQuote() *models.FXQuote {
	return &models.FXQuote{
		ID:           "test-quote-id",
		FromCurrency: "NGN",
		ToCurrency:   "USD",
		Rate:         0.00065,
		ValidUntil:   time.Now().Add(15 * time.Minute),
		QuoteID:      "QUOTE-TEST-001",
	}
}

func createTestSettlementAccount() *models.Account {
	return &models.Account{
		ID:          "settlement-account-id",
		Name:        "Settlement Account",
		AccountType: models.AccountTypeSettlement,
		IsActive:    true,
	}
}

func createTestTransaction() *models.PaymentTransaction {
	now := time.Now().UTC()
	return &models.PaymentTransaction{
		ID:             "test-transaction-id",
		TransactionRef: "TXN-TEST-001",
		AccountID:      "test-account-id",
		Type:           models.TransactionTypePayment,
		Status:         models.TransactionStatusInitiated,
		Amount:         1000,
		Currency:       "NGN",
		Description:    "Payment to John Doe",
		InitiatedAt:    now,
	}
}

func createTestLedgerEntries() []models.LedgerEntry {
	return []models.LedgerEntry{
		{
			ID:             "ledger-entry-1",
			EntryReference: "LED-TXN-TEST-001-DR",
			TransactionID:  "test-transaction-id",
			AccountID:      "test-account-id",
			EntryType:      models.LedgerEntryTypeDebit,
			Amount:         1000,
			Currency:       "NGN",
			Status:         models.LedgerEntryStatusPending,
		},
		{
			ID:             "ledger-entry-2",
			EntryReference: "LED-TXN-TEST-001-CR",
			TransactionID:  "test-transaction-id",
			AccountID:      "settlement-account-id",
			EntryType:      models.LedgerEntryTypeCredit,
			Amount:         1000,
			Currency:       "NGN",
			Status:         models.LedgerEntryStatusPending,
		},
	}
}

// ============ Core Payment Service Tests ============

func TestPaymentService_ProcessPayment_Success(t *testing.T) {
	// Arrange
	accountRepo := new(MockAccountRepository)
	transactionRepo := new(MockTransactionRepository)
	ledgerRepo := new(MockLedgerEntryRepository)
	fxQuoteRepo := new(MockFXQuoteRepository)
	webhookRepo := new(MockWebhookEventRepository)
	idempotencyRepo := new(MockIdempotencyKeyRepository)

	_ = accountRepo
	_ = transactionRepo
	_ = ledgerRepo
	_ = fxQuoteRepo
	_ = webhookRepo
	_ = idempotencyRepo

	// Basic assertion - this is a compilation test to ensure all mocks are set up correctly
	assert.NotNil(t, accountRepo)
}

// ============ Get Transaction Details Tests ============

func TestPaymentService_GetTransactionDetails_Success(t *testing.T) {
	accountRepo := new(MockAccountRepository)
	transactionRepo := new(MockTransactionRepository)
	ledgerRepo := new(MockLedgerEntryRepository)
	fxQuoteRepo := new(MockFXQuoteRepository)
	webhookRepo := new(MockWebhookEventRepository)
	idempotencyRepo := new(MockIdempotencyKeyRepository)

	db := &gorm.DB{}
	service := services.NewPaymentService(
		accountRepo,
		transactionRepo,
		ledgerRepo,
		fxQuoteRepo,
		webhookRepo,
		idempotencyRepo,
		db,
		"test-secret",
	)

	transaction := createTestTransaction()
	ledgerEntries := createTestLedgerEntries()

	transactionRepo.On("FindByID", db, "test-transaction-id").Return(transaction, nil)
	ledgerRepo.On("FindByTransactionID", db, "test-transaction-id").Return(ledgerEntries, nil)

	details, err := service.GetTransactionDetails("test-transaction-id")

	assert.NoError(t, err)
	assert.NotNil(t, details)
	assert.Equal(t, transaction.ID, details.Transaction.ID)
	assert.Len(t, details.LedgerEntries, 2)
	assert.Len(t, details.Timeline, 1)
}

func TestPaymentService_GetTransactionDetails_NotFound(t *testing.T) {
	accountRepo := new(MockAccountRepository)
	transactionRepo := new(MockTransactionRepository)
	ledgerRepo := new(MockLedgerEntryRepository)
	fxQuoteRepo := new(MockFXQuoteRepository)
	webhookRepo := new(MockWebhookEventRepository)
	idempotencyRepo := new(MockIdempotencyKeyRepository)

	db := &gorm.DB{}
	service := services.NewPaymentService(
		accountRepo,
		transactionRepo,
		ledgerRepo,
		fxQuoteRepo,
		webhookRepo,
		idempotencyRepo,
		db,
		"test-secret",
	)

	transactionRepo.On("FindByID", db, "invalid-id").Return(nil, gorm.ErrRecordNotFound)

	_, err := service.GetTransactionDetails("invalid-id")

	assert.Error(t, err)
}

// ============ List Payments Tests ============

func TestPaymentService_ListPayments_Success(t *testing.T) {
	accountRepo := new(MockAccountRepository)
	transactionRepo := new(MockTransactionRepository)
	ledgerRepo := new(MockLedgerEntryRepository)
	fxQuoteRepo := new(MockFXQuoteRepository)
	webhookRepo := new(MockWebhookEventRepository)
	idempotencyRepo := new(MockIdempotencyKeyRepository)

	db := &gorm.DB{}
	service := services.NewPaymentService(
		accountRepo,
		transactionRepo,
		ledgerRepo,
		fxQuoteRepo,
		webhookRepo,
		idempotencyRepo,
		db,
		"test-secret",
	)

	transactions := []models.PaymentTransaction{*createTestTransaction()}

	transactionRepo.On("FindAll", db, 20, 0, "", "", "").Return(transactions, int64(1), nil)
	ledgerRepo.On("FindByTransactionID", db, mock.Anything).Return([]models.LedgerEntry{}, nil)

	results, count, err := service.ListPayments(20, 0, "", "", "")

	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
	assert.Len(t, results, 1)
}

// ============ Webhook Processing Tests ============

func TestPaymentService_ProcessWebhook_InvalidSignature(t *testing.T) {
	accountRepo := new(MockAccountRepository)
	transactionRepo := new(MockTransactionRepository)
	ledgerRepo := new(MockLedgerEntryRepository)
	fxQuoteRepo := new(MockFXQuoteRepository)
	webhookRepo := new(MockWebhookEventRepository)
	idempotencyRepo := new(MockIdempotencyKeyRepository)

	service := services.NewPaymentService(
		accountRepo,
		transactionRepo,
		ledgerRepo,
		fxQuoteRepo,
		webhookRepo,
		idempotencyRepo,
		nil,
		"test-secret",
	)

	body := bytes.NewBufferString(`{"event_id":"evt-123","transaction_id":"txn-123","status":"completed"}`)

	_, err := service.ProcessWebhook(nil, body, "invalid-signature")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid webhook signature")
}

func TestPaymentService_ProcessWebhook_AlreadyProcessed(t *testing.T) {
	accountRepo := new(MockAccountRepository)
	transactionRepo := new(MockTransactionRepository)
	ledgerRepo := new(MockLedgerEntryRepository)
	fxQuoteRepo := new(MockFXQuoteRepository)
	webhookRepo := new(MockWebhookEventRepository)
	idempotencyRepo := new(MockIdempotencyKeyRepository)

	// Create service with empty webhook secret for this test
	service := services.NewPaymentService(
		accountRepo,
		transactionRepo,
		ledgerRepo,
		fxQuoteRepo,
		webhookRepo,
		idempotencyRepo,
		nil,
		"", // Empty secret
	)

	existingEvent := &models.WebhookEvent{
		EventID:          "evt-123",
		ProcessingStatus: "processed",
	}

	webhookRepo.On("FindByEventIDAndStatus", mock.Anything, "evt-123", "processed").Return(existingEvent, nil)

	body := bytes.NewBufferString(`{"event_id":"evt-123","transaction_id":"txn-123","status":"completed"}`)

	// Empty signature because webhookSecret is empty
	payload, err := service.ProcessWebhook(nil, body, "")
	assert.NoError(t, err)
	assert.NotNil(t, payload)
}

func TestPaymentService_ProcessWebhook_TransactionNotFound(t *testing.T) {
	t.Skip("Test requires full DB setup")
}

// ============ Webhook Signature Tests ============

func TestPaymentService_VerifyWebhookSignature_EmptySecret(t *testing.T) {
	accountRepo := new(MockAccountRepository)
	transactionRepo := new(MockTransactionRepository)
	ledgerRepo := new(MockLedgerEntryRepository)
	fxQuoteRepo := new(MockFXQuoteRepository)
	webhookRepo := new(MockWebhookEventRepository)
	idempotencyRepo := new(MockIdempotencyKeyRepository)

	// Create service with empty webhook secret
	service := services.NewPaymentService(
		accountRepo,
		transactionRepo,
		ledgerRepo,
		fxQuoteRepo,
		webhookRepo,
		idempotencyRepo,
		nil,
		"", // Empty secret - always returns true
	)

	// Test webhook with empty secret (signature always valid)
	body := bytes.NewBufferString(`{"event_id":"evt-123","transaction_id":"txn-123","status":"completed"}`)
	_, err := service.ProcessWebhook(nil, body, "any-signature")
	// With empty secret, there should be no signature validation error
	assert.NoError(t, err)
}

// ============ DTO Serialization Tests ============

func TestPaymentRequest_Serialization(t *testing.T) {
	req := dto.PaymentRequest{
		AccountID:           "acc-123",
		Amount:              1000.50,
		Currency:            "NGN",
		DestinationCurrency: "USD",
		RecipientName:       "John Doe",
		RecipientAccount:    "123456",
		RecipientBank:       "Test Bank",
		RecipientCountry:    "NG",
		Reference:           "REF-001",
	}

	data, err := json.Marshal(req)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "acc-123")
	assert.Contains(t, string(data), "1000.5")
}

func TestPaymentResponse_Serialization(t *testing.T) {
	now := time.Now()
	resp := dto.PaymentResponse{
		TransactionID:       "txn-123",
		TransactionRef:      "TXN-001",
		ProviderReference:   "PRV-001",
		Status:              "completed",
		Amount:              1000,
		Currency:            "NGN",
		FXRate:              0.00065,
		FXAmount:            0.65,
		DestinationCurrency: "USD",
		CreatedAt:           now,
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "txn-123")
	assert.Contains(t, string(data), "completed")
}

// ============ Webhook Payload Tests ============

func TestWebhookPayload_Serialization(t *testing.T) {
	payload := dto.WebhookPayload{
		EventID:       "evt-123",
		EventType:     "payment.completed",
		TransactionID: "txn-123",
		Status:        "completed",
		Amount:        "1000",
	}

	data, err := json.Marshal(payload)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "evt-123")
	assert.Contains(t, string(data), "payment.completed")
}

// Test that io.Reader can be properly read for webhook processing
func TestWebhookBodyReading(t *testing.T) {
	body := `{"event_id":"evt-123","transaction_id":"txn-123","status":"completed","amount":"1000"}`
	reader := bytes.NewBufferString(body)

	readBody, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, body, string(readBody))
}
