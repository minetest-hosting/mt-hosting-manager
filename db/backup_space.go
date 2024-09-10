package db

import (
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BackupSpaceRepository struct {
	g *gorm.DB
}

func (r *BackupSpaceRepository) Insert(n *types.BackupSpace) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.g.Create(n).Error
}

func (r *BackupSpaceRepository) Update(n *types.BackupSpace) error {
	return r.g.Model(n).Updates(n).Error
}

func (r *BackupSpaceRepository) GetByID(id string) (*types.BackupSpace, error) {
	var list []*types.BackupSpace
	err := r.g.Where(types.BackupSpace{ID: id}).Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *BackupSpaceRepository) GetByUserID(user_id string) ([]*types.BackupSpace, error) {
	var list []*types.BackupSpace
	err := r.g.Where(types.BackupSpace{UserID: user_id}).Find(&list).Error
	return list, err
}

func (r *BackupSpaceRepository) Delete(id string) error {
	return r.g.Delete(types.BackupSpace{ID: id}).Error
}
