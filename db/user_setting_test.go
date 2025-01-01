package db_test

import (
	"mt-hosting-manager/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserSetting(t *testing.T) {
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

	// create or update
	assert.NoError(t, repos.UserSettingRepo.Set(&types.UserSetting{UserID: user.ID, Key: "mykey", Value: "myvalue"}))

	// get
	list, err := repos.UserSettingRepo.GetByUserID(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "myvalue", list[0].Value)

	// delete
	assert.NoError(t, repos.UserSettingRepo.Delete(user.ID, "mykey"))

	// get empty
	list, err = repos.UserSettingRepo.GetByUserID(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 0, len(list))

}
