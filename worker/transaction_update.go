package worker

import (
	"mt-hosting-manager/types"
	"time"

	"github.com/sirupsen/logrus"
)

func (w *Worker) TransactionUpdateJob() {
	for w.running.Load() {
		pending := types.PaymentStatePending
		q := &types.PaymentTransactionSearch{
			FromTimestamp: time.Now().Add(1 * time.Hour).Unix(),
			ToTimestamp:   time.Now().Unix(),
			State:         &pending,
		}
		list, err := w.repos.PaymentTransactionRepo.Search(q)
		if err != nil {
			logrus.WithError(err).Error("tx search")
		}

		for _, tx := range list {
			_, err = w.core.CheckTransaction(tx.ID)
			if err != nil {
				logrus.WithError(err).Error("CheckTransaction")
			}
		}

		time.Sleep(1 * time.Minute)
	}
}
