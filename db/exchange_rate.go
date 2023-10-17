package db

import (
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
