package worker

import (
	"mt-hosting-manager/core"
	"time"

	"github.com/sirupsen/logrus"
)

func (w *Worker) CollectJob() {
	for w.running.Load() {
		ts := time.Now().Unix()
		w.wg.Add(1)
		err := w.core.Collect(ts - core.SECONDS_IN_A_DAY)
		if err != nil {
			logrus.WithError(err).Error("collect error")
		}
		w.wg.Done()

		time.Sleep(time.Minute)
	}
}
