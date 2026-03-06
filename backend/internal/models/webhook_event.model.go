package models

import (
	"time"

	"gorm.io/gorm"
)

type WebhookEvent struct {
	ID               string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Source           string     `gorm:"size:100;not null;index" json:"source"`
	EventType        string     `gorm:"size:100;not null;index" json:"event_type"`
	EventID          string     `gorm:"uniqueIndex;size:255" json:"event_id"`
	Payload          *string    `gorm:"type:jsonb" json:"payload"`
	Headers          *string    `gorm:"type:jsonb" json:"headers"`
	ProcessingStatus string     `gorm:"size:50;not null;default:'received';index" json:"processing_status"`
	ProcessingError  string     `gorm:"type:text" json:"processing_error"`
	ProcessedAt      *time.Time `json:"processed_at"`
	CreatedAt        time.Time  `gorm:"default:CURRENT_TIMESTAMP;index" json:"created_at"`
}

func (WebhookEvent) TableName() string { return "webhook_events" }

func (w *WebhookEvent) BeforeCreate(tx *gorm.DB) error {
	w.CreatedAt = time.Now().UTC()
	return nil
}
