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

	// public
	apir := r.PathPrefix("/api").Subrouter()
	apir.HandleFunc("/info", api.GetInfo)
	apir.HandleFunc("/login", api.Logout).Methods(http.MethodDelete)
	apir.HandleFunc("/login", api.GetLogin).Methods(http.MethodGet)

	// user api
	user_api := apir.NewRoute().Subrouter()
	user_api.Use(SecureHandler(api.LoginCheck()))
	user_api.HandleFunc("/node", api.Secure(api.GetNodes)).Methods(http.MethodGet)
	user_api.HandleFunc("/node", api.Secure(api.CreateNode)).Methods(http.MethodPost)
	user_api.HandleFunc("/node/{id}", api.Secure(api.DeleteNode)).Methods(http.MethodDelete)
	user_api.HandleFunc("/node/{id}", api.Secure(api.UpdateNode)).Methods(http.MethodPost)
	user_api.HandleFunc("/transaction/create", api.Secure(api.CreateTransaction)).Methods(http.MethodPost)
	user_api.HandleFunc("/transaction/callback", api.Secure(api.TransactionCallback)).Methods(http.MethodGet)

	// admin api
	admin_api := apir.NewRoute().Subrouter()
	admin_api.Use(SecureHandler(api.RoleCheck(types.UserRoleAdmin)))
	admin_api.HandleFunc("/nodetype", api.Secure(api.GetNodeTypes)).Methods(http.MethodGet)
	admin_api.HandleFunc("/nodetype", api.Secure(api.CreateNodeType)).Methods(http.MethodPost)
	admin_api.HandleFunc("/nodetype/{id}", api.Secure(api.UpdateNodeType)).Methods(http.MethodPost)
	admin_api.HandleFunc("/nodetype/{id}", api.Secure(api.DeleteNodeType)).Methods(http.MethodDelete)
	admin_api.HandleFunc("/job", api.Secure(api.GetJobs)).Methods(http.MethodGet)
	admin_api.HandleFunc("/job/{id}", api.Secure(api.DeleteJob)).Methods(http.MethodDelete)
	admin_api.HandleFunc("/job/{id}", api.Secure(api.RetryJob)).Methods(http.MethodPost)

	// oauth
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
