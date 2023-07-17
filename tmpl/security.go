package tmpl

import (
	"errors"
	"mt-hosting-manager/types"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type SecureHandlerFunc func(http.ResponseWriter, *http.Request, *types.Claims)
type ClaimsCheck func(*types.Claims) (bool, error)

func RoleCheck(req_role types.UserRole) ClaimsCheck {
	return func(c *types.Claims) (bool, error) {
		return c.Role == req_role, nil
	}
}

var err_unauthorized = errors.New("unauthorized")
var err_forbidden = errors.New("forbidden")

func (tu *TemplateUtil) OptionalSecure(h SecureHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := tu.GetClaims(r)
		if err != nil {
			tu.RenderError(w, r, 500, err)
			return
		}

		h(w, r, claims)
	}
}

func (tu *TemplateUtil) Secure(h SecureHandlerFunc, checks ...ClaimsCheck) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := tu.GetClaims(r)
		if err != nil {
			tu.RenderError(w, r, 500, err)
			return
		}
		if claims == nil {
			tu.RenderError(w, r, 401, err_unauthorized)
			return
		}

		for _, c := range checks {
			ok, err := c(claims)
			if err != nil {
				tu.RenderError(w, r, 500, err)
				return
			}
			if !ok {
				tu.RenderError(w, r, 403, err_forbidden)
				return
			}
		}

		h(w, r, claims)
	}
}

func (tu *TemplateUtil) CreateToken(c *types.Claims, d time.Duration) (string, error) {
	c.RegisteredClaims = &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(d)),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString([]byte(tu.JWTKey))
}

func (tu *TemplateUtil) ParseToken(token string) (*types.Claims, error) {
	t, err := jwt.ParseWithClaims(token, &types.Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(tu.JWTKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, err_unauthorized
	}

	claims, ok := t.Claims.(*types.Claims)
	if !ok {
		return nil, errors.New("internal error")
	}

	return claims, nil
}

func (tu *TemplateUtil) GetClaims(r *http.Request) (*types.Claims, error) {
	co, err := r.Cookie(tu.CookieName)
	if err == http.ErrNoCookie {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	c, err := tu.ParseToken(co.Value)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (tu *TemplateUtil) SetToken(w http.ResponseWriter, token string, d time.Duration) {
	c := &http.Cookie{
		Name:     tu.CookieName,
		Value:    token,
		Path:     tu.CookiePath,
		Expires:  time.Now().Add(d),
		Domain:   tu.CookieDomain,
		HttpOnly: true,
		Secure:   tu.CookieSecure,
	}
	http.SetCookie(w, c)
}

func (tu *TemplateUtil) ClearToken(w http.ResponseWriter) {
	tu.SetToken(w, "", time.Duration(0))
}
