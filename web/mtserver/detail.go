package mtserver

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"

	"github.com/gorilla/mux"
)

type ServerDetailModel struct {
	Server *types.MinetestServer
}

func (ctx *Context) Detail(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	server, err := ctx.repos.MinetestServerRepo.GetByID(id)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("get server error: %v", err))
		return
	}
	if server == nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("server not found: %s", id))
		return
	}

	m := &ServerDetailModel{
		Server: server,
	}

	ctx.tu.ExecuteTemplate(w, r, "mtserver/detail.html", m)
}
