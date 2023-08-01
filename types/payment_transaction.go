package types

type PaymentTransaction struct {
	ID            string `json:"id"`
	TransactionID string `json:"transaction_id"`
	Created       int64  `json:"created"`
	NodeTypeID    string `json:"node_type_id"`
	Months        int    `json:"months"`
}

func (m *PaymentTransaction) Columns(action string) []string {
	return []string{
		"id",
		"transaction_id",
		"created",
		"node_type_id",
		"months",
	}
}

func (m *PaymentTransaction) Table() string {
	return "payment_transaction"
}

func (m *PaymentTransaction) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.TransactionID, &m.Created, &m.NodeTypeID, &m.Months)
}

func (m *PaymentTransaction) Values(action string) []any {
	return []any{m.ID, m.TransactionID, m.Created, m.NodeTypeID, m.Months}
}
