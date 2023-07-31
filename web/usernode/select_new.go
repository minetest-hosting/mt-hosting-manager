package usernode

import (
	"mt-hosting-manager/types"
	"net/http"
	"strconv"
	"strings"
)

type TimeSelection struct {
	Months   int64
	Cost     float64
	Currency string
}

type NodeTypeInfo struct {
	*types.NodeType
	TimeSelection         []*TimeSelection
	DescriptionParagraphs []string
}

type SelectNewModel struct {
	NodeTypes []*NodeTypeInfo
}

func (ctx *Context) SelectNew(w http.ResponseWriter, r *http.Request, c *types.Claims) {

	ntl, err := ctx.repos.NodeTypeRepo.GetByState(types.NodeTypeStateActive)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	m := &SelectNewModel{
		NodeTypes: make([]*NodeTypeInfo, len(ntl)),
	}

	for i, nt := range ntl {

		months_str := strings.Split(nt.MonthChoices, ";")
		costs_str := strings.Split(nt.Cost, ";")

		tsl := make([]*TimeSelection, len(months_str))
		for j, month_str := range months_str {

			months, _ := strconv.ParseInt(month_str, 10, 64)
			cost, _ := strconv.ParseFloat(costs_str[j], 64)

			tsl[j] = &TimeSelection{
				Currency: "€",
				Months:   months,
				Cost:     cost,
			}
		}

		m.NodeTypes[i] = &NodeTypeInfo{
			NodeType:              nt,
			TimeSelection:         tsl,
			DescriptionParagraphs: strings.Split(nt.Description, "\n"),
		}
	}

	ctx.tu.ExecuteTemplate(w, r, "usernode/select_new.html", m)
}
