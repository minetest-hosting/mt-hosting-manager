package db

import (
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
	Lock                   *DBLock
	g                      *gorm.DB
}

func NewRepositories(g *gorm.DB) *Repositories {
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
		Lock:                   &DBLock{g: g},
		g:                      g,
	}
}

func (r *Repositories) Gorm() *gorm.DB {
	return r.g
}
