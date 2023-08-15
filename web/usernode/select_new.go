package usernode

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
	"strconv"
	"strings"
)

type NodeTypeInfo struct {
	*types.NodeType
	DescriptionParagraphs []string
	MonthlyCost           string
}

type SelectNewModel struct {
	NodeTypes []*NodeTypeInfo
}

func (ctx *Context) SelectNew(w http.ResponseWriter, r *http.Request, c *types.Claims) {

	if r.Method == http.MethodPost {
		url := fmt.Sprintf("%s/nodes/order/%s/%s",
			ctx.tu.BaseURL, r.FormValue("node_type_id"), r.FormValue("days"),
		)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	}

	ntl, err := ctx.repos.NodeTypeRepo.GetByState(types.NodeTypeStateActive)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	m := &SelectNewModel{
		NodeTypes: make([]*NodeTypeInfo, len(ntl)),
	}

	for i, nt := range ntl {

		daily_cost, _ := strconv.ParseFloat(nt.DailyCost, 64)

		m.NodeTypes[i] = &NodeTypeInfo{
			NodeType:              nt,
			DescriptionParagraphs: strings.Split(nt.Description, "\n"),
			MonthlyCost:           fmt.Sprintf("%.2f", daily_cost*30),
		}
	}

	ctx.tu.ExecuteTemplate(w, r, "usernode/select_new.html", m)
}
