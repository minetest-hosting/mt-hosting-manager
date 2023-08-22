package web

import (
	"errors"
	"mt-hosting-manager/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) CreateTransaction(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	//TODO
}

func (a *Api) TransactionCallback(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	//TODO
}

func (a *Api) GetTransactions(w http.ResponseWriter, r *http.Request, c *types.Claims) {
}

func (a *Api) GetTransaction(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	tx, err := a.repos.PaymentTransactionRepo.GetByID(vars["id"])
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if tx.UserID != c.UserID {
		SendError(w, 403, errors.New("not allowed"))
		return
	}

	// TODO: only show needed fields
	Send(w, tx, err)
}
