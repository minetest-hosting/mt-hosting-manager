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

func (w *Worker) NodeProvision(job *types.Job) error {
	node, err := w.repos.UserNodeRepo.GetByID(*job.UserNodeID)
	if err != nil {
		return err
	}
	if node == nil {
		return errors.New("node not found")
	}

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

	records, err := w.hdc.GetRecords()
	if err != nil {
		return fmt.Errorf("fetch records error: %v", err)
	}

	err = w.UpdateDNSRecord(records, hetzner_dns.RecordA, node.Name, node.IPv4)
	if err != nil {
		return err
	}

	err = w.UpdateDNSRecord(records, hetzner_dns.RecordAAAA, node.Name, node.IPv6)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"external_id": node.ExternalID,
		"ipv4":        node.IPv4,
		"ipv6":        node.IPv6,
	}).Info("Executing provisioning")

	client, err := TrySSHConnection(node)
	if err != nil {
		return err
	}

	err = provision.Provision(client)
	if err != nil {
		return fmt.Errorf("provision error: %v", err)
	}
	err = client.Close()
	if err != nil {
		return err
	}

	node.State = types.UserNodeStateRunning
	return w.repos.UserNodeRepo.Update(node)
}
