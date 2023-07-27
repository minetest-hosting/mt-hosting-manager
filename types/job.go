package types

type JobState string

const (
	JobStateCreated JobState = "CREATED"
)

type JobType string

const (
	JobTypeNodeSetup JobType = "NODE_SETUP"
)

type Job struct {
	ID               string   `json:"id"`
	Type             JobType  `json:"type"`
	State            JobState `json:"state"`
	Started          int64    `json:"started"`
	Finished         int64    `json:"finished"`
	UserNodeID       string   `json:"user_node_id"`
	MinetestServerID string   `json:"minetest_server_id"`
	ProgressPercent  float64  `json:"progress_percent"`
	Message          string   `json:"message"`
	Data             []byte   `json:"data"`
}

func (m *Job) Columns(action string) []string {
	return []string{
		"id",
		"type",
		"state",
		"started",
		"finished",
		"user_node_id",
		"minetest_server_id",
		"progress_percent",
		"message",
		"data",
	}
}

func (m *Job) Table() string {
	return "job"
}

func (m *Job) Scan(action string, r func(dest ...any) error) error {
	return r(
		&m.ID,
		&m.Type,
		&m.State,
		&m.Started,
		&m.Finished,
		&m.UserNodeID,
		&m.MinetestServerID,
		&m.ProgressPercent,
		&m.Message,
		&m.Data,
	)
}

func (m *Job) Values(action string) []any {
	return []any{
		m.ID,
		m.Type,
		m.State,
		m.Started,
		m.Finished,
		m.UserNodeID,
		m.MinetestServerID,
		m.ProgressPercent,
		m.Message,
		m.Data,
	}
}
