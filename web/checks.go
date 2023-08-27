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

func (a *Api) CheckedGetUserNode(user_node_id string, c *types.Claims) (*types.UserNode, int, error) {
	node, err := a.repos.UserNodeRepo.GetByID(user_node_id)
	if err != nil {
		return nil, 500, fmt.Errorf("fetch node error: %v", err)
	}
	if node == nil {
		return nil, 404, fmt.Errorf("node not found '%s'", user_node_id)
	}
	if c.Role != types.UserRoleAdmin && node.UserID != c.UserID {
		return nil, 403, fmt.Errorf("not authorized to access node '%s'", user_node_id)
	}
	return node, 0, nil
}

func (a *Api) CheckedGetMTServer(mt_server_id string, c *types.Claims) (*types.UserNode, *types.MinetestServer, int, error) {
	mtserver, err := a.repos.MinetestServerRepo.GetByID(mt_server_id)
	if err != nil {
		return nil, nil, 500, fmt.Errorf("fetch mtserver error: %v", err)
	}
	if mtserver == nil {
		return nil, nil, 404, fmt.Errorf("mtserver not found '%s'", mt_server_id)
	}

	node, status, err := a.CheckedGetUserNode(mtserver.UserNodeID, c)
	return node, mtserver, status, err
}
