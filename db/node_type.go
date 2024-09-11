package db

import (
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NodeTypeRepository struct {
	g *gorm.DB
}

func (r *NodeTypeRepository) Insert(n *types.NodeType) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.g.Create(n).Error
}

func (r *NodeTypeRepository) Update(n *types.NodeType) error {
	return r.g.Model(n).Updates(n).Error
}

func (r *NodeTypeRepository) GetByID(id string) (*types.NodeType, error) {
	var list []*types.NodeType
	err := r.g.Where(types.NodeType{ID: id}).Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *NodeTypeRepository) GetByState(t types.NodeTypeState) ([]*types.NodeType, error) {
	var list []*types.NodeType
	err := r.g.Where(types.NodeType{State: t}).Find(&list).Error
	return list, err
}

func (r *NodeTypeRepository) GetAll() ([]*types.NodeType, error) {
	var list []*types.NodeType
	err := r.g.Where(types.NodeType{}).Order("order_id ASC").Find(&list).Error
	return list, err
}

func (r *NodeTypeRepository) Delete(node_type_id string) error {
	return r.g.Delete(types.NodeType{ID: node_type_id}).Error
}

func (r *NodeTypeRepository) DeleteAll() error {
	return r.g.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(types.NodeType{}).Error
}
