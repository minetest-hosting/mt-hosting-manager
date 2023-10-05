package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func (a *Api) Logout(w http.ResponseWriter, r *http.Request) {
	a.RemoveClaims(w)
}

func (a *Api) GetLogin(w http.ResponseWriter, r *http.Request) {
	claims, err := a.GetClaims(r)
	if err == err_unauthorized {
		w.WriteHeader(401)
		w.Write([]byte("unauthorized"))
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

		a.core.AddAuditLog(&types.AuditLog{
			Type:   types.AuditLogUserLoggedIn,
			UserID: claims.UserID,
		})

		err = a.SetClaims(w, claims)
		Send(w, new_claims, err)
	}
}

type LoginRequest struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

func (a *Api) Login(w http.ResponseWriter, r *http.Request) {
	lr := &LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(lr)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	user, err := a.repos.UserRepo.GetByMail(lr.Mail)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if user == nil {
		SendError(w, 404, fmt.Errorf("user with mail '%s' not found", lr.Mail))
		return
	}
	if user.Hash == "" {
		SendError(w, 405, fmt.Errorf("not allowed"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(lr.Password))
	if err != nil {
		SendError(w, 500, err)
		return
	}

	err = a.loginUser(w, user)
	user.RemoveSensitiveFields()
	Send(w, user, err)
}

func (a *Api) loginUser(w http.ResponseWriter, user *types.User) error {
	if a.cfg.SignupWhitelist[0] != "" {
		// check whitelist
		found := false
		for _, mail := range a.cfg.SignupWhitelist {
			if user.Mail == mail {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("login currently restricted, sorry")
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

	a.core.AddAuditLog(&types.AuditLog{
		Type:   types.AuditLogUserLoggedIn,
		UserID: claims.UserID,
	})

	return a.SetClaims(w, claims)
}
