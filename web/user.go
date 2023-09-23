package web

import (
	"encoding/json"
	"mt-hosting-manager/types"
	"net/http"
)

func (a *Api) SearchUser(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	s := &types.UserSearch{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	list, err := a.repos.UserRepo.Search(s)
	Send(w, list, err)
}
