package worker

import (
	"fmt"
	"mt-hosting-manager/api/hetzner_dns"
	"mt-hosting-manager/types"
	"time"
)

func (w *Worker) ServerPrepareSetup(job *types.Job, node *types.UserNode, server *types.MinetestServer) error {

	server.State = types.MinetestServerStateProvisioning
	err := w.repos.MinetestServerRepo.Update(server)
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
	time.Sleep(20 * time.Second)

	return nil
}
