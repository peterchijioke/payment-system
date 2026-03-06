package dto

import "time"

type PaymentResponse struct {
	TransactionID       string    `json:"transaction_id"`
	TransactionRef      string    `json:"transaction_reference"`
	ProviderReference   string    `json:"provider_reference"`
	Status              string    `json:"status"`
	Amount              float64   `json:"amount"`
	Currency            string    `json:"currency"`
	FXRate              float64   `json:"fx_rate,omitempty"`
	FXAmount            float64   `json:"fx_amount,omitempty"`
	DestinationCurrency string    `json:"destination_currency,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
}
