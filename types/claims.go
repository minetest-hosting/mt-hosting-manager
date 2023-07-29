package types

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	*jwt.RegisteredClaims
	UserID string   `json:"user_id"`
	Mail   string   `json:"mail"`
	Role   UserRole `json:"role"`
}
