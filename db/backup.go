package db

import (
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BackupRepository struct {
	g *gorm.DB
}

func (r *BackupRepository) Insert(n *types.Backup) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.g.Create(n).Error
}

func (r *BackupRepository) Update(n *types.Backup) error {
	return r.g.Model(n).Updates(n).Error
}

func (r *BackupRepository) GetByID(id string) (*types.Backup, error) {
	var list []*types.Backup
	err := r.g.Where(types.Backup{ID: id}).Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *BackupRepository) GetByBackupSpaceID(backup_space_id string) ([]*types.Backup, error) {
	var list []*types.Backup
	err := r.g.Where(types.Backup{BackupSpaceID: backup_space_id}).Find(&list).Error
	return list, err
}

func (r *BackupRepository) Delete(id string) error {
	return r.g.Delete(types.Backup{ID: id}).Error
}
