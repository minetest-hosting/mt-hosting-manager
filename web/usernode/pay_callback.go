package usernode

import (
	"fmt"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/types"
	"net/http"
	"time"

	"github.com/google/uuid"
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
		ctx.tu.RenderError(w, r, 404, fmt.Errorf("transaction not found: %s", txid))
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
	txr := &wallee.TransactionSearchRequest{
		Filter: &wallee.TransactionSearchFilter{
			FieldName: "id",
			Operator:  wallee.FilterOperatorEquals,
			Type:      wallee.FilterTypeLeaf,
			Value:     tx.TransactionID,
		},
	}
	tx_list, err := ctx.wc.SearchTransaction(txr)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("failed to fetch transaction %s: %v", tx.ID, err))
		return
	}
	if tx_list == nil || len(tx_list) != 1 {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("transaction not found %s: %v", tx.ID, err))
		return
	}
	verfifed_tx := tx_list[0]
	if verfifed_tx.State != wallee.TransactionStateFulfilled {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("invalid transaction state %s", verfifed_tx.State))
		return
	}

	//TODO: increment expiration time properly
	node.Expires = time.Unix(node.Expires, 0).Add(time.Hour * 24 * 31).Unix()
	err = ctx.repos.UserNodeRepo.Update(node)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("failed to update expiration time on node %s: %v", node.ID, err))
		return
	}

	// start node provisioning
	job := &types.Job{
		ID:         uuid.NewString(),
		Type:       types.JobTypeNodeSetup,
		State:      types.JobStateCreated,
		Started:    time.Now().Unix(),
		UserNodeID: &node.ID,
	}
	err = ctx.repos.JobRepo.Insert(job)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	ctx.tu.ExecuteTemplate(w, r, "usernode/pay_callback.html", m)
}