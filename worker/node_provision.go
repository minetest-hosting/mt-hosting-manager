package worker

import (
	"errors"
	"mt-hosting-manager/types"
	"time"
)

func (w *Worker) NodeProvision(job *types.Job) error {
	node, err := w.repos.UserNodeRepo.GetByID(*job.UserNodeID)
	if err != nil {
		return err
	}
	if node == nil {
		return errors.New("node not found")
	}

	//TODO: provision stuff
	time.Sleep(time.Second * 10)

	node.State = types.UserNodeStateRunning
	return w.repos.UserNodeRepo.Update(node)
}
