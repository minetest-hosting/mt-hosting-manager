package usernode

import (
	"mt-hosting-manager/db"
	"mt-hosting-manager/tmpl"

	"github.com/gorilla/mux"
)

type Context struct {
	tu    *tmpl.TemplateUtil
	repos *db.Repositories
}

func New(tu *tmpl.TemplateUtil, repos *db.Repositories) *Context {
	return &Context{
		tu:    tu,
		repos: repos,
	}
}

func (ctx *Context) Setup(r *mux.Router) {
	r.HandleFunc("/nodes", ctx.tu.Secure(ctx.List))
	r.HandleFunc("/nodes/new", ctx.tu.Secure(ctx.Create))
	r.HandleFunc("/nodes/{id}", ctx.tu.Secure(ctx.Detail))
	r.HandleFunc("/nodes/{id}/delete", ctx.tu.Secure(ctx.Delete))
}
