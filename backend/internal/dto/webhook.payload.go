package dto

import (
	"encoding/json"
	"time"
)

type WebhookPayload struct {
	EventType     string          `json:"event_type"`
	EventID       string          `json:"event_id"`
	TransactionID string          `json:"transaction_id"`
	Status        string          `json:"status"`
	Amount        string          `json:"amount,omitempty"`
	Currency      string          `json:"currency,omitempty"`
	FailureReason string          `json:"failure_reason,omitempty"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	Timestamp     time.Time       `json:"timestamp"`
}
