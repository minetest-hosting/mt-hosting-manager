package worker

import (
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"
	"time"

	"github.com/sirupsen/logrus"
)

type Worker struct {
	repos *db.Repositories
}

func NewWorker(repos *db.Repositories) *Worker {
	return &Worker{
		repos: repos,
	}
}

func (w *Worker) Run() {

	// execute previously running jobs
	jobs, err := w.repos.JobRepo.GetByState(types.JobStateRunning)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Running jobs fetch error")
	}

	for _, job := range jobs {
		go w.ExecuteJob(job)
	}

	for {
		//Execute pending (created) jobs
		jobs, err := w.repos.JobRepo.GetByState(types.JobStateCreated)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("Created jobs fetch error")
			time.Sleep(10 * time.Second)
			continue
		}

		for _, job := range jobs {
			job.Started = time.Now().Unix()
			job.State = types.JobStateRunning
			err := w.repos.JobRepo.Update(job)
			if err != nil {
				fields := job.LogrusFields()
				fields["err"] = err
				logrus.WithFields(fields).Error("job update failed (running)")
				continue
			}

			go w.ExecuteJob(job)
		}

		//TODO: remove old jobs

		time.Sleep(500 * time.Millisecond)
	}
}

func (w *Worker) ExecuteJob(job *types.Job) {
	logrus.WithFields(job.LogrusFields()).Debug("Executing job")

	if job.Type == types.JobTypeNodeDestroy {
		node, err := w.repos.UserNodeRepo.GetByID(job.UserNodeID)
		if err != nil {
			fields := job.LogrusFields()
			fields["err"] = err
			logrus.WithFields(fields).Error("node fetch failed")
			return
		}
		if node == nil {
			logrus.WithFields(job.LogrusFields()).Warn("node not found")
			job.Finished = time.Now().Unix()
			job.State = types.JobStateDoneFailure
			job.Message = "Node not found"
			err := w.repos.JobRepo.Update(job)
			if err != nil {
				fields := job.LogrusFields()
				fields["err"] = err
				logrus.WithFields(fields).Error("job update failed")
			}
			return
		}

		//TODO: resource removals
		err = w.repos.UserNodeRepo.Delete(node.ID)
		if err != nil {
			fields := job.LogrusFields()
			fields["err"] = err
			logrus.WithFields(fields).Error("job update failed")
			return
		}
	}

	//TODO: node provision

	job.Finished = time.Now().Unix()
	job.State = types.JobStateDoneSuccess
	err := w.repos.JobRepo.Update(job)
	if err != nil {
		fields := job.LogrusFields()
		fields["err"] = err
		logrus.WithFields(fields).Error("job update failed")
		return
	}

	logrus.WithFields(job.LogrusFields()).Debug("Job finished (success)")
}
