package db

import (
	"mt-hosting-manager/types"
	"time"

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
	if l.Timestamp == 0 {
		l.Timestamp = time.Now().Unix()
	}
	return r.dbu.Insert(l)
}

func (r *AuditLogRepository) Search(s *types.AuditLogSearch) ([]*types.AuditLog, error) {

	q := "where timestamp > %s and timestamp < %s"
	params := []any{s.FromTimestamp, s.ToTimestamp}

	if s.Type != nil {
		q += " and type = %s"
		params = append(params, *s.Type)
	}

	if s.UserID != nil {
		q += " and user_id = %s"
		params = append(params, *s.UserID)
	}

	return r.dbu.SelectMulti(q, params...)
}
