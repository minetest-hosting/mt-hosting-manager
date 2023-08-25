package web

import (
	"errors"
	"mt-hosting-manager/types"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func (a *Api) Logout(w http.ResponseWriter, r *http.Request) {
	a.RemoveClaims(w)
}

func (a *Api) GetLogin(w http.ResponseWriter, r *http.Request) {
	claims, err := a.GetClaims(r)
	if err == err_unauthorized {
		SendError(w, 401, errors.New("unauthorized"))
	} else if err != nil {
		SendError(w, 500, err)
	} else {
		// refresh token
		auth_entry, err := a.repos.UserRepo.GetByMail(claims.Mail)
		if err != nil {
			SendError(w, 500, err)
			return
		}
		if auth_entry == nil {
			SendError(w, 404, errors.New("auth entry not found"))
			return
		}

		expires := time.Now().Add(7 * 24 * time.Hour)
		new_claims := &types.Claims{
			RegisteredClaims: &jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expires),
			},
			Mail:   claims.Mail,
			UserID: claims.ID,
			Role:   claims.Role,
		}

		err = a.SetClaims(w, claims)
		Send(w, new_claims, err)
	}
}
