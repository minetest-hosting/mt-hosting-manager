package mtserver

import (
	"mt-hosting-manager/db"
	"mt-hosting-manager/tmpl"
	"mt-hosting-manager/types"

	"github.com/gorilla/mux"
)

type Context struct {
	tu    *tmpl.TemplateUtil
	repos *db.Repositories
	cfg   *types.Config
}

func New(tu *tmpl.TemplateUtil, repos *db.Repositories, cfg *types.Config) *Context {
	return &Context{
		tu:    tu,
		repos: repos,
		cfg:   cfg,
	}
}

func (ctx *Context) Setup(r *mux.Router) {
	r.HandleFunc("/mtserver", ctx.tu.Secure(ctx.List))
	r.HandleFunc("/mtserver/create", ctx.tu.Secure(ctx.Create))
}
