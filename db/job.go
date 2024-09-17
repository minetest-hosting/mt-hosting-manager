package db

import (
	"mt-hosting-manager/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JobRepository struct {
	g *gorm.DB
}

func (r *JobRepository) Insert(n *types.Job) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.g.Create(n).Error
}

func (r *JobRepository) Update(n *types.Job) error {
	return r.g.Model(n).Updates(n).Error
}

func (r *JobRepository) UpdateWithTx(tx *gorm.DB, n *types.Job) error {
	return tx.Model(n).Updates(n).Error
}

func (r *JobRepository) GetByID(id string) (*types.Job, error) {
	var list []*types.Job
	err := r.g.Where(types.Job{ID: id}).Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *JobRepository) GetByState(state types.JobState) ([]*types.Job, error) {
	var list []*types.Job
	err := r.g.Where(types.Job{State: state}).Find(&list).Error
	return list, err
}

func (r *JobRepository) GetByStateAndNextRun(state types.JobState, nextrun int64) ([]*types.Job, error) {
	var list []*types.Job
	q := r.g.Where(types.Job{State: state})
	q = q.Where("next_run < ?", nextrun)
	q = q.Find(&list)
	err := q.Error
	return list, err
}

func (r *JobRepository) GetNextJob(tx *gorm.DB, state types.JobState, nextrun int64) (*types.Job, error) {
	var list []*types.Job
	q := tx.Clauses(clause.Locking{
		Strength: clause.LockingStrengthUpdate,
		Options:  clause.LockingOptionsSkipLocked,
	})
	q = q.Where(types.Job{State: state})
	q = q.Limit(1)
	q = q.Find(&list)
	err := q.Error

	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *JobRepository) GetLatestByUserNodeID(usernodeID string) (*types.Job, error) {
	var list []*types.Job
	err := r.g.Where(types.Job{UserNodeID: &usernodeID}).Order("created desc").Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *JobRepository) GetLatestByMinetestServerID(minetestserverID string) (*types.Job, error) {
	var list []*types.Job
	err := r.g.Where(types.Job{MinetestServerID: &minetestserverID}).Order("created desc").Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *JobRepository) GetLatestByBackupID(backupID string) (*types.Job, error) {
	var list []*types.Job
	err := r.g.Where(types.Job{BackupID: &backupID}).Order("created desc").Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *JobRepository) GetAll() ([]*types.Job, error) {
	var list []*types.Job
	err := r.g.Where(types.Job{}).Order("created desc").Find(&list).Error
	return list, err
}

func (r *JobRepository) Delete(id string) error {
	return r.g.Delete(types.Job{ID: id}).Error
}

func (r *JobRepository) DeleteBefore(t time.Time) error {
	return r.g.Where("created < ?", t.Unix()).Delete(types.Job{}).Error
}

func (r *JobRepository) DeleteAll() error {
	return r.g.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(types.Job{}).Error
}
