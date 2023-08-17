package mtserver

import (
	"fmt"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker"
	"net/http"

	"github.com/gorilla/mux"
)

type DeleteModel struct {
	Server        *types.MinetestServer
	ConfirmFailed bool
	ServerDeleted bool
}

func (ctx *Context) Delete(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	server, err := ctx.repos.MinetestServerRepo.GetByID(id)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}
	if server == nil {
		ctx.tu.RenderError(w, r, 404, fmt.Errorf("server not found: %s", id))
		return
	}

	node, err := ctx.repos.UserNodeRepo.GetByID(server.UserNodeID)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}
	if node.UserID != c.UserID {
		ctx.tu.RenderError(w, r, 403, fmt.Errorf("unauthorized"))
		return
	}

	m := &DeleteModel{
		Server:        server,
		ConfirmFailed: false,
		ServerDeleted: false,
	}

	if r.Method == http.MethodPost {
		switch r.FormValue("action") {
		case "delete":
			// mark for removal
			if r.FormValue("confirm_name") != server.Name {
				m.ConfirmFailed = true
				break
			}

			server.State = types.MinetestServerStateRemoving
			err = ctx.repos.MinetestServerRepo.Update(server)
			if err != nil {
				ctx.tu.RenderError(w, r, 500, fmt.Errorf("node-update failed: %v", err))
				return
			}

			// dispatch removal job
			job := worker.RemoveServerJob(node, server)
			err = ctx.repos.JobRepo.Insert(job)
			if err != nil {
				ctx.tu.RenderError(w, r, 500, fmt.Errorf("job insert failed: %v", err))
				return
			}

			m.ServerDeleted = true
		}
	}

	ctx.tu.ExecuteTemplate(w, r, "mtserver/delete.html", m)
}
