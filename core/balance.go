package core

import (
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"
	"time"
)

func GetNodeCost(nt *types.NodeType, d time.Duration) float64 {
	return nt.DailyCost * (d.Hours() / 24)
}

func GetNodesCost(repos *db.Repositories, nodes []*types.UserNode, d time.Duration) (float64, error) {
	total := 0.0

	for _, node := range nodes {
		nt, err := repos.NodeTypeRepo.GetByID(node.NodeTypeID)
		if err != nil {
			return 0, err
		}

		total += GetNodeCost(nt, d)
	}

	return total, nil
}

func GetBalanceDuration(nt *types.NodeType, balance float64) time.Duration {
	return time.Duration(balance/nt.DailyCost) * time.Hour * 24
}
