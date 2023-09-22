package types

type CreateTransactionRequest struct {
	Amount int64 `json:"amount"`
}

type CreateTransactionResponse struct {
	URL string `json:"url"`
}
