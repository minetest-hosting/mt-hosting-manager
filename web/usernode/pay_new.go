package usernode

import (
	"fmt"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/types"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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

	if r.Method == http.MethodGet {
		// show details
		ctx.tu.ExecuteTemplate(w, r, "usernode/pay_new.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		// create tx and redirect to payment site
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

		payment_tx_id := uuid.NewString()

		tx, err := ctx.wc.CreateTransaction(&wallee.TransactionRequest{
			Currency:   "EUR",
			LineItems:  []*wallee.LineItem{item},
			SuccessURL: fmt.Sprintf("%s/nodes/pay-callback/%s", ctx.tu.BaseURL, payment_tx_id),
			//TODO: failedURL
		})
		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}

		url, err := ctx.wc.CreatePaymentPageURL(tx.ID)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}

		// create usernode
		node := &types.UserNode{
			ID:         uuid.NewString(),
			UserID:     c.UserID,
			NodeTypeID: nodetype.ID,
			Created:    time.Now().Unix(),
			Expires:    time.Now().Unix(),
			State:      types.UserNodeStateCreated,
			Name:       uuid.NewString(), //TODO: better hostnames
		}
		err = ctx.repos.UserNodeRepo.Insert(node)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}

		payment_tx := &types.PaymentTransaction{
			ID:            payment_tx_id,
			TransactionID: fmt.Sprintf("%d", tx.ID),
			Created:       time.Now().Unix(),
			UserNodeID:    node.ID,
		}
		err = ctx.repos.PaymentTransactionRepo.Insert(payment_tx)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}

		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}
