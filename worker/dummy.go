package worker

import (
	"math/rand"
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

		nodes, err := repos.UserNodeRepo.GetAll()
		if err != nil {
			panic(err)
		}

		for _, node := range nodes {
			if node.State != types.UserNodeStateRunning {
				continue
			}

			node.DiskSize = 1000 * 1000 * 1000 * 10
			node.DiskUsed = 1000 * 1000 * 1000 * 2.5
			node.MemorySize = 1000 * 1000 * 1000 * 2
			node.MemoryUsed = 1000 * 1000 * 1000 * 0.2
			node.LoadPercent = rand.Intn(20)
			err = repos.UserNodeRepo.Update(node)
			if err != nil {
				panic(err)
			}
		}

		time.Sleep(1 * time.Second)
	}
}
