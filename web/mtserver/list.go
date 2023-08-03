package mtserver

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
)

type ServerListModel struct {
	Servers []*types.MinetestServer
}

func (ctx *Context) List(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	servers, err := ctx.repos.MinetestServerRepo.GetByUserID(c.UserID)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("fetch servers error: %v", err))
		return
	}

	m := &ServerListModel{
		Servers: servers,
	}

	ctx.tu.ExecuteTemplate(w, r, "mtserver/list.html", m)
}
