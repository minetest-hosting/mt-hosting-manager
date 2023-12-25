package worker

import (
	"errors"
	"fmt"
	"mt-hosting-manager/api/hetzner_dns"
	"mt-hosting-manager/notify"
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

	nt, err := w.repos.NodeTypeRepo.GetByID(node.NodeTypeID)
	if err != nil {
		return err
	}
	if nt == nil {
		return errors.New("node-type not found")
	}

	user, err := w.repos.UserRepo.GetByID(node.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	servers, err := w.repos.MinetestServerRepo.GetByNodeID(node.ID)
	if err != nil {
		return fmt.Errorf("could not fetch servers: %v", err)
	}

	for _, server := range servers {
		// remove CNAME record
		if server.ExternalCNAMEDNSID != "" {
			err = w.hdc.DeleteRecord(server.ExternalCNAMEDNSID)
			if err != nil && err != hetzner_dns.ErrRecordNotFound {
				return fmt.Errorf("could not remove cname (id: %s) of server %s: %v", server.ExternalCNAMEDNSID, server.DNSName, err)
			}

			server.ExternalCNAMEDNSID = ""
			server.State = types.MinetestServerStateDecommissioned
			err = w.repos.MinetestServerRepo.Update(server)
			if err != nil {
				return fmt.Errorf("could not update server entry '%s': %v", server.ID, err)
			}
		}
	}

	if node.ExternalID != "" {
		err = w.hcc.DeleteServer(node.ExternalID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"ExternalID": node.ExternalID,
			}).Warn("Server instance not found, not deleting anything")
		}
		node.ExternalID = ""
		err = w.repos.UserNodeRepo.Update(node)
		if err != nil {
			return fmt.Errorf("could not update usernode: %v", err)
		}
	}

	if node.ExternalIPv4DNSID != "" {
		err = w.hdc.DeleteRecord(node.ExternalIPv4DNSID)
		if err != nil && err != hetzner_dns.ErrRecordNotFound {
			return fmt.Errorf("could not remove A-record: %v", err)
		}
		node.ExternalIPv4DNSID = ""
		err = w.repos.UserNodeRepo.Update(node)
		if err != nil {
			return fmt.Errorf("could not update usernode: %v", err)
		}
	}

	if node.ExternalIPv6DNSID != "" {
		err = w.hdc.DeleteRecord(node.ExternalIPv6DNSID)
		if err != nil && err != hetzner_dns.ErrRecordNotFound {
			return fmt.Errorf("could not remove AAAA-record: %v", err)
		}
		node.ExternalIPv6DNSID = ""
		err = w.repos.UserNodeRepo.Update(node)
		if err != nil {
			return fmt.Errorf("could not update usernode: %v", err)
		}
	}

	w.core.AddAuditLog(&types.AuditLog{
		Type:       types.AuditLogNodeRemoved,
		UserID:     node.UserID,
		UserNodeID: &node.ID,
	})

	notify.Send(&notify.NtfyNotification{
		Title:    fmt.Sprintf("Node removed by %s (Type: %s)", user.Name, nt.Name),
		Message:  fmt.Sprintf("User: %s, Node-type %s, Name: %s", user.Name, nt.Name, node.Name),
		Priority: 3,
		Tags:     []string{"computer", "x"},
	}, true)

	node.State = types.UserNodeStateDecommissioned

	return w.repos.UserNodeRepo.Update(node)
}
