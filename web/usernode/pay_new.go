package usernode

import (
	"mt-hosting-manager/types"
	"net/http"
)

func (ctx *Context) PayNew(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	ctx.tu.ExecuteTemplate(w, r, "usernode/pay_new.html", nil)
}
