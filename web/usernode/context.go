package usernode

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
	r.HandleFunc("/nodes", ctx.tu.Secure(ctx.List))
	r.HandleFunc("/nodes/select-new", ctx.tu.Secure(ctx.SelectNew))
	r.HandleFunc("/nodes/new/{nodetype-id}/{months}", ctx.tu.Secure(ctx.PayNew))
	r.HandleFunc("/nodes/pay-callback", ctx.tu.Secure(ctx.PayCallback))
	r.HandleFunc("/nodes/{id}", ctx.tu.Secure(ctx.Detail))
	r.HandleFunc("/nodes/{id}/delete", ctx.tu.Secure(ctx.Delete))
}
