package worker

import (
	"errors"
	"fmt"
	"mt-hosting-manager/types"
)

func (w *Worker) NodeDestroy(job *types.Job) error {
	node, err := w.repos.UserNodeRepo.GetByID(*job.UserNodeID)
	if err != nil {
		return err
	}
	if node == nil {
		return errors.New("node not found")
	}

	err = w.hcc.DeleteServer(node.ExternalID)
	if err != nil {
		return fmt.Errorf("delete server failed: %v", err)
	}

	return w.repos.UserNodeRepo.Delete(node.ID)
}
