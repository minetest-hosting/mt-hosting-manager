package usernode

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
	"strings"
)

type NodeTypeInfo struct {
	*types.NodeType
	DescriptionParagraphs []string
}

type SelectNewModel struct {
	NodeTypes []*NodeTypeInfo
}

func (ctx *Context) SelectNew(w http.ResponseWriter, r *http.Request, c *types.Claims) {

	if r.Method == http.MethodPost {
		url := fmt.Sprintf("%s/nodes/order/%s/%s",
			ctx.tu.BaseURL, r.FormValue("node_type_id"), r.FormValue("months"),
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

		m.NodeTypes[i] = &NodeTypeInfo{
			NodeType:              nt,
			DescriptionParagraphs: strings.Split(nt.Description, "\n"),
		}
	}

	ctx.tu.ExecuteTemplate(w, r, "usernode/select_new.html", m)
}
