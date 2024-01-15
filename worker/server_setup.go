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

	w.core.AddAuditLog(&types.AuditLog{
		Type:             types.AuditLogServerSetupStarted,
		UserID:           node.UserID,
		UserNodeID:       &node.ID,
		MinetestServerID: &server.ID,
	})

	server.State = types.MinetestServerStateProvisioning
	err = w.repos.MinetestServerRepo.Update(server)
	if err != nil {
		return fmt.Errorf("server entity update error: %v", err)
	}

	record_name := fmt.Sprintf("%s%s", server.DNSName, w.cfg.DNSRecordSuffix)
	record_value := fmt.Sprintf("%s%s", node.Name, w.cfg.DNSRecordSuffix)
	if server.ExternalCNAMEDNSID == "" {
		// create new record
		record, err := w.CreateDNSRecord(hetzner_dns.RecordCNAME, record_name, record_value)
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
		if record.Record.Name != record_name || record.Record.Value != record_value {
			// values changed, remove and recreate
			err = w.hdc.DeleteRecord(server.ExternalCNAMEDNSID)
			if err != nil {
				return fmt.Errorf("could not remove record with id '%s', %v", server.ExternalCNAMEDNSID, err)
			}
			created_record, err := w.CreateDNSRecord(hetzner_dns.RecordCNAME, record_name, record_value)
			if err != nil {
				return fmt.Errorf("could not re-create CNAME record: %v", err)
			}
			server.ExternalCNAMEDNSID = created_record.ID
		}
	}

	// save external dns id
	err = w.repos.MinetestServerRepo.Update(server)
	if err != nil {
		return fmt.Errorf("mid-setup update failed: %v", err)
	}

	// dns propagation time (LE has issues with really _fresh_ records)
	time.Sleep(10 * time.Second)

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

	w.core.AddAuditLog(&types.AuditLog{
		Type:             types.AuditLogServerSetupFinished,
		UserID:           node.UserID,
		UserNodeID:       &node.ID,
		MinetestServerID: &server.ID,
	})

	return nil
}
