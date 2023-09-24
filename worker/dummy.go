package worker

import (
	"mt-hosting-manager/core"
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"
	"time"

	"github.com/sirupsen/logrus"
)

func executeDummyJob(repos *db.Repositories, job *types.Job) {
	logrus.WithFields(job.LogrusFields()).Debug("executing job")

	job.Started = time.Now().Unix()
	job.State = types.JobStateRunning
	err := repos.JobRepo.Update(job)
	if err != nil {
		panic(err)
	}

	var node *types.UserNode
	var server *types.MinetestServer

	if job.UserNodeID != nil {
		node, err = repos.UserNodeRepo.GetByID(*job.UserNodeID)
		if err != nil {
			panic(err)
		}
	}
	if job.MinetestServerID != nil {
		server, err = repos.MinetestServerRepo.GetByID(*job.MinetestServerID)
		if err != nil {
			panic(err)
		}
	}

	switch job.Type {
	case types.JobTypeNodeSetup:
		node.State = types.UserNodeStateProvisioning
		err = repos.UserNodeRepo.Update(node)
		if err != nil {
			panic(err)
		}
		time.Sleep(10 * time.Second)

		node.IPv4 = "127.0.0.1"
		node.IPv6 = "::1"
		node.Fingerprint = "dummy"
		node.State = types.UserNodeStateRunning
		err = repos.UserNodeRepo.Update(node)
		if err != nil {
			panic(err)
		}

	case types.JobTypeNodeDestroy:
		node.State = types.UserNodeStateRemoving
		err = repos.UserNodeRepo.Update(node)
		if err != nil {
			panic(err)
		}
		time.Sleep(10 * time.Second)

		err = repos.UserNodeRepo.Delete(node.ID)
		if err != nil {
			panic(err)
		}

	case types.JobTypeServerSetup:
		server.State = types.MinetestServerStateProvisioning
		err = repos.MinetestServerRepo.Update(server)
		if err != nil {
			panic(err)
		}
		time.Sleep(10 * time.Second)

		server.State = types.MinetestServerStateRunning
		err = repos.MinetestServerRepo.Update(server)
		if err != nil {
			panic(err)
		}

	case types.JobTypeServerDestroy:
		server.State = types.MinetestServerStateRemoving
		err = repos.MinetestServerRepo.Update(server)
		if err != nil {
			panic(err)
		}
		time.Sleep(10 * time.Second)

		err = repos.MinetestServerRepo.Delete(server.ID)
		if err != nil {
			panic(err)
		}
	}

	job.Finished = time.Now().Unix()
	job.State = types.JobStateDoneSuccess
	err = repos.JobRepo.Update(job)
	if err != nil {
		panic(err)
	}

	logrus.WithFields(job.LogrusFields()).Debug("job done")
}

func DummyWorker(repos *db.Repositories, cfg *types.Config) {

	c := core.New(repos, cfg)
	go func() {
		for {
			ts := time.Now().Unix()
			err := c.Collect(ts - core.SECONDS_IN_A_DAY)
			if err != nil {
				logrus.WithError(err).Error("collect error")
			}

			time.Sleep(time.Minute)
		}
	}()

	jobs, err := repos.JobRepo.GetByState(types.JobStateRunning)
	if err != nil {
		panic(err)
	}

	for _, job := range jobs {
		executeDummyJob(repos, job)
	}

	for {
		jobs, err := repos.JobRepo.GetByState(types.JobStateCreated)
		if err != nil {
			panic(err)
		}

		for _, job := range jobs {
			executeDummyJob(repos, job)
		}

		time.Sleep(1 * time.Second)
	}
}
