package worker

import (
	"mt-hosting-manager/types"

	"github.com/google/uuid"
)

func SetupNodeJob(node *types.UserNode) *types.Job {
	return &types.Job{
		ID:         uuid.NewString(),
		Type:       types.JobTypeNodeSetup,
		State:      types.JobStateCreated,
		UserNodeID: &node.ID,
	}
}

func RemoveNodeJob(node *types.UserNode) *types.Job {
	return &types.Job{
		ID:         uuid.NewString(),
		Type:       types.JobTypeNodeDestroy,
		State:      types.JobStateCreated,
		UserNodeID: &node.ID,
	}
}

func SetupServerJob(node *types.UserNode, server *types.MinetestServer) *types.Job {
	return &types.Job{
		ID:               uuid.NewString(),
		Type:             types.JobTypeServerSetup,
		State:            types.JobStateCreated,
		UserNodeID:       &node.ID,
		MinetestServerID: &server.ID,
	}
}

func RemoveServerJob(node *types.UserNode, server *types.MinetestServer) *types.Job {
	return &types.Job{
		ID:               uuid.NewString(),
		Type:             types.JobTypeServerDestroy,
		State:            types.JobStateCreated,
		UserNodeID:       &node.ID,
		MinetestServerID: &server.ID,
	}
}
