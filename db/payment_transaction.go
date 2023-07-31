package db

import (
	"database/sql"
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type PaymentTransactionRepository struct {
	DB dbutil.DBTx
}

func (r *PaymentTransactionRepository) Insert(n *types.PaymentTransaction) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return dbutil.Insert(r.DB, n)
}

func (r *PaymentTransactionRepository) GetByID(id string) (*types.PaymentTransaction, error) {
	nt, err := dbutil.Select(r.DB, &types.PaymentTransaction{}, "where id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *PaymentTransactionRepository) Delete(id string) error {
	return dbutil.Delete(r.DB, &types.PaymentTransaction{}, "where id = $1", id)
}
