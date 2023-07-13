package db

import (
	"github.com/minetest-go/dbutil"
)

type Repositories struct {
	UserRepo     *UserRepository
	NodeTypeRepo *NodeTypeRepository
}

func NewRepositories(db dbutil.DBTx) *Repositories {
	return &Repositories{
		UserRepo:     &UserRepository{DB: db},
		NodeTypeRepo: &NodeTypeRepository{DB: db},
	}
}
