package db

import (
	"database/sql"
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type BackupRepository struct {
	dbu *dbutil.DBUtil[*types.Backup]
}

func (r *BackupRepository) Insert(n *types.Backup) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.dbu.Insert(n)
}

func (r *BackupRepository) Update(n *types.Backup) error {
	return r.dbu.Update(n, "where id = %s", n.ID)
}

func (r *BackupRepository) GetByID(id string) (*types.Backup, error) {
	nt, err := r.dbu.Select("where id = %s", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *BackupRepository) GetByBackupSpaceID(backup_space_id string) ([]*types.Backup, error) {
	return r.dbu.SelectMulti("where backup_space_id = %s", backup_space_id)
}

func (r *BackupRepository) Delete(id string) error {
	return r.dbu.Delete("where id = %s", id)
}
