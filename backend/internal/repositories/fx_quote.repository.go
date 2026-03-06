package repositories

import (
	"take-Home-assignment/internal/models"
	"time"

	"gorm.io/gorm"
)

type FXQuoteRepository interface {
	FindValidQuote(tx *gorm.DB, fromCurrency, toCurrency string, now time.Time) (*models.FXQuote, error)
	Create(tx *gorm.DB, quote *models.FXQuote) error
	LockQuote(tx *gorm.DB, quoteID string, lockedAt time.Time) error
}

type fxQuoteRepository struct{}

func NewFXQuoteRepository() FXQuoteRepository {
	return &fxQuoteRepository{}
}

func (r *fxQuoteRepository) FindValidQuote(tx *gorm.DB, fromCurrency, toCurrency string, now time.Time) (*models.FXQuote, error) {
	var quote models.FXQuote
	err := tx.Where("from_currency = ? AND to_currency = ? AND valid_until > ? AND is_locked = ?",
		fromCurrency, toCurrency, now, false).
		Order("created_at DESC").
		First(&quote).Error
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

func (r *fxQuoteRepository) Create(tx *gorm.DB, quote *models.FXQuote) error {
	return tx.Create(quote).Error
}

func (r *fxQuoteRepository) LockQuote(tx *gorm.DB, quoteID string, lockedAt time.Time) error {
	return tx.Model(&models.FXQuote{}).
		Where("id = ?", quoteID).
		Updates(map[string]interface{}{
			"is_locked": true,
			"locked_at": lockedAt,
		}).Error
}
