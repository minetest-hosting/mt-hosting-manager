package db

import (
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type AuditLogRepository struct {
	dbu *dbutil.DBUtil[*types.AuditLog]
}

func (r *AuditLogRepository) Insert(l *types.AuditLog) error {
	if l.ID == "" {
		l.ID = uuid.NewString()
	}
	return r.dbu.Insert(l)
}
