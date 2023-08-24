package types

type PaymentStateType string

const (
	PaymentStatePending PaymentStateType = "PENDING"
	PaymentStateSuccess PaymentStateType = "SUCCESS"
	PaymentStateError   PaymentStateType = "ERROR"
)

func PaymentTransactionProvider() *PaymentTransaction { return &PaymentTransaction{} }

type PaymentTransaction struct {
	ID             string           `json:"id"`
	TransactionID  string           `json:"transaction_id"`
	Created        int64            `json:"created"`
	UserID         string           `json:"user_id"`
	Amount         string           `json:"amount"`
	AmountRefunded string           `json:"amount_refunded"`
	State          PaymentStateType `json:"state"`
}

func (m *PaymentTransaction) Columns(action string) []string {
	return []string{
		"id",
		"transaction_id",
		"created",
		"user_id",
		"amount",
		"amount_refunded",
		"state",
	}
}

func (m *PaymentTransaction) Table() string {
	return "payment_transaction"
}

func (m *PaymentTransaction) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.TransactionID, &m.Created, &m.UserID, &m.Amount, &m.AmountRefunded, &m.State)
}

func (m *PaymentTransaction) Values(action string) []any {
	return []any{m.ID, m.TransactionID, m.Created, m.UserID, m.Amount, m.AmountRefunded, m.State}
}
