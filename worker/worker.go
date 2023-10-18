package worker

import (
	"mt-hosting-manager/api/coinbase"
	"mt-hosting-manager/api/hetzner_cloud"
	"mt-hosting-manager/api/hetzner_dns"
	"mt-hosting-manager/core"
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type Worker struct {
	repos   *db.Repositories
	cfg     *types.Config
	hcc     *hetzner_cloud.HetznerCloudClient
	hdc     *hetzner_dns.HetznerDNSClient
	cbc     *coinbase.CoinbaseClient
	running *atomic.Bool
	core    *core.Core
	wg      *sync.WaitGroup
}

func NewWorker(repos *db.Repositories, cfg *types.Config) *Worker {
	return &Worker{
		repos:   repos,
		cfg:     cfg,
		hcc:     hetzner_cloud.New(cfg.HetznerCloudKey),
		hdc:     hetzner_dns.New(cfg.HetznerApiKey, cfg.HetznerApiZoneID),
		cbc:     coinbase.New(cfg.CoinbaseKey),
		running: &atomic.Bool{},
		core:    core.New(repos, cfg),
		wg:      &sync.WaitGroup{},
	}
}

func (w *Worker) Stop() {
	w.running.Store(false)
	w.wg.Wait()
}

func (w *Worker) Start() {
	if w.running.CompareAndSwap(false, true) {
		go w.Run()
	}
}

func (w *Worker) Run() {
	// start collector job
	go w.CollectJob()

	// start exchange rate update job
	go w.ExchangeRateUpdateJob()

	// start mail job
	go w.MailJob()

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
