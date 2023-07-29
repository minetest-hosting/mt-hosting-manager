package web

import (
	"fmt"
	"mt-hosting-manager/types"
	"mt-hosting-manager/web/oauth"
	"net/http"
)

type LoginModel struct {
	GithubOauth *oauth.OAuthConfig
}

func (ctx *Context) Login(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	lm := &LoginModel{
		GithubOauth: ctx.GithubOauth,
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}

		switch r.Form.Get("action") {
		case "logout":
			ctx.tu.ClearToken(w)
			http.Redirect(w, r, fmt.Sprintf("%s/login", ctx.tu.BaseURL), http.StatusSeeOther)
			return
		}
	}

	ctx.tu.ExecuteTemplate(w, r, "login.html", lm)
}
