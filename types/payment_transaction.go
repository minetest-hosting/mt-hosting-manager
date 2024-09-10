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
	ID             string           `json:"id" gorm:"primarykey;column:id"`
	Type           PaymentType      `json:"type" gorm:"column:type"`
	TransactionID  string           `json:"transaction_id" gorm:"column:transaction_id"`
	PaymentURL     string           `json:"payment_url" gorm:"column:payment_url"`
	Created        int64            `json:"created" gorm:"column:created"`
	Expires        int64            `json:"expires" gorm:"column:expires"`
	UserID         string           `json:"user_id" gorm:"column:user_id"`
	Amount         int64            `json:"amount" gorm:"column:amount"`
	AmountRefunded int64            `json:"amount_refunded" gorm:"column:amount_refunded"`
	State          PaymentStateType `json:"state" gorm:"column:state"`
}

type PaymentTransactionSearch struct {
	FromTimestamp int64             `json:"from_timestamp"`
	ToTimestamp   int64             `json:"to_timestamp"`
	UserID        *string           `json:"user_id"`
	State         *PaymentStateType `json:"state"`
}

func (m *PaymentTransaction) TableName() string {
	return "payment_transaction"
}
