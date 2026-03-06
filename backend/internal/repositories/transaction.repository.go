package repositories

import (
	"take-Home-assignment/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionRepository interface {
	Create(tx *gorm.DB, transaction *models.PaymentTransaction) error
	Update(tx *gorm.DB, transaction *models.PaymentTransaction) error
	FindByID(tx *gorm.DB, id string) (*models.PaymentTransaction, error)
	FindByIDForUpdate(tx *gorm.DB, id string) (*models.PaymentTransaction, error)
	FindByIdempotencyKey(tx *gorm.DB, key string) (*models.PaymentTransaction, error)
	FindAll(tx *gorm.DB, limit, offset int, status, startDate, endDate string) ([]models.PaymentTransaction, int64, error)
	FindByAccountID(tx *gorm.DB, accountID string, limit, offset int) ([]models.PaymentTransaction, int64, error)
}

type transactionRepository struct{}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{}
}

func (r *transactionRepository) Create(tx *gorm.DB, transaction *models.PaymentTransaction) error {
	return tx.Create(transaction).Error
}

func (r *transactionRepository) Update(tx *gorm.DB, transaction *models.PaymentTransaction) error {
	return tx.Save(transaction).Error
}

func (r *transactionRepository) FindByID(tx *gorm.DB, id string) (*models.PaymentTransaction, error) {
	var transaction models.PaymentTransaction
	err := tx.Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindByIDForUpdate(tx *gorm.DB, id string) (*models.PaymentTransaction, error) {
	var transaction models.PaymentTransaction
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindByIdempotencyKey(tx *gorm.DB, key string) (*models.PaymentTransaction, error) {
	var transaction models.PaymentTransaction
	err := tx.Where("idempotency_key = ?", key).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindAll(tx *gorm.DB, limit, offset int, status, startDate, endDate string) ([]models.PaymentTransaction, int64, error) {
	var transactions []models.PaymentTransaction
	var count int64

	query := tx.Model(&models.PaymentTransaction{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != "" {
		query = query.Where("initiated_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("initiated_at <= ?", endDate)
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("initiated_at DESC").Limit(limit).Offset(offset).Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, count, nil
}

func (r *transactionRepository) FindByAccountID(tx *gorm.DB, accountID string, limit, offset int) ([]models.PaymentTransaction, int64, error) {
	var transactions []models.PaymentTransaction
	var count int64

	err := tx.Model(&models.PaymentTransaction{}).
		Where("account_id = ?", accountID).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = tx.Where("account_id = ?", accountID).
		Order("initiated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, count, nil
}
