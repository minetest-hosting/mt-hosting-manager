package db

import (
	"database/sql"
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type NodeTypeRepository struct {
	dbu *dbutil.DBUtil[*types.NodeType]
}

func (r *NodeTypeRepository) Insert(n *types.NodeType) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.dbu.Insert(n)
}

func (r *NodeTypeRepository) Update(n *types.NodeType) error {
	return r.dbu.Update(n, "where id = %s", n.ID)
}

func (r *NodeTypeRepository) GetByID(id string) (*types.NodeType, error) {
	nt, err := r.dbu.Select("where id = %s", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *NodeTypeRepository) GetByState(t types.NodeTypeState) ([]*types.NodeType, error) {
	return r.dbu.SelectMulti("where state = %s order by order_id asc", t)
}

func (r *NodeTypeRepository) GetAll() ([]*types.NodeType, error) {
	return r.dbu.SelectMulti("order by order_id asc")
}

func (r *NodeTypeRepository) Delete(node_type_id string) error {
	return r.dbu.Delete("where id = %s", node_type_id)
}

func (r *NodeTypeRepository) DeleteAll() error {
	return r.dbu.Delete("")
}
