package web

import (
	"embed"
	"mt-hosting-manager/tmpl"

	"github.com/gorilla/mux"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

//go:embed *
var Files embed.FS

type Context struct {
	tu *tmpl.TemplateUtil
}

func (ctx *Context) Setup(r *mux.Router) {
	r.HandleFunc("/", ctx.Index)
	r.PathPrefix("/assets").Handler(statigz.FileServer(Files, brotli.AddEncoding))

	r.NotFoundHandler = ctx.NotFound()
}
