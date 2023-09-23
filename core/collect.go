package core

import (
	"fmt"
	"mt-hosting-manager/types"
	"time"

	"github.com/sirupsen/logrus"
)

const SECONDS_IN_A_DAY = 3600 * 24

func (c *Core) Collect(last_collected_time int64) error {
	now := time.Now().Unix()
	list, err := c.repos.UserNodeRepo.GetByLastCollectedTime(last_collected_time)
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
		delta := now - node.LastCollectedTime
		delta_days := delta / SECONDS_IN_A_DAY
		if delta_days > 0 {
			nt := nt_map[node.NodeTypeID]
			if nt == nil {
				return fmt.Errorf("nodetype not found: %s", node.NodeTypeID)
			}

			// cost in eurocents
			cost := nt.DailyCost * delta_days
			err = c.repos.UserRepo.AddBalance(node.UserID, cost*-1)
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
				"Days":     delta_days,
			}).Debug("Usernode collected")

			// update last collected time
			node.LastCollectedTime += (SECONDS_IN_A_DAY * delta_days)
			err = c.repos.UserNodeRepo.Update(node)
			if err != nil {
				return fmt.Errorf("could not update usernode '%s': %v", node.ID, err)
			}
		}
	}

	return nil
}
