package worker

import (
	"fmt"
	"mt-hosting-manager/api/mtui"
	"mt-hosting-manager/types"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/sirupsen/logrus"
)

var mtui_cache = expirable.NewLRU[string, *mtui.MtuiClient](5, nil, time.Hour*2)

func (w *Worker) updateBackupJob(job *types.Job) error {
	server, err := w.repos.MinetestServerRepo.GetByID(*job.MinetestServerID)
	if err != nil {
		return fmt.Errorf("get server error: %v", err)
	}
	if server == nil {
		return fmt.Errorf("server not found")
	}

	url := fmt.Sprintf("https://%s.%s/ui", server.DNSName, w.cfg.HostingDomainSuffix)
	client, found := mtui_cache.Get(url)
	if !found {
		// create a new client and log in
		client = mtui.New(url)
		err = client.Login(server.Admin, server.JWTKey)
		if err != nil {
			return fmt.Errorf("login error: %v", err)
		}
		mtui_cache.Add(url, client)
	}

	info, err := client.GetBackupJobInfo(string(job.Data))
	if err != nil {
		return fmt.Errorf("GetBackupJobInfo error: %v", err)
	}

	backup, err := w.repos.BackupRepo.GetByID(*job.BackupID)
	if err != nil {
		return fmt.Errorf("get backup error: %v", err)
	}
	if backup == nil {
		return fmt.Errorf("backup not found: '%s'", *job.BackupID)
	}

	switch info.Status {
	case mtui.BackupJobRunning:
		job.Message = info.Message
	case mtui.BackupJobFailure:
		job.State = types.JobStateDoneFailure
		job.Message = info.Message
		backup.State = types.BackupStateError
	case mtui.BackupJobSuccess:
		// get size from storage
		size, err := w.core.GetBackupSize(backup)
		if err != nil {
			job.State = types.JobStateDoneFailure
			job.Message = fmt.Sprintf("backup-file stat failed: %v", err)
			backup.State = types.BackupStateError
		} else {
			// everything checks out
			job.State = types.JobStateDoneSuccess
			job.Message = info.Message
			backup.State = types.BackupStateComplete
			backup.Size = size
		}
	}

	err = w.repos.BackupRepo.Update(backup)
	if err != nil {
		return fmt.Errorf("update backup error: %v", err)
	}

	return w.repos.JobRepo.Update(job)
}

func (w *Worker) UpdateBackupProgressJob() {
	for w.running.Load() {
		time.Sleep(1 * time.Second)

		jobs, err := w.repos.JobRepo.GetByTypeAndState(types.JobTypeServerBackup, types.JobStateRunning)
		if err != nil {
			logrus.WithError(err).Error("JobRepo.GetByTypeAndState error")
			continue
		}

		for _, job := range jobs {
			err = w.updateBackupJob(job)
			if err != nil {
				logrus.WithError(err).Error("updateBackupJob error")
				continue
			}
		}
	}
}
