package worker

import (
	"errors"
	"fmt"
	"mt-hosting-manager/api/hetzner_cloud"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker/provision"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func (w *Worker) NodeProvision(job *types.Job) error {
	node, err := w.repos.UserNodeRepo.GetByID(*job.UserNodeID)
	if err != nil {
		return err
	}
	if node == nil {
		return errors.New("node not found")
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
			},
			Location: hetzner_cloud.LocationNuernberg,
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
		node.IPv6 = resp.Server.PublicNet.IPv6.IP
		err = w.repos.UserNodeRepo.Update(node)
		if err != nil {
			return fmt.Errorf("external-id update failed: %v", err)
		}
	}

	logrus.WithFields(logrus.Fields{
		"external_id": node.ExternalID,
		"ipv4":        node.IPv4,
	}).Info("Executing provisioning")

	var client *ssh.Client
	try_count := 0
	for {
		client, err = provision.CreateClient(fmt.Sprintf("%s:22", node.IPv4))
		if err != nil {
			if try_count > 5 {
				return fmt.Errorf("ssh-client connection failed: %v", err)
			} else {
				logrus.WithFields(logrus.Fields{
					"err": err,
				}).Warn("ssh-client failed")
				try_count++
			}
		} else {
			break
		}
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
