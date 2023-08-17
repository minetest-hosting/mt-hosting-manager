package worker

import (
	"errors"
	"fmt"
	"mt-hosting-manager/api/hetzner_dns"
	"mt-hosting-manager/types"

	"github.com/sirupsen/logrus"
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
		logrus.WithFields(logrus.Fields{
			"ExternalID": node.ExternalID,
		}).Warn("Server instance not found, not deleting anything")
	}

	records, err := w.hdc.GetRecords()
	if err != nil {
		return fmt.Errorf("fetch records error: %v", err)
	}

	err = w.RemoveDNSRecord(records, hetzner_dns.RecordA, node.Name)
	if err != nil {
		return err
	}

	err = w.RemoveDNSRecord(records, hetzner_dns.RecordAAAA, node.Name)
	if err != nil {
		return err
	}

	return w.repos.UserNodeRepo.Delete(node.ID)
}
