package db

import (
	"database/sql"
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type UserNodeRepository struct {
	dbu *dbutil.DBUtil[*types.UserNode]
}

func (r *UserNodeRepository) Insert(n *types.UserNode) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.dbu.Insert(n)
}

func (r *UserNodeRepository) Update(n *types.UserNode) error {
	return r.dbu.Update(n, "where id = %s", n.ID)
}

func (r *UserNodeRepository) GetByID(id string) (*types.UserNode, error) {
	nt, err := r.dbu.Select("where id = %s", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *UserNodeRepository) GetByName(name string) (*types.UserNode, error) {
	nt, err := r.dbu.Select("where name = %s", name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *UserNodeRepository) GetByUserID(user_id string) ([]*types.UserNode, error) {
	return r.dbu.SelectMulti("where user_id = %s", user_id)
}

func (r *UserNodeRepository) GetByUserIDAndState(user_id string, state types.UserNodeState) ([]*types.UserNode, error) {
	return r.dbu.SelectMulti("where user_id = %s and state = %s", user_id, state)
}

func (r *UserNodeRepository) GetAll() ([]*types.UserNode, error) {
	return r.dbu.SelectMulti("")
}

func (r *UserNodeRepository) Delete(id string) error {
	return r.dbu.Delete("where id = %s", id)
}

func (r *UserNodeRepository) GetByLastCollectedTime(last_collected_time int64) ([]*types.UserNode, error) {
	return r.dbu.SelectMulti("where last_collected_time < %s", last_collected_time)
}
