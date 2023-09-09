package worker

import (
	"errors"
	"mt-hosting-manager/types"
	"time"

	"github.com/sirupsen/logrus"
)

func (w *Worker) ExecuteJob(job *types.Job) {
	logrus.WithFields(job.LogrusFields()).Debug("Executing job")

	var err error
	switch job.Type {
	case types.JobTypeNodeDestroy:
		err = w.NodeDestroy(job)
	case types.JobTypeNodeSetup:
		err = w.NodeProvision(job)
	case types.JobTypeServerSetup:
		err = w.ServerSetup(job)
	case types.JobTypeServerDestroy:
		err = w.ServerDestroy(job)
	default:
		err = errors.New("type not implemented")
	}

	if err != nil {
		job.State = types.JobStateDoneFailure
		job.Message = err.Error()

		fields := job.LogrusFields()
		fields["err"] = err
		logrus.WithFields(fields).Error("job failed")

	} else {
		job.State = types.JobStateDoneSuccess

	}

	job.Finished = time.Now().Unix()
	err = w.repos.JobRepo.Update(job)
	if err != nil {
		fields := job.LogrusFields()
		fields["err"] = err
		logrus.WithFields(fields).Error("job update failed")
		return
	}

	logrus.WithFields(job.LogrusFields()).Debug("Job finished")
}
