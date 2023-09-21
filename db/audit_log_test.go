package db_test

import (
	"mt-hosting-manager/types"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuditLog(t *testing.T) {
	repos := SetupRepos(t)

	user := &types.User{ID: uuid.NewString()}
	assert.NoError(t, repos.UserRepo.Insert(user))

	node := "456"
	log := &types.AuditLog{
		Type:       types.AuditLogNodeCreated,
		UserID:     user.ID,
		UserNodeID: &node,
	}
	err := repos.AuditLogRepo.Insert(log)
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), log.Timestamp)

	userid1 := "1234"
	list, err := repos.AuditLogRepo.Search(&types.AuditLogSearch{
		UserID:        &userid1,
		FromTimestamp: log.Timestamp - 1,
		ToTimestamp:   log.Timestamp + 1,
	})
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 0, len(list))

	list, err = repos.AuditLogRepo.Search(&types.AuditLogSearch{
		UserID:        &user.ID,
		FromTimestamp: log.Timestamp - 1,
		ToTimestamp:   log.Timestamp + 1,
	})
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
}
