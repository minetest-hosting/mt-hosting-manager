package db_test

import (
	"mt-hosting-manager/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository(t *testing.T) {
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
	user.ExternalID = "def"
	assert.NoError(t, repos.UserRepo.Update(user))

	// existing user
	u, err := repos.UserRepo.GetByName("somedude")
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, types.UserStateActive, u.State)

	// balance
	assert.Equal(t, int64(0), u.Balance)
	// add
	assert.NoError(t, repos.UserRepo.AddBalance(u.ID, 100))
	u, err = repos.UserRepo.GetByID(u.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), u.Balance)
	// subtract
	assert.NoError(t, repos.UserRepo.AddBalance(u.ID, -100))
	u, err = repos.UserRepo.GetByID(u.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), u.Balance)

	// non existent user
	u, err = repos.UserRepo.GetByName("non@existent")
	assert.NoError(t, err)
	assert.Nil(t, u)

	// get all users
	list, err := repos.UserRepo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))

	// get all users by role
	list, err = repos.UserRepo.GetAllByRole(types.UserRoleUser)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))

	list, err = repos.UserRepo.GetAllByRole(types.UserRoleAdmin)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 0, len(list))

	// count
	count, err := repos.UserRepo.CountAll()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// search existing
	limit := 10
	name_like := "some%"
	list, err = repos.UserRepo.Search(&types.UserSearch{NameLike: &name_like, Limit: &limit})
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))

	// search non-existing
	name_like = "another%"
	list, err = repos.UserRepo.Search(&types.UserSearch{NameLike: &name_like})
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 0, len(list))

	// delete
	assert.NoError(t, repos.UserRepo.Delete(user.ID))
}
