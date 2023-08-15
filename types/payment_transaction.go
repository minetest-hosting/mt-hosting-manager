package types

func PaymentTransactionProvider() *PaymentTransaction { return &PaymentTransaction{} }

type PaymentTransaction struct {
	ID            string `json:"id"`
	TransactionID string `json:"transaction_id"`
	Created       int64  `json:"created"`
	NodeTypeID    string `json:"node_type_id"`
	NodeID        string `json:"node_id"`
	Months        int    `json:"months"`
	State         string `json:"state"`
}

func (m *PaymentTransaction) Columns(action string) []string {
	return []string{
		"id",
		"transaction_id",
		"created",
		"node_type_id",
		"node_id",
		"months",
		"state",
	}
}

func (m *PaymentTransaction) Table() string {
	return "payment_transaction"
}

func (m *PaymentTransaction) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.TransactionID, &m.Created, &m.NodeTypeID, &m.NodeID, &m.Months, &m.State)
}

func (m *PaymentTransaction) Values(action string) []any {
	return []any{m.ID, m.TransactionID, m.Created, m.NodeTypeID, m.NodeID, m.Months, m.State}
}
