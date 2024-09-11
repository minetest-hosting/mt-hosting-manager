package worker

import (
	"fmt"
	"mt-hosting-manager/api/mtui"
	"mt-hosting-manager/types"
	"time"

	"github.com/sirupsen/logrus"
)

func (w *Worker) updateBackupJob(job *types.Job) error {
	server, err := w.repos.MinetestServerRepo.GetByID(*job.MinetestServerID)
	if err != nil {
		return fmt.Errorf("get server error: %v", err)
	}
	if server == nil {
		return fmt.Errorf("server not found")
	}

	// TODO: cache logged in client
	url := fmt.Sprintf("https://%s.%s/ui", server.DNSName, w.cfg.HostingDomainSuffix)
	client := mtui.New(url)
	err = client.Login(server.Admin, server.JWTKey)
	if err != nil {
		return fmt.Errorf("login error: %v", err)
	}

	info, err := client.GetBackupJobInfo(string(job.Data))
	if err != nil {
		return fmt.Errorf("GetBackupJobInfo error: %v", err)
	}

	// TODO: update backup entry
	switch info.Status {
	case mtui.BackupJobRunning:
		job.Message = info.Message
	case mtui.BackupJobFailure:
		job.State = types.JobStateDoneFailure
		job.Message = info.Message
	case mtui.BackupJobSuccess:
		job.State = types.JobStateDoneSuccess
		job.Message = info.Message
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
