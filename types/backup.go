package types

type BackupState string

const (
	BackupStateCreated  BackupState = "CREATED"
	BackupStateProgress BackupState = "PROGRESS"
	BackupStateComplete BackupState = "COMPLETE"
	BackupStateError    BackupState = "ERROR"
)

type Backup struct {
	ID               string      `json:"id" gorm:"primarykey;column:id"`
	State            BackupState `json:"state" gorm:"column:state"`
	Passphrase       string      `json:"passphrase" gorm:"column:passphrase"`
	UserID           string      `json:"user_id" gorm:"column:user_id"`
	MinetestServerID string      `json:"minetest_server_id" gorm:"column:minetest_server_id"`
	Created          int64       `json:"created" gorm:"column:created"`
	Expires          int64       `json:"expires" gorm:"column:expires"`
	Size             int64       `json:"size" gorm:"column:size"`
	Comment          string      `json:"comment" gorm:"column:comment"`
}

func (m *Backup) TableName() string {
	return "backup"
}
