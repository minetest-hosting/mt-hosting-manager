package worker

import (
	"errors"
	"mt-hosting-manager/api/hetzner_cloud"
	"mt-hosting-manager/api/hetzner_dns"
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Worker struct {
	repos *db.Repositories
	hcc   *hetzner_cloud.HetznerCloudClient
	hdc   *hetzner_dns.HetznerDNSClient
}

func NewWorker(repos *db.Repositories) *Worker {
	return &Worker{
		repos: repos,
		hcc:   hetzner_cloud.New(os.Getenv("HETZNER_CLOUD_KEY")),
		hdc:   hetzner_dns.New(os.Getenv("HETZNER_API_KEY"), os.Getenv("HETZNER_API_ZONEID")),
	}
}

func (w *Worker) Run() {

	// start background jobs
	go w.MetricsCollector()

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
