package usernode

import (
	"fmt"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker"
	"net/http"

	"github.com/gorilla/mux"
)

type DeleteModel struct {
	UserNode      *types.UserNode
	ConfirmFailed bool
	NodeDeleted   bool
}

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
	if node.UserID != c.UserID {
		ctx.tu.RenderError(w, r, 403, fmt.Errorf("unauthorized"))
		return
	}

	m := &DeleteModel{
		UserNode:      node,
		ConfirmFailed: false,
		NodeDeleted:   false,
	}

	if r.Method == http.MethodPost {
		switch r.FormValue("action") {
		case "delete":
			// mark for removal
			if r.FormValue("confirm_name") != node.Name {
				m.ConfirmFailed = true
				break
			}

			node.State = types.UserNodeStateRemoving
			err = ctx.repos.UserNodeRepo.Update(node)
			if err != nil {
				ctx.tu.RenderError(w, r, 500, fmt.Errorf("node-update failed: %v", err))
				return
			}

			// dispatch removal job
			job := worker.RemoveNodeJob(node)
			err = ctx.repos.JobRepo.Insert(job)
			if err != nil {
				ctx.tu.RenderError(w, r, 500, fmt.Errorf("job insert failed: %v", err))
				return
			}

			m.NodeDeleted = true
		}
	}

	ctx.tu.ExecuteTemplate(w, r, "usernode/delete.html", m)
}
