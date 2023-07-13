package db

import (
	"github.com/minetest-go/dbutil"
)

type Repositories struct {
	UserRepo *UserRepository
	NodeRepo *NodeRepository
}

func NewRepositories(db dbutil.DBTx) *Repositories {
	return &Repositories{
		UserRepo: &UserRepository{DB: db},
		NodeRepo: &NodeRepository{DB: db},
	}
}
