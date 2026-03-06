package database

import (
	"fmt"
	"log"
	"time"

	"take-Home-assignment/internal/models"

	"gorm.io/gorm"
)

func SeedFinancialData(db *gorm.DB) {
	log.Println("Seeding financial data...")


	settlementAccount := models.Account{
		AccountNumber: "ACC-SETTLEMENT-001",
		AccountType:   models.AccountTypeSettlement,
		Name:          "Settlement Account",
		Description:   "Main settlement account for cross-border payments",
		IsActive:      true,
		IsVerified:    true,
	}

	result := db.Where("account_number = ?", settlementAccount.AccountNumber).First(&models.Account{})
	if result.Error == gorm.ErrRecordNotFound {
		if err := db.Create(&settlementAccount).Error; err != nil {
			log.Printf("Failed to create settlement account: %v", err)
		}
	}


	reserveAccount := models.Account{
		AccountNumber: "ACC-RESERVE-001",
		AccountType:   models.AccountTypeReserve,
		Name:          "Reserve Account",
		Description:   "Reserve account for pending payments",
		IsActive:      true,
		IsVerified:    true,
	}

	result = db.Where("account_number = ?", reserveAccount.AccountNumber).First(&models.Account{})
	if result.Error == gorm.ErrRecordNotFound {
		if err := db.Create(&reserveAccount).Error; err != nil {
			log.Printf("Failed to create reserve account: %v", err)
		}
	}


	internalAccount := models.Account{
		AccountNumber:  "ACC-INTERNAL-001",
		AccountType:    models.AccountTypeInternal,
		Name:           "Operations Account",
		Description:    "Internal operational account",
		IsActive:       true,
		IsVerified:     true,
		DailyLimit:     10000000,
		DailyLimitCurr: "NGN",
	}

	result = db.Where("account_number = ?", internalAccount.AccountNumber).First(&models.Account{})
	if result.Error == gorm.ErrRecordNotFound {
		if err := db.Create(&internalAccount).Error; err != nil {
			log.Printf("Failed to create internal account: %v", err)
		}
	}


	merchantAccount := models.Account{
		AccountNumber:  "ACC-MERCHANT-001",
		AccountType:    models.AccountTypeExternal,
		Name:           "Test Merchant Account",
		Description:    "Test merchant for payments",
		IsActive:       true,
		IsVerified:     true,
		DailyLimit:     5000000,
		DailyLimitCurr: "NGN",
	}

	result = db.Where("account_number = ?", merchantAccount.AccountNumber).First(&models.Account{})
	if result.Error == gorm.ErrRecordNotFound {
		if err := db.Create(&merchantAccount).Error; err != nil {
			log.Printf("Failed to create merchant account: %v", err)
		}
	}


	var accounts []models.Account
	db.Find(&accounts)

	for _, account := range accounts {

		ngnBalance := models.AccountBalance{
			AccountID:         account.ID,
			Currency:          "NGN",
			AvailableBalance:  10000000, // 10 million NGN
			PendingBalance:    0,
			ReservedBalance:   0,
			TotalCredited:     10000000,
			TotalDebited:      0,
			LastTransactionAt: func() *time.Time { t := time.Now().UTC(); return &t }(),
		}

		result := db.Where("account_id = ? AND currency = ?", account.ID, "NGN").First(&models.AccountBalance{})
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&ngnBalance).Error; err != nil {
				log.Printf("Failed to create NGN balance for %s: %v", account.AccountNumber, err)
			}
		}


		usdBalance := models.AccountBalance{
			AccountID:         account.ID,
			Currency:          "USD",
			AvailableBalance:  50000, // 50k USD
			PendingBalance:    0,
			ReservedBalance:   0,
			TotalCredited:     50000,
			TotalDebited:      0,
			LastTransactionAt: func() *time.Time { t := time.Now().UTC(); return &t }(),
		}

		result = db.Where("account_id = ? AND currency = ?", account.ID, "USD").First(&models.AccountBalance{})
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&usdBalance).Error; err != nil {
				log.Printf("Failed to create USD balance for %s: %v", account.AccountNumber, err)
			}
		}


		eurBalance := models.AccountBalance{
			AccountID:         account.ID,
			Currency:          "EUR",
			AvailableBalance:  25000, // 25k EUR
			PendingBalance:    0,
			ReservedBalance:   0,
			TotalCredited:     25000,
			TotalDebited:      0,
			LastTransactionAt: func() *time.Time { t := time.Now().UTC(); return &t }(),
		}

		result = db.Where("account_id = ? AND currency = ?", account.ID, "EUR").First(&models.AccountBalance{})
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&eurBalance).Error; err != nil {
				log.Printf("Failed to create EUR balance for %s: %v", account.AccountNumber, err)
			}
		}
	}


	currencyPairs := []struct {
		From string
		To   string
		Rate float64
	}{
		{"NGN", "USD", 0.00065},
		{"NGN", "EUR", 0.00060},
		{"NGN", "GBP", 0.00052},
		{"USD", "EUR", 0.92},
		{"USD", "GBP", 0.79},
		{"EUR", "GBP", 0.86},
	}

	for _, pair := range currencyPairs {
		quote := models.FXQuote{
			FromCurrency: pair.From,
			ToCurrency:   pair.To,
			Rate:         pair.Rate,
			ValidUntil:   time.Now().UTC().Add(15 * time.Minute),
			QuoteID:      fmt.Sprintf("QUOTE-%d-%s-%s", time.Now().Unix(), pair.From, pair.To),
			IsLocked:     false,
		}

		result := db.Where("quote_id = ?", quote.QuoteID).First(&models.FXQuote{})
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&quote).Error; err != nil {
				log.Printf("Failed to create FX quote: %v", err)
			}
		}
	}

	log.Println("Financial data seeding completed")
}
