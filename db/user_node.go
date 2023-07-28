package db

import (
	"database/sql"
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type UserNodeRepository struct {
	DB dbutil.DBTx
}

func UserNodeFactory() *types.UserNode { return &types.UserNode{} }

func (r *UserNodeRepository) Insert(n *types.UserNode) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return dbutil.Insert(r.DB, n)
}

func (r *UserNodeRepository) Update(n *types.UserNode) error {
	return dbutil.Update(r.DB, n, "where id = $1", n.ID)
}

func (r *UserNodeRepository) GetByID(id string) (*types.UserNode, error) {
	nt, err := dbutil.Select(r.DB, &types.UserNode{}, "where id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *UserNodeRepository) GetByName(name string) (*types.UserNode, error) {
	nt, err := dbutil.Select(r.DB, &types.UserNode{}, "where name = $1", name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *UserNodeRepository) GetByUserID(user_id string) ([]*types.UserNode, error) {
	return dbutil.SelectMulti(r.DB, UserNodeFactory, "where user_id = $1", user_id)
}

func (r *UserNodeRepository) GetAll() ([]*types.UserNode, error) {
	return dbutil.SelectMulti(r.DB, UserNodeFactory, "")
}

func (r *UserNodeRepository) Delete(id string) error {
	return dbutil.Delete(r.DB, &types.UserNode{}, "where id = $1", id)
}
