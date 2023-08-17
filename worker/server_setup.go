package worker

import (
	"fmt"
	"mt-hosting-manager/api/hetzner_dns"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker/server_setup"
)

func (w *Worker) ServerSetup(job *types.Job) error {
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
	server.State = types.MinetestServerStateProvisioning
	err = w.repos.MinetestServerRepo.Update(server)
	if err != nil {
		return fmt.Errorf("server entity update error: %v", err)
	}

	records, err := w.hdc.GetRecords()
	if err != nil {
		return fmt.Errorf("fetch records error: %v", err)
	}

	err = w.UpdateDNSRecord(records, hetzner_dns.RecordCNAME, server.DNSName, node.Name)
	if err != nil {
		return err
	}

	client, err := TrySSHConnection(node)
	if err != nil {
		return err
	}

	err = server_setup.Setup(client, node, server)
	if err != nil {
		return err
	}

	server.State = types.MinetestServerStateRunning
	err = w.repos.MinetestServerRepo.Update(server)
	if err != nil {
		return fmt.Errorf("server entity update error: %v", err)
	}

	return nil
}
