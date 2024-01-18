package db

import (
	"database/sql"
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type BackupSpaceRepository struct {
	dbu *dbutil.DBUtil[*types.BackupSpace]
}

func (r *BackupSpaceRepository) Insert(n *types.BackupSpace) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.dbu.Insert(n)
}

func (r *BackupSpaceRepository) Update(n *types.BackupSpace) error {
	return r.dbu.Update(n, "where id = %s", n.ID)
}

func (r *BackupSpaceRepository) GetByID(id string) (*types.BackupSpace, error) {
	nt, err := r.dbu.Select("where id = %s", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *BackupSpaceRepository) GetByUserID(user_id string) ([]*types.BackupSpace, error) {
	return r.dbu.SelectMulti("where user_id = %s", user_id)
}

func (r *BackupSpaceRepository) Delete(id string) error {
	return r.dbu.Delete("where id = %s", id)
}
