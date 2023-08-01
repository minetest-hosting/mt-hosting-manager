package worker

import (
	"errors"
	"fmt"
	"mt-hosting-manager/api/hetzner_cloud"
	"mt-hosting-manager/types"

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

	nodetype, err := w.repos.NodeTypeRepo.GetByID(node.NodeTypeID)
	if err != nil {
		return err
	}
	if nodetype == nil {
		return errors.New("nodetype not found")
	}

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
		"hetzner_id": resp.ID,
		"status":     resp.Status,
	}).Info("Server created")

	node.ExternalID = fmt.Sprintf("%d", resp.ID)
	err = w.repos.UserNodeRepo.Update(node)
	if err != nil {
		return fmt.Errorf("external-id update failed: %v", err)
	}

	node.State = types.UserNodeStateRunning
	return w.repos.UserNodeRepo.Update(node)
}
