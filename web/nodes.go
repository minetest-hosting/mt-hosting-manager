package web

import (
	"mt-hosting-manager/types"
	"net/http"
)

func (ctx *Context) Nodes(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	model := make(map[string]any)

	model["Nodes"] = []*types.UserNode{}

	ctx.tu.ExecuteTemplate(w, r, "nodes.html", model)
}
