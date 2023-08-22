package types

type CreateTransactionRequest struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type CreateTransactionResponse struct {
	URL string `json:"url"`
}
