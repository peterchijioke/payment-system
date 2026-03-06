package repositories

import (
	"take-Home-assignment/internal/models"
	"time"

	"gorm.io/gorm"
)

type IdempotencyKeyRepository interface {
	FindValid(tx *gorm.DB, key, requestHash string, now time.Time) (*models.IdempotencyKey, error)
	Create(tx *gorm.DB, idempotencyKey *models.IdempotencyKey) error
}

type idempotencyKeyRepository struct{}

func NewIdempotencyKeyRepository() IdempotencyKeyRepository {
	return &idempotencyKeyRepository{}
}

func (r *idempotencyKeyRepository) FindValid(tx *gorm.DB, key, requestHash string, now time.Time) (*models.IdempotencyKey, error) {
	var idempotency models.IdempotencyKey
	err := tx.Where("key = ? AND request_hash = ? AND expires_at > ?", key, requestHash, now).
		First(&idempotency).Error
	if err != nil {
		return nil, err
	}
	return &idempotency, nil
}

func (r *idempotencyKeyRepository) Create(tx *gorm.DB, idempotencyKey *models.IdempotencyKey) error {
	return tx.Create(idempotencyKey).Error
}
