package models

import (
	"time"

	"gorm.io/gorm"
)

type FXQuote struct {
	ID           string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	FromCurrency string     `gorm:"size:3;not null" json:"from_currency"`
	ToCurrency   string     `gorm:"size:3;not null" json:"to_currency"`
	Rate         float64    `gorm:"type:decimal(38,12);not null" json:"rate"`
	ValidFrom    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"valid_from"`
	ValidUntil   time.Time  `gorm:"not null" json:"valid_until"`
	QuoteID      string     `gorm:"uniqueIndex;size:100" json:"quote_id"`
	IsLocked     bool       `gorm:"default:false" json:"is_locked"`
	LockedAt     *time.Time `json:"locked_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (FXQuote) TableName() string { return "fx_quotes" }

func (f *FXQuote) BeforeCreate(tx *gorm.DB) error {
	f.ValidFrom = time.Now().UTC()
	return nil
}

func (f *FXQuote) IsExpired() bool {
	return time.Now().UTC().After(f.ValidUntil)
}

func (f *FXQuote) IsValid() bool {
	return !f.IsExpired() && !f.IsLocked
}
