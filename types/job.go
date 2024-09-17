package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type JobState string

const (
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

type Job struct {
	ID               string   `json:"id" gorm:"primarykey;column:id"`
	Type             JobType  `json:"type" gorm:"column:type"`
	State            JobState `json:"state" gorm:"column:state"`
	Created          int64    `json:"created" gorm:"column:created"`
	NextRun          int64    `json:"next_run" gorm:"column:next_run"`
	Step             int      `json:"step" gorm:"column:step"`
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
		Created:    time.Now().Unix(),
		NextRun:    time.Now().Unix(),
		Type:       JobTypeNodeSetup,
		State:      JobStateRunning,
		UserNodeID: &node.ID,
	}
}

func RemoveNodeJob(node *UserNode) *Job {
	return &Job{
		ID:         uuid.NewString(),
		Created:    time.Now().Unix(),
		NextRun:    time.Now().Unix(),
		Type:       JobTypeNodeDestroy,
		State:      JobStateRunning,
		UserNodeID: &node.ID,
	}
}

func SetupServerJob(node *UserNode, server *MinetestServer) *Job {
	return &Job{
		ID:               uuid.NewString(),
		Created:          time.Now().Unix(),
		NextRun:          time.Now().Unix(),
		Type:             JobTypeServerSetup,
		State:            JobStateRunning,
		UserNodeID:       &node.ID,
		MinetestServerID: &server.ID,
	}
}

func RemoveServerJob(node *UserNode, server *MinetestServer) *Job {
	return &Job{
		ID:               uuid.NewString(),
		Created:          time.Now().Unix(),
		NextRun:          time.Now().Unix(),
		Type:             JobTypeServerDestroy,
		State:            JobStateRunning,
		UserNodeID:       &node.ID,
		MinetestServerID: &server.ID,
	}
}

func BackupServerJob(node *UserNode, server *MinetestServer, backup *Backup) *Job {
	return &Job{
		ID:               uuid.NewString(),
		Created:          time.Now().Unix(),
		NextRun:          time.Now().Unix(),
		Type:             JobTypeServerBackup,
		State:            JobStateRunning,
		UserNodeID:       &node.ID,
		MinetestServerID: &server.ID,
		BackupID:         &backup.ID,
	}
}

func (job *Job) LogrusFields() logrus.Fields {
	return logrus.Fields{
		"jobid":              job.ID,
		"type":               job.Type,
		"state":              job.State,
		"user_node_id":       job.UserNodeID,
		"backup_id":          job.BackupID,
		"next_run":           job.NextRun,
		"minetest_server_id": job.MinetestServerID,
		"message":            job.Message,
		"created":            job.Created,
	}
}
