package db_test

import (
	"fmt"
	"mt-hosting-manager/types"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestJobRepo(t *testing.T) {
	repos := SetupRepos(t)

	j := &types.Job{
		Type:    types.JobTypeNodeSetup,
		State:   types.JobStateRunning,
		Created: time.Now().Unix(),
		NextRun: 100,
	}
	assert.NoError(t, repos.JobRepo.Insert(j))

	j.Data = []byte{0x00, 0x01}
	assert.NoError(t, repos.JobRepo.Update(j))

	list, err := repos.JobRepo.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))

	list, err = repos.JobRepo.GetByStateAndNextRun(types.JobStateRunning, 20)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(list))

	list, err = repos.JobRepo.GetByStateAndNextRun(types.JobStateRunning, 200)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))

	list, err = repos.JobRepo.GetByStateAndNextRun(types.JobStateDoneFailure, 200)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(list))

	assert.NoError(t, repos.JobRepo.Delete(j.ID))
	assert.NoError(t, repos.JobRepo.DeleteBefore(time.Now()))
}

func TestJobQueue(t *testing.T) {
	repos := SetupRepos(t)

	assert.NoError(t, repos.JobRepo.Insert(&types.Job{
		Type:    types.JobTypeNodeSetup,
		State:   types.JobStateRunning,
		Created: time.Now().Unix(),
		NextRun: 100,
	}))

	wg := &sync.WaitGroup{}
	wg.Add(1)

	// fetch first job async with block
	go func() {
		err := repos.Gorm().Transaction(func(tx *gorm.DB) error {
			job, err := repos.JobRepo.GetNextJob(tx, types.JobStateRunning, 200)
			if err != nil {
				return err
			}
			if job == nil {
				return fmt.Errorf("no job fetched")
			}
			time.Sleep(100 * time.Millisecond)
			job.State = types.JobStateDoneSuccess

			return repos.JobRepo.UpdateWithTx(tx, job)
		})
		assert.NoError(t, err)
		wg.Done()
	}()

	// try to fetch next job sync
	time.Sleep(10 * time.Millisecond)
	err := repos.Gorm().Transaction(func(tx *gorm.DB) error {
		job, err := repos.JobRepo.GetNextJob(tx, types.JobStateRunning, 200)
		if err != nil {
			return err
		}
		if job != nil {
			// no job should be available here
			return fmt.Errorf("job fetched: %v", job)
		}
		return nil
	})
	assert.NoError(t, err)

	wg.Wait()

	jobs, err := repos.JobRepo.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(jobs))
	assert.Equal(t, types.JobStateDoneSuccess, jobs[0].State)
}
