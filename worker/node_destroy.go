package worker

import (
	"errors"
	"fmt"
	"mt-hosting-manager/api/hetzner_dns"
	"mt-hosting-manager/types"
)

func (w *Worker) NodeDestroy(job *types.Job) error {
	node, err := w.repos.UserNodeRepo.GetByID(*job.UserNodeID)
	if err != nil {
		return err
	}
	if node == nil {
		return errors.New("node not found")
	}

	err = w.hcc.DeleteServer(node.ExternalID)
	if err != nil {
		return fmt.Errorf("delete server failed: %v", err)
	}

	var a_record *hetzner_dns.Record
	var aaaa_record *hetzner_dns.Record

	records, err := w.hdc.GetRecords()
	if err != nil {
		return fmt.Errorf("fetch records error: %v", err)
	}
	for _, r := range records.Records {
		if r.Type == hetzner_dns.RecordA && r.Name == node.Name {
			a_record = r
			continue
		}
		if r.Type == hetzner_dns.RecordAAAA && r.Name == node.Name {
			aaaa_record = r
			continue
		}
	}

	if a_record != nil {
		err = w.hdc.DeleteRecord(a_record.ID)
		if err != nil {
			return fmt.Errorf("delete a-record error: %v", err)
		}
	}
	if aaaa_record != nil {
		err = w.hdc.DeleteRecord(aaaa_record.ID)
		if err != nil {
			return fmt.Errorf("delete aaaa-record error: %v", err)
		}
	}

	return w.repos.UserNodeRepo.Delete(node.ID)
}
