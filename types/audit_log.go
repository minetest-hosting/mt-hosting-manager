package types

func AuditLogProvider() *AuditLog { return &AuditLog{} }

type AuditLogType string

const (
	AuditLogUserCreated  AuditLogType = "user_created"
	AuditLogUserLoggedIn AuditLogType = "user_logged_in"

	AuditLogNodeCreated AuditLogType = "node_created"
	AuditLogNodeRemoved AuditLogType = "node_removed"

	AuditLogNodeProvisioningStarted  AuditLogType = "node_provisioning_started"
	AuditLogNodeProvisioningFinished AuditLogType = "node_provisioning_finished"

	AuditLogServerCreated       AuditLogType = "server_created"
	AuditLogServerRemoved       AuditLogType = "server_removed"
	AuditLogServerSetupStarted  AuditLogType = "server_setup_started"
	AuditLogServerSetupFinished AuditLogType = "server_setup_finished"

	AuditLogPaymentCreated  AuditLogType = "payment_created"
	AuditLogPaymentReceived AuditLogType = "payment_received"
	AuditLogPaymentRefunded AuditLogType = "payment_refunded"
)

type AuditLogSearch struct {
	FromTimestamp int64         `json:"from_timestamp"`
	ToTimestamp   int64         `json:"to_timestamp"`
	Type          *AuditLogType `json:"type"`
	UserID        *string       `json:"user_id"`
}

type AuditLog struct {
	ID                   string       `json:"id"`
	Type                 AuditLogType `json:"type"`
	Timestamp            int64        `json:"timestamp"`
	UserID               string       `json:"user_id"`
	UserNodeID           *string      `json:"user_node_id"`
	MinetestServerID     *string      `json:"minetest_server_id"`
	PaymentTransactionID *string      `json:"payment_transaction_id"`
	Amount               *string      `json:"amount"`
	Currency             *string      `json:"currency"`
}

func (m *AuditLog) Columns(action string) []string {
	return []string{"id", "type", "timestamp", "user_id", "user_node_id", "minetest_server_id", "payment_transaction_id", "amount", "currency"}
}

func (m *AuditLog) Table() string {
	return "audit_log"
}

func (m *AuditLog) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.Type, &m.Timestamp, &m.UserID, &m.UserNodeID, &m.MinetestServerID, &m.PaymentTransactionID, &m.Amount, &m.Currency)
}

func (m *AuditLog) Values(action string) []any {
	return []any{m.ID, m.Type, m.Timestamp, m.UserID, m.UserNodeID, m.MinetestServerID, m.PaymentTransactionID, m.Amount, m.Currency}
}
