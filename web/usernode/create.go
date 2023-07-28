package usernode

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CreateModel struct {
	NodeTypes []*types.NodeType
	HasError  bool
	Name      string
	NameErr   string
}

// create new node
func (ctx *Context) Create(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	nodetypes, err := ctx.repos.NodeTypeRepo.GetByState(types.NodeTypeStateActive)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	model := &CreateModel{
		NodeTypes: nodetypes,
	}

	if r.Method == http.MethodPost {
		// new node
		err := r.ParseForm()
		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}

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

		// check for valid name
		model.Name = r.FormValue("name")
		if !types.ValidUserNodeName.MatchString(model.Name) {
			model.HasError = true
			model.NameErr = "invalid-server-name"
		}

		// check for duplicate name
		existing_node, err := ctx.repos.UserNodeRepo.GetByName(model.Name)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}
		if existing_node != nil {
			model.HasError = true
			model.NameErr = "duplicate-server-name"
		}

		if !model.HasError {
			// everything ok, add server and redirect to detail page
			user_node := &types.UserNode{
				ID:         uuid.NewString(),
				UserID:     c.ID,
				NodeTypeID: node_type.ID,
				Created:    time.Now().Unix(),
				State:      types.UserNodeStateCreated,
				Name:       model.Name,
			}
			err = ctx.repos.UserNodeRepo.Insert(user_node)
			if err != nil {
				ctx.tu.RenderError(w, r, 500, err)
				return
			}
			//TODO: add provisioning job

			// redirect to detail page
			http.Redirect(w, r, fmt.Sprintf("%s/nodes/%s", ctx.tu.BaseURL, user_node.ID), http.StatusSeeOther)
		}
	}

	ctx.tu.ExecuteTemplate(w, r, "usernode/create.html", model)
}
