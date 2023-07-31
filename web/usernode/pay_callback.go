package usernode

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"

	"github.com/gorilla/mux"
)

type PayCallbackModel struct {
	Transaction *types.PaymentTransaction
	Node        *types.UserNode
}

func (ctx *Context) PayCallback(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	txid := vars["tx-id"]
	tx, err := ctx.repos.PaymentTransactionRepo.GetByID(txid)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}
	if tx == nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("transaction not found: %s", txid))
		return
	}

	// remove tx entry after success
	err = ctx.repos.PaymentTransactionRepo.Delete(tx.ID)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	node, err := ctx.repos.UserNodeRepo.GetByID(tx.UserNodeID)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	m := &PayCallbackModel{
		Transaction: tx,
		Node:        node,
	}

	//TODO: verify tx success
	//TODO: start node provisioning

	ctx.tu.ExecuteTemplate(w, r, "usernode/pay_callback.html", m)
}
