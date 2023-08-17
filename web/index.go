package web

import (
	"mt-hosting-manager/web/components"
	"net/http"
)

type StartModel struct {
	Breadcrumb *components.Breadcrumb
}

func (ctx *Context) Index(w http.ResponseWriter, r *http.Request) {
	m := &StartModel{
		Breadcrumb: &components.Breadcrumb{
			Entries: []*components.BreadcrumbEntry{
				{
					Name:   "Home",
					Link:   "/",
					Active: true,
				},
			},
		},
	}

	ctx.tu.ExecuteTemplate(w, r, "index.html", m)
}
