package core

import (
	"mt-hosting-manager/types"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (c *Core) CreateUser(mail string, t types.UserType, role types.UserRole, mailVerified bool) (*types.User, error) {
	var balance int64 = 0
	var warnBalance int64 = 200

	if c.cfg.InitialBalance != "" {
		b, err := strconv.ParseInt(c.cfg.InitialBalance, 10, 32)
		if err == nil {
			balance = b
		}
	}

	user := &types.User{
		ID:           uuid.NewString(),
		Name:         mail,
		Mail:         mail,
		MailVerified: mailVerified,
		State:        types.UserStateActive,
		Created:      time.Now().Unix(),
		Balance:      balance,
		WarnBalance:  warnBalance,
		Currency:     "EUR",
		Type:         t,
		Role:         role,
	}
	return user, c.repos.UserRepo.Insert(user)
}
