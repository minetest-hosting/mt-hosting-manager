package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/types"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type SetPasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (a *Api) SetPassword(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	spr := &SetPasswordRequest{}
	err := json.NewDecoder(r.Body).Decode(spr)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	if spr.OldPassword == "" || spr.NewPassword == "" {
		SendError(w, 401, fmt.Errorf("empty password"))
		return
	}
	if len(spr.NewPassword) < 8 {
		SendError(w, 405, fmt.Errorf("password min-length not sufficient"))
		return
	}

	user, err := a.repos.UserRepo.GetByID(c.UserID)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	if user.Hash == "" {
		SendError(w, 401, fmt.Errorf("no old password found"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(spr.OldPassword))
	if err != nil {
		SendError(w, 401, fmt.Errorf("old password invalid"))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(spr.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	user.Hash = string(hash)
	err = a.repos.UserRepo.Update(user)
	Send(w, true, err)
}
