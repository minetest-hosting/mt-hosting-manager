package db_test

import (
	"mt-hosting-manager/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserNode(t *testing.T) {
	repos := SetupRepos(t)

	user := &types.User{
		Name:       "somedude",
		State:      types.UserStateActive,
		Created:    time.Now().Unix(),
		ExternalID: "abc",
		Type:       types.UserTypeGithub,
		Role:       types.UserRoleUser,
	}
	assert.NoError(t, repos.UserRepo.Insert(user))

	nt := &types.NodeType{
		Provider:   types.ProviderHetzner,
		ServerType: "cx11",
		Name:       "",
	}
	assert.NoError(t, repos.NodeTypeRepo.Insert(nt))

	un := &types.UserNode{
		UserID:     user.ID,
		NodeTypeID: nt.ID,
		State:      types.UserNodeStateCreated,
	}
	assert.NoError(t, repos.UserNodeRepo.Insert(un))

	un.Alias = "xy"
	assert.NoError(t, repos.UserNodeRepo.Update(un))

	list, err := repos.UserNodeRepo.Search(&types.UserNodeSearch{
		ID: &un.ID,
	})
	assert.NoError(t, err)
	assert.True(t, len(list) > 0)

	zero := int64(0)
	list, err = repos.UserNodeRepo.Search(&types.UserNodeSearch{
		ValidUntil: &zero,
	})
	assert.NoError(t, err)
	assert.True(t, len(list) == 0)

	assert.NoError(t, repos.UserNodeRepo.Delete(un.ID))
}
