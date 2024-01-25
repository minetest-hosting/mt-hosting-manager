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

	if s.MinetestServerID != nil {
		q += " and minetest_server_id = %s"
		params = append(params, *s.MinetestServerID)
	}

	if s.UserNodeID != nil {
		q += " and user_node_id = %s"
		params = append(params, *s.UserNodeID)
	}

	if s.BackupID != nil {
		q += " and backup_id = %s"
		params = append(params, *s.BackupID)
	}

	if s.PaymentTransactionID != nil {
		q += " and payment_transaction_id = %s"
		params = append(params, *s.PaymentTransactionID)
	}

	q += " order by timestamp desc limit 1000"

	return r.dbu.SelectMulti(q, params...)
}
