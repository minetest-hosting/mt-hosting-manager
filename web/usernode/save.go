package usernode

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// POST for edit / create
func (ctx *Context) Save(w http.ResponseWriter, r *http.Request, c *types.Claims) {

	err := r.ParseForm()
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	id := r.FormValue("id")
	user_node, err := ctx.repos.UserNodeRepo.GetByID(id)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	// owner check
	if user_node != nil && user_node.UserID != c.ID {
		ctx.tu.RenderError(w, r, 403, fmt.Errorf("Unauthorized"))
	}

	if user_node == nil {
		// new node
		nodetypeid := r.FormValue("nodetype")
		node_type, err := ctx.repos.NodeTypeRepo.GetByID(nodetypeid)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}
		if node_type == nil {
			ctx.tu.RenderError(w, r, 404, fmt.Errorf("nodetype not found: '%s'", nodetypeid))
			return
		}

		user_node = &types.UserNode{
			ID:         uuid.NewString(),
			UserID:     c.ID,
			NodeTypeID: node_type.ID,
			Created:    time.Now().Unix(),
			State:      types.UserNodeStateCreated,
			Name:       r.FormValue("name"),
		}
		err = ctx.repos.UserNodeRepo.Insert(user_node)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}
		//TODO: add provisioning job

		// redirect to detail page
		http.Redirect(w, r, fmt.Sprintf("%s/nodes/%s", ctx.tu.BaseURL, user_node.ID), http.StatusSeeOther)
		return

	} else {
		// existing node
		//TODO
	}

	http.Redirect(w, r, fmt.Sprintf("%s/nodes", ctx.tu.BaseURL), http.StatusSeeOther)
}