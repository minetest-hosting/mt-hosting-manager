package db

import (
	"database/sql"
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type NodeTypeRepository struct {
	DB dbutil.DBTx
}

func (r *NodeTypeRepository) Insert(n *types.NodeType) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return dbutil.Insert(r.DB, n)
}

func (r *NodeTypeRepository) Update(n *types.NodeType) error {
	return dbutil.Update(r.DB, n, "where id = $1", n.ID)
}

func (r *NodeTypeRepository) GetByID(id string) (*types.NodeType, error) {
	nt, err := dbutil.Select(r.DB, &types.NodeType{}, "where id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *NodeTypeRepository) GetByState(t types.NodeTypeState) ([]*types.NodeType, error) {
	return dbutil.SelectMulti(r.DB, func() *types.NodeType { return &types.NodeType{} }, "where state = $1", t)
}

func (r *NodeTypeRepository) GetAll() ([]*types.NodeType, error) {
	return dbutil.SelectMulti(r.DB, func() *types.NodeType { return &types.NodeType{} }, "")
}

func (r *NodeTypeRepository) Delete(node_type_id string) error {
	return dbutil.Delete(r.DB, &types.NodeType{}, "where id = $1", node_type_id)
}
