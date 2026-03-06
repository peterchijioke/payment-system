package models

import (
	"time"

	"gorm.io/gorm"
)

type LedgerEntryType string

const (
	LedgerEntryTypeDebit                    LedgerEntryType = "debit"
	LedgerEntryTypeCredit                   LedgerEntryType = "credit"
	LedgerEntryTypeSettlementDebit          LedgerEntryType = "settlement_debit"
	LedgerEntryTypeSettlementCredit         LedgerEntryType = "settlement_credit"
	LedgerEntryTypeSettlementReversalDebit  LedgerEntryType = "settlement_reversal_debit"
	LedgerEntryTypeSettlementReversalCredit LedgerEntryType = "settlement_reversal_credit"
)

type LedgerEntryStatus string

const (
	LedgerEntryStatusPending  LedgerEntryStatus = "pending"
	LedgerEntryStatusPosted   LedgerEntryStatus = "posted"
	LedgerEntryStatusReversed LedgerEntryStatus = "reversed"
	LedgerEntryStatusVoided   LedgerEntryStatus = "voided"
)

type LedgerEntry struct {
	ID                 string             `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	EntryReference     string             `gorm:"uniqueIndex;size:100" json:"entry_reference"`
	TransactionID      string             `gorm:"type:uuid;not null;index" json:"transaction_id"`
	AccountID          string             `gorm:"type:uuid;not null;index" json:"account_id"`
	EntryType          LedgerEntryType    `gorm:"type:varchar(30);not null" json:"entry_type"`
	Amount             float64            `gorm:"type:decimal(38,12);not null" json:"amount"`
	Currency           string             `gorm:"size:3;not null" json:"currency"`
	CounterpartEntryID *string            `gorm:"type:uuid" json:"counterpart_entry_id"`
	OriginalEntryID    *string            `gorm:"type:uuid" json:"original_entry_id"`
	Status             LedgerEntryStatus  `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	ReversalReason     string             `gorm:"type:text" json:"reversal_reason"`
	Description        string             `gorm:"size:500" json:"description"`
	EffectiveDate      time.Time          `gorm:"type:date;not null" json:"effective_date"`
	PostedAt           *time.Time         `json:"posted_at"`
	ReversedByID       *string            `gorm:"type:uuid" json:"reversed_by_id"`
	CreatedAt          time.Time          `json:"created_at"`
	CreatedBy          *string            `gorm:"type:uuid" json:"created_by"`
	Transaction        PaymentTransaction `gorm:"foreignKey:TransactionID" json:"-"`
	Account            Account            `gorm:"foreignKey:AccountID" json:"-"`
}

func (LedgerEntry) TableName() string { return "ledger_entries" }

func (l *LedgerEntry) BeforeCreate(tx *gorm.DB) error {
	if l.Status == "" {
		l.Status = LedgerEntryStatusPending
	}
	if l.EffectiveDate.IsZero() {
		l.EffectiveDate = time.Now().UTC().Truncate(24 * time.Hour)
	}
	return nil
}
