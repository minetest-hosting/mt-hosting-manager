package db_test

import (
	"mt-hosting-manager/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJobRepo(t *testing.T) {
	repos := SetupRepos(t)

	err := repos.JobRepo.Insert(&types.Job{
		Type:    types.JobTypeNodeSetup,
		State:   types.JobStateCreated,
		Started: time.Now().Unix(),
	})
	assert.NoError(t, err)

}
