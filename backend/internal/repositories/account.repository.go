package repositories

import (
	"take-Home-assignment/internal/models"

	"gorm.io/gorm"
)

type AccountRepository interface {
	FindByID(tx *gorm.DB, id string) (*models.Account, error)
	FindByIDAndActive(tx *gorm.DB, id string) (*models.Account, error)
	FindAll(tx *gorm.DB) ([]models.Account, error)
	FindBalance(tx *gorm.DB, accountID, currency string) (*models.AccountBalance, error)
	LockFunds(tx *gorm.DB, accountID, currency string, amount float64) error
	ReleaseFunds(tx *gorm.DB, accountID, currency string, amount float64) error
	CreditFunds(tx *gorm.DB, accountID, currency string, amount float64) error
	UpdateBalance(tx *gorm.DB, accountID, currency string, updates map[string]interface{}) error
	GetSettlementAccount(tx *gorm.DB) (*models.Account, error)
	GetDailyTransactionTotal(tx *gorm.DB, accountID, currency string, since interface{}) (float64, error)
}

type accountRepository struct{}

func NewAccountRepository() AccountRepository {
	return &accountRepository{}
}

func (r *accountRepository) FindByID(tx *gorm.DB, id string) (*models.Account, error) {
	var account models.Account
	err := tx.Where("id = ?", id).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) FindByIDAndActive(tx *gorm.DB, id string) (*models.Account, error) {
	var account models.Account
	err := tx.Where("id = ? AND is_active = ?", id, true).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) FindAll(tx *gorm.DB) ([]models.Account, error) {
	var accounts []models.Account
	err := tx.Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) FindBalance(tx *gorm.DB, accountID, currency string) (*models.AccountBalance, error) {
	var balance models.AccountBalance
	err := tx.Where("account_id = ? AND currency = ?", accountID, currency).First(&balance).Error
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

func (r *accountRepository) LockFunds(tx *gorm.DB, accountID, currency string, amount float64) error {
	result := tx.Model(&models.AccountBalance{}).
		Where("account_id = ? AND currency = ? AND available_balance - reserved_balance >= ?",
			accountID, currency, amount).
		Updates(map[string]interface{}{
			"available_balance":   gorm.Expr("available_balance - ?", amount),
			"reserved_balance":    gorm.Expr("reserved_balance + ?", amount),
			"last_transaction_at": tx.NowFunc(),
			"version":             gorm.Expr("version + 1"),
		})

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *accountRepository) ReleaseFunds(tx *gorm.DB, accountID, currency string, amount float64) error {
	return tx.Model(&models.AccountBalance{}).
		Where("account_id = ? AND currency = ?", accountID, currency).
		Updates(map[string]interface{}{
			"available_balance":   gorm.Expr("available_balance + ?", amount),
			"reserved_balance":    gorm.Expr("reserved_balance - ?", amount),
			"last_transaction_at": tx.NowFunc(),
			"version":             gorm.Expr("version + 1"),
		}).Error
}

func (r *accountRepository) CreditFunds(tx *gorm.DB, accountID, currency string, amount float64) error {
	return tx.Model(&models.AccountBalance{}).
		Where("account_id = ? AND currency = ?", accountID, currency).
		Updates(map[string]interface{}{
			"available_balance":   gorm.Expr("available_balance + ?", amount),
			"last_transaction_at": tx.NowFunc(),
			"version":             gorm.Expr("version + 1"),
		}).Error
}

func (r *accountRepository) UpdateBalance(tx *gorm.DB, accountID, currency string, updates map[string]interface{}) error {
	return tx.Model(&models.AccountBalance{}).
		Where("account_id = ? AND currency = ?", accountID, currency).
		Updates(updates).Error
}

func (r *accountRepository) GetSettlementAccount(tx *gorm.DB) (*models.Account, error) {
	var account models.Account
	err := tx.Where("account_type = ?", models.AccountTypeSettlement).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) GetDailyTransactionTotal(tx *gorm.DB, accountID, currency string, since interface{}) (float64, error) {
	var total float64
	err := tx.Model(&models.PaymentTransaction{}).
		Where("account_id = ? AND currency = ? AND initiated_at >= ? AND status IN (?, ?)",
			accountID, currency, since,
			models.TransactionStatusCompleted,
			models.TransactionStatusProcessing).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}
