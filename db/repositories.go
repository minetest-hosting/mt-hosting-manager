package db

import (
	"github.com/minetest-go/dbutil"
	"gorm.io/gorm"
)

type Repositories struct {
	UserRepo               *UserRepository
	NodeTypeRepo           *NodeTypeRepository
	UserNodeRepo           *UserNodeRepository
	MinetestServerRepo     *MinetestServerRepository
	JobRepo                *JobRepository
	PaymentTransactionRepo *PaymentTransactionRepository
	AuditLogRepo           *AuditLogRepository
	BackupRepo             *BackupRepository
	BackupSpaceRepo        *BackupSpaceRepository
	ExchangeRateRepo       *ExchangeRateRepository
}

func NewRepositories(db dbutil.DBTx, g *gorm.DB) *Repositories {
	return &Repositories{
		UserRepo:               &UserRepository{g: g},
		NodeTypeRepo:           &NodeTypeRepository{g: g},
		UserNodeRepo:           &UserNodeRepository{g: g},
		MinetestServerRepo:     &MinetestServerRepository{g: g},
		JobRepo:                &JobRepository{g: g},
		PaymentTransactionRepo: &PaymentTransactionRepository{g: g},
		AuditLogRepo:           &AuditLogRepository{g: g},
		BackupRepo:             &BackupRepository{g: g},
		BackupSpaceRepo:        &BackupSpaceRepository{g: g},
		ExchangeRateRepo:       &ExchangeRateRepository{g: g},
	}
}
