package types

type CreateTransactionRequest struct {
	Amount int64       `json:"amount"`
	Type   PaymentType `json:"type"`
}
