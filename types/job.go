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
	JobTypeServerRestore JobType = "SERVER_RESTORE"
	JobTypeServerDestroy JobType = "SERVER_DESTROY"
	JobTypeServerBackup  JobType = "SERVER_BACKUP"
)

func JobProvider() *Job { return &Job{} }

type Job struct {
	ID               string   `json:"id" gorm:"primarykey;column:id"`
	Type             JobType  `json:"type" gorm:"column:type"`
	State            JobState `json:"state" gorm:"column:state"`
	Started          int64    `json:"started" gorm:"column:started"`
	Finished         int64    `json:"finished" gorm:"column:finished"`
	UserNodeID       *string  `json:"user_node_id" gorm:"column:user_node_id"`
	MinetestServerID *string  `json:"minetest_server_id" gorm:"column:minetest_server_id"`
	BackupID         *string  `json:"backup_id" gorm:"column:backup_id"`
	ProgressPercent  float64  `json:"progress_percent" gorm:"column:progress_percent"`
	Message          string   `json:"message" gorm:"column:message"`
	Data             []byte   `json:"data" gorm:"column:data"`
}

func (m *Job) TableName() string {
	return "job"
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
