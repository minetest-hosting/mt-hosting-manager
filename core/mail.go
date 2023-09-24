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
