package core

import (
	"fmt"
	"mt-hosting-manager/types"

	"github.com/sirupsen/logrus"
)

func (c *Core) SubtractBalance(user_id string, eurocents int64) error {
	before_user, err := c.repos.UserRepo.GetByID(user_id)
	if err != nil {
		return fmt.Errorf("could not fetch before user: %v", err)
	}

	err = c.repos.UserRepo.SubtractBalance(user_id, eurocents)
	if err != nil {
		return fmt.Errorf("could not subtract balance: %v", err)
	}

	after_user, err := c.repos.UserRepo.GetByID(user_id)
	if err != nil {
		return fmt.Errorf("could not fetch after user: %v", err)
	}

	if before_user.Balance >= before_user.WarnBalance && after_user.Balance < before_user.WarnBalance {
		// crossed warning threshold
		c.AddAuditLog(&types.AuditLog{
			Type:   types.AuditLogPaymentWarning,
			UserID: user_id,
			Amount: &after_user.Balance,
		})

		err = c.SendBalanceWarning(before_user)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":            err,
				"user_id":        user_id,
				"before_balance": before_user.Balance,
				"after_balance":  after_user.Balance,
				"warn_balance":   before_user.WarnBalance,
			}).Error("could not send balance warning")
		}
	}

	if before_user.Balance >= 0 && after_user.Balance < 0 {
		// crossed zero threshold
		c.AddAuditLog(&types.AuditLog{
			Type:   types.AuditLogPaymentZero,
			UserID: user_id,
			Amount: &after_user.Balance,
		})

		err = c.SendRemovalNotice(before_user)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":            err,
				"user_id":        user_id,
				"before_balance": before_user.Balance,
				"after_balance":  after_user.Balance,
				"warn_balance":   before_user.WarnBalance,
			}).Error("could not send removal notice")
		}

		nodes, err := c.repos.UserNodeRepo.GetByUserID(before_user.ID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":     err,
				"user_id": before_user.ID,
			}).Error("could not fetch usernodes")
		} else if len(nodes) > 0 {

			for _, node := range nodes {
				j := types.RemoveNodeJob(node)
				err = c.repos.JobRepo.Insert(j)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"err":     err,
						"node_id": node.ID,
					}).Error("could not schedule removal job")
				}
			}
		}

	}

	return nil
}
