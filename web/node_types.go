package web

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"

	"github.com/google/uuid"
)

func (ctx *Context) NodeTypes(w http.ResponseWriter, r *http.Request, c *types.Claims) {

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}

		switch r.Form.Get("action") {
		case "new":
			http.Redirect(w, r, fmt.Sprintf("%s/node_types/%s", ctx.BaseURL, uuid.NewString()), http.StatusSeeOther)
			return
		}
	}

	ctx.tu.ExecuteTemplate(w, r, "node_types.html", nil)
}
