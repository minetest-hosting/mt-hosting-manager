package db

import (
	"database/sql"
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type UserRepository struct {
	DB dbutil.DBTx
}

func (r *UserRepository) Insert(u *types.User) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return dbutil.Insert(r.DB, u)
}

func (r *UserRepository) GetByMail(mail string) (*types.User, error) {
	u, err := dbutil.Select(r.DB, &types.User{}, "where mail = $1", mail)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return u, err
	}
}

func (r *UserRepository) GetAll() ([]*types.User, error) {
	return dbutil.SelectMulti(r.DB, func() *types.User { return &types.User{} }, "")
}

func (r *UserRepository) Delete(user_id string) error {
	return dbutil.Delete(r.DB, &types.User{}, "where id = $1", user_id)
}
