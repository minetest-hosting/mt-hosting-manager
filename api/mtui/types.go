package mtui

type BackupJobState string

const (
	BackupJobRunning BackupJobState = "running"
	BackupJobSuccess BackupJobState = "success"
	BackupJobFailure BackupJobState = "failure"
)

type BackupJobInfo struct {
	ID      string         `json:"id"`
	Status  BackupJobState `json:"state"`
	Message string         `json:"message"`
}

type BackupJobType string

const (
	BackupJobTypeSCP BackupJobType = "scp"
)

type CreateBackupJob struct {
	ID       string        `json:"id"`
	Type     BackupJobType `json:"type"`
	Host     string        `json:"host"`
	Port     int           `json:"port"`
	Filename string        `json:"filename"`
	Username string        `json:"username"`
	Password string        `json:"password"`
	Key      string        `json:"key"`
}

type RestoreJobState string

const (
	RestoreJobRunning RestoreJobState = "running"
	RestoreJobSuccess RestoreJobState = "success"
	RestoreJobFailure RestoreJobState = "failure"
)

type RestoreJobInfo struct {
	ID      string          `json:"id"`
	Status  RestoreJobState `json:"state"`
	Message string          `json:"message"`
}

type RestoreJobType string

const (
	RestoreJobTypeSCP RestoreJobType = "scp"
)

type CreateRestoreJob struct {
	ID       string         `json:"id"`
	Type     RestoreJobType `json:"type"`
	Host     string         `json:"host"`
	Port     int            `json:"port"`
	Filename string         `json:"filename"`
	Username string         `json:"username"`
	Password string         `json:"password"`
	Key      string         `json:"key"`
}
