package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SendActivationRequest struct {
	Mail string `json:"mail"`
}

func (a *Api) SendActivationMail(w http.ResponseWriter, r *http.Request) {
	sar := &SendActivationRequest{}
	err := json.NewDecoder(r.Body).Decode(sar)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	user, err := a.repos.UserRepo.GetByMail(sar.Mail)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if user == nil {
		// create new user
		user = &types.User{
			ID:           uuid.NewString(),
			Name:         sar.Mail,
			Mail:         sar.Mail,
			MailVerified: false,
			State:        types.UserStateActive,
			Created:      time.Now().Unix(),
			Balance:      0,
			WarnBalance:  500,
			Type:         types.UserTypeLocal,
			Role:         types.UserRoleUser,
		}
		err = a.repos.UserRepo.Insert(user)
		if err != nil {
			SendError(w, 500, err)
			return
		}
	}

	err = a.core.SendActivationMail(user)
	Send(w, true, err)
}

type ActivationCallbackRequest struct {
	Mail           string `json:"mail"`
	ActivationCode string `json:"activation_code"`
	NewPassword    string `json:"new_password"`
}

func (a *Api) ActivationCallback(w http.ResponseWriter, r *http.Request) {
	acr := &ActivationCallbackRequest{}
	err := json.NewDecoder(r.Body).Decode(acr)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	user, err := a.repos.UserRepo.GetByMail(acr.Mail)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if user == nil {
		SendError(w, 404, fmt.Errorf("user with mail '%s' not found", acr.Mail))
		return
	}
	if user.ActivationCode != acr.ActivationCode {
		SendError(w, 403, fmt.Errorf("activationcode does not match"))
		return
	}
	if acr.NewPassword == "" {
		SendError(w, 405, fmt.Errorf("empty password"))
		return
	}
	if len(acr.NewPassword) < 8 {
		SendError(w, 405, fmt.Errorf("password min-length not sufficient"))
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(acr.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	user.Hash = string(hash)
	user.MailVerified = true
	err = a.repos.UserRepo.Update(user)
	user.RemoveSensitiveFields()
	Send(w, user, err)

	a.core.AddAuditLog(&types.AuditLog{
		Type:   types.AuditLogUserActivated,
		UserID: user.ID,
	})
}
