package types

import (
	"github.com/golang-jwt/jwt/v4"
)

type JWTPermission string

const (
	JWTPermissionAdmin JWTPermission = "ADMIN"
)

type Claims struct {
	*jwt.RegisteredClaims
	Username    string          `json:"username"`
	UserID      int64           `json:"user_id"`
	Mail        string          `json:"mail"`
	Permissions []JWTPermission `json:"permissions"`
}

func (c *Claims) HasPermission(perm JWTPermission) bool {
	for _, p := range c.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}

func (c *Claims) IsAdmin() bool {
	return c.HasPermission(JWTPermissionAdmin)
}
