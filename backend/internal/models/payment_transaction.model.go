package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionStatus string

const (
	TransactionStatusInitiated     TransactionStatus = "initiated"
	TransactionStatusProcessing    TransactionStatus = "processing"
	TransactionStatusSettled       TransactionStatus = "settled"
	TransactionStatusCompleted     TransactionStatus = "completed"
	TransactionStatusFailed        TransactionStatus = "failed"
	TransactionStatusReversed      TransactionStatus = "reversed"
	TransactionStatusPendingReview TransactionStatus = "pending_review"
)

type TransactionType string

const (
	TransactionTypeDeposit      TransactionType = "deposit"
	TransactionTypeWithdrawal   TransactionType = "withdrawal"
	TransactionTypeTransfer     TransactionType = "transfer"
	TransactionTypePayment      TransactionType = "payment"
	TransactionTypeRefund       TransactionType = "refund"
	TransactionTypeFXConversion TransactionType = "fx_conversion"
	TransactionTypeFee          TransactionType = "fee"
	TransactionTypeAdjustment   TransactionType = "adjustment"
)

type PaymentTransaction struct {
	ID             string            `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TransactionRef string            `gorm:"uniqueIndex;size:100" json:"transaction_reference"`
	IdempotencyKey string            `gorm:"index;size:255" json:"idempotency_key"`
	AccountID      string            `gorm:"type:uuid;not null;index" json:"account_id"`
	CounterpartyID *string           `gorm:"type:uuid" json:"counterparty_id"`
	Type           TransactionType   `gorm:"type:varchar(20);not null" json:"type"`
	Status         TransactionStatus `gorm:"type:varchar(20);not null;default:'initiated'" json:"status"`
	Amount         float64           `gorm:"type:decimal(38,12);not null" json:"amount"`
	Currency       string            `gorm:"size:3;not null" json:"currency"`
	SettledAmount  *float64          `gorm:"type:decimal(38,12)" json:"settled_amount"`
	FXQuoteID      *string           `gorm:"type:uuid" json:"fx_quote_id"`
	FXRate         *float64          `gorm:"type:decimal(38,12)" json:"fx_rate"`
	FXAmount       *float64          `gorm:"type:decimal(38,12)" json:"fx_amount"`
	FXCurrency     *string           `gorm:"size:3" json:"fx_currency"`
	Description    string            `gorm:"type:text" json:"description"`
	Reference      string            `gorm:"size:255" json:"reference"`
	Metadata       *string           `gorm:"type:jsonb" json:"metadata"`
	InitiatedAt    time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"initiated_at"`
	ProcessedAt    *time.Time        `json:"processed_at"`
	SettledAt      *time.Time        `json:"settled_at"`
	CompletedAt    *time.Time        `json:"completed_at"`
	FailedAt       *time.Time        `json:"failed_at"`
	FailureReason  string            `gorm:"type:text" json:"failure_reason"`
	ReversalReason string            `gorm:"type:text" json:"reversal_reason"`
	ReversedByID   *string           `gorm:"type:uuid" json:"reversed_by_id"`
	Version        int               `gorm:"default:0" json:"version"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

func (PaymentTransaction) TableName() string { return "transactions" }

func (p *PaymentTransaction) BeforeCreate(tx *gorm.DB) error {
	if p.Status == "" {
		p.Status = TransactionStatusInitiated
	}
	p.InitiatedAt = time.Now().UTC()
	return nil
}
