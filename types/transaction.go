package types

type CreateTransactionRequest struct {
	Amount int64       `json:"amount"`
	Type   PaymentType `json:"type"`
}

type CreateTransactionResponse struct {
	URL         string              `json:"url"`
	Transaction *PaymentTransaction `json:"transaction"`
}
