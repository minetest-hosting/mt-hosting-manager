package usernode

import (
	"fmt"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"mt-hosting-manager/web/components"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type OrderNewModel struct {
	NodeType   *types.NodeType
	Days       int
	TotalCost  string
	Expiration int64
	Breadcrumb *components.Breadcrumb
}

func (ctx *Context) OrderNew(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)

	nodetype_id := vars["nodetype-id"]
	days_str := vars["days"]

	nodetype, err := ctx.repos.NodeTypeRepo.GetByID(nodetype_id)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}
	if nodetype == nil {
		ctx.tu.RenderError(w, r, 404, fmt.Errorf("nodetype not found: %s", nodetype_id))
		return
	}

	days, err := strconv.ParseInt(days_str, 10, 32)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("days parse error: %v", err))
		return
	}
	if days < 1 || days > int64(nodetype.MaxDays) {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("invalid day choice: %s", days_str))
		return
	}

	cost_per_day, err := strconv.ParseFloat(nodetype.DailyCost, 64)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("DailyCost parse error: %v", err))
		return
	}

	total_cost := cost_per_day * float64(days)

	if r.Method == http.MethodGet {
		// show details
		m := &OrderNewModel{
			NodeType:   nodetype,
			Days:       int(days),
			TotalCost:  fmt.Sprintf("%.2f", total_cost),
			Expiration: core.AddDays(time.Now(), int(days)).Unix(),
			Breadcrumb: &components.Breadcrumb{
				Entries: []*components.BreadcrumbEntry{
					components.HomeBreadcrumb, {
						Name:   "Order overview",
						FAIcon: "shopping-cart",
					}},
			},
		}

		ctx.tu.ExecuteTemplate(w, r, "usernode/order_new.html", m)
		return
	}

	if r.Method == http.MethodPost {
		// create tx and redirect to payment site
		item := &wallee.LineItem{
			Name:               nodetype.Name,
			Quantity:           float64(days),
			AmountIncludingTax: total_cost,
			Type:               wallee.LineItemTypeProduct,
			UniqueID:           nodetype.ID,
		}

		payment_tx_id := uuid.NewString()

		tx, err := ctx.wc.CreateTransaction(&wallee.TransactionRequest{
			Currency:   "EUR",
			LineItems:  []*wallee.LineItem{item},
			SuccessURL: fmt.Sprintf("%s/nodes/payment/%s", ctx.tu.BaseURL, payment_tx_id),
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
			StartDate:     time.Now().Unix(),
			UntilDate:     core.AddDays(time.Now(), int(days)).Unix(),
			State:         types.PaymentStatePending,
		}
		err = ctx.repos.PaymentTransactionRepo.Insert(payment_tx)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}

		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}
