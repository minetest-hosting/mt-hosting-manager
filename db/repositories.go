package db

import (
	"mt-hosting-manager/types"

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
		UserRepo:               &UserRepository{dbu: dbutil.New[*types.User](db, dbutil.DialectSQLite, func() *types.User { return &types.User{} })},
		NodeTypeRepo:           &NodeTypeRepository{dbu: dbutil.New[*types.NodeType](db, dbutil.DialectSQLite, func() *types.NodeType { return &types.NodeType{} })},
		UserNodeRepo:           &UserNodeRepository{dbu: dbutil.New[*types.UserNode](db, dbutil.DialectSQLite, func() *types.UserNode { return &types.UserNode{} })},
		MinetestServerRepo:     &MinetestServerRepository{dbu: dbutil.New[*types.MinetestServer](db, dbutil.DialectSQLite, func() *types.MinetestServer { return &types.MinetestServer{} })},
		JobRepo:                &JobRepository{dbu: dbutil.New[*types.Job](db, dbutil.DialectSQLite, func() *types.Job { return &types.Job{} })},
		PaymentTransactionRepo: &PaymentTransactionRepository{dbu: dbutil.New[*types.PaymentTransaction](db, dbutil.DialectSQLite, func() *types.PaymentTransaction { return &types.PaymentTransaction{} })},
	}
}
