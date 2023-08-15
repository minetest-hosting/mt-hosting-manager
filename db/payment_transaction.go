package db

import (
	"database/sql"
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type PaymentTransactionRepository struct {
	dbu *dbutil.DBUtil[*types.PaymentTransaction]
}

func (r *PaymentTransactionRepository) Insert(n *types.PaymentTransaction) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.dbu.Insert(n)
}

func (r *PaymentTransactionRepository) GetByID(id string) (*types.PaymentTransaction, error) {
	nt, err := r.dbu.Select("where id = %s", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *PaymentTransactionRepository) Delete(id string) error {
	return r.dbu.Delete("where id = %s", id)
}
