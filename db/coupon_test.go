package db_test

import (
	"mt-hosting-manager/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCoupon(t *testing.T) {
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

	c := &types.Coupon{
		Code:       "ABCD",
		Name:       "My coupon",
		ValidFrom:  time.Now().Unix(),
		ValidUntil: time.Now().Add(2 * time.Hour).Unix(),
	}

	// create

	assert.NoError(t, repos.CouponRepo.Insert(c))

	// read

	list, err := repos.CouponRepo.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))

	r, err := repos.CouponRepo.IsRedeemed(c.ID, user.ID)
	assert.NoError(t, err)
	assert.False(t, r)

	// redeem

	assert.NoError(t, repos.CouponRepo.Redeem(c.ID, user.ID))

	// read

	r, err = repos.CouponRepo.IsRedeemed(c.ID, user.ID)
	assert.NoError(t, err)
	assert.True(t, r)

	// update

	c.Name = "xy"
	assert.NoError(t, repos.CouponRepo.Update(c))

	// read

	c1, err := repos.CouponRepo.GetByCode(c.Code)
	assert.NoError(t, err)
	assert.Equal(t, c, c1)

	// delete

	assert.NoError(t, repos.CouponRepo.Delete(c.ID))
}
