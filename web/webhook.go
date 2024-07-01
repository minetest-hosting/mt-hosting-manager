package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/api/zahlsch"
	"mt-hosting-manager/types"
	"net/http"
)

func (a *Api) ZahlschWebhook(w http.ResponseWriter, r *http.Request) {
	if !a.cfg.ZahlschEnabled {
		SendError(w, 500, fmt.Errorf("webhook not enabled"))
		return
	}
	key := r.URL.Query().Get("key")
	if key != a.cfg.ZahlschWebhookKey {
		SendError(w, 401, fmt.Errorf("key invalid"))
		return
	}

	payload := &zahlsch.WebhookPayload{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		SendError(w, 500, fmt.Errorf("json parse error: %v", err))
		return
	}

	if payload.Transaction == nil {
		SendError(w, 500, fmt.Errorf("no transaction found"))
		return
	}

	if payload.Transaction.PageID != a.cfg.ZahlschPageID {
		// ignore
		return
	}

	if payload.Transaction.Status != zahlsch.TransactionConfirmed {
		// not confirmed, ignore
		return
	}

	if payload.Transaction.Invoice == nil {
		SendError(w, 500, fmt.Errorf("invoice not found"))
		return
	}

	user_id := ""
	tx_id := ""

	for _, field := range payload.Transaction.Invoice.CustomFields {
		switch field.Name {
		case "user_id":
			user_id = field.Value
		case "transaction_id":
			tx_id = field.Value
		}
	}

	user, err := a.repos.UserRepo.GetByID(user_id)
	if err != nil {
		SendError(w, 500, fmt.Errorf("user fetch error: %v", err))
		return
	}
	if user == nil {
		SendError(w, 500, fmt.Errorf("user not found: '%s'", user_id))
		return
	}

	tx, err := a.repos.PaymentTransactionRepo.GetByID(tx_id)
	if err != nil {
		SendError(w, 500, fmt.Errorf("tx fetch error: %v", err))
		return
	}
	if tx == nil {
		SendError(w, 500, fmt.Errorf("tx not found: '%s'", tx_id))
		return
	}
	if tx.UserID != user.ID {
		SendError(w, 500, fmt.Errorf("user mismatch: found '%s' expected: '%s'", tx.UserID, user.ID))
		return
	}
	if tx.Amount != payload.Transaction.Amount {
		SendError(w, 500, fmt.Errorf("amount mismatch: found '%d' expected: '%d'", payload.Transaction.Amount, tx.Amount))
		return
	}

	err = a.repos.UserRepo.AddBalance(tx.UserID, tx.Amount)
	if err != nil {
		SendError(w, 500, fmt.Errorf("could not add balance '%d' to user '%s': %v", tx.Amount, tx.UserID, err))
		return
	}

	// everything ok, update status
	tx.State = types.PaymentStateSuccess
	err = a.repos.PaymentTransactionRepo.Update(tx)

	Send(w, true, err)
}
