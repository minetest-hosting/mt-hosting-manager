package web

import (
	"mt-hosting-manager/api/coinbase"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/core"
	"mt-hosting-manager/db"
	"mt-hosting-manager/public"
	"mt-hosting-manager/types"
	"mt-hosting-manager/web/middleware"
	"mt-hosting-manager/web/oauth"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

type Api struct {
	repos   *db.Repositories
	cfg     *types.Config
	core    *core.Core
	wc      *wallee.WalleeClient
	cbc     *coinbase.CoinbaseClient
	running *atomic.Bool
}

func NewApi(repos *db.Repositories, cfg *types.Config) *Api {
	return &Api{
		repos:   repos,
		cfg:     cfg,
		wc:      wallee.New(cfg.WalleeUserID, cfg.WalleeSpaceID, cfg.WalleeKey),
		running: &atomic.Bool{},
		core:    core.New(repos, cfg),
		cbc:     coinbase.New(cfg.CoinbaseKey),
	}
}

func (api *Api) Setup() {
	api.running.Store(true)

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.PrometheusMiddleware)

	// setup routes

	// public
	apir := r.PathPrefix("/api").Subrouter()
	apir.HandleFunc("/info", api.GetInfo).Methods(http.MethodGet)
	apir.HandleFunc("/healthcheck", api.Healthcheck).Methods(http.MethodGet)
	apir.HandleFunc("/login", api.Logout).Methods(http.MethodDelete)
	apir.HandleFunc("/login", api.GetLogin).Methods(http.MethodGet)
	apir.HandleFunc("/login", api.Login).Methods(http.MethodPost)
	apir.HandleFunc("/nodetype", api.GetNodeTypes).Methods(http.MethodGet)
	apir.HandleFunc("/nodetype/{id}", api.GetNodeType).Methods(http.MethodGet)
	apir.HandleFunc("/logstream/{id}", api.LogStream).Methods(http.MethodPost)
	apir.HandleFunc("/send_activation", api.SendActivationMail).Methods(http.MethodPost)
	apir.HandleFunc("/activate", api.ActivationCallback).Methods(http.MethodPost)
	apir.HandleFunc("/exchange_rate", api.GetExchangeRates)

	// user api
	user_api := apir.NewRoute().Subrouter()
	user_api.Use(SecureHandler(api.LoginCheck()))
	user_api.HandleFunc("/set_password", api.Secure(api.SetPassword)).Methods(http.MethodPost)
	user_api.HandleFunc("/profile", api.Secure(api.GetUserProfile)).Methods(http.MethodGet)
	user_api.HandleFunc("/profile", api.Secure(api.UpdateUserProfile)).Methods(http.MethodPost)
	user_api.HandleFunc("/audit_log", api.Secure(api.SearchAuditLog)).Methods(http.MethodPost)
	user_api.HandleFunc("/node", api.Secure(api.GetNodes)).Methods(http.MethodGet)
	user_api.HandleFunc("/node", api.Secure(api.CreateNode)).Methods(http.MethodPost)
	user_api.HandleFunc("/node/{id}", api.Secure(api.GetNode)).Methods(http.MethodGet)
	user_api.HandleFunc("/node/{id}/stats", api.Secure(api.GetNodeStats)).Methods(http.MethodGet)
	user_api.HandleFunc("/node/{id}/job", api.Secure(api.GetLatestNodeJob)).Methods(http.MethodGet)
	user_api.HandleFunc("/node/{id}/mtservers", api.Secure(api.GetNodeServers)).Methods(http.MethodGet)
	user_api.HandleFunc("/node/{id}", api.Secure(api.DeleteNode)).Methods(http.MethodDelete)
	user_api.HandleFunc("/node/{id}", api.Secure(api.UpdateNode)).Methods(http.MethodPost)
	user_api.HandleFunc("/mtserver", api.Secure(api.GetMTServers)).Methods(http.MethodGet)
	user_api.HandleFunc("/mtserver", api.Secure(api.CreateMTServer)).Methods(http.MethodPost)
	user_api.HandleFunc("/mtserver/validate", api.Secure(api.ValidateCreateMTServer)).Methods(http.MethodPost)
	user_api.HandleFunc("/mtserver/{id}", api.Secure(api.GetMTServer)).Methods(http.MethodGet)
	user_api.HandleFunc("/mtserver/{id}", api.Secure(api.UpdateMTServer)).Methods(http.MethodPost)
	user_api.HandleFunc("/mtserver/{id}/setup", api.Secure(api.SetupMTServer)).Methods(http.MethodPost)
	user_api.HandleFunc("/mtserver/{id}/job", api.Secure(api.GetLatestMTServerJob)).Methods(http.MethodGet)
	user_api.HandleFunc("/mtserver/{id}", api.Secure(api.DeleteMTServer)).Methods(http.MethodDelete)
	user_api.HandleFunc("/transaction", api.Secure(api.GetTransactions)).Methods(http.MethodGet)
	user_api.HandleFunc("/transaction/create", api.Secure(api.CreateTransaction)).Methods(http.MethodPost)
	user_api.HandleFunc("/transaction/search", api.Secure(api.SearchTransaction)).Methods(http.MethodPost)
	user_api.HandleFunc("/transaction/{id}", api.Secure(api.GetTransaction)).Methods(http.MethodGet)
	user_api.HandleFunc("/transaction/{id}/check", api.Secure(api.CheckTransaction)).Methods(http.MethodGet)
	user_api.HandleFunc("/transaction/{id}/refund", api.Secure(api.RefundTransaction)).Methods(http.MethodPost)
	user_api.HandleFunc("/backup", api.Secure(api.GetBackups)).Methods(http.MethodGet)

	// semi public, only with known identifiers (user_id and minetest_server_id)
	apir.HandleFunc("/backup/create", api.CreateBackup).Methods(http.MethodPost)
	apir.HandleFunc("/backup/{id}/complete", api.CompleteBackup).Methods(http.MethodPost)
	apir.HandleFunc("/backup/{id}/error", api.MarkBackupError).Methods(http.MethodPost)

	// admin api
	admin_api := apir.NewRoute().Subrouter()
	admin_api.Use(SecureHandler(api.RoleCheck(types.UserRoleAdmin)))
	admin_api.HandleFunc("/user", api.Secure(api.GetUsers)).Methods(http.MethodGet)
	admin_api.HandleFunc("/user/search", api.Secure(api.SearchUser)).Methods(http.MethodPost)
	admin_api.HandleFunc("/nodetype", api.Secure(api.CreateNodeType)).Methods(http.MethodPost)
	admin_api.HandleFunc("/nodetype/{id}", api.Secure(api.UpdateNodeType)).Methods(http.MethodPost)
	admin_api.HandleFunc("/nodetype/{id}", api.Secure(api.DeleteNodeType)).Methods(http.MethodDelete)
	admin_api.HandleFunc("/job", api.Secure(api.GetJobs)).Methods(http.MethodGet)
	admin_api.HandleFunc("/job/{id}", api.Secure(api.DeleteJob)).Methods(http.MethodDelete)
	admin_api.HandleFunc("/job/{id}", api.Secure(api.RetryJob)).Methods(http.MethodPost)
	admin_api.HandleFunc("/mail/send/{user_id}", api.Secure(api.SendMail)).Methods(http.MethodPost)

	// oauth
	if api.cfg.GithubOauthConfig.ClientID != "" {
		oauth_handler := &oauth.OauthHandler{
			Core:     api.core,
			Impl:     &oauth.GithubOauth{},
			UserRepo: api.repos.UserRepo,
			Config:   api.cfg.GithubOauthConfig,
			BaseURL:  api.cfg.BaseURL,
			Type:     types.UserTypeGithub,
			Callback: api.OauthCallback,
		}
		r.Handle("/oauth_callback/github", oauth_handler)
	}

	if api.cfg.DiscordOauthConfig.ClientID != "" {
		oauth_handler := &oauth.OauthHandler{
			Core:     api.core,
			Impl:     &oauth.DiscordOauth{},
			UserRepo: api.repos.UserRepo,
			Config:   api.cfg.DiscordOauthConfig,
			BaseURL:  api.cfg.BaseURL,
			Type:     types.UserTypeDiscord,
			Callback: api.OauthCallback,
		}
		r.Handle("/oauth_callback/discord", oauth_handler)
	}

	if api.cfg.MesehubOauthConfig.ClientID != "" {
		oauth_handler := &oauth.OauthHandler{
			Core:     api.core,
			Impl:     &oauth.MesehubOauth{},
			UserRepo: api.repos.UserRepo,
			Config:   api.cfg.MesehubOauthConfig,
			BaseURL:  api.cfg.BaseURL,
			Type:     types.UserTypeMesehub,
			Callback: api.OauthCallback,
		}
		r.Handle("/oauth_callback/mesehub", oauth_handler)
	}

	// static files
	if api.cfg.Webdev {
		logrus.WithFields(logrus.Fields{"dir": "public"}).Info("Using live mode")
		fs := http.FileServer(http.FS(os.DirFS("public")))
		r.PathPrefix("/").HandlerFunc(fs.ServeHTTP)

	} else {
		logrus.Info("Using embed mode")
		r.PathPrefix("/").Handler(statigz.FileServer(public.Webapp, brotli.AddEncoding))
	}

	http.Handle("/", r)
}

func (api *Api) Stop() {
	api.running.Store(false)
}
