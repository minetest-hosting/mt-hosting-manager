package usernode

import (
	"mt-hosting-manager/types"
	"net/http"
)

// view details
func (ctx *Context) Detail(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	ctx.tu.ExecuteTemplate(w, r, "usernode/detail.html", nil)
}
