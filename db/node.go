package db

import (
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type NodeRepository struct {
	DB dbutil.DBTx
}

func (r *NodeRepository) Insert(n *types.Node) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return dbutil.Insert(r.DB, n)
}

func (r *NodeRepository) Update(n *types.Node) error {
	return dbutil.Update(r.DB, n, "where id = $1", n.ID)
}

func (r *NodeRepository) GetAll() ([]*types.Node, error) {
	return dbutil.SelectMulti(r.DB, func() *types.Node { return &types.Node{} }, "")
}

func (r *NodeRepository) Delete(node_id string) error {
	return dbutil.Delete(r.DB, &types.Node{}, "where id = $1", node_id)
}
