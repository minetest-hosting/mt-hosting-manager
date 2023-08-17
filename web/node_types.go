package web

import (
	"mt-hosting-manager/types"
	"mt-hosting-manager/web/components"
	"net/http"
)

type NodeTypeModel struct {
	NodeTypes  []*types.NodeType
	Breadcrumb *components.Breadcrumb
}

func (ctx *Context) NodeTypes(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	list, err := ctx.repos.NodeTypeRepo.GetAll()
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	m := &NodeTypeModel{
		NodeTypes: list,
		Breadcrumb: &components.Breadcrumb{
			Entries: []*components.BreadcrumbEntry{
				{
					Name: "Home",
					Link: "/",
				}, {
					Name:   "Node types",
					Link:   "/node_types",
					Active: true,
				},
			},
		},
	}

	ctx.tu.ExecuteTemplate(w, r, "node_types.html", m)
}
