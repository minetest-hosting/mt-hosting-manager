package db

import (
	"database/sql"
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type MinetestServerRepository struct {
	DB dbutil.DBTx
}

func (r *MinetestServerRepository) Insert(n *types.MinetestServer) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return dbutil.Insert(r.DB, n)
}

func (r *MinetestServerRepository) Update(n *types.MinetestServer) error {
	return dbutil.Update(r.DB, n, "where id = $1", n.ID)
}

func (r *MinetestServerRepository) GetByID(id string) (*types.MinetestServer, error) {
	nt, err := dbutil.Select(r.DB, &types.MinetestServer{}, "where id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *MinetestServerRepository) GetAll() ([]*types.MinetestServer, error) {
	return dbutil.SelectMulti(r.DB, func() *types.MinetestServer { return &types.MinetestServer{} }, "")
}

func (r *MinetestServerRepository) Delete(id string) error {
	return dbutil.Delete(r.DB, &types.MinetestServer{}, "where id = $1", id)
}
