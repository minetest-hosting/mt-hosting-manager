package usernode

import (
	"mt-hosting-manager/types"
	"net/http"

	"github.com/gorilla/mux"
)

type DetailModel struct {
	UserNode *types.UserNode
}

// view details
func (ctx *Context) Detail(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	node, err := ctx.repos.UserNodeRepo.GetByID(id)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	m := &DetailModel{
		UserNode: node,
	}

	ctx.tu.ExecuteTemplate(w, r, "usernode/detail.html", m)
}
