package db_test

import (
	"mt-hosting-manager/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository(t *testing.T) {
	repos := setupRepos(t)

	assert.NoError(t, repos.UserRepo.Insert(&types.User{
		Name:       "Some dude",
		State:      types.UserStateActive,
		Mail:       "x@y.ch",
		Created:    time.Now().Unix(),
		ExternalID: "abc",
		Type:       types.UserTypeGithub,
		Role:       types.UserRoleUser,
	}))

	// existing user
	u, err := repos.UserRepo.GetByMail("x@y.ch")
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, types.UserStateActive, u.State)

	// non existent user
	u, err = repos.UserRepo.GetByMail("non@existent")
	assert.NoError(t, err)
	assert.Nil(t, u)

	// all users
	list, err := repos.UserRepo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
}
