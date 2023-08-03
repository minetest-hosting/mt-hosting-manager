package usernode

import (
	"fmt"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type OrderNewModel struct {
	NodeType   *types.NodeType
	Months     int
	TotalCost  string
	Expiration int64
}

func (ctx *Context) OrderNew(w http.ResponseWriter, r *http.Request, c *types.Claims) {
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

	months, err := strconv.ParseInt(months_str, 10, 32)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("month parse error: %v", err))
		return
	}
	if months < 1 || months > int64(nodetype.MaxMonths) {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("invalid month choice: %s", months_str))
		return
	}

	cost_per_month, err := strconv.ParseFloat(nodetype.MonthlyCost, 64)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("monthlycost parse error: %v", err))
		return
	}

	total_cost := cost_per_month * float64(months)

	if r.Method == http.MethodGet {
		// show details
		m := &OrderNewModel{
			NodeType:   nodetype,
			Months:     int(months),
			TotalCost:  fmt.Sprintf("%.2f", total_cost),
			Expiration: core.AddMonths(time.Now(), int(months)).Unix(),
		}

		ctx.tu.ExecuteTemplate(w, r, "usernode/order_new.html", m)
		return
	}

	if r.Method == http.MethodPost {
		// create tx and redirect to payment site
		item := &wallee.LineItem{
			Name:               nodetype.Name,
			Quantity:           float64(months),
			AmountIncludingTax: total_cost,
			Type:               wallee.LineItemTypeProduct,
			UniqueID:           nodetype.ID,
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

		payment_tx := &types.PaymentTransaction{
			ID:            payment_tx_id,
			TransactionID: fmt.Sprintf("%d", tx.ID),
			Created:       time.Now().Unix(),
			NodeTypeID:    nodetype.ID,
			Months:        int(months),
		}
		err = ctx.repos.PaymentTransactionRepo.Insert(payment_tx)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}

		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}
