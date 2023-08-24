package core_test

import (
	"mt-hosting-manager/core"
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
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
	now := time.Now()

	user := &types.User{
		ID:      uuid.NewString(),
		State:   types.UserStateActive,
		Name:    "testuser",
		Mail:    "test@user",
		Balance: "0",
	}
	assert.NoError(t, repos.UserRepo.Insert(user))

	nt := &types.NodeType{
		ID:         uuid.NewString(),
		State:      types.NodeTypeStateActive,
		Provider:   types.ProviderHetzner,
		ServerType: "test",
		DailyCost:  "0.4",
	}
	assert.NoError(t, repos.NodeTypeRepo.Insert(nt))

	node := &types.UserNode{
		ID:                uuid.NewString(),
		UserID:            user.ID,
		NodeTypeID:        nt.ID,
		LastCollectedTime: now.Unix(),
		State:             types.UserNodeStateRunning,
	}
	assert.NoError(t, repos.UserNodeRepo.Insert(node))

	// 1 day later
	then := now.Add(time.Hour * 25)
	a, err := core.Collect(repos, user.ID, then)
	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, types.DEFAULT_CURRENCY, a.CurrencyCode())
	assert.Equal(t, "0.4", a.Number())

	// 1 day and an hour
	then = now.Add(time.Hour * 26)
	a, err = core.Collect(repos, user.ID, then)
	assert.NoError(t, err)
	assert.Nil(t, a)

	// 2 days and an hour
	then = now.Add(time.Hour * 49)
	a, err = core.Collect(repos, user.ID, then)
	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, "0.4", a.Number())

	// 1 week and an hour
	then = now.Add(time.Hour * ((24 * 7) + 1))
	a, err = core.Collect(repos, user.ID, then)
	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, "2.0", a.Number())

}
