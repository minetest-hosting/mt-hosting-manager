package worker

import (
	"errors"
	"fmt"
	"mt-hosting-manager/api/hetzner_cloud"
	"mt-hosting-manager/api/hetzner_dns"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker/provision"
	"strings"

	"github.com/sirupsen/logrus"
)

func (w *Worker) NodeProvision(job *types.Job, status func(string, int)) error {
	node, err := w.repos.UserNodeRepo.GetByID(*job.UserNodeID)
	if err != nil {
		return err
	}
	if node == nil {
		return errors.New("node not found")
	}

	w.core.AddAuditLog(&types.AuditLog{
		Type:       types.AuditLogNodeProvisioningStarted,
		UserID:     node.UserID,
		UserNodeID: &node.ID,
	})

	node.State = types.UserNodeStateProvisioning
	err = w.repos.UserNodeRepo.Update(node)
	if err != nil {
		return err
	}

	nodetype, err := w.repos.NodeTypeRepo.GetByID(node.NodeTypeID)
	if err != nil {
		return err
	}
	if nodetype == nil {
		return errors.New("nodetype not found")
	}

	if node.ExternalID == "" {
		// create new instance
		logrus.WithFields(logrus.Fields{
			"node_id": node.ID,
		}).Info("Creating new hetzner server instance")
		status("creating new server", 10)

		csr := &hetzner_cloud.CreateServerRequest{
			Image: "ubuntu-22.04",
			Labels: map[string]string{
				"node_id": node.ID,
				"stage":   w.cfg.Stage,
			},
			Location: hetzner_cloud.LocationType(nodetype.Location),
			Name:     node.Name,
			PublicNet: &hetzner_cloud.PublicNet{
				EnableIPv4: true,
				EnableIPv6: true,
			},
			ServerType:       nodetype.ServerType,
			SSHKeys:          []string{"minetest@keymaster", "thomas@keymaster"},
			StartAfterCreate: true,
		}
		resp, err := w.hcc.CreateServer(csr)
		if err != nil {
			return fmt.Errorf("create server failed: %v", err)
		}
		logrus.WithFields(logrus.Fields{
			"hetzner_id": resp.Server.ID,
			"status":     resp.Server.Status,
		}).Info("Server created")

		node.ExternalID = fmt.Sprintf("%d", resp.Server.ID)
		node.IPv4 = resp.Server.PublicNet.IPv4.IP

		// replace trailing /64 with 1
		ipv6_parts := strings.Split(resp.Server.PublicNet.IPv6.IP, "/")
		node.IPv6 = fmt.Sprintf("%s1", ipv6_parts[0])
		err = w.repos.UserNodeRepo.Update(node)
		if err != nil {
			return fmt.Errorf("external-id update failed: %v", err)
		}
	}

	record_name := fmt.Sprintf("%s%s", node.Name, w.cfg.DNSRecordSuffix)

	if node.ExternalIPv4DNSID == "" {
		status("creating dns A-record", 20)

		record, err := w.CreateDNSRecord(hetzner_dns.RecordA, record_name, node.IPv4)
		if err != nil {
			return fmt.Errorf("could not create A-record: %v", err)
		}
		node.ExternalIPv4DNSID = record.ID
		err = w.repos.UserNodeRepo.Update(node)
		if err != nil {
			return fmt.Errorf("ipv4 record update failed: %v", err)
		}
	}

	if node.ExternalIPv6DNSID == "" {
		status("creating dns AAAA-record", 30)

		record, err := w.CreateDNSRecord(hetzner_dns.RecordAAAA, record_name, node.IPv6)
		if err != nil {
			return fmt.Errorf("could not create AAAA-record: %v", err)
		}
		node.ExternalIPv6DNSID = record.ID
		err = w.repos.UserNodeRepo.Update(node)
		if err != nil {
			return fmt.Errorf("ipv6 record update failed: %v", err)
		}
	}

	logrus.WithFields(logrus.Fields{
		"external_id": node.ExternalID,
		"ipv4":        node.IPv4,
		"ipv6":        node.IPv6,
	}).Info("Executing provisioning")

	status("waiting for node-startup", 40)
	client, err := TrySSHConnection(node)
	if err != nil {
		return err
	}

	err = provision.Provision(client, w.cfg, node.UserID, status)
	if err != nil {
		return fmt.Errorf("provision error: %v", err)
	}
	err = client.Close()
	if err != nil {
		return err
	}

	status("done", 100)

	node.State = types.UserNodeStateRunning

	w.core.AddAuditLog(&types.AuditLog{
		Type:       types.AuditLogNodeProvisioningFinished,
		UserID:     node.UserID,
		UserNodeID: &node.ID,
	})

	return w.repos.UserNodeRepo.Update(node)
}
