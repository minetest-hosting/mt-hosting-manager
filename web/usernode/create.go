package usernode

import (
	"mt-hosting-manager/types"
	"net/http"
)

// create new node
func (ctx *Context) Create(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	nodetypes, err := ctx.repos.NodeTypeRepo.GetByState(types.NodeTypeStateActive)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	model := make(map[string]any)
	model["NodeTypes"] = nodetypes

	ctx.tu.ExecuteTemplate(w, r, "usernode/create.html", model)
}
