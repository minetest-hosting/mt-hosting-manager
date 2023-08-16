package mtserver

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"

	"github.com/gorilla/mux"
)

type ServerDetailModel struct {
	Server *types.MinetestServer
	Node   *types.UserNode
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

	node, err := ctx.repos.UserNodeRepo.GetByID(server.UserNodeID)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("get node error: %v", err))
		return
	}

	m := &ServerDetailModel{
		Server: server,
		Node:   node,
	}

	ctx.tu.ExecuteTemplate(w, r, "mtserver/detail.html", m)
}
