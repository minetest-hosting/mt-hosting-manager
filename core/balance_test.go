package core_test

import (
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetBalanceDuration(t *testing.T) {
	nt := &types.NodeType{DailyCost: 0.4}
	d := core.GetBalanceDuration(nt, 10)
	assert.Equal(t, time.Hour*24*25, d)
}

func TestGetNodeCost(t *testing.T) {
	nt := &types.NodeType{DailyCost: 0.4}
	cost := core.GetNodeCost(nt, time.Hour*24*5)
	assert.Equal(t, 2.0, cost)
}
