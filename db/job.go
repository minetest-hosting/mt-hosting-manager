package db

import (
	"database/sql"
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type JobRepository struct {
	DB dbutil.DBTx
}

func (r *JobRepository) Insert(n *types.Job) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return dbutil.Insert(r.DB, n)
}

func (r *JobRepository) Update(n *types.Job) error {
	return dbutil.Update(r.DB, n, "where id = $1", n.ID)
}

func (r *JobRepository) GetByID(id string) (*types.Job, error) {
	nt, err := dbutil.Select(r.DB, &types.Job{}, "where id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *JobRepository) GetAll() ([]*types.Job, error) {
	return dbutil.SelectMulti(r.DB, func() *types.Job { return &types.Job{} }, "")
}

func (r *JobRepository) Delete(id string) error {
	return dbutil.Delete(r.DB, &types.Job{}, "where id = $1", id)
}
