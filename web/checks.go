package web

import (
	"errors"
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
)

type SecureHandlerFunc func(http.ResponseWriter, *http.Request, *types.Claims)

// check for login only (UI)
func (api *Api) Secure(fn SecureHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := api.GetClaims(r)
		if err == err_unauthorized {
			SendError(w, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		} else if err != nil {
			SendError(w, http.StatusInternalServerError, err)
			return
		}
		fn(w, r, claims)
	}
}

type Check func(w http.ResponseWriter, r *http.Request) bool

func (api *Api) RoleCheck(required_role types.UserRole) Check {
	return func(w http.ResponseWriter, r *http.Request) bool {
		claims, err := api.GetClaims(r)
		if err == err_unauthorized {
			SendError(w, http.StatusUnauthorized, errors.New("unauthorized"))
			return false
		} else if err != nil {
			SendError(w, http.StatusInternalServerError, err)
			return false
		}
		if claims.Role != required_role {
			SendError(w, http.StatusForbidden, fmt.Errorf("forbidden, missing role: %s", required_role))
			return false
		}
		return true
	}
}

func (api *Api) LoginCheck() Check {
	return func(w http.ResponseWriter, r *http.Request) bool {
		claims, err := api.GetClaims(r)
		if err == err_unauthorized {
			SendError(w, http.StatusUnauthorized, errors.New("unauthorized"))
			return false
		} else if err != nil {
			SendError(w, http.StatusInternalServerError, err)
			return false
		}
		return claims != nil
	}
}

type SecureHandlerImpl struct {
	checks  []Check
	handler http.Handler
}

func (sh SecureHandlerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, check := range sh.checks {
		success := check(w, r)
		if !success {
			return
		}
	}
	sh.handler.ServeHTTP(w, r)
}

func SecureHandler(checks ...Check) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return SecureHandlerImpl{checks: checks, handler: h}
	}
}
