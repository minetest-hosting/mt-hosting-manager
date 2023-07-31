package usernode

import (
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/db"
	"mt-hosting-manager/tmpl"
	"mt-hosting-manager/types"
	"os"

	"github.com/gorilla/mux"
)

type Context struct {
	tu    *tmpl.TemplateUtil
	repos *db.Repositories
	cfg   *types.Config
	wc    *wallee.WalleeClient
}

func New(tu *tmpl.TemplateUtil, repos *db.Repositories, cfg *types.Config) *Context {
	return &Context{
		tu:    tu,
		repos: repos,
		cfg:   cfg,
		wc: wallee.New(
			os.Getenv("WALLEE_USERID"),
			os.Getenv("WALLEE_SPACEID"),
			os.Getenv("WALLEE_KEY"),
		),
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
