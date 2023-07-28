package usernode

import (
	"mt-hosting-manager/types"
	"net/http"
)

// show all nodes by the user
func (ctx *Context) List(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	nodes, err := ctx.repos.UserNodeRepo.GetByUserID(c.ID)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	model := make(map[string]any)
	model["Nodes"] = nodes

	ctx.tu.ExecuteTemplate(w, r, "usernode/list.html", model)
}

//TODO: remove (and confirm remove)
