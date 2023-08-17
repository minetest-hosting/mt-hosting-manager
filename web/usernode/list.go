package usernode

import (
	"fmt"
	"mt-hosting-manager/types"
	"mt-hosting-manager/web/components"
	"net/http"
)

type UserNodeListModel struct {
	Nodes      []*UserNodeInfo
	Breadcrumb *components.Breadcrumb
}

type UserNodeInfo struct {
	*types.UserNode
	*types.NodeType
	ExpirationWarning bool
	ServerCount       int
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
		if n.State == types.UserNodeStateRemoving {
			continue
		}

		nt, err := ctx.repos.NodeTypeRepo.GetByID(n.NodeTypeID)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}
		if nt == nil {
			ctx.tu.RenderError(w, r, 404, fmt.Errorf("nodetype not found: '%s'", n.NodeTypeID))
			return
		}

		servers, err := ctx.repos.MinetestServerRepo.GetByNodeID(n.ID)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, fmt.Errorf("fetch servers error: %v", err))
		}

		nodeinfos[i] = &UserNodeInfo{
			UserNode:    n,
			NodeType:    nt,
			ServerCount: len(servers),
		}
	}

	model := &UserNodeListModel{
		Nodes: nodeinfos,
		Breadcrumb: &components.Breadcrumb{
			Entries: []*components.BreadcrumbEntry{
				{
					Name: "Home",
					Link: "/",
				}, {
					Name: "Nodes",
					Link: "/nodes",
				},
			},
		},
	}

	ctx.tu.ExecuteTemplate(w, r, "usernode/list.html", model)
}

//TODO: remove (and confirm remove)
