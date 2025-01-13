package core

import (
	"fmt"
	"mt-hosting-manager/types"
	"strings"

	"github.com/wneessen/go-mail"
)

func (c *Core) SendMail(to, subject, body string) error {
	m := mail.NewMsg()
	m.From(c.cfg.MailAddress)
	m.To(to)
	m.Subject(subject)
	m.SetBodyString(mail.TypeTextPlain, body)

	client, err := mail.NewClient(c.cfg.MailHost,
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithPort(587),
		mail.WithUsername(c.cfg.MailAddress), mail.WithPassword(c.cfg.MailPassword))
	if err != nil {
		return fmt.Errorf("new smtp client error: %v", err)
	}

	err = client.DialAndSend(m)
	if err != nil {
		return fmt.Errorf("smtp send error: %v", err)
	}

	return nil
}

func (c *Core) SendBalanceWarningMail(user *types.User) error {
	body := strings.Builder{}
	body.WriteString(fmt.Sprintf("Warning: your luanti hosting balance is low (EUR %.2f)\n", float64(user.Balance)/100))
	body.WriteString("Your nodes will be automatically removed when the balance is reaching zero.")

	return c.SendMail(user.Mail, "Luanti hosting balance warning", body.String())
}

func (c *Core) SendBalanceZeroMail(user *types.User) error {
	body := strings.Builder{}
	body.WriteString(fmt.Sprintf("Warning: your luanti hosting balance has reached EUR %.2f\n", float64(user.Balance)/100))
	body.WriteString("Your nodes will now be automatically removed.")

	return c.SendMail(user.Mail, "Luanti hosting removal notice", body.String())
}
