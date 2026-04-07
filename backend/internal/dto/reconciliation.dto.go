package dto

type ReconciliationRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
	AccountID string `form:"account_id"`
	Status    string `form:"status"`
}

type ReconciliationResult struct {
	TransactionID     string `json:"transaction_id"`
	TransactionRef    string `json:"transaction_reference"`
	Status            string `json:"status"`
	Issue             string `json:"issue,omitempty"`
	LocalUpdatedAt    string `json:"local_updated_at"`
	ProviderStatus    string `json:"provider_status,omitempty"`
	ProviderUpdatedAt string `json:"provider_updated_at,omitempty"`
}

type ReconciliationResponse struct {
	Success           bool                   `json:"success"`
	TotalTransactions int64                  `json:"total_transactions"`
	MismatchedCount   int64                  `json:"mismatched_count"`
	MissingCount      int64                  `json:"missing_count"`
	Results           []ReconciliationResult `json:"results"`
}
