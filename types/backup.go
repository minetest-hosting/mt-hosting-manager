package types

type BackupState string

const (
	BackupStateCreated  BackupState = "CREATED"
	BackupStateProgress BackupState = "PROGRESS"
	BackupStateComplete BackupState = "COMPLETE"
	BackupStateError    BackupState = "ERROR"
)

func BackupProvider() *Backup { return &Backup{} }

type Backup struct {
	ID               string      `json:"id" gorm:"primarykey;column:id"`
	State            BackupState `json:"state" gorm:"column:state"`
	Passphrase       string      `json:"passphrase" gorm:"column:passphrase"`
	BackupSpaceID    string      `json:"backup_space_id" gorm:"column:backup_space_id"`
	MinetestServerID string      `json:"minetest_server_id" gorm:"column:minetest_server_id"`
	Created          int64       `json:"created" gorm:"column:created"`
	Size             int64       `json:"size" gorm:"column:size"`
}

func (m *Backup) TableName() string {
	return "backup"
}
