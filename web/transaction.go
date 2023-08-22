package web

import (
	"mt-hosting-manager/types"
	"net/http"
)

func (a *Api) CreateTransaction(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	//TODO
}

func (a *Api) TransactionCallback(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	//TODO
}

func (a *Api) GetTransactions(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	list, err := a.repos.PaymentTransactionRepo.GetByUserID(c.UserID)
	// TODO: only show needed fields
	Send(w, list, err)
}
