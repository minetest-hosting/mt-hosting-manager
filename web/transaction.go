package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/types"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
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

	payment_tx_id := uuid.NewString()

	amount, err := strconv.ParseFloat(create_tx_req.Amount, 64)
	if err != nil {
		SendError(w, 500, fmt.Errorf("parse amount failed: %v", err))
		return
	}

	item := &wallee.LineItem{
		Name:               "Minetest hosting credits",
		Quantity:           1,
		AmountIncludingTax: amount,
		Type:               wallee.LineItemTypeProduct,
		UniqueID:           payment_tx_id,
	}

	wc := wallee.New(
		os.Getenv("WALLEE_USERID"),
		os.Getenv("WALLEE_SPACEID"),
		os.Getenv("WALLEE_KEY"),
	)

	tx, err := wc.CreateTransaction(&wallee.TransactionRequest{
		Currency:   user.Currency,
		LineItems:  []*wallee.LineItem{item},
		SuccessURL: fmt.Sprintf("%s/transactions/%s", a.cfg.BaseURL, payment_tx_id),
	})
	if err != nil {
		SendError(w, 500, fmt.Errorf("create transaction failed: %v", err))
		return
	}

	url, err := wc.CreatePaymentPageURL(tx.ID)
	if err != nil {
		SendError(w, 500, fmt.Errorf("create payment url failed: %v", err))
		return
	}

	payment_tx := &types.PaymentTransaction{
		ID:            payment_tx_id,
		TransactionID: fmt.Sprintf("%d", tx.ID),
		Created:       time.Now().Unix(),
		UserID:        c.UserID,
		Amount:        create_tx_req.Amount,
		Currency:      user.Currency,
		State:         types.PaymentStatePending,
	}
	err = a.repos.PaymentTransactionRepo.Insert(payment_tx)
	if err != nil {
		SendError(w, 500, fmt.Errorf("payment tx insert failed: %v", err))
		return
	}

	create_tx_resp := &types.CreateTransactionResponse{
		URL: url,
	}

	Send(w, create_tx_resp, nil)
}

func (a *Api) TransactionCallback(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	//TODO convert amount in EUR currency
}

func (a *Api) GetTransactions(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	list, err := a.repos.PaymentTransactionRepo.GetByUserID(c.UserID)
	// TODO: only show needed fields
	Send(w, list, err)
}
