package dto

type PaymentRequest struct {
	AccountID           string  `json:"account_id" binding:"required,uuid"`
	Amount              float64 `json:"amount" binding:"required,gt=0"`
	Currency            string  `json:"currency" binding:"required,len=3"`
	DestinationCurrency string  `json:"destination_currency" binding:"required,len=3"`
	RecipientName       string  `json:"recipient_name" binding:"required"`
	RecipientAccount    string  `json:"recipient_account" binding:"required"`
	RecipientBank       string  `json:"recipient_bank" binding:"required"`
	RecipientCountry    string  `json:"recipient_country" binding:"required"`
	Reference           string  `json:"reference"`
	IdempotencyKey      string  `json:"-"`
}
