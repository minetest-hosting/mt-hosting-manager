package db

import (
	"database/sql"
	"mt-hosting-manager/types"

	"github.com/minetest-go/dbutil"
)

type ExchangeRateRepository struct {
	dbu *dbutil.DBUtil[*types.ExchangeRate]
}

func (r *ExchangeRateRepository) Insert(u *types.ExchangeRate) error {
	return r.dbu.Insert(u)
}

func (r *ExchangeRateRepository) Update(u *types.ExchangeRate) error {
	return r.dbu.Update(u, "where currency =%s", u.Currency)
}

func (r *ExchangeRateRepository) GetByCurrency(currency string) (*types.ExchangeRate, error) {
	u, err := r.dbu.Select("where currency = %s", currency)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return u, err
	}
}

func (r *ExchangeRateRepository) GetAll() ([]*types.ExchangeRate, error) {
	return r.dbu.SelectMulti("")
}

func (r *ExchangeRateRepository) Delete(currency string) error {
	return r.dbu.Delete("where currency = %s", currency)
}
