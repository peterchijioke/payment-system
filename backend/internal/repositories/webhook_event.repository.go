package repositories

import (
	"take-Home-assignment/internal/models"
	"time"

	"gorm.io/gorm"
)

type WebhookEventRepository interface {
	Create(tx *gorm.DB, event *models.WebhookEvent) error
	Update(tx *gorm.DB, event *models.WebhookEvent) error
	FindByEventID(tx *gorm.DB, eventID string) (*models.WebhookEvent, error)
	FindByEventIDAndStatus(tx *gorm.DB, eventID, status string) (*models.WebhookEvent, error)
	MarkProcessed(tx *gorm.DB, eventID string) error
	MarkFailed(tx *gorm.DB, eventID, errorMsg string) error
}

type webhookEventRepository struct{}

func NewWebhookEventRepository() WebhookEventRepository {
	return &webhookEventRepository{}
}

func (r *webhookEventRepository) Create(tx *gorm.DB, event *models.WebhookEvent) error {
	return tx.Create(event).Error
}

func (r *webhookEventRepository) Update(tx *gorm.DB, event *models.WebhookEvent) error {
	return tx.Save(event).Error
}

func (r *webhookEventRepository) FindByEventID(tx *gorm.DB, eventID string) (*models.WebhookEvent, error) {
	var event models.WebhookEvent
	err := tx.Where("event_id = ?", eventID).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *webhookEventRepository) FindByEventIDAndStatus(tx *gorm.DB, eventID, status string) (*models.WebhookEvent, error) {
	var event models.WebhookEvent
	err := tx.Where("event_id = ? AND processing_status = ?", eventID, status).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *webhookEventRepository) MarkProcessed(tx *gorm.DB, eventID string) error {
	now := time.Now().UTC()
	return tx.Model(&models.WebhookEvent{}).
		Where("event_id = ?", eventID).
		Updates(map[string]interface{}{
			"processing_status": "processed",
			"processed_at":      now,
		}).Error
}

func (r *webhookEventRepository) MarkFailed(tx *gorm.DB, eventID, errorMsg string) error {
	return tx.Model(&models.WebhookEvent{}).
		Where("event_id = ?", eventID).
		Updates(map[string]interface{}{
			"processing_status": "failed",
			"processing_error":  errorMsg,
		}).Error
}
