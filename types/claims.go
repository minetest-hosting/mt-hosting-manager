package types

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	*jwt.RegisteredClaims
	UserID string   `json:"user_id"`
	Name   string   `json:"name"`
	Role   UserRole `json:"role"`
}
