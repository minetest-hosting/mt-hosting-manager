package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/types"
	"net/http"
	"strconv"
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

	back_url := fmt.Sprintf("%s/#/finance", a.cfg.BaseURL)
	tx, err := a.wc.CreateTransaction(&wallee.TransactionRequest{
		Currency:   user.Currency,
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

func (a *Api) CheckTransaction(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	tx, err := a.repos.PaymentTransactionRepo.GetByID(id)
	if err != nil {
		SendError(w, 500, fmt.Errorf("fetch payment tx failed: %v", err))
		return
	}
	if tx == nil {
		SendError(w, 404, fmt.Errorf("payment tx not found: %s", id))
		return
	}

	if tx.State == types.PaymentStatePending {
		// verify tx success
		txr := &wallee.TransactionSearchRequest{
			Filter: &wallee.TransactionSearchFilter{
				FieldName: "id",
				Operator:  wallee.FilterOperatorEquals,
				Type:      wallee.FilterTypeLeaf,
				Value:     tx.TransactionID,
			},
		}
		tx_list, err := a.wc.SearchTransaction(txr)
		if err != nil {
			SendError(w, 500, fmt.Errorf("failed to fetch transaction %s: %v", tx.ID, err))
			return
		}
		if tx_list == nil || len(tx_list) != 1 {
			SendError(w, 500, fmt.Errorf("transaction not found %s: %v", tx.ID, err))
			return
		}
		verfifed_tx := tx_list[0]
		if verfifed_tx.State == wallee.TransactionStateFulfilled {
			tx.State = types.PaymentStateSuccess
			err = a.repos.PaymentTransactionRepo.Update(tx)
			if err != nil {
				SendError(w, 500, fmt.Errorf("failed to save transaction: %v", err))
				return
			}
		}
	}

	Send(w, tx, nil)
}

func (a *Api) GetTransactions(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	list, err := a.repos.PaymentTransactionRepo.GetByUserID(c.UserID)
	// TODO: only show needed fields
	Send(w, list, err)
}
