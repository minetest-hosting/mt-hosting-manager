package worker

import (
	"mt-hosting-manager/api/hetzner_cloud"
	"mt-hosting-manager/api/hetzner_dns"
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type Worker struct {
	repos   *db.Repositories
	cfg     *types.Config
	hcc     *hetzner_cloud.HetznerCloudClient
	hdc     *hetzner_dns.HetznerDNSClient
	running *atomic.Bool
}

func NewWorker(repos *db.Repositories, cfg *types.Config) *Worker {
	return &Worker{
		repos:   repos,
		cfg:     cfg,
		hcc:     hetzner_cloud.New(cfg.HetznerCloudKey),
		hdc:     hetzner_dns.New(cfg.HetznerApiKey, cfg.HetznerApiZoneID),
		running: &atomic.Bool{},
	}
}

func (w *Worker) Stop() {
	w.running.Store(false)
}

func (w *Worker) Start() {
	if w.running.CompareAndSwap(false, true) {
		go w.Run()
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

	for w.running.Load() {
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
