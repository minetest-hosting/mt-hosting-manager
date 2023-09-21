package core

import (
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"
)

type Core struct {
	repos *db.Repositories
	cfg   *types.Config
	wc    *wallee.WalleeClient
}

func New(repos *db.Repositories, cfg *types.Config) *Core {
	return &Core{repos: repos, cfg: cfg, wc: wallee.New()}
}
