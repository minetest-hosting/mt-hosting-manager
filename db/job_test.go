package db_test

import (
	"mt-hosting-manager/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJobRepo(t *testing.T) {
	repos := SetupRepos(t)

	j := &types.Job{
		Type:    types.JobTypeNodeSetup,
		State:   types.JobStateCreated,
		Created: time.Now().Unix(),
	}
	assert.NoError(t, repos.JobRepo.Insert(j))

	j.Data = []byte{0x00, 0x01}
	assert.NoError(t, repos.JobRepo.Update(j))

	assert.NoError(t, repos.JobRepo.Delete(j.ID))
	assert.NoError(t, repos.JobRepo.DeleteBefore(time.Now()))

}
