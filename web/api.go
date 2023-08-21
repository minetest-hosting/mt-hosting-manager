package web

import (
	"mt-hosting-manager/db"
	"mt-hosting-manager/public"
	"mt-hosting-manager/types"
	"mt-hosting-manager/web/middleware"
	"mt-hosting-manager/web/oauth"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

type Api struct {
	repos *db.Repositories
	cfg   *types.Config
}

func NewApi(repos *db.Repositories, cfg *types.Config) *Api {
	return &Api{
		repos: repos,
		cfg:   cfg,
	}
}

func (api *Api) Setup() {
	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.PrometheusMiddleware)

	// setup routes
	apir := r.PathPrefix("/api").Subrouter()
	apir.HandleFunc("/info", api.GetInfo)
	apir.HandleFunc("/login", api.Logout).Methods(http.MethodDelete)
	apir.HandleFunc("/login", api.GetLogin).Methods(http.MethodGet)

	if api.cfg.GithubOauthConfig.ClientID != "" {
		oauth_handler := &oauth.OauthHandler{
			Impl:     &oauth.GithubOauth{},
			UserRepo: api.repos.UserRepo,
			Config:   api.cfg.GithubOauthConfig,
			BaseURL:  api.cfg.BaseURL,
			Type:     types.UserTypeGithub,
			Callback: api.OauthCallback,
		}
		r.Handle("/oauth_callback/github", oauth_handler)
	}

	// static files
	if os.Getenv("WEBDEV") == "true" {
		logrus.WithFields(logrus.Fields{"dir": "public"}).Info("Using live mode")
		fs := http.FileServer(http.FS(os.DirFS("public")))
		r.PathPrefix("/").HandlerFunc(fs.ServeHTTP)

	} else {
		logrus.Info("Using embed mode")
		r.PathPrefix("/").Handler(statigz.FileServer(public.Webapp, brotli.AddEncoding))
	}

	http.Handle("/", r)
}
