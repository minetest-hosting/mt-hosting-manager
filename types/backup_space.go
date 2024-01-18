package types

func BackupSpaceProvider() *BackupSpace { return &BackupSpace{} }

type BackupSpace struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	UserID        string `json:"user_id"`
	RetentionDays int    `json:"retention_days"`
	Created       int64  `json:"created"`
	ValidUntil    int64  `json:"valid_until"`
}

func (m *BackupSpace) Columns(action string) []string {
	return []string{
		"id",
		"name",
		"user_id",
		"retention_days",
		"created",
		"valid_until",
	}
}

func (m *BackupSpace) Table() string {
	return "backup_space"
}

func (m *BackupSpace) Scan(action string, r func(dest ...any) error) error {
	return r(
		&m.ID,
		&m.Name,
		&m.UserID,
		&m.RetentionDays,
		&m.Created,
		&m.ValidUntil,
	)
}

func (m *BackupSpace) Values(action string) []any {
	return []any{
		m.ID,
		m.Name,
		m.UserID,
		m.RetentionDays,
		m.Created,
		m.ValidUntil,
	}
}
