package services

import (
	"take-Home-assignment/internal/config"
	"take-Home-assignment/internal/repositories"

	"gorm.io/gorm"
)

func InitServices(db *gorm.DB, cfg *config.ServerConfig) *Container {
	accountRepo := repositories.NewAccountRepository()
	transactionRepo := repositories.NewTransactionRepository()
	ledgerRepo := repositories.NewLedgerEntryRepository()
	fxQuoteRepo := repositories.NewFXQuoteRepository()
	webhookRepo := repositories.NewWebhookEventRepository()
	idempotencyRepo := repositories.NewIdempotencyKeyRepository()

	paymentService := NewPaymentService(
		accountRepo,
		transactionRepo,
		ledgerRepo,
		fxQuoteRepo,
		webhookRepo,
		idempotencyRepo,
		db,
		cfg.WebhookSecret,
	)
	return &Container{
		Payment: paymentService,
	}
}
