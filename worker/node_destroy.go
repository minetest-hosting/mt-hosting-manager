package worker

import (
	"errors"
	"fmt"
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

	if node.ExternalIPv4DNSID != "" {
		err = w.RemoveDNSRecord(node.ExternalIPv4DNSID)
		if err != nil {
			return fmt.Errorf("could not remove A-record: %v", err)
		}
	}

	if node.ExternalIPv6DNSID != "" {
		err = w.RemoveDNSRecord(node.ExternalIPv6DNSID)
		if err != nil {
			return fmt.Errorf("could not remove AAAA-record: %v", err)
		}
	}

	return w.repos.UserNodeRepo.Delete(node.ID)
}
