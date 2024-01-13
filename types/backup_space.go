package types

func BackupSpaceProvider() *BackupSpace { return &BackupSpace{} }

type BackupSpace struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	UserID  string `json:"user_id"`
	Created int64  `json:"created"`
}

func (m *BackupSpace) Columns(action string) []string {
	return []string{
		"id",
		"name",
		"user_id",
		"created",
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
		&m.Created,
	)
}

func (m *BackupSpace) Values(action string) []any {
	return []any{
		m.ID,
		m.Name,
		m.UserID,
		m.Created,
	}
}
