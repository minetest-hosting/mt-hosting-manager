package web

import (
	"mt-hosting-manager/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (ctx *Context) NodeTypes(w http.ResponseWriter, r *http.Request, c *types.Claims) {

	ctx.tu.ExecuteTemplate(w, r, "node_types.html", nil)
}

func (ctx *Context) NodeTypeEdit(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx.tu.ExecuteTemplate(w, r, "node_type_edit.html", id)
}
