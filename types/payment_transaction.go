package types

type PaymentStateType string

const (
	PaymentStatePending PaymentStateType = "PENDING"
	PaymentStateSuccess PaymentStateType = "SUCCESS"
	PaymentStateExpired PaymentStateType = "EXPIRED"
	PaymentStateError   PaymentStateType = "ERROR"
)

type PaymentType string

const (
	PaymentTypeWallee   PaymentType = "WALLEE"
	PaymentTypeCoinbase PaymentType = "COINBASE"
	PaymentTypeZahlsch  PaymentType = "ZAHLSCH"
)

func PaymentTransactionProvider() *PaymentTransaction { return &PaymentTransaction{} }

type PaymentTransaction struct {
	ID             string           `json:"id"`
	Type           PaymentType      `json:"type"`
	TransactionID  string           `json:"transaction_id"`
	PaymentURL     string           `json:"payment_url"`
	Created        int64            `json:"created"`
	Expires        int64            `json:"expires"`
	UserID         string           `json:"user_id"`
	Amount         int64            `json:"amount"`
	AmountRefunded int64            `json:"amount_refunded"`
	State          PaymentStateType `json:"state"`
}

type PaymentTransactionSearch struct {
	FromTimestamp int64   `json:"from_timestamp"`
	ToTimestamp   int64   `json:"to_timestamp"`
	UserID        *string `json:"user_id"`
}

func (m *PaymentTransaction) Columns(action string) []string {
	return []string{
		"id",
		"type",
		"transaction_id",
		"payment_url",
		"created",
		"expires",
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
	return r(
		&m.ID,
		&m.Type,
		&m.TransactionID,
		&m.PaymentURL,
		&m.Created,
		&m.Expires,
		&m.UserID,
		&m.Amount,
		&m.AmountRefunded,
		&m.State,
	)
}

func (m *PaymentTransaction) Values(action string) []any {
	return []any{
		m.ID,
		m.Type,
		m.TransactionID,
		m.PaymentURL,
		m.Created,
		m.Expires,
		m.UserID,
		m.Amount,
		m.AmountRefunded,
		m.State,
	}
}
