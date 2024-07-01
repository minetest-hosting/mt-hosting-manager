package zahlsch

type TransactionStatus string

const (
	TransactionConfirmed TransactionStatus = "confirmed"
)

type CustomField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Invoice struct {
	CustomFields []*CustomField `json:"custom_fields"`
}

type Transaction struct {
	ID      int64             `json:"id"`
	Amount  int64             `json:"amount"`
	Time    string            `json:"time"`
	Status  TransactionStatus `json:"status"`
	PageID  string            `json:"pageUuid"`
	Invoice *Invoice          `json:"invoice"`
}

type WebhookPayload struct {
	Transaction *Transaction `json:"transaction"`
}
