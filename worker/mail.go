package worker

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
)

func (w *Worker) MailJob() {
	for w.running.Load() {

		mails, err := w.repos.MailQueueRepo.GetByState(types.MailQueueStateCreated)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("mail queue fetch error")
			time.Sleep(10 * time.Second)
			continue
		}

		for _, mail := range mails {
			lf := logrus.Fields{
				"receiver": mail.Receiver,
				"subject":  mail.Subject,
			}
			logrus.WithFields(lf).Info("sending mail")

			err = w.SendMail(mail)
			if err != nil {
				lf["err"] = err
				logrus.WithFields(lf).Error("send-mail error")

				mail.State = types.MailQueueStateDoneFailure
			} else {
				mail.State = types.MailQueueStateDoneSuccess
			}
			err = w.repos.MailQueueRepo.Update(mail)
			if err != nil {
				logrus.WithError(err).Error("mail-update error")
			}
		}

		time.Sleep(time.Second)
	}
}

func (w *Worker) SendMail(m *types.MailQueue) error {
	e := &email.Email{
		To:      []string{m.Receiver},
		From:    fmt.Sprintf("Minetest hosting <%s>", w.cfg.MailAddress),
		Subject: m.Subject,
		Text:    []byte(m.Content),
	}

	err := e.Send(fmt.Sprintf("%s:587", w.cfg.MailHost), smtp.PlainAuth("", w.cfg.MailAddress, w.cfg.MailPassword, w.cfg.MailHost))
	if err != nil {
		return fmt.Errorf("could not send mail: %v", err)
	}

	return nil
}
