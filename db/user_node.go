package db

import (
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserNodeRepository struct {
	g *gorm.DB
}

func (r *UserNodeRepository) Insert(n *types.UserNode) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.g.Create(n).Error
}

func (r *UserNodeRepository) Update(n *types.UserNode) error {
	return r.g.Model(n).Updates(n).Error
}

func (r *UserNodeRepository) GetByID(id string) (*types.UserNode, error) {
	var list []*types.UserNode
	err := r.g.Where(types.UserNode{ID: id}).Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *UserNodeRepository) GetAll() ([]*types.UserNode, error) {
	var list []*types.UserNode
	err := r.g.Where(types.UserNode{}).Find(&list).Error
	return list, err
}

func (r *UserNodeRepository) Delete(id string) error {
	return r.g.Delete(types.UserNode{ID: id}).Error
}

func (r *UserNodeRepository) Search(search *types.UserNodeSearch) ([]*types.UserNode, error) {
	q := r.g

	if search.ID != nil {
		q = q.Where(types.UserNode{ID: *search.ID})
	}

	if search.Name != nil {
		q = q.Where(types.UserNode{Name: *search.Name})
	}

	if search.UserID != nil {
		q = q.Where(types.UserNode{UserID: *search.UserID})
	}

	if search.State != nil {
		q = q.Where(types.UserNode{State: *search.State})
	}

	if search.ValidUntil != nil {
		q = q.Where("valid_until < ?", *search.ValidUntil)
	}

	var list []*types.UserNode
	err := q.Find(&list).Error
	return list, err
}
