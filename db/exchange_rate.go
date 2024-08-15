package db

import (
	"database/sql"
	"mt-hosting-manager/types"

	"github.com/minetest-go/dbutil"
)

type ExchangeRateRepository struct {
	dbu *dbutil.DBUtil[*types.ExchangeRate]
}

func (r *ExchangeRateRepository) Insert(n *types.ExchangeRate) error {
	return r.dbu.Insert(n)
}

func (r *ExchangeRateRepository) Update(n *types.ExchangeRate) error {
	return r.dbu.Update(n, "where currency = %s", n.Currency)
}

func (r *ExchangeRateRepository) GetAll() ([]*types.ExchangeRate, error) {
	return r.dbu.SelectMulti("")
}

func (r *ExchangeRateRepository) GetByCurrency(currency string) (*types.ExchangeRate, error) {
	rate, err := r.dbu.Select("where currency = %s", currency)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return rate, err
}

func (r *ExchangeRateRepository) DeleteByCurrency(currency string) error {
	return r.dbu.Delete("where currency = %s", currency)
}

func (r *ExchangeRateRepository) DeleteAll() error {
	return r.dbu.Delete("")
}
