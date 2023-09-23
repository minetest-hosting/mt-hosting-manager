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
	AuditLogRepo           *AuditLogRepository
}

func NewRepositories(db dbutil.DBTx) *Repositories {
	dialect := dbutil.DialectSQLite
	return &Repositories{
		UserRepo:               &UserRepository{dbu: dbutil.New[*types.User](db, dialect, types.UserProvider), db: db},
		NodeTypeRepo:           &NodeTypeRepository{dbu: dbutil.New[*types.NodeType](db, dialect, types.NodeTypeProvider)},
		UserNodeRepo:           &UserNodeRepository{dbu: dbutil.New[*types.UserNode](db, dialect, types.UserNodeProvider)},
		MinetestServerRepo:     &MinetestServerRepository{dbu: dbutil.New[*types.MinetestServer](db, dialect, types.MinetestServerProvider)},
		JobRepo:                &JobRepository{dbu: dbutil.New[*types.Job](db, dialect, types.JobProvider)},
		PaymentTransactionRepo: &PaymentTransactionRepository{dbu: dbutil.New[*types.PaymentTransaction](db, dialect, types.PaymentTransactionProvider)},
		AuditLogRepo:           &AuditLogRepository{dbu: dbutil.New[*types.AuditLog](db, dialect, types.AuditLogProvider)},
	}
}
