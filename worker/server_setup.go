package worker

import (
	"fmt"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker/server_setup"
)

func (w *Worker) ServerSetup(job *types.Job) error {
	node, server, err := w.GetJobContext(job)
	if err != nil {
		return err
	}

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
	return nil
}
