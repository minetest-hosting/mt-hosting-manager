package worker

import (
	"fmt"
	"mt-hosting-manager/types"
)

func (w *Worker) ServerDestroy(job *types.Job) error {
	node, server, err := w.GetJobContext(job)
	if err != nil {
		return err
	}

	err = w.removeServer(node, server, true)
	if err != nil {
		return fmt.Errorf("server remove error: %v", err)
	}

	job.State = types.JobStateDoneSuccess
	return nil
}
