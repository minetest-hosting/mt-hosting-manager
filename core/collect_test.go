package core_test

import (
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCollect(t *testing.T) {
	repos := SetupRepos(t)
	c := core.New(repos, types.NewConfig())

	user := &types.User{
		Balance: 100,
	}
	assert.NoError(t, repos.UserRepo.Insert(user))

	nt := &types.NodeType{DailyCost: 20}
	assert.NoError(t, repos.NodeTypeRepo.Insert(nt))

	ts := time.Now().Unix()

	un := &types.UserNode{
		UserID:     user.ID,
		NodeTypeID: nt.ID,
		ValidUntil: ts + (core.SECONDS_IN_A_DAY * 1),
		State:      types.UserNodeStateRunning,
	}
	assert.NoError(t, repos.UserNodeRepo.Insert(un))

	// no updates
	assert.NoError(t, c.Collect(ts+300))
	user, err := repos.UserRepo.GetByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), user.Balance)

	// 1 day
	assert.NoError(t, c.Collect(ts+300+core.SECONDS_IN_A_DAY))
	user, err = repos.UserRepo.GetByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(80), user.Balance)
}
