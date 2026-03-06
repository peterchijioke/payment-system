package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountType string

const (
	AccountTypeInternal   AccountType = "internal"
	AccountTypeExternal   AccountType = "external"
	AccountTypeSettlement AccountType = "settlement"
	AccountTypeReserve    AccountType = "reserve"
	AccountTypeFee        AccountType = "fee"
)

type Account struct {
	ID             string           `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	AccountNumber  string           `gorm:"uniqueIndex;size:50" json:"account_number"`
	AccountType    AccountType      `gorm:"type:varchar(20);not null" json:"account_type"`
	OwnerID        *uuid.UUID       `gorm:"type:uuid" json:"owner_id"`
	OwnerType      string           `gorm:"size:50" json:"owner_type"`
	Name           string           `gorm:"size:255;not null" json:"name"`
	Description    string           `gorm:"type:text" json:"description"`
	IsActive       bool             `gorm:"default:true" json:"is_active"`
	IsVerified     bool             `gorm:"default:false" json:"is_verified"`
	DailyLimit     float64          `json:"daily_limit"`
	DailyLimitCurr string           `gorm:"size:3" json:"daily_limit_currency"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	ClosedAt       *time.Time       `json:"closed_at"`
	Balances       []AccountBalance `gorm:"foreignKey:AccountID" json:"balances,omitempty"`
}

func (Account) TableName() string { return "accounts" }

func (a *Account) BeforeCreate(tx *gorm.DB) error {
	return nil
}
