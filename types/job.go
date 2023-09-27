package types

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type JobState string

const (
	JobStateCreated     JobState = "CREATED"
	JobStateRunning     JobState = "RUNNING"
	JobStateDoneSuccess JobState = "DONE_SUCCESS"
	JobStateDoneFailure JobState = "DONE_FAILURE"
)

// Created -> Running -> Done_Success -> (Created)
//                    -> Done_Failure -> (Created)

type JobType string

const (
	JobTypeNodeSetup     JobType = "NODE_SETUP"
	JobTypeNodeDestroy   JobType = "NODE_DESTROY"
	JobTypeServerSetup   JobType = "SERVER_SETUP"
	JobTypeServerDestroy JobType = "SERVER_DESTROY"
)

func JobProvider() *Job { return &Job{} }

type Job struct {
	ID               string   `json:"id"`
	Type             JobType  `json:"type"`
	State            JobState `json:"state"`
	Started          int64    `json:"started"`
	Finished         int64    `json:"finished"`
	UserNodeID       *string  `json:"user_node_id"`
	MinetestServerID *string  `json:"minetest_server_id"`
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

func (job *Job) LogrusFields() logrus.Fields {
	return logrus.Fields{
		"jobid":              job.ID,
		"type":               job.Type,
		"state":              job.State,
		"user_node_id":       job.UserNodeID,
		"minetest_server_id": job.MinetestServerID,
		"message":            job.Message,
		"started":            job.Started,
	}
}

func SetupNodeJob(node *UserNode) *Job {
	return &Job{
		ID:         uuid.NewString(),
		Type:       JobTypeNodeSetup,
		State:      JobStateCreated,
		UserNodeID: &node.ID,
	}
}

func RemoveNodeJob(node *UserNode) *Job {
	return &Job{
		ID:         uuid.NewString(),
		Type:       JobTypeNodeDestroy,
		State:      JobStateCreated,
		UserNodeID: &node.ID,
	}
}

func SetupServerJob(node *UserNode, server *MinetestServer) *Job {
	return &Job{
		ID:               uuid.NewString(),
		Type:             JobTypeServerSetup,
		State:            JobStateCreated,
		UserNodeID:       &node.ID,
		MinetestServerID: &server.ID,
	}
}

func RemoveServerJob(node *UserNode, server *MinetestServer) *Job {
	return &Job{
		ID:               uuid.NewString(),
		Type:             JobTypeServerDestroy,
		State:            JobStateCreated,
		UserNodeID:       &node.ID,
		MinetestServerID: &server.ID,
	}
}
