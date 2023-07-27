package web

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
)

func (ctx *Context) ShowUserNodes(w http.ResponseWriter, r *http.Request, c *types.Claims) {

	nodes, err := ctx.repos.UserNodeRepo.GetByUserID(c.ID)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	model := make(map[string]any)
	model["Nodes"] = nodes

	ctx.tu.ExecuteTemplate(w, r, "user_node.html", model)
}

func (ctx *Context) UserNodeCreate(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	// TODO: create new node
	ctx.tu.ExecuteTemplate(w, r, "user_node_create.html", nil)
}

func (ctx *Context) UserNodeDetail(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	// TODO: view details
	ctx.tu.ExecuteTemplate(w, r, "user_node_detail.html", nil)
}

func (ctx *Context) UserNodeSave(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	// TODO: POST for edit and create
	http.Redirect(w, r, fmt.Sprintf("%s/nodes", ctx.BaseURL), http.StatusSeeOther)
}

//TODO: remove (and confirm remove)
