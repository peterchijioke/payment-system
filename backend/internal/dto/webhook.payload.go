package dto

import (
	"encoding/json"
	"time"
)

type WebhookPayload struct {
	EventType     string          `json:"event_type" binding:"required"`
	EventID       string          `json:"event_id" binding:"required,uuid"`
	TransactionID string          `json:"transaction_id" binding:"required,uuid"`
	Status        string          `json:"status" binding:"required"`
	Amount        string          `json:"amount,omitempty"`
	Currency      string          `json:"currency,omitempty" binding:"omitempty,len=3"`
	FailureReason string          `json:"failure_reason,omitempty"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	Timestamp     time.Time       `json:"timestamp" binding:"required"`
}
