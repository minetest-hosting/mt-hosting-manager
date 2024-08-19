package core

import (
	"fmt"
	"mt-hosting-manager/types"

	"github.com/sirupsen/logrus"
)

const SECONDS_IN_A_DAY = 3600 * 24

func (c *Core) Collect(valid_until int64) error {
	list, err := c.repos.UserNodeRepo.Search(&types.UserNodeSearch{ValidUntil: &valid_until})
	if err != nil {
		return fmt.Errorf("could not fetch usernodes: %v", err)
	}
	if len(list) == 0 {
		// nothing to do
		return nil
	}

	// cache nodetypes
	nts, err := c.repos.NodeTypeRepo.GetAll()
	if err != nil {
		return fmt.Errorf("could not fetch nodetypes: %v", err)
	}
	nt_map := map[string]*types.NodeType{}
	for _, nt := range nts {
		nt_map[nt.ID] = nt
	}

	for _, node := range list {
		if node.State != types.UserNodeStateRunning {
			// don't bill stopped nodes
			continue
		}

		nt := nt_map[node.NodeTypeID]
		if nt == nil {
			return fmt.Errorf("nodetype not found: %s", node.NodeTypeID)
		}

		// cost in eurocents
		cost := nt.DailyCost
		err = c.SubtractBalance(node.UserID, cost)
		if err != nil {
			return fmt.Errorf("could not subtract cost '%d' from user '%s': %v", cost, node.UserID, err)
		}

		c.AddAuditLog(&types.AuditLog{
			Type:       types.AuditLogNodeBilled,
			UserID:     node.UserID,
			UserNodeID: &node.ID,
			Amount:     &cost,
		})

		logrus.WithFields(logrus.Fields{
			"UserID":   node.UserID,
			"UserNode": node.ID,
			"Cost":     cost,
		}).Debug("Usernode collected")

		// update valid until time
		node.ValidUntil += SECONDS_IN_A_DAY
		err = c.repos.UserNodeRepo.Update(node)
		if err != nil {
			return fmt.Errorf("could not update usernode '%s': %v", node.ID, err)
		}

	}

	return nil
}
