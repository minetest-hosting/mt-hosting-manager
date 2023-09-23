package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/types"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (a *Api) CreateTransaction(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	create_tx_req := &types.CreateTransactionRequest{}
	err := json.NewDecoder(r.Body).Decode(create_tx_req)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	user, err := a.repos.UserRepo.GetByMail(c.Mail)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if user == nil {
		SendError(w, 404, fmt.Errorf("user not found: '%s'", c.Mail))
		return
	}

	if user.Balance+create_tx_req.Amount > int64(a.cfg.MaxBalance) {
		SendError(w, 405, fmt.Errorf("max balance of %d exceeded", a.cfg.MaxBalance))
	}

	payment_tx_id := uuid.NewString()

	item := &wallee.LineItem{
		Name:               "Minetest hosting credits",
		Quantity:           1,
		AmountIncludingTax: float64(create_tx_req.Amount) / 100,
		Type:               wallee.LineItemTypeProduct,
		UniqueID:           payment_tx_id,
	}

	back_url := fmt.Sprintf("%s/#/finance/detail/%s", a.cfg.BaseURL, payment_tx_id)
	tx, err := a.wc.CreateTransaction(&wallee.TransactionRequest{
		Currency:   "EUR",
		LineItems:  []*wallee.LineItem{item},
		SuccessURL: back_url,
		FailedURL:  back_url,
	})
	if err != nil {
		SendError(w, 500, fmt.Errorf("create transaction failed: %v", err))
		return
	}

	url, err := a.wc.CreatePaymentPageURL(tx.ID)
	if err != nil {
		SendError(w, 500, fmt.Errorf("create payment url failed: %v", err))
		return
	}

	payment_tx := &types.PaymentTransaction{
		ID:             payment_tx_id,
		TransactionID:  fmt.Sprintf("%d", tx.ID),
		Created:        time.Now().Unix(),
		UserID:         c.UserID,
		Amount:         create_tx_req.Amount,
		AmountRefunded: 0,
		State:          types.PaymentStatePending,
	}
	err = a.repos.PaymentTransactionRepo.Insert(payment_tx)
	if err != nil {
		SendError(w, 500, fmt.Errorf("payment tx insert failed: %v", err))
		return
	}

	create_tx_resp := &types.CreateTransactionResponse{
		URL: url,
	}

	a.core.AddAuditLog(&types.AuditLog{
		Type:                 types.AuditLogPaymentCreated,
		UserID:               c.UserID,
		PaymentTransactionID: &payment_tx_id,
		Amount:               &create_tx_req.Amount,
	})

	Send(w, create_tx_resp, nil)
}

func (a *Api) CheckTransaction(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	tx, err := a.core.CheckTransaction(id)
	Send(w, tx, err)
}

func (a *Api) RefundTransaction(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	tx, err := a.core.RefundTransaction(id)
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
	if tx.UserID != c.UserID {
		SendError(w, 403, fmt.Errorf("not authorized to fetch %s", id))
		return
	}

	Send(w, tx, err)
}
