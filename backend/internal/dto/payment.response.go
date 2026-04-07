package dto

import "time"

type PaymentResponse struct {
	TransactionID       string    `json:"transaction_id" binding:"required,uuid"`
	TransactionRef      string    `json:"transaction_reference" binding:"required"`
	ProviderReference   string    `json:"provider_reference"`
	Status              string    `json:"status" binding:"required"`
	Amount              float64   `json:"amount" binding:"required,gt=0"`
	Currency            string    `json:"currency" binding:"required,len=3"`
	FXRate              float64   `json:"fx_rate,omitempty" binding:"omitempty,gt=0"`
	FXAmount            float64   `json:"fx_amount,omitempty" binding:"omitempty,gt=0"`
	DestinationCurrency string    `json:"destination_currency,omitempty" binding:"omitempty,len=3"`
	CreatedAt           time.Time `json:"created_at" binding:"required"`
}
