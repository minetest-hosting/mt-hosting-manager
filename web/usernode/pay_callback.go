package usernode

import (
	"fmt"
	"math/rand"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/types"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

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

	// verify tx success
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

	randstr := RandStringRunes(7)

	// create usernode
	node := &types.UserNode{
		ID:         uuid.NewString(),
		UserID:     c.UserID,
		NodeTypeID: tx.NodeTypeID,
		Created:    time.Now().Unix(),
		Expires:    tx.UntilDate,
		State:      types.UserNodeStateCreated,
		Name:       fmt.Sprintf("node-%s-%s", os.Getenv("STAGE"), randstr),
		Alias:      randstr,
	}
	err = ctx.repos.UserNodeRepo.Insert(node)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("usernode insert failed: %v", err))
		return
	}

	// mark tx as successful after payment
	tx.State = types.PaymentStateSuccess
	tx.NodeID = node.ID
	err = ctx.repos.PaymentTransactionRepo.Update(tx)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("payment-tx update failed: %v", err))
		return
	}

	m := &PayCallbackModel{
		Transaction: tx,
		Node:        node,
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
