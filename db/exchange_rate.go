package db

import (
	"mt-hosting-manager/types"

	"gorm.io/gorm"
)

type ExchangeRateRepository struct {
	g *gorm.DB
}

func (r *ExchangeRateRepository) Insert(n *types.ExchangeRate) error {
	return r.g.Create(n).Error
}

func (r *ExchangeRateRepository) Update(n *types.ExchangeRate) error {
	return r.g.Model(n).Updates(n).Error
}

func (r *ExchangeRateRepository) GetAll() ([]*types.ExchangeRate, error) {
	var list []*types.ExchangeRate
	err := r.g.Where(types.ExchangeRate{}).Find(&list).Error
	return list, err
}

func (r *ExchangeRateRepository) GetByCurrency(currency string) (*types.ExchangeRate, error) {
	var list []*types.ExchangeRate
	err := r.g.Where(types.ExchangeRate{Currency: currency}).Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *ExchangeRateRepository) DeleteByCurrency(currency string) error {
	return r.g.Delete(types.ExchangeRate{Currency: currency}).Error
}

func (r *ExchangeRateRepository) DeleteAll() error {
	return r.g.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(types.ExchangeRate{}).Error
}
