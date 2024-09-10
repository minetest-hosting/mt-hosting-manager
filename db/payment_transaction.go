package db

import (
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentTransactionRepository struct {
	g *gorm.DB
}

func (r *PaymentTransactionRepository) Insert(n *types.PaymentTransaction) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.g.Create(n).Error
}

func (r *PaymentTransactionRepository) Update(tx *types.PaymentTransaction) error {
	return r.g.Model(tx).Updates(tx).Error
}

func (r *PaymentTransactionRepository) GetByID(id string) (*types.PaymentTransaction, error) {
	var tx *types.PaymentTransaction
	err := r.g.Where(types.Job{ID: id}).First(&tx).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return tx, err
}

func (r *PaymentTransactionRepository) GetByUserID(user_id string) ([]*types.PaymentTransaction, error) {
	var list []*types.PaymentTransaction
	err := r.g.Where(types.PaymentTransaction{UserID: user_id}).Find(&list).Error
	return list, err
}

func (r *PaymentTransactionRepository) Delete(id string) error {
	return r.g.Delete(types.PaymentTransaction{ID: id}).Error
}

func (r *PaymentTransactionRepository) Search(s *types.PaymentTransactionSearch) ([]*types.PaymentTransaction, error) {
	q := r.g.Where("created > ?", s.FromTimestamp)
	q = q.Where("created < ?", s.ToTimestamp)

	if s.UserID != nil {
		q = q.Where(types.PaymentTransaction{UserID: *s.UserID})
	}

	if s.State != nil {
		q = q.Where(types.PaymentTransaction{State: *s.State})
	}

	q = q.Order("created desc").Limit(1000)

	var list []*types.PaymentTransaction
	err := q.Find(&list).Error
	return list, err
}
