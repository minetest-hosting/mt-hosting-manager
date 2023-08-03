package worker

import (
	"errors"
	"fmt"
	"mt-hosting-manager/api/hetzner_cloud"
	"mt-hosting-manager/api/hetzner_dns"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker/provision"
	"os"
	"strings"
	"time"

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
				"stage":   os.Getenv("STAGE"),
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

		// replace trailing /64 with 1
		ipv6_parts := strings.Split(resp.Server.PublicNet.IPv6.IP, "/")
		node.IPv6 = fmt.Sprintf("%s1", ipv6_parts[0])
		err = w.repos.UserNodeRepo.Update(node)
		if err != nil {
			return fmt.Errorf("external-id update failed: %v", err)
		}
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

	if a_record == nil {
		// create new record
		logrus.WithFields(logrus.Fields{
			"ipv4": node.IPv4,
			"name": node.Name,
		}).Info("Creating A-Record")

		a_record = &hetzner_dns.Record{
			Type:  hetzner_dns.RecordA,
			Name:  node.Name,
			Value: node.IPv4,
			TTL:   300,
		}
		err := w.hdc.CreateRecord(a_record)
		if err != nil {
			return fmt.Errorf("create a-record error: %v", err)
		}

	} else if a_record.Value != node.IPv4 {
		// update record
		a_record.Value = node.IPv4
		err = w.hdc.UpdateRecord(a_record)
		if err != nil {
			return fmt.Errorf("update a-record error: %v", err)
		}

	}

	if aaaa_record == nil {
		logrus.WithFields(logrus.Fields{
			"ipv6": node.IPv6,
			"name": node.Name,
		}).Info("Creating AAAA-Record")

		a_record = &hetzner_dns.Record{
			Type:  hetzner_dns.RecordAAAA,
			Name:  node.Name,
			Value: node.IPv6,
			TTL:   300,
		}
		err := w.hdc.CreateRecord(a_record)
		if err != nil {
			return fmt.Errorf("create aaaa-record error: %v", err)
		}

	} else if aaaa_record.Value != node.IPv6 {
		// update record
		a_record.Value = node.IPv6
		err = w.hdc.UpdateRecord(a_record)
		if err != nil {
			return fmt.Errorf("update aaaa-record error: %v", err)
		}
	}

	logrus.WithFields(logrus.Fields{
		"external_id": node.ExternalID,
		"ipv4":        node.IPv4,
		"ipv6":        node.IPv6,
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
					"err":       err,
					"try_count": try_count,
				}).Warn("ssh-client failed")
				try_count++
				time.Sleep(10 * time.Second)
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
