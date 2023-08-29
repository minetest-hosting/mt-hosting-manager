package worker

import (
	"fmt"
	"mt-hosting-manager/api/hetzner_dns"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker/server_setup"
	"time"
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

	if server.ExternalCNAMEDNSID == "" {
		// create new record
		record, err := w.CreateDNSRecord(hetzner_dns.RecordCNAME, server.DNSName, node.Name)
		if err != nil {
			return fmt.Errorf("could not create CNAME record: %v", err)
		}
		server.ExternalCNAMEDNSID = record.ID

	} else {
		// check if record matches config
		record, err := w.hdc.GetRecord(server.ExternalCNAMEDNSID)
		if err != nil {
			return fmt.Errorf("could not fetch current cname with id: '%s': %v", server.ExternalCNAMEDNSID, err)
		}
		if record.Record.Name != server.DNSName || record.Record.Value != node.Name {
			// values changed, remove and recreate
			err = w.RemoveDNSRecord(server.ExternalCNAMEDNSID)
			if err != nil {
				return fmt.Errorf("could not remove record with id '%s', %v", server.ExternalCNAMEDNSID, err)
			}
			created_record, err := w.CreateDNSRecord(hetzner_dns.RecordCNAME, server.DNSName, node.Name)
			if err != nil {
				return fmt.Errorf("could not re-create CNAME record: %v", err)
			}
			server.ExternalCNAMEDNSID = created_record.ID
		}
	}

	// dns propagation time (LE has issues with really _fresh_ records)
	time.Sleep(30 * time.Second)

	client, err := TrySSHConnection(node)
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

	return nil
}
