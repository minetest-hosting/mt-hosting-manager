package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
	"strings"
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
		a.RemoveClaims(w)
		w.WriteHeader(401)
		w.Write([]byte("unauthorized"))
	} else if err != nil {
		a.RemoveClaims(w)
		SendError(w, 500, err)
	} else {
		// refresh token
		auth_entry, err := a.repos.UserRepo.GetByID(claims.UserID)
		if err != nil {
			a.RemoveClaims(w)
			SendError(w, 500, err)
			return
		}
		if auth_entry == nil {
			a.RemoveClaims(w)
			SendError(w, 404, fmt.Errorf("auth entry not found: userid='%s' username='%s'", claims.UserID, claims.Name))
			return
		}

		new_claims, err := a.loginUser(w, r, auth_entry)
		Send(w, new_claims, err)
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *Api) Login(w http.ResponseWriter, r *http.Request) {
	lr := &LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(lr)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	user, err := a.repos.UserRepo.GetByName(lr.Username)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if user == nil {
		SendError(w, 404, fmt.Errorf("user with name '%s' not found", lr.Username))
		return
	}
	if user.Type != types.UserTypeLocal {
		SendError(w, 405, fmt.Errorf("non-local user"))
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

	_, err = a.loginUser(w, r, user)
	user.RemoveSensitiveFields()
	Send(w, user, err)
}

func (a *Api) loginUser(w http.ResponseWriter, req *http.Request, user *types.User) (*types.Claims, error) {
	dur := time.Duration(24 * 180 * time.Hour)
	claims := &types.Claims{
		Name:   user.Name,
		Role:   user.Role,
		UserID: user.ID,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(dur)),
		},
	}

	l := &types.AuditLog{
		Type:   types.AuditLogUserLoggedIn,
		UserID: claims.UserID,
	}
	if req != nil {
		// web request
		fwdfor := req.Header.Get("X-Forwarded-For")
		if fwdfor != "" {
			// behind reverse proxy
			parts := strings.Split(fwdfor, ",")
			l.IPAddress = &parts[0]
		} else {
			// direct access
			parts := strings.Split(req.RemoteAddr, ":")
			l.IPAddress = &parts[0]
		}
	}

	a.core.AddAuditLog(l)

	return claims, a.SetClaims(w, claims)
}
