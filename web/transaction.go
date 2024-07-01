package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) CreateTransaction(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	create_tx_req := &types.CreateTransactionRequest{}
	err := json.NewDecoder(r.Body).Decode(create_tx_req)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	create_tx_resp, err := a.core.CreateTransaction(c.UserID, create_tx_req)
	Send(w, create_tx_resp, err)
}

func (a *Api) CheckTransaction(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	tx, err := a.core.CheckTransaction(id)
	Send(w, tx, err)
}

func (a *Api) GetTransactions(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	list, err := a.repos.PaymentTransactionRepo.GetByUserID(c.UserID)
	Send(w, list, err)
}

func (a *Api) GetTransaction(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	tx, err := a.repos.PaymentTransactionRepo.GetByID(id)
	if err != nil {
		SendError(w, 500, fmt.Errorf("failed to fetch transaction %s: %v", tx.ID, err))
		return
	}
	if tx == nil {
		SendError(w, 404, fmt.Errorf("transaction not found %s", id))
		return
	}
	if tx.UserID != c.UserID && c.Role != types.UserRoleAdmin {
		SendError(w, 403, fmt.Errorf("not authorized to fetch %s", id))
		return
	}

	Send(w, tx, err)
}

func (a *Api) SearchTransaction(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	s := &types.PaymentTransactionSearch{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	if c.Role != types.UserRoleAdmin {
		// non-admins can only search their own transactions
		s.UserID = &c.UserID
	}

	list, err := a.repos.PaymentTransactionRepo.Search(s)
	Send(w, list, err)
}
