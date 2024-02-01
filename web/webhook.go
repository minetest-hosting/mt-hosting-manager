package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/api/zahlsch"
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

	//TODO: get tx and update status
}
