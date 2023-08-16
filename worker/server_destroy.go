package worker

import (
	"fmt"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker/server_setup"
)

func (w *Worker) ServerDestroy(job *types.Job) error {
	node, err := w.repos.UserNodeRepo.GetByID(*job.UserNodeID)
	if err != nil {
		return fmt.Errorf("usernode fetch error: %v", err)
	}
	if node == nil {
		return fmt.Errorf("usernode not found: %s", *job.UserNodeID)
	}

	server, err := w.repos.MinetestServerRepo.GetByID(*job.MinetestServerID)
	if err != nil {
		return fmt.Errorf("usernode fetch error: %v", err)
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", *job.MinetestServerID)
	}
	server.State = types.MinetestServerStateRemoving
	err = w.repos.MinetestServerRepo.Update(server)
	if err != nil {
		return fmt.Errorf("server entity update error: %v", err)
	}

	//TODO: dns teardown

	client, err := TrySSHConnection(node)
	if err != nil {
		return err
	}

	basedir := server_setup.GetBaseDir(server)
	_, _, err = core.SSHExecute(client, fmt.Sprintf("cd %s && docker-compose down -v", basedir))
	if err != nil {
		return fmt.Errorf("could not run docker-compose down: %v", err)
	}

	_, _, err = core.SSHExecute(client, fmt.Sprintf("rm -rf %s", basedir))
	if err != nil {
		return fmt.Errorf("could not run remove data-dir '%s': %v", basedir, err)
	}

	return nil
}
