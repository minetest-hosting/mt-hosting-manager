package core

import (
	"fmt"
	"mt-hosting-manager/notify"
	"mt-hosting-manager/types"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (c *Core) CreateUser(name, external_id, hash string, t types.UserType, role types.UserRole) (*types.User, error) {
	var balance int64 = 0

	if c.cfg.InitialBalance != "" {
		b, err := strconv.ParseInt(c.cfg.InitialBalance, 10, 32)
		if err == nil {
			balance = b
		}
	}

	user := &types.User{
		ID:         uuid.NewString(),
		ExternalID: external_id,
		Hash:       hash,
		Name:       name,
		State:      types.UserStateActive,
		Created:    time.Now().Unix(),
		LastLogin:  time.Now().Unix(),
		Balance:    balance,
		Currency:   "EUR",
		Type:       t,
		Role:       role,
	}

	logrus.WithFields(logrus.Fields{
		"name":        user.Name,
		"type":        user.Type,
		"external_id": user.ExternalID,
	}).Debug("created new user")

	notify.Send(&notify.NtfyNotification{
		Title:    fmt.Sprintf("New user signed up: %s", user.Name),
		Message:  fmt.Sprintf("Name: %s, Auth: %s", user.Name, user.Type),
		Priority: 3,
		Tags:     []string{"new"},
	}, true)

	c.AddAuditLog(&types.AuditLog{
		Type:   types.AuditLogUserCreated,
		UserID: user.ID,
	})

	return user, c.repos.UserRepo.Insert(user)
}
