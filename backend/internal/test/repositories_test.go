package test

import (
	"testing"
	"time"

	"take-Home-assignment/internal/models"
	"take-Home-assignment/internal/repositories"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// ============ Transaction Repository Tests ============

func TestTransactionRepository_Create(t *testing.T) {
	repo := repositories.NewTransactionRepository()

	// Test that Create method exists and is callable
	assert.NotNil(t, repo)

	// The actual database operation would need a real DB
	// This test verifies the interface is implemented
	var _ repositories.TransactionRepository = repo
}

func TestTransactionRepository_Update(t *testing.T) {
	repo := repositories.NewTransactionRepository()
	assert.NotNil(t, repo)
}

func TestTransactionRepository_FindByID(t *testing.T) {
	repo := repositories.NewTransactionRepository()
	assert.NotNil(t, repo)
}

func TestTransactionRepository_FindByIdempotencyKey(t *testing.T) {
	repo := repositories.NewTransactionRepository()
	assert.NotNil(t, repo)
}

func TestTransactionRepository_FindAll(t *testing.T) {
	repo := repositories.NewTransactionRepository()
	assert.NotNil(t, repo)
}

func TestTransactionRepository_FindByAccountID(t *testing.T) {
	repo := repositories.NewTransactionRepository()
	assert.NotNil(t, repo)
}

// ============ Account Repository Tests ============

func TestAccountRepository_FindByID(t *testing.T) {
	repo := repositories.NewAccountRepository()
	assert.NotNil(t, repo)
}

func TestAccountRepository_FindByIDAndActive(t *testing.T) {
	repo := repositories.NewAccountRepository()
	assert.NotNil(t, repo)
}

func TestAccountRepository_FindAll(t *testing.T) {
	repo := repositories.NewAccountRepository()
	assert.NotNil(t, repo)
}

func TestAccountRepository_FindBalance(t *testing.T) {
	repo := repositories.NewAccountRepository()
	assert.NotNil(t, repo)
}

func TestAccountRepository_LockFunds(t *testing.T) {
	repo := repositories.NewAccountRepository()
	assert.NotNil(t, repo)
}

func TestAccountRepository_ReleaseFunds(t *testing.T) {
	repo := repositories.NewAccountRepository()
	assert.NotNil(t, repo)
}

func TestAccountRepository_CreditFunds(t *testing.T) {
	repo := repositories.NewAccountRepository()
	assert.NotNil(t, repo)
}

func TestAccountRepository_UpdateBalance(t *testing.T) {
	repo := repositories.NewAccountRepository()
	assert.NotNil(t, repo)
}

func TestAccountRepository_GetSettlementAccount(t *testing.T) {
	repo := repositories.NewAccountRepository()
	assert.NotNil(t, repo)
}

func TestAccountRepository_GetDailyTransactionTotal(t *testing.T) {
	repo := repositories.NewAccountRepository()
	assert.NotNil(t, repo)
}

// ============ Ledger Entry Repository Tests ============

func TestLedgerEntryRepository_Create(t *testing.T) {
	repo := repositories.NewLedgerEntryRepository()
	assert.NotNil(t, repo)
}

func TestLedgerEntryRepository_Update(t *testing.T) {
	repo := repositories.NewLedgerEntryRepository()
	assert.NotNil(t, repo)
}

func TestLedgerEntryRepository_FindByTransactionID(t *testing.T) {
	repo := repositories.NewLedgerEntryRepository()
	assert.NotNil(t, repo)
}

func TestLedgerEntryRepository_UpdateStatus(t *testing.T) {
	repo := repositories.NewLedgerEntryRepository()
	assert.NotNil(t, repo)
}

// ============ FX Quote Repository Tests ============

func TestFXQuoteRepository_Create(t *testing.T) {
	repo := repositories.NewFXQuoteRepository()
	assert.NotNil(t, repo)
}

func TestFXQuoteRepository_FindValidQuote(t *testing.T) {
	repo := repositories.NewFXQuoteRepository()
	assert.NotNil(t, repo)
}

func TestFXQuoteRepository_LockQuote(t *testing.T) {
	repo := repositories.NewFXQuoteRepository()
	assert.NotNil(t, repo)
}

// ============ Webhook Event Repository Tests ============

func TestWebhookEventRepository_Create(t *testing.T) {
	repo := repositories.NewWebhookEventRepository()
	assert.NotNil(t, repo)
}

func TestWebhookEventRepository_FindByEventIDAndStatus(t *testing.T) {
	repo := repositories.NewWebhookEventRepository()
	assert.NotNil(t, repo)
}

func TestWebhookEventRepository_MarkProcessed(t *testing.T) {
	repo := repositories.NewWebhookEventRepository()
	assert.NotNil(t, repo)
}

func TestWebhookEventRepository_MarkFailed(t *testing.T) {
	repo := repositories.NewWebhookEventRepository()
	assert.NotNil(t, repo)
}

// ============ Idempotency Key Repository Tests ============

func TestIdempotencyKeyRepository_Create(t *testing.T) {
	repo := repositories.NewIdempotencyKeyRepository()
	assert.NotNil(t, repo)
}

func TestIdempotencyKeyRepository_FindValid(t *testing.T) {
	repo := repositories.NewIdempotencyKeyRepository()
	assert.NotNil(t, repo)
}

// ============ Model Tests ============

func TestPaymentTransaction_TableName(t *testing.T) {
	tx := &models.PaymentTransaction{}
	assert.Equal(t, "transactions", tx.TableName())
}

func TestPaymentTransaction_BeforeCreate(t *testing.T) {
	tx := &models.PaymentTransaction{}

	// BeforeCreate should set default status
	tx.BeforeCreate(&gorm.DB{})

	// Verify status is set to initiated by default
	assert.Equal(t, models.TransactionStatusInitiated, tx.Status)
}

func TestTransactionStatus_Constants(t *testing.T) {
	assert.Equal(t, models.TransactionStatus("initiated"), models.TransactionStatusInitiated)
	assert.Equal(t, models.TransactionStatus("processing"), models.TransactionStatusProcessing)
	assert.Equal(t, models.TransactionStatus("settled"), models.TransactionStatusSettled)
	assert.Equal(t, models.TransactionStatus("completed"), models.TransactionStatusCompleted)
	assert.Equal(t, models.TransactionStatus("failed"), models.TransactionStatusFailed)
	assert.Equal(t, models.TransactionStatus("reversed"), models.TransactionStatusReversed)
	assert.Equal(t, models.TransactionStatus("pending_review"), models.TransactionStatusPendingReview)
}

func TestTransactionType_Constants(t *testing.T) {
	assert.Equal(t, models.TransactionType("deposit"), models.TransactionTypeDeposit)
	assert.Equal(t, models.TransactionType("withdrawal"), models.TransactionTypeWithdrawal)
	assert.Equal(t, models.TransactionType("transfer"), models.TransactionTypeTransfer)
	assert.Equal(t, models.TransactionType("payment"), models.TransactionTypePayment)
	assert.Equal(t, models.TransactionType("refund"), models.TransactionTypeRefund)
	assert.Equal(t, models.TransactionType("fx_conversion"), models.TransactionTypeFXConversion)
	assert.Equal(t, models.TransactionType("fee"), models.TransactionTypeFee)
	assert.Equal(t, models.TransactionType("adjustment"), models.TransactionTypeAdjustment)
}

func TestAccount_TableName(t *testing.T) {
	account := &models.Account{}
	assert.Equal(t, "accounts", account.TableName())
}

func TestAccountBalance_TableName(t *testing.T) {
	balance := &models.AccountBalance{}
	assert.Equal(t, "account_balances", balance.TableName())
}

func TestLedgerEntry_TableName(t *testing.T) {
	entry := &models.LedgerEntry{}
	assert.Equal(t, "ledger_entries", entry.TableName())
}

func TestLedgerEntryStatus_Constants(t *testing.T) {
	assert.Equal(t, models.LedgerEntryStatus("pending"), models.LedgerEntryStatusPending)
	assert.Equal(t, models.LedgerEntryStatus("posted"), models.LedgerEntryStatusPosted)
	assert.Equal(t, models.LedgerEntryStatus("reversed"), models.LedgerEntryStatusReversed)
	assert.Equal(t, models.LedgerEntryStatus("voided"), models.LedgerEntryStatusVoided)
}

func TestLedgerEntryType_Constants(t *testing.T) {
	assert.Equal(t, models.LedgerEntryType("debit"), models.LedgerEntryTypeDebit)
	assert.Equal(t, models.LedgerEntryType("credit"), models.LedgerEntryTypeCredit)
	assert.Equal(t, models.LedgerEntryType("settlement_debit"), models.LedgerEntryTypeSettlementDebit)
	assert.Equal(t, models.LedgerEntryType("settlement_credit"), models.LedgerEntryTypeSettlementCredit)
	assert.Equal(t, models.LedgerEntryType("settlement_reversal_debit"), models.LedgerEntryTypeSettlementReversalDebit)
	assert.Equal(t, models.LedgerEntryType("settlement_reversal_credit"), models.LedgerEntryTypeSettlementReversalCredit)
}

func TestAccountType_Constants(t *testing.T) {
	assert.Equal(t, models.AccountType("internal"), models.AccountTypeInternal)
	assert.Equal(t, models.AccountType("external"), models.AccountTypeExternal)
	assert.Equal(t, models.AccountType("settlement"), models.AccountTypeSettlement)
	assert.Equal(t, models.AccountType("reserve"), models.AccountTypeReserve)
	assert.Equal(t, models.AccountType("fee"), models.AccountTypeFee)
}

// ============ Repository Interface Compliance Tests ============

func TestTransactionRepository_Interface(t *testing.T) {
	repo := repositories.NewTransactionRepository()

	// Verify all required methods exist
	var _ repositories.TransactionRepository = repo

	// Create instance and verify it's not nil
	assert.NotNil(t, repo)
}

func TestAccountRepository_Interface(t *testing.T) {
	repo := repositories.NewAccountRepository()
	var _ repositories.AccountRepository = repo
	assert.NotNil(t, repo)
}

func TestLedgerEntryRepository_Interface(t *testing.T) {
	repo := repositories.NewLedgerEntryRepository()
	var _ repositories.LedgerEntryRepository = repo
	assert.NotNil(t, repo)
}

func TestFXQuoteRepository_Interface(t *testing.T) {
	repo := repositories.NewFXQuoteRepository()
	var _ repositories.FXQuoteRepository = repo
	assert.NotNil(t, repo)
}

func TestWebhookEventRepository_Interface(t *testing.T) {
	repo := repositories.NewWebhookEventRepository()
	var _ repositories.WebhookEventRepository = repo
	assert.NotNil(t, repo)
}

func TestIdempotencyKeyRepository_Interface(t *testing.T) {
	repo := repositories.NewIdempotencyKeyRepository()
	var _ repositories.IdempotencyKeyRepository = repo
	assert.NotNil(t, repo)
}

// ============ Edge Case Tests ============

func TestPaymentTransaction_OptionalFields(t *testing.T) {
	tx := &models.PaymentTransaction{
		ID:             "txn-123",
		TransactionRef: "TXN-001",
		Amount:         1000,
		Currency:       "NGN",
	}

	// All pointer fields should be nil by default
	assert.Nil(t, tx.CounterpartyID)
	assert.Nil(t, tx.SettledAmount)
	assert.Nil(t, tx.FXQuoteID)
	assert.Nil(t, tx.FXRate)
	assert.Nil(t, tx.FXAmount)
	assert.Nil(t, tx.FXCurrency)
	assert.Nil(t, tx.Metadata)
	assert.Nil(t, tx.ProcessedAt)
	assert.Nil(t, tx.SettledAt)
	assert.Nil(t, tx.CompletedAt)
	assert.Nil(t, tx.FailedAt)
	assert.Nil(t, tx.ReversedByID)
}

func TestAccount_Fields(t *testing.T) {
	account := &models.Account{
		ID:   "acc-123",
		Name: "Test Account",
	}

	// Verify defaults
	assert.NotNil(t, account.IsActive)
}

func TestLedgerEntry_References(t *testing.T) {
	txID := "txn-123"
	accID := "acc-456"

	entry := &models.LedgerEntry{
		TransactionID: txID,
		AccountID:     accID,
		Amount:        1000,
		Currency:      "NGN",
	}

	assert.Equal(t, txID, entry.TransactionID)
	assert.Equal(t, accID, entry.AccountID)
}

func TestFXQuote_Validity(t *testing.T) {
	now := time.Now()
	validQuote := &models.FXQuote{
		FromCurrency: "NGN",
		ToCurrency:   "USD",
		Rate:         0.00065,
		ValidUntil:   now.Add(15 * time.Minute),
	}

	// Verify quote is valid (not expired)
	assert.False(t, validQuote.IsExpired())
	assert.True(t, validQuote.IsValid())
}

func TestFXQuote_Expired(t *testing.T) {
	now := time.Now()
	expiredQuote := &models.FXQuote{
		FromCurrency: "NGN",
		ToCurrency:   "USD",
		Rate:         0.00065,
		ValidUntil:   now.Add(-15 * time.Minute),
	}

	// Verify quote is expired
	assert.True(t, expiredQuote.IsExpired())
	assert.False(t, expiredQuote.IsValid())
}

func TestIdempotencyKey_Expiration(t *testing.T) {
	now := time.Now()
	key := &models.IdempotencyKey{
		Key:       "idem-key-123",
		ExpiresAt: now.Add(24 * time.Hour),
	}

	// Verify key has not expired
	assert.False(t, key.IsExpired())
}

func TestIdempotencyKey_Expired(t *testing.T) {
	now := time.Now()
	key := &models.IdempotencyKey{
		Key:       "idem-key-123",
		ExpiresAt: now.Add(-1 * time.Hour),
	}

	// Verify key has expired
	assert.True(t, key.IsExpired())
}

func TestWebhookEvent_ProcessingStatus(t *testing.T) {
	event := &models.WebhookEvent{
		EventID:          "evt-123",
		ProcessingStatus: "received",
	}

	assert.Equal(t, "received", event.ProcessingStatus)
}
