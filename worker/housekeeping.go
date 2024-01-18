package worker

import (
	"time"

	"github.com/sirupsen/logrus"
)

func (w *Worker) HousekeepingJob() {
	for w.running.Load() {
		err := w.core.UpdateExchangeRates()
		if err != nil {
			logrus.WithError(err).Error("exchange rate update")
		}

		err = w.repos.JobRepo.DeleteBefore(time.Now().Add(time.Hour * -24))
		if err != nil {
			logrus.WithError(err).Error("job cleanup")
		}

		time.Sleep(1 * time.Hour)
	}
}
