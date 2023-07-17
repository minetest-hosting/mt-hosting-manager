package web

import (
	"mt-hosting-manager/types"
	"mt-hosting-manager/web/oauth"
	"net/http"
)

type LoginModel struct {
	GithubOauth *oauth.OAuthConfig
}

func (ctx *Context) LoginGet(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	lm := &LoginModel{
		GithubOauth: ctx.GithubOauth,
	}

	ctx.tu.ExecuteTemplate(w, r, "login.html", lm)
}

func (ctx *Context) LoginPost(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	err := r.ParseForm()
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	switch r.Form.Get("action") {
	case "logout":
		ctx.tu.ClearToken(w)
	}

	http.Redirect(w, r, "login", http.StatusSeeOther)

}
