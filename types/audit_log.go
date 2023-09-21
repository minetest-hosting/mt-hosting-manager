package types

func AuditLogProvider() *AuditLog { return &AuditLog{} }

type AuditLogType string

const (
	AuditLogUserCreated  AuditLogType = "user_created"
	AuditLogUserLoggedIn AuditLogType = "user_logged_in"
)

type AuditLog struct {
	ID                   string       `json:"id"`
	Type                 AuditLogType `json:"type"`
	Timestamp            int64        `json:"timestamp"`
	UserID               string       `json:"user_id"`
	UserNodeID           *string      `json:"user_node_id"`
	MinetestServerID     *string      `json:"minetest_server_id"`
	PaymentTransactionID *string      `json:"payment_transaction_id"`
}

func (m *AuditLog) Columns(action string) []string {
	return []string{"id", "type", "timestamp", "user_id", "user_node_id", "minetest_server_id", "payment_transaction_id"}
}

func (m *AuditLog) Table() string {
	return "audit_log"
}

func (m *AuditLog) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.Type, &m.Timestamp, &m.UserID, &m.UserNodeID, &m.MinetestServerID, &m.PaymentTransactionID)
}

func (m *AuditLog) Values(action string) []any {
	return []any{m.ID, m.Type, m.Timestamp, m.UserID, m.UserNodeID, m.MinetestServerID, m.PaymentTransactionID}
}
