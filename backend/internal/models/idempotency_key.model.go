package models

import (
	"time"
)

type IdempotencyKey struct {
	Key              string    `gorm:"primaryKey;size:255" json:"key"`
	AccountID        string    `gorm:"type:uuid;not null;index" json:"account_id"`
	RequestHash      string    `gorm:"size:64;not null" json:"request_hash"`
	RequestMethod    string    `gorm:"size:10;not null" json:"request_method"`
	RequestPath      string    `gorm:"size:500;not null" json:"request_path"`
	OriginalAmount   float64   `gorm:"type:decimal(38,12);not null" json:"original_amount"`
	OriginalCurrency string    `gorm:"size:3;not null" json:"original_currency"`
	ResponseStatus   *int      `json:"response_status"`
	ResponseBody     string    `gorm:"type:jsonb" json:"response_body"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	ExpiresAt        time.Time `gorm:"not null;index" json:"expires_at"`
}

func (IdempotencyKey) TableName() string { return "idempotency_keys" }

func (i *IdempotencyKey) IsExpired() bool {
	return time.Now().UTC().After(i.ExpiresAt)
}
