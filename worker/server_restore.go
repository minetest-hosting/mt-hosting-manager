package worker

import (
	"fmt"
	"mt-hosting-manager/types"
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

	w.core.AddAuditLog(&types.AuditLog{
		Type:             types.AuditLogServerRestoreStarted,
		UserID:           node.UserID,
		UserNodeID:       &node.ID,
		MinetestServerID: &server.ID,
		BackupID:         job.BackupID,
	})

	//TODO
	job.State = types.JobStateDoneSuccess

	server.State = types.MinetestServerStateProvisioning
	err = w.repos.MinetestServerRepo.Update(server)
	if err != nil {
		return fmt.Errorf("server entity update error: %v", err)
	}

	return nil
}
