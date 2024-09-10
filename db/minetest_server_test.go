package db_test

import (
	"mt-hosting-manager/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMinetestServer(t *testing.T) {
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

	mt := &types.MinetestServer{
		UserNodeID: un.ID,
	}
	assert.NoError(t, repos.MinetestServerRepo.Insert(mt))

	mt.Admin = "xy"
	assert.NoError(t, repos.MinetestServerRepo.Update(mt))

	list, err := repos.MinetestServerRepo.Search(&types.MinetestServerSearch{
		NodeID: &un.ID,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))

	assert.NoError(t, repos.MinetestServerRepo.Delete(mt.ID))
}
