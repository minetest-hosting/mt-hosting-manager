package web

import (
	"fmt"
	"mt-hosting-manager/notify"
	"mt-hosting-manager/types"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func (api *Api) OauthCallback(w http.ResponseWriter, user *types.User, new_user bool) error {

	if new_user {
		notify.Send(&notify.NtfyNotification{
			Title:    fmt.Sprintf("New user signed up: %s", user.Name),
			Message:  fmt.Sprintf("Name: %s, Mail: %s, Auth: %s", user.Name, user.Mail, user.Type),
			Priority: 3,
			Tags:     []string{"new"},
		}, true)

		api.core.AddAuditLog(&types.AuditLog{
			Type:   types.AuditLogUserCreated,
			UserID: user.ID,
		})
	}

	if api.cfg.SignupWhitelist[0] != "" {
		// check whitelist
		found := false
		for _, mail := range api.cfg.SignupWhitelist {
			if user.Mail == mail {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("signup currently restricted, sorry")
		}
	}

	dur := time.Duration(24 * 180 * time.Hour)
	claims := &types.Claims{
		Mail:   user.Mail,
		Role:   user.Role,
		UserID: user.ID,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(dur)),
		},
	}

	api.core.AddAuditLog(&types.AuditLog{
		Type:   types.AuditLogUserLoggedIn,
		UserID: claims.UserID,
	})

	return api.SetClaims(w, claims)
}
