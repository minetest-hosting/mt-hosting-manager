package core

import (
	"fmt"
	"mt-hosting-manager/types"
	"time"
)

func (c *Core) SendActivationMail(user *types.User) error {
	latest_mail, err := c.repos.MailQueueRepo.GetLatestByReceiver(user.Mail)
	if err != nil {
		return fmt.Errorf("could not fetch latest mail: %v", err)
	}

	dur := time.Since(time.Unix(latest_mail.Timestamp, 0))
	if dur.Minutes() < 5 {
		return fmt.Errorf("cooldown duration error: %s", dur)
	}

	// TODO: generate code and send mail

	return nil
}

func (c *Core) SendBalanceWarning(user *types.User) error {
	euros := float64(user.Balance) / 100
	return c.repos.MailQueueRepo.Insert(&types.MailQueue{
		Receiver: user.Mail,
		Subject:  fmt.Sprintf("Low balance warning (EUR %.2f)", euros),
		Content:  fmt.Sprintf("Your balance just dropped below the warning limit of EUR %.2f, services will be interrupted if it reaches zero!", euros),
	})
}
