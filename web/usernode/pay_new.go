package usernode

import (
	"fmt"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/types"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func (ctx *Context) PayNew(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)

	nodetype_id := vars["nodetype-id"]
	months_str := vars["months"]

	nodetype, err := ctx.repos.NodeTypeRepo.GetByID(nodetype_id)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}
	if nodetype == nil {
		ctx.tu.RenderError(w, r, 404, fmt.Errorf("nodetype not found: %s", nodetype_id))
		return
	}

	costs := strings.Split(nodetype.Cost, ";")
	months := strings.Split(nodetype.MonthChoices, ";")

	item := &wallee.LineItem{
		Name:     nodetype.Name,
		Quantity: 1,
		Type:     wallee.LineItemTypeProduct,
		UniqueID: nodetype.ID,
	}

	cost_found := false
	for i, m := range months {
		if m == months_str {
			cost, err := strconv.ParseFloat(costs[i], 64)
			if err != nil {
				ctx.tu.RenderError(w, r, 500, err)
				return
			}
			item.AmountIncludingTax = cost
			cost_found = true
		}
	}

	if !cost_found {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("month selection not found: '%s'", months_str))
		return
	}

	tx, err := ctx.wc.CreateTransaction(&wallee.TransactionRequest{
		Currency:  "EUR",
		LineItems: []*wallee.LineItem{item},
	})
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	fmt.Printf("Tx-id: %d\n", tx.ID)

	ctx.tu.ExecuteTemplate(w, r, "usernode/pay_new.html", nil)
}
