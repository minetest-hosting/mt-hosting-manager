package db_test

import (
	"mt-hosting-manager/types"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPaymentTransaction(t *testing.T) {
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

	tx := &types.PaymentTransaction{
		TransactionID: uuid.NewString(),
		UserID:        user.ID,
		Created:       50,
	}
	assert.NoError(t, repos.PaymentTransactionRepo.Insert(tx))

	tx.Amount = 100
	assert.NoError(t, repos.PaymentTransactionRepo.Update(tx))

	list, err := repos.PaymentTransactionRepo.GetByUserID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, int64(100), list[0].Amount)

	list, err = repos.PaymentTransactionRepo.Search(&types.PaymentTransactionSearch{
		FromTimestamp: 0,
		ToTimestamp:   100,
		UserID:        &user.ID,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, int64(100), list[0].Amount)

	assert.NoError(t, repos.PaymentTransactionRepo.Delete(tx.ID))
}
