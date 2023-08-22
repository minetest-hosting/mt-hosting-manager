package db

import (
	"database/sql"
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type UserRepository struct {
	dbu *dbutil.DBUtil[*types.User]
}

func (r *UserRepository) Insert(u *types.User) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return r.dbu.Insert(u)
}

func (r *UserRepository) Update(u *types.User) error {
	return r.dbu.Update(u, "where mail =%s", u.Mail)
}

func (r *UserRepository) GetByMail(mail string) (*types.User, error) {
	u, err := r.dbu.Select("where mail = %s", mail)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return u, err
	}
}

func (r *UserRepository) GetAll() ([]*types.User, error) {
	return r.dbu.SelectMulti("")
}

func (r *UserRepository) Delete(user_id string) error {
	return r.dbu.Delete("where id = %s", user_id)
}
