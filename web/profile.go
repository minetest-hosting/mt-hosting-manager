package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
)

func (a *Api) GetUserProfile(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	user, err := a.repos.UserRepo.GetByID(c.UserID)
	user.RemoveSensitiveFields()
	Send(w, user, err)
}

func (a *Api) UpdateUserProfile(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	user, err := a.repos.UserRepo.GetByID(c.UserID)
	if err != nil {
		SendError(w, 500, fmt.Errorf("user fetch error: %v", err))
		return
	}
	if user == nil {
		SendError(w, 404, fmt.Errorf("user not found: '%s'", c.UserID))
		return
	}

	updated_user := &types.User{}
	err = json.NewDecoder(r.Body).Decode(updated_user)
	if err != nil {
		SendError(w, 500, fmt.Errorf("json error: %v", err))
		return
	}

	// update allowed fields
	user.Mail = updated_user.Mail
	user.Currency = updated_user.Currency

	err = a.repos.UserRepo.Update(user)
	Send(w, user, err)
}
