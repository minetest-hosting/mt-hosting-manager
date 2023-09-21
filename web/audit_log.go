package web

import (
	"encoding/json"
	"mt-hosting-manager/types"
	"net/http"
)

func (a *Api) SearchAuditLog(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	s := &types.AuditLogSearch{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	if c.Role != types.UserRoleAdmin {
		// non-admins can only search their own audit-log
		s.UserID = &c.UserID
	}

	list, err := a.repos.AuditLogRepo.Search(s)
	Send(w, list, err)
}
