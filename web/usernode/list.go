package usernode

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
)

type UserNodeInfo struct {
	*types.UserNode
	*types.NodeType
}

// show all nodes by the user
func (ctx *Context) List(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	nodes, err := ctx.repos.UserNodeRepo.GetByUserID(c.ID)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	nodeinfos := make([]*UserNodeInfo, len(nodes))
	for i, n := range nodes {
		nt, err := ctx.repos.NodeTypeRepo.GetByID(n.NodeTypeID)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}
		if nt == nil {
			ctx.tu.RenderError(w, r, 404, fmt.Errorf("nodetype not found: '%s'", n.NodeTypeID))
			return
		}

		nodeinfos[i] = &UserNodeInfo{
			UserNode: n,
			NodeType: nt,
		}
	}

	model := map[string]any{
		"Nodes": nodeinfos,
	}

	ctx.tu.ExecuteTemplate(w, r, "usernode/list.html", model)
}

//TODO: remove (and confirm remove)
