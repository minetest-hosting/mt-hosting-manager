package core

import (
	"mt-hosting-manager/types"
	"time"

	"github.com/sirupsen/logrus"
)

func (c *Core) AddAuditLog(l *types.AuditLog) {
	l.Timestamp = time.Now().Unix()

	logrus.WithFields(logrus.Fields{
		"type":                   l.Type,
		"user_id":                l.UserID,
		"user_node_id":           l.UserNodeID,
		"minetest_server_id":     l.MinetestServerID,
		"payment_transaction_id": l.PaymentTransactionID,
	}).Info("audit-log")

	err := c.repos.AuditLogRepo.Insert(l)
	if err != nil {
		logrus.WithError(err).Error("Error inserting log")
	}
}
