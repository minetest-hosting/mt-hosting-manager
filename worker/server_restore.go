package worker

import (
	"fmt"
	"mt-hosting-manager/api/mtui"
	"mt-hosting-manager/types"
	"time"
)

func (w *Worker) ServerRestore(job *types.Job) error {
	node, server, err := w.GetJobContext(job)
	if err != nil {
		return err
	}
	backup, err := w.repos.BackupRepo.GetByID(*job.BackupID)
	if err != nil {
		return err
	}
	if backup == nil {
		return fmt.Errorf("backup not found: %s", *job.BackupID)
	}

	client, err := w.core.GetMTUIClient(server)
	if err != nil {
		return fmt.Errorf("get client error: %v", err)
	}

	switch job.Step {
	case 0:
		w.core.AddAuditLog(&types.AuditLog{
			Type:             types.AuditLogServerRestoreStarted,
			UserID:           node.UserID,
			UserNodeID:       &node.ID,
			MinetestServerID: &server.ID,
			BackupID:         job.BackupID,
		})

		info, err := client.CreateRestoreJob(&mtui.CreateRestoreJob{
			ID:       backup.ID,
			Type:     mtui.RestoreJobTypeSCP,
			Host:     w.cfg.StorageHostname,
			Username: w.cfg.StorageUsername,
			Password: w.cfg.StoragePassword,
			Port:     22,
			Filename: fmt.Sprintf("%s.tar.gz", backup.ID),
		})
		if err != nil {
			return fmt.Errorf("create restore job error: %v", err)
		}

		server.State = types.MinetestServerStateProvisioning
		err = w.repos.MinetestServerRepo.Update(server)
		if err != nil {
			return fmt.Errorf("server entity update error: %v", err)
		}

		job.Message = info.Message
		job.Data = []byte(info.ID)
		job.Step = 1
		job.NextRun = time.Now().Add(5 * time.Second).Unix()

	case 1:
		// check backup
		info, err := client.GetRestoreJobInfo(string(job.Data))
		if err != nil {
			return fmt.Errorf("get restore job error: %v", err)
		}

		switch info.Status {
		case mtui.RestoreJobRunning:
			// still running
			job.Message = info.Message
			job.NextRun = time.Now().Add(5 * time.Second).Unix()

		case mtui.RestoreJobSuccess:
			// all done
			job.Message = info.Message
			job.State = types.JobStateDoneSuccess

		case mtui.RestoreJobFailure:
			// restore failed
			job.Message = info.Message
			job.State = types.JobStateDoneFailure
		}
	}

	return nil
}
