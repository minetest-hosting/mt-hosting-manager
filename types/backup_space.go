package types

type BackupSpace struct {
	ID            string `json:"id" gorm:"primarykey;column:id"`
	Name          string `json:"name" gorm:"column:name"`
	UserID        string `json:"user_id" gorm:"column:user_id"`
	RetentionDays int    `json:"retention_days" gorm:"column:retention_days"`
	Created       int64  `json:"created" gorm:"column:created"`
	ValidUntil    int64  `json:"valid_until" gorm:"column:valid_until"`
}

func (m *BackupSpace) TableName() string {
	return "backup_space"
}
