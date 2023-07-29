package usernode

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (ctx *Context) Delete(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	node, err := ctx.repos.UserNodeRepo.GetByID(id)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}
	if node == nil {
		ctx.tu.RenderError(w, r, 404, fmt.Errorf("node not found: %s", id))
		return
	}
	if node.UserID != c.ID {
		ctx.tu.RenderError(w, r, 403, fmt.Errorf("unauthorized"))
		return
	}

	if r.Method == http.MethodPost {
		switch r.FormValue("action") {
		case "abort":
			http.Redirect(w, r, fmt.Sprintf("%s/nodes/%s", ctx.tu.BaseURL, id), http.StatusSeeOther)
		case "delete":
			// mark for removal
			node.State = types.UserNodeStateRemoving
			err = ctx.repos.UserNodeRepo.Update(node)
			if err != nil {
				ctx.tu.RenderError(w, r, 500, err)
				return
			}

			// dispatch removal job
			job := &types.Job{
				ID:         uuid.NewString(),
				Type:       types.JobTypeNodeDestroy,
				State:      types.JobStateCreated,
				UserNodeID: id,
			}
			err = ctx.repos.JobRepo.Insert(job)
			if err != nil {
				ctx.tu.RenderError(w, r, 500, err)
				return
			}

			http.Redirect(w, r, fmt.Sprintf("%s/nodes", ctx.tu.BaseURL), http.StatusSeeOther)
		}
		return
	}

	m := &DetailModel{
		UserNode: node,
	}

	ctx.tu.ExecuteTemplate(w, r, "usernode/delete.html", m)
}
