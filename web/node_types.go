package web

import (
	"mt-hosting-manager/types"
	"net/http"
)

func (ctx *Context) NodeTypes(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	model := make(map[string]any)

	list, err := ctx.repos.NodeTypeRepo.GetAll()
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}
	model["NodeTypes"] = list

	ctx.tu.ExecuteTemplate(w, r, "node_types.html", model)
}
