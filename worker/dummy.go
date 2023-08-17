package worker

import (
	"mt-hosting-manager/types"
	"time"
)

// dummy job for testing locally
func (w *Worker) Dummy(job *types.Job) error {
	time.Sleep(15 * time.Second)
	return nil
}
