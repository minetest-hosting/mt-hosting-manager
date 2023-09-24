package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) SendMail(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	user_id := vars["user_id"]

	m := &types.MailQueue{}
	err := json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	user, err := a.repos.UserRepo.GetByID(user_id)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if user == nil {
		SendError(w, 404, fmt.Errorf("user not found: '%s'", user_id))
		return
	}

	m.Receiver = user.Mail

	err = a.repos.MailQueueRepo.Insert(m)
	Send(w, m, err)
}
