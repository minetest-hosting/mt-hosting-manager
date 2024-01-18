package db

import (
	"database/sql"
	"mt-hosting-manager/types"
	"time"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type JobRepository struct {
	dbu *dbutil.DBUtil[*types.Job]
}

func (r *JobRepository) Insert(n *types.Job) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.dbu.Insert(n)
}

func (r *JobRepository) Update(n *types.Job) error {
	return r.dbu.Update(n, "where id = %s", n.ID)
}

func (r *JobRepository) GetByID(id string) (*types.Job, error) {
	nt, err := r.dbu.Select("where id = %s", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *JobRepository) GetByState(state types.JobState) ([]*types.Job, error) {
	return r.dbu.SelectMulti("where state = %s", state)
}

func (r *JobRepository) GetLatestByUserNodeID(usernodeID string) (*types.Job, error) {
	nt, err := r.dbu.Select("where user_node_id = %s order by started desc limit 1", usernodeID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *JobRepository) GetLatestByMinetestServerID(minetestserverID string) (*types.Job, error) {
	nt, err := r.dbu.Select("where minetest_server_id = %s order by started desc limit 1", minetestserverID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *JobRepository) GetAll() ([]*types.Job, error) {
	return r.dbu.SelectMulti("")
}

func (r *JobRepository) Delete(id string) error {
	return r.dbu.Delete("where id = %s", id)
}

func (r *JobRepository) DeleteBefore(t time.Time) error {
	return r.dbu.Delete("where started < %s", t.Unix())
}
