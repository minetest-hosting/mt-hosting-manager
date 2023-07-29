package worker

import (
	"errors"
	"mt-hosting-manager/types"
	"time"
)

func (w *Worker) NodeDestroy(job *types.Job) error {
	node, err := w.repos.UserNodeRepo.GetByID(*job.UserNodeID)
	if err != nil {
		return err
	}
	if node == nil {
		return errors.New("node not found")
	}

	//TODO: resource removals
	time.Sleep(time.Second * 10)

	return w.repos.UserNodeRepo.Delete(node.ID)
}
