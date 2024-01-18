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
	ID               string      `json:"id"`
	State            BackupState `json:"state"`
	Passphrase       string      `json:"passphrase"`
	BackupSpaceID    string      `json:"backup_space_id"`
	MinetestServerID string      `json:"minetest_server_id"`
	Created          int64       `json:"created"`
	Size             int64       `json:"size"`
}

func (m *Backup) Columns(action string) []string {
	return []string{
		"id",
		"state",
		"passphrase",
		"backup_space_id",
		"minetest_server_id",
		"created",
		"size",
	}
}

func (m *Backup) Table() string {
	return "backup"
}

func (m *Backup) Scan(action string, r func(dest ...any) error) error {
	return r(
		&m.ID,
		&m.State,
		&m.Passphrase,
		&m.BackupSpaceID,
		&m.MinetestServerID,
		&m.Created,
		&m.Size,
	)
}

func (m *Backup) Values(action string) []any {
	return []any{
		m.ID,
		m.State,
		m.Passphrase,
		m.BackupSpaceID,
		m.MinetestServerID,
		m.Created,
		m.Size,
	}
}
