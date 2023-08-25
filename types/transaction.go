package types

type CreateTransactionRequest struct {
	Amount string `json:"amount"`
}

type CreateTransactionResponse struct {
	URL string `json:"url"`
}
