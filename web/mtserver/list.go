package mtserver

import (
	"fmt"
	"mt-hosting-manager/types"
	"mt-hosting-manager/web/components"
	"net/http"
)

type ServerListModel struct {
	Servers    []*types.MinetestServer
	Breadcrumb *components.Breadcrumb
}

func (ctx *Context) List(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	servers, err := ctx.repos.MinetestServerRepo.GetByUserID(c.UserID)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("fetch servers error: %v", err))
		return
	}

	m := &ServerListModel{
		Servers: servers,
		Breadcrumb: &components.Breadcrumb{
			Entries: []*components.BreadcrumbEntry{
				components.HomeBreadcrumb,
				components.ServersBreadcrumb,
			},
		},
	}

	ctx.tu.ExecuteTemplate(w, r, "mtserver/list.html", m)
}
