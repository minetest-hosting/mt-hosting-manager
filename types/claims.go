package types

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	*jwt.RegisteredClaims
	Mail string   `json:"mail"`
	Role UserRole `json:"role"`
}
