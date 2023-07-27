package web

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
)

func (ctx *Context) UserNodes(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	model := make(map[string]any)

	model["Nodes"] = []*types.UserNode{}

	ctx.tu.ExecuteTemplate(w, r, "user_node.html", model)
}

func (ctx *Context) UserNodeEdit(w http.ResponseWriter, r *http.Request, c *types.Claims) {

	// TODO
	ctx.tu.ExecuteTemplate(w, r, "user_node_edit.html", nil)
}

func (ctx *Context) UserNodeSave(w http.ResponseWriter, r *http.Request, c *types.Claims) {

	// TODO
	http.Redirect(w, r, fmt.Sprintf("%s/nodes", ctx.BaseURL), http.StatusSeeOther)
}
