package core_test

import (
	"mt-hosting-manager/core"
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func SetupRepos(t *testing.T) *db.Repositories {
	tmpdir, err := os.MkdirTemp(os.TempDir(), "mt-hosting")
	assert.NoError(t, err)

	db_, err := db.Init(tmpdir)
	assert.NoError(t, err)
	assert.NotNil(t, db_)

	err = db.Migrate(db_)
	assert.NoError(t, err)

	return db.NewRepositories(db_)
}

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
		UserID:            user.ID,
		NodeTypeID:        nt.ID,
		LastCollectedTime: ts - (core.SECONDS_IN_A_DAY * 2),
	}
	assert.NoError(t, repos.UserNodeRepo.Insert(un))

	assert.NoError(t, c.Collect(ts-core.SECONDS_IN_A_DAY))

	user, err := repos.UserRepo.GetByID(user.ID)
	assert.NoError(t, err)

	assert.Equal(t, int64(60), user.Balance)
}
