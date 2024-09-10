package db

import (
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	g *gorm.DB
}

func (r *UserRepository) Insert(u *types.User) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return r.g.Create(u).Error
}

func (r *UserRepository) Update(u *types.User) error {
	return r.g.Model(u).Updates(u).Error
}

func (r *UserRepository) GetByID(id string) (*types.User, error) {
	var user *types.User
	err := r.g.Where(types.User{ID: id}).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetByName(name string) (*types.User, error) {
	var user *types.User
	err := r.g.Where(types.User{Name: name}).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetByTypeAndExternalID(t types.UserType, external_id string) (*types.User, error) {
	var user *types.User
	err := r.g.Where(types.User{Type: t, ExternalID: external_id}).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetAll() ([]*types.User, error) {
	var users []*types.User
	err := r.g.Find(&users).Error
	return users, err
}

func (r *UserRepository) Delete(user_id string) error {
	return r.g.Delete(types.User{ID: user_id}).Error
}

func (r *UserRepository) DeleteAll() error {
	return r.g.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(types.User{}).Error
}

func (r *UserRepository) AddBalance(user_id string, eurocents int64) error {
	return r.g.Exec("update public.user set balance = balance + $1 where id = $2", eurocents, user_id).Error
}

func (r *UserRepository) SubtractBalance(user_id string, eurocents int64) error {
	return r.g.Exec("update public.user set balance = balance - $1 where id = $2", eurocents, user_id).Error
}

func (r *UserRepository) Search(s *types.UserSearch) ([]*types.User, error) {
	q := r.g

	if s.NameLike != nil {
		q = q.Where("name like ?", *s.NameLike)
	}

	if s.UserID != nil {
		q = q.Where(types.User{ID: *s.UserID})
	}

	if s.Limit != nil && *s.Limit > 0 && *s.Limit < 100 {
		q = q.Limit(*s.Limit)
	} else {
		q = q.Limit(100)
	}

	var users []*types.User
	err := q.Find(&users).Error
	return users, err
}
