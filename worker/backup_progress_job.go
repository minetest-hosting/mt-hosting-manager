package worker

import (
	"time"
)

func (w *Worker) UpdateBackupProgressJob() {
	for w.running.Load() {
		//TODO
		time.Sleep(1 * time.Second)
	}
}
