package web

import (
	"fmt"
	"mt-hosting-manager/types"
	"mt-hosting-manager/web/components"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type NodeTypeEditModel struct {
	NodeType   *types.NodeType
	Breadcrumb *components.Breadcrumb
	HasError   bool
}

func (ctx *Context) NodeTypeEdit(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	model := &NodeTypeEditModel{}

	nt, err := ctx.repos.NodeTypeRepo.GetByID(id)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	if id == "new" {
		nt = &types.NodeType{
			ID: uuid.NewString(),
		}
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}

		// remove action
		if r.FormValue("action") == "remove" {
			err = ctx.repos.NodeTypeRepo.Delete(id)
			if err != nil {
				ctx.tu.RenderError(w, r, 404, err)
				return
			}

			http.Redirect(w, r, fmt.Sprintf("%s/node_types", ctx.tu.BaseURL), http.StatusSeeOther)
			return
		}

		if nt == nil {
			nt = &types.NodeType{}
		}

		nt.Name = r.FormValue("name")
		nt.Description = r.FormValue("description")
		nt.State = r.FormValue("state")
		nt.Provider = types.ProviderType(r.FormValue("provider"))
		nt.ServerType = r.FormValue("server_type")
		nt.MonthChoices = r.FormValue("month_choices")
		nt.Cost = r.FormValue("cost")
		order_id, _ := strconv.ParseInt(r.FormValue("order_id"), 10, 32)
		nt.OrderID = int(order_id)
		num, _ := strconv.ParseInt(r.FormValue("max_recommended_instances"), 10, 32)
		nt.MaxRecommendedInstances = int(num)
		num, _ = strconv.ParseInt(r.FormValue("max_instances"), 10, 32)
		nt.MaxInstances = int(num)

		//TODO: validate input

		if nt.ID == "" {
			// insert new record
			nt.ID = r.FormValue("id")
			err = ctx.repos.NodeTypeRepo.Insert(nt)
		} else {
			// update
			err = ctx.repos.NodeTypeRepo.Update(nt)
		}

		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("%s/node_types", ctx.tu.BaseURL), http.StatusSeeOther)

	}

	if nt == nil {
		ctx.tu.RenderError(w, r, 404, err)
		return
	}

	model.NodeType = nt
	model.Breadcrumb = &components.Breadcrumb{
		Entries: []*components.BreadcrumbEntry{
			{
				Name: "Start",
				Link: "",
			}, {
				Name: "Node-Types",
				Link: "",
			}, {
				Name:   nt.ID,
				Link:   "",
				Active: true,
			},
		},
	}

	ctx.tu.ExecuteTemplate(w, r, "node_type_edit.html", model)
}
