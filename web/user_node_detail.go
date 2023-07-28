package web

import (
	"mt-hosting-manager/types"
	"net/http"
)

// view details
func (ctx *Context) UserNodeDetail(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	ctx.tu.ExecuteTemplate(w, r, "user_node_detail.html", nil)
}
