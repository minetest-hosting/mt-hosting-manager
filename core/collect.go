package core

import (
	"fmt"
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"
	"time"

	"github.com/bojanz/currency"
)

func Collect(repos *db.Repositories, userID string, now time.Time) (*currency.Amount, error) {
	nodes, err := repos.UserNodeRepo.GetByUserIDAndState(userID, types.UserNodeStateRunning)
	if err != nil {
		return nil, fmt.Errorf("usernode fetch failed: %v", err)
	}

	if len(nodes) == 0 {
		return nil, nil
	}

	cost_eur, err := currency.NewAmount("0", types.DEFAULT_CURRENCY)
	if err != nil {
		return nil, fmt.Errorf("currency parse failed: %v", err)
	}

	cost_changed := false

	for _, node := range nodes {
		last_collected := time.Unix(node.LastCollectedTime, 0)
		dur := now.Sub(last_collected)
		hours := dur.Hours()

		if hours <= 24 {
			// running less than a day, skip collection
			continue
		}

		// at least one node collected
		cost_changed = true

		nt, err := repos.NodeTypeRepo.GetByID(node.NodeTypeID)
		if err != nil {
			return nil, fmt.Errorf("nodetype fetch failed: %v", err)
		}
		if nt == nil {
			return nil, fmt.Errorf("nodetype not found for node '%s'", node.NodeTypeID)
		}

		dailycost, err := currency.NewAmount(nt.DailyCost, types.DEFAULT_CURRENCY)
		if err != nil {
			return nil, fmt.Errorf("currency parse failed: %v", err)
		}

		for hours > 24 {
			cost_eur, err = cost_eur.Add(dailycost)
			if err != nil {
				return nil, fmt.Errorf("currency add failed: %v", err)
			}

			hours -= 24
		}

		// reset last collected time
		node.LastCollectedTime = now.Unix()
	}

	if !cost_changed {
		// nothing changed
		return nil, nil
	}

	for _, node := range nodes {
		err = repos.UserNodeRepo.Update(node)
		if err != nil {
			return nil, fmt.Errorf("node update failed for '%s': %v", node.ID, err)
		}
	}

	return &cost_eur, nil
}
