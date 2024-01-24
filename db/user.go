package db

import (
	"database/sql"
	"fmt"
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type UserRepository struct {
	dbu *dbutil.DBUtil[*types.User]
	db  dbutil.DBTx
}

func (r *UserRepository) Insert(u *types.User) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return r.dbu.Insert(u)
}

func (r *UserRepository) Update(u *types.User) error {
	return r.dbu.Update(u, "where id = %s", u.ID)
}

func (r *UserRepository) GetByID(id string) (*types.User, error) {
	u, err := r.dbu.Select("where id = %s", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return u, err
	}
}

func (r *UserRepository) GetByName(name string) (*types.User, error) {
	u, err := r.dbu.Select("where name = %s", name)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return u, err
	}
}

func (r *UserRepository) GetByTypeAndExternalID(t types.UserType, external_id string) (*types.User, error) {
	u, err := r.dbu.Select("where type = %s and external_id = %s", t, external_id)
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

func (r *UserRepository) DeleteAll() error {
	return r.dbu.Delete("")
}

func (r *UserRepository) AddBalance(user_id string, eurocents int64) error {
	_, err := r.db.Exec("update public.user set balance = balance + $1 where id = $2", eurocents, user_id)
	return err
}

func (r *UserRepository) SubtractBalance(user_id string, eurocents int64) error {
	_, err := r.db.Exec("update public.user set balance = balance - $1 where id = $2", eurocents, user_id)
	return err
}

func (r *UserRepository) Search(s *types.UserSearch) ([]*types.User, error) {
	q := "where true=true"
	params := []any{}

	if s.NameLike != nil {
		q += " and name like %s"
		params = append(params, *s.NameLike)
	}

	if s.Limit != nil && *s.Limit > 0 && *s.Limit < 100 {
		q += fmt.Sprintf(" limit %d", *s.Limit)
	} else {
		q += " limit 100"
	}

	return r.dbu.SelectMulti(q, params...)
}
