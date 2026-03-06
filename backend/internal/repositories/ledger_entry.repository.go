package repositories

import (
	"take-Home-assignment/internal/models"

	"gorm.io/gorm"
)

type LedgerEntryRepository interface {
	Create(tx *gorm.DB, entry *models.LedgerEntry) error
	Update(tx *gorm.DB, entry *models.LedgerEntry) error
	UpdateStatus(tx *gorm.DB, transactionID string, status models.LedgerEntryStatus, postedAt interface{}, reversedByID interface{}) error
	FindByTransactionID(tx *gorm.DB, transactionID string) ([]models.LedgerEntry, error)
}

type ledgerEntryRepository struct{}

func NewLedgerEntryRepository() LedgerEntryRepository {
	return &ledgerEntryRepository{}
}

func (r *ledgerEntryRepository) Create(tx *gorm.DB, entry *models.LedgerEntry) error {
	return tx.Create(entry).Error
}

func (r *ledgerEntryRepository) Update(tx *gorm.DB, entry *models.LedgerEntry) error {
	return tx.Save(entry).Error
}

func (r *ledgerEntryRepository) UpdateStatus(tx *gorm.DB, transactionID string, status models.LedgerEntryStatus, postedAt interface{}, reversedByID interface{}) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if postedAt != nil {
		updates["posted_at"] = postedAt
	}
	if reversedByID != nil {
		updates["reversed_by_id"] = reversedByID
	}
	return tx.Model(&models.LedgerEntry{}).
		Where("transaction_id = ?", transactionID).
		Updates(updates).Error
}

func (r *ledgerEntryRepository) FindByTransactionID(tx *gorm.DB, transactionID string) ([]models.LedgerEntry, error) {
	var entries []models.LedgerEntry
	err := tx.Where("transaction_id = ?", transactionID).Find(&entries).Error
	return entries, err
}
