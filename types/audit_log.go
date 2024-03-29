package types

func AuditLogProvider() *AuditLog { return &AuditLog{} }

type AuditLogType string

const (
	AuditLogUserActivated AuditLogType = "user_activated"
	AuditLogUserCreated   AuditLogType = "user_created"
	AuditLogUserLoggedIn  AuditLogType = "user_logged_in"

	AuditLogNodeCreated AuditLogType = "node_created"
	AuditLogNodeRemoved AuditLogType = "node_removed"

	AuditLogNodeProvisioningStarted  AuditLogType = "node_provisioning_started"
	AuditLogNodeProvisioningFinished AuditLogType = "node_provisioning_finished"
	AuditLogNodeBilled               AuditLogType = "node_billed"

	AuditLogServerCreated        AuditLogType = "server_created"
	AuditLogServerRemoved        AuditLogType = "server_removed"
	AuditLogServerSetupStarted   AuditLogType = "server_setup_started"
	AuditLogServerRestoreStarted AuditLogType = "server_restore_started"
	AuditLogServerSetupFinished  AuditLogType = "server_setup_finished"

	AuditLogServerBackupStarted  AuditLogType = "server_backup_started"
	AuditLogServerBackupFinished AuditLogType = "server_backup_finished"

	AuditLogPaymentCreated  AuditLogType = "payment_created"
	AuditLogPaymentReceived AuditLogType = "payment_received"
	AuditLogPaymentRefunded AuditLogType = "payment_refunded"
	AuditLogPaymentWarning  AuditLogType = "payment_warning"
	AuditLogPaymentZero     AuditLogType = "payment_zero"
)

type AuditLogSearch struct {
	FromTimestamp        int64         `json:"from_timestamp"`
	ToTimestamp          int64         `json:"to_timestamp"`
	Type                 *AuditLogType `json:"type"`
	UserID               *string       `json:"user_id"`
	UserNodeID           *string       `json:"user_node_id"`
	MinetestServerID     *string       `json:"minetest_server_id"`
	BackupID             *string       `json:"backup_id"`
	PaymentTransactionID *string       `json:"payment_transaction_id"`
}

type AuditLog struct {
	ID                   string       `json:"id"`
	Type                 AuditLogType `json:"type"`
	Timestamp            int64        `json:"timestamp"`
	UserID               string       `json:"user_id"`
	IPAddress            *string      `json:"ip_address"`
	UserNodeID           *string      `json:"user_node_id"`
	MinetestServerID     *string      `json:"minetest_server_id"`
	BackupID             *string      `json:"backup_id"`
	PaymentTransactionID *string      `json:"payment_transaction_id"`
	Amount               *int64       `json:"amount"`
}

func (m *AuditLog) Columns(action string) []string {
	return []string{
		"id",
		"type",
		"timestamp",
		"user_id",
		"ip_address",
		"user_node_id",
		"minetest_server_id",
		"backup_id",
		"payment_transaction_id",
		"amount",
	}
}

func (m *AuditLog) Table() string {
	return "audit_log"
}

func (m *AuditLog) Scan(action string, r func(dest ...any) error) error {
	return r(
		&m.ID,
		&m.Type,
		&m.Timestamp,
		&m.UserID,
		&m.IPAddress,
		&m.UserNodeID,
		&m.MinetestServerID,
		&m.BackupID,
		&m.PaymentTransactionID,
		&m.Amount,
	)
}

func (m *AuditLog) Values(action string) []any {
	return []any{
		m.ID,
		m.Type,
		m.Timestamp,
		m.UserID,
		m.IPAddress,
		m.UserNodeID,
		m.MinetestServerID,
		m.BackupID,
		m.PaymentTransactionID,
		m.Amount,
	}
}
