package web

import (
	"fmt"
	"mt-hosting-manager/notify"
	"mt-hosting-manager/types"
	"net/http"
)

func (api *Api) OauthCallback(w http.ResponseWriter, user *types.User, new_user bool) error {

	if new_user {
		notify.Send(&notify.NtfyNotification{
			Title:    fmt.Sprintf("New user signed up: %s", user.Name),
			Message:  fmt.Sprintf("Name: %s, Auth: %s", user.Name, user.Type),
			Priority: 3,
			Tags:     []string{"new"},
		}, true)

		api.core.AddAuditLog(&types.AuditLog{
			Type:   types.AuditLogUserCreated,
			UserID: user.ID,
		})
	}

	return api.loginUser(w, user)
}
