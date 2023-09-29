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

	if latest_mail != nil {
		// rate-limit per mail-receiver
		dur := time.Since(time.Unix(latest_mail.Timestamp, 0))
		if dur.Hours() < 1 {
			return fmt.Errorf("cooldown duration error: %s", dur)
		}
	}

	if user.ActivationCode == "" {
		//create activation code
		user.ActivationCode = RandSeq(8)
		err = c.repos.UserRepo.Update(user)
		if err != nil {
			return err
		}
	}

	url := fmt.Sprintf("%s/#/activate/%s/%s", c.cfg.BaseURL, user.Mail, user.ActivationCode)

	return c.repos.MailQueueRepo.Insert(&types.MailQueue{
		Receiver: user.Mail,
		Subject:  "Minetest hosting activation",
		Content:  fmt.Sprintf("Please visit %s to activate your minetest-hosting account", url),
	})
}

func (c *Core) SendBalanceWarning(user *types.User) error {
	euros := float64(user.Balance) / 100
	return c.repos.MailQueueRepo.Insert(&types.MailQueue{
		Receiver: user.Mail,
		Subject:  fmt.Sprintf("Low balance warning (EUR %.2f)", euros),
		Content:  fmt.Sprintf("Your balance just dropped below the warning limit of EUR %.2f, services will be interrupted if it reaches zero!", euros),
	})
}

func (c *Core) SendRemovalNotice(user *types.User) error {
	return c.repos.MailQueueRepo.Insert(&types.MailQueue{
		Receiver: user.Mail,
		Subject:  "Server removal notice",
		Content:  "Due to insufficient balance on your account all your servers and nodes are now being removed!",
	})
}
