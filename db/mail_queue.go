package db

import (
	"database/sql"
	"mt-hosting-manager/types"
	"time"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type MailQueueRepository struct {
	dbu *dbutil.DBUtil[*types.MailQueue]
}

func (r *MailQueueRepository) Insert(n *types.MailQueue) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	if n.State == "" {
		n.State = types.MailQueueStateCreated
	}
	if n.Timestamp == 0 {
		n.Timestamp = time.Now().Unix()
	}
	return r.dbu.Insert(n)
}

func (r *MailQueueRepository) Update(n *types.MailQueue) error {
	return r.dbu.Update(n, "where id = %s", n.ID)
}

func (r *MailQueueRepository) GetByID(id string) (*types.MailQueue, error) {
	m, err := r.dbu.Select("where id = %s", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return m, err
}

func (r *MailQueueRepository) GetLatestByReceiver(receiver string) (*types.MailQueue, error) {
	m, err := r.dbu.Select("where receiver = %s order by timestamp desc limit 1")
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return m, err
}

func (r *MailQueueRepository) GetByState(state types.MailQueueState) ([]*types.MailQueue, error) {
	return r.dbu.SelectMulti("where state = %s", state)
}
