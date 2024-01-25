package worker

import (
	"fmt"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker/server_setup"
)

func (w *Worker) ServerRestore(job *types.Job, status StatusCallback) error {
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

	w.core.AddAuditLog(&types.AuditLog{
		Type:             types.AuditLogServerRestoreStarted,
		UserID:           node.UserID,
		UserNodeID:       &node.ID,
		MinetestServerID: &server.ID,
		BackupID:         job.BackupID,
	})

	err = w.ServerPrepareSetup(job, node, server)
	if err != nil {
		return err
	}

	client, err := core.TrySSHConnection(node)
	if err != nil {
		return err
	}

	err = server_setup.Restore(client, w.cfg, node, server, backup)
	if err != nil {
		return err
	}

	server.State = types.MinetestServerStateProvisioning
	err = w.repos.MinetestServerRepo.Update(server)
	if err != nil {
		return fmt.Errorf("server entity update error: %v", err)
	}

	return ErrJobStillRunning
}
