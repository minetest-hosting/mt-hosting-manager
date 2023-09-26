package worker

import (
	"time"

	"github.com/sirupsen/logrus"
)

func (w *Worker) CollectJob() {
	for w.running.Load() {
		w.wg.Add(1)
		err := w.core.Collect(time.Now().Unix())
		if err != nil {
			logrus.WithError(err).Error("collect error")
		}
		w.wg.Done()

		time.Sleep(time.Minute)
	}
}
