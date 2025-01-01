package db

import (
	"mt-hosting-manager/types"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserSettingRepository struct {
	g *gorm.DB
}

func (r *UserSettingRepository) Set(n *types.UserSetting) error {
	return r.g.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(n).Error
}

func (r *UserSettingRepository) GetByUserID(user_id string) ([]*types.UserSetting, error) {
	return FindMulti[types.UserSetting](r.g.Where(types.UserSetting{UserID: user_id}))
}

func (r *UserSettingRepository) Delete(user_id, key string) error {
	return r.g.Delete(types.UserSetting{UserID: user_id, Key: key}).Error
}
