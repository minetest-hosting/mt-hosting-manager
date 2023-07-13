package web

import (
	"embed"
	"mt-hosting-manager/db"
	"mt-hosting-manager/tmpl"
	"mt-hosting-manager/types"
	"mt-hosting-manager/web/oauth"
	"os"

	"github.com/gorilla/mux"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

//go:embed *
var Files embed.FS

type Context struct {
	tu          *tmpl.TemplateUtil
	repos       *db.Repositories
	GithubOauth *oauth.OAuthConfig
}

func (ctx *Context) Setup(r *mux.Router) {
	r.HandleFunc("/", ctx.Index)
	r.HandleFunc("/login", ctx.tu.OptionalSecure(ctx.Login))

	r.PathPrefix("/assets").Handler(statigz.FileServer(Files, brotli.AddEncoding))

	if os.Getenv("GITHUB_CLIENTID") != "" {
		ctx.GithubOauth = &oauth.OAuthConfig{
			ClientID: os.Getenv("GITHUB_CLIENTID"),
			Secret:   os.Getenv("GITHUB_SECRET"),
		}

		r.Handle("/oauth_callback/github", &OauthHandler{
			Impl:     &oauth.GithubOauth{},
			BaseURL:  os.Getenv("BASEURL"),
			Type:     types.UserTypeGithub,
			Tu:       ctx.tu,
			UserRepo: ctx.repos.UserRepo,
			Config:   ctx.GithubOauth,
		})
	}

	r.NotFoundHandler = ctx.NotFound()
}
