package db

import (
	"github.com/minetest-go/dbutil"
)

type Repositories struct {
	UserRepo               *UserRepository
	NodeTypeRepo           *NodeTypeRepository
	UserNodeRepo           *UserNodeRepository
	MinetestServerRepo     *MinetestServerRepository
	JobRepo                *JobRepository
	PaymentTransactionRepo *PaymentTransactionRepository
}

func NewRepositories(db dbutil.DBTx) *Repositories {
	return &Repositories{
		UserRepo:               &UserRepository{DB: db},
		NodeTypeRepo:           &NodeTypeRepository{DB: db},
		UserNodeRepo:           &UserNodeRepository{DB: db},
		MinetestServerRepo:     &MinetestServerRepository{DB: db},
		JobRepo:                &JobRepository{DB: db},
		PaymentTransactionRepo: &PaymentTransactionRepository{DB: db},
	}
}
