package worker

import (
	"fmt"
	"mt-hosting-manager/types"
)

func (w *Worker) GetJobContext(job *types.Job) (*types.UserNode, *types.MinetestServer, error) {
	node, err := w.repos.UserNodeRepo.GetByID(*job.UserNodeID)
	if err != nil {
		return nil, nil, fmt.Errorf("usernode fetch error: %v", err)
	}
	if node == nil {
		return nil, nil, fmt.Errorf("usernode not found: %s", *job.UserNodeID)
	}

	server, err := w.repos.MinetestServerRepo.GetByID(*job.MinetestServerID)
	if err != nil {
		return nil, nil, fmt.Errorf("server fetch error: %v", err)
	}
	if server == nil {
		return nil, nil, fmt.Errorf("server not found: %s", *job.MinetestServerID)
	}

	return node, server, nil
}
