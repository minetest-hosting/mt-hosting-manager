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

func (r *PaymentTransactionRepository) GetByUserID(user_id string) ([]*types.PaymentTransaction, error) {
	return r.dbu.SelectMulti("where user_id = %s", user_id)
}

func (r *PaymentTransactionRepository) Update(tx *types.PaymentTransaction) error {
	return r.dbu.Update(tx, "where id = %s", tx.ID)
}

func (r *PaymentTransactionRepository) Delete(id string) error {
	return r.dbu.Delete("where id = %s", id)
}

func (r *PaymentTransactionRepository) Search(s *types.PaymentTransactionSearch) ([]*types.PaymentTransaction, error) {
	q := "where created > %s and created < %s"
	params := []any{s.FromTimestamp, s.ToTimestamp}

	if s.UserID != nil {
		q += " and user_id = %s"
		params = append(params, *s.UserID)
	}

	if s.State != nil {
		q += " and state = %s"
		params = append(params, *s.State)
	}

	q += " order by created desc limit 1000"

	return r.dbu.SelectMulti(q, params...)
}
