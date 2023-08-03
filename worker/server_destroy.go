package worker

import "mt-hosting-manager/types"

func (w *Worker) ServerDestroy(job *types.Job) error {
	_, err := w.repos.UserNodeRepo.GetByID(*job.UserNodeID)
	//TODO
	return err
}
