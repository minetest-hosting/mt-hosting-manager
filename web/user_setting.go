package web

import (
	"fmt"
	"io"
	"mt-hosting-manager/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) GetUserSettings(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	list, err := a.repos.UserSettingRepo.GetByUserID(c.UserID)
	Send(w, list, err)
}

func (a *Api) SetUserSetting(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	key := vars["key"]
	body, err := io.ReadAll(r.Body)
	if err != nil {
		SendError(w, 500, fmt.Errorf("readall error: %v", err))
		return
	}

	err = a.repos.UserSettingRepo.Set(&types.UserSetting{
		UserID: c.UserID,
		Key:    key,
		Value:  string(body),
	})

	Send(w, true, err)
}

func (a *Api) DeleteUserSetting(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	key := vars["key"]
	err := a.repos.UserSettingRepo.Delete(c.UserID, key)
	Send(w, true, err)
}
