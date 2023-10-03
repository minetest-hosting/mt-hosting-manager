package worker

import (
	"errors"
	"fmt"
	"mt-hosting-manager/notify"
	"mt-hosting-manager/types"
	"time"

	"github.com/sirupsen/logrus"
)

func (w *Worker) ExecuteJob(job *types.Job) {
	logrus.WithFields(job.LogrusFields()).Debug("Executing job")
	w.wg.Add(1)
	defer w.wg.Done()

	status_callback := func(msg string, progress_percent int) {
		job.Message = msg
		job.ProgressPercent = float64(progress_percent)
		// update and ignore errors
		err := w.repos.JobRepo.Update(job)
		if err != nil {
			logrus.WithError(err).Error("failed to update job progress")
		}
	}

	var err error
	switch job.Type {
	case types.JobTypeNodeDestroy:
		err = w.NodeDestroy(job)
	case types.JobTypeNodeSetup:
		err = w.NodeProvision(job, status_callback)
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

		job_url := fmt.Sprintf("%s/#/jobs", w.cfg.BaseURL)
		notify.Send(&notify.NtfyNotification{
			Title:    fmt.Sprintf("Job failed: %s", job.Type),
			Message:  fmt.Sprintf("Type: %s, ID %s, progress %.2f, message: '%s'", job.Type, job.ID, job.ProgressPercent, job.Message),
			Priority: 3,
			Click:    &job_url,
			Tags:     []string{"arrow_forward", "warning"},
		}, true)

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
