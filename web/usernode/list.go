package usernode

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
)

type UserNodeListModel struct {
	Nodes []*UserNodeInfo
}

type UserNodeInfo struct {
	*types.UserNode
	*types.NodeType
	ExpirationWarning bool
}

// show all nodes by the user
func (ctx *Context) List(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	nodes, err := ctx.repos.UserNodeRepo.GetByUserID(c.UserID)
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

	model := &UserNodeListModel{
		Nodes: nodeinfos,
	}

	ctx.tu.ExecuteTemplate(w, r, "usernode/list.html", model)
}

//TODO: remove (and confirm remove)
