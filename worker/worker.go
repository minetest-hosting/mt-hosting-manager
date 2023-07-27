package worker

import (
	"mt-hosting-manager/db"
)

type Worker struct {
	repos *db.Repositories
}

func NewWorker(repos *db.Repositories) *Worker {
	return &Worker{
		repos: repos,
	}
}

func (w *Worker) Start() {}
