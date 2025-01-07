package types

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

	AuditLogCouponRedeemed AuditLogType = "coupon_redeemed"

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
	ID                   string       `json:"id" gorm:"primarykey;column:id"`
	Type                 AuditLogType `json:"type" gorm:"column:type"`
	Timestamp            int64        `json:"timestamp" gorm:"column:timestamp"`
	UserID               string       `json:"user_id" gorm:"column:user_id"`
	IPAddress            *string      `json:"ip_address" gorm:"column:ip_address"`
	UserNodeID           *string      `json:"user_node_id" gorm:"column:user_node_id"`
	MinetestServerID     *string      `json:"minetest_server_id" gorm:"column:minetest_server_id"`
	BackupID             *string      `json:"backup_id" gorm:"column:backup_id"`
	PaymentTransactionID *string      `json:"payment_transaction_id" gorm:"column:payment_transaction_id"`
	Amount               *int64       `json:"amount" gorm:"column:amount"`
}

func (m *AuditLog) TableName() string {
	return "audit_log"
}
