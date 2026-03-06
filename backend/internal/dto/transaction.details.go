package dto

import (
	"take-Home-assignment/internal/models"
)

type TransactionDetails struct {
	Transaction   *models.PaymentTransaction `json:"transaction"`
	LedgerEntries []models.LedgerEntry       `json:"ledger_entries"`
	Timeline      []StatusChange             `json:"timeline"`
}

type StatusChange struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}
