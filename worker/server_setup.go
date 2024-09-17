package worker

import (
	"fmt"
	"mt-hosting-manager/api/mtui"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker/server_setup"
	"time"
)

func (w *Worker) ServerSetup(job *types.Job) error {
	node, server, err := w.GetJobContext(job)
	if err != nil {
		return err
	}

	switch job.Step {
	case 0:
		w.core.AddAuditLog(&types.AuditLog{
			Type:             types.AuditLogServerSetupStarted,
			UserID:           node.UserID,
			UserNodeID:       &node.ID,
			MinetestServerID: &server.ID,
		})

		err = w.ServerPrepareSetup(job, node, server)
		if err != nil {
			return err
		}

		client, err := core.TrySSHConnection(node)
		if err != nil {
			return err
		}

		err = server_setup.Setup(client, w.cfg, node, server)
		if err != nil {
			return err
		}

		job.Step = 1

	case 1:
		if job.BackupID == nil {
			// skip restore steps
			job.Step = 3
			return nil
		}

		client, err := w.core.GetMTUIClient(server)
		if err != nil {
			return fmt.Errorf("get client error: %v", err)
		}

		info, err := client.CreateRestoreJob(&mtui.CreateRestoreJob{
			ID:       *job.BackupID,
			Type:     mtui.RestoreJobTypeSCP,
			Host:     w.cfg.StorageHostname,
			Username: w.cfg.StorageUsername,
			Password: w.cfg.StoragePassword,
			Port:     22,
			Filename: fmt.Sprintf("%s.tar.gz", *job.BackupID),
		})
		if err != nil {
			return fmt.Errorf("create restore job error: %v", err)
		}

		job.Message = info.Message
		job.Data = []byte(info.ID)
		job.Step = 2
		job.NextRun = time.Now().Add(5 * time.Second).Unix()

	case 2:

		client, err := w.core.GetMTUIClient(server)
		if err != nil {
			return fmt.Errorf("get client error: %v", err)
		}

		// check restore job
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
			job.Step = 3

		case mtui.RestoreJobFailure:
			// restore failed
			job.Message = info.Message
			job.State = types.JobStateDoneFailure
		}

	case 3:
		server.State = types.MinetestServerStateRunning
		err = w.repos.MinetestServerRepo.Update(server)
		if err != nil {
			return fmt.Errorf("server entity update error: %v", err)
		}

		w.core.AddAuditLog(&types.AuditLog{
			Type:             types.AuditLogServerSetupFinished,
			UserID:           node.UserID,
			UserNodeID:       &node.ID,
			MinetestServerID: &server.ID,
		})

		job.State = types.JobStateDoneSuccess
	}

	return nil
}
