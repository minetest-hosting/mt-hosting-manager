package web

import (
	"net/http"
)

func (ctx *Context) Index(w http.ResponseWriter, r *http.Request) {

	ctx.tu.ExecuteTemplate(w, r, "index.html", nil)
}
