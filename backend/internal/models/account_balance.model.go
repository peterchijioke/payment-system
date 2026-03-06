package models

import (
	"time"
)

type AccountBalance struct {
	ID                string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	AccountID         string     `gorm:"type:uuid;not null;index" json:"account_id"`
	Currency          string     `gorm:"size:3;not null" json:"currency"`
	AvailableBalance  float64    `gorm:"type:decimal(38,12);default:0" json:"available_balance"`
	PendingBalance    float64    `gorm:"type:decimal(38,12);default:0" json:"pending_balance"`
	ReservedBalance   float64    `gorm:"type:decimal(38,12);default:0" json:"reserved_balance"`
	TotalCredited     float64    `gorm:"type:decimal(38,12);default:0" json:"total_credited"`
	TotalDebited      float64    `gorm:"type:decimal(38,12);default:0" json:"total_debited"`
	LastTransactionAt *time.Time `json:"last_transaction_at"`
	Version           int        `gorm:"default:0" json:"version"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	Account           Account    `gorm:"foreignKey:AccountID" json:"-"`
}

func (AccountBalance) TableName() string { return "account_balances" }
