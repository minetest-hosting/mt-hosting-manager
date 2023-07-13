package db

import (
	"github.com/minetest-go/dbutil"
)

type Repositories struct {
	UserRepo *UserRepository
}

func NewRepositories(db dbutil.DBTx) *Repositories {
	return &Repositories{
		UserRepo: &UserRepository{DB: db},
	}
}
