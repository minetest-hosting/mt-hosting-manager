package db

import (
	"mt-hosting-manager/types"

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
	dialect := dbutil.DialectPostgres
	return &Repositories{
		UserRepo:               &UserRepository{g: g},
		NodeTypeRepo:           &NodeTypeRepository{g: g},
		UserNodeRepo:           &UserNodeRepository{g: g},
		MinetestServerRepo:     &MinetestServerRepository{g: g},
		JobRepo:                &JobRepository{g: g},
		PaymentTransactionRepo: &PaymentTransactionRepository{g: g},
		AuditLogRepo:           &AuditLogRepository{dbu: dbutil.New[*types.AuditLog](db, dialect, types.AuditLogProvider)},
		BackupRepo:             &BackupRepository{dbu: dbutil.New[*types.Backup](db, dialect, types.BackupProvider)},
		BackupSpaceRepo:        &BackupSpaceRepository{dbu: dbutil.New[*types.BackupSpace](db, dialect, types.BackupSpaceProvider)},
		ExchangeRateRepo:       &ExchangeRateRepository{dbu: dbutil.New[*types.ExchangeRate](db, dialect, types.ExchangeRateProvider)},
	}
}
