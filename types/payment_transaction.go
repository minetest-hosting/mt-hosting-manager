package types

type PaymentTransaction struct {
	ID            string `json:"id"`
	TransactionID string `json:"transaction_id"`
	Created       int64  `json:"created"`
	UserNodeID    string `json:"user_node_id"`
}

func (m *PaymentTransaction) Columns(action string) []string {
	return []string{
		"id",
		"transaction_id",
		"created",
		"user_node_id",
	}
}

func (m *PaymentTransaction) Table() string {
	return "payment_transaction"
}

func (m *PaymentTransaction) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.TransactionID, &m.Created, &m.UserNodeID)
}

func (m *PaymentTransaction) Values(action string) []any {
	return []any{m.ID, m.TransactionID, m.Created, m.UserNodeID}
}
