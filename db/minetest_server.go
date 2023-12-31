package db

import (
	"database/sql"
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type MinetestServerRepository struct {
	dbu *dbutil.DBUtil[*types.MinetestServer]
}

func (r *MinetestServerRepository) Insert(n *types.MinetestServer) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.dbu.Insert(n)
}

func (r *MinetestServerRepository) Update(n *types.MinetestServer) error {
	return r.dbu.Update(n, "where id = %s", n.ID)
}

func (r *MinetestServerRepository) GetByName(name string) (*types.MinetestServer, error) {
	nt, err := r.dbu.Select("where name = %s", name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *MinetestServerRepository) GetByID(id string) (*types.MinetestServer, error) {
	nt, err := r.dbu.Select("where id = %s", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *MinetestServerRepository) GetAll() ([]*types.MinetestServer, error) {
	return r.dbu.SelectMulti("")
}

func (r *MinetestServerRepository) GetByUserID(userID string) ([]*types.MinetestServer, error) {
	return r.dbu.SelectMulti(
		"where user_node_id in (select id from user_node where user_id = %s)",
		userID,
	)
}

func (r *MinetestServerRepository) GetByNodeID(nodeID string) ([]*types.MinetestServer, error) {
	return r.dbu.SelectMulti(
		"where user_node_id = %s",
		nodeID,
	)
}

func (r *MinetestServerRepository) Delete(id string) error {
	return r.dbu.Delete("where id = %s", id)
}
