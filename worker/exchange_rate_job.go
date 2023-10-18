package worker

import (
	"time"

	"github.com/sirupsen/logrus"
)

func (w *Worker) ExchangeRateUpdateJob() {
	for w.running.Load() {
		err := w.core.UpdateExchangeRates()
		if err != nil {
			logrus.WithError(err).Error("exchange rate update")
		}

		time.Sleep(1 * time.Hour)
	}
}
