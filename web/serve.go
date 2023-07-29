package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"time"

	"mt-hosting-manager/db"
	"mt-hosting-manager/tmpl"
	"mt-hosting-manager/web/middleware"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
	"golang.org/x/text/language"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func prettysize(num int) string {
	if num > (1000 * 1000) {
		return fmt.Sprintf("%d MB", num/(1000*1000))
	} else if num > 1000 {
		return fmt.Sprintf("%d kB", num/(1000))
	} else {
		return fmt.Sprintf("%d bytes", num)
	}
}

func formattime(ts int64) string {
	t := time.Unix(ts, 0)
	return t.Format(time.RFC3339)
}

func Serve(repos *db.Repositories) error {

	r := mux.NewRouter()
	// static assets
	r.PathPrefix("/assets").Handler(statigz.FileServer(Files, brotli.AddEncoding))

	tmplRoute := r.NewRoute().Subrouter()
	tmplRoute.Use(csrf.Protect([]byte(os.Getenv("CSRF_KEY"))))
	tmplRoute.Use(middleware.PrometheusMiddleware)
	tmplRoute.Use(middleware.LoggingMiddleware)

	var files fs.FS
	if os.Getenv("WEBDEV") == "true" {
		logrus.Warn("Webdev mode enabled")
		files = os.DirFS("web")
	} else {
		files = Files
	}

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.LoadMessageFileFS(files, "locale/en.json")
	bundle.LoadMessageFileFS(files, "locale/de.json")

	tu := &tmpl.TemplateUtil{
		Files: files,
		AddFuncs: func(funcs template.FuncMap, r *http.Request) {
			funcs["prettysize"] = prettysize
			funcs["formattime"] = formattime
			funcs["CSRFField"] = func() template.HTML { return csrf.TemplateField(r) }
			funcs["T"] = func(msgId string) (string, error) {
				localizer := i18n.NewLocalizer(bundle, r.Header.Get("Accept-Language"))
				return localizer.Localize(&i18n.LocalizeConfig{MessageID: msgId})
			}
		},
		JWTKey:       os.Getenv("JWT_KEY"),
		BaseURL:      os.Getenv("BASEURL"),
		CookieName:   "mt-hosting-manager",
		CookieDomain: os.Getenv("COOKIE_DOMAIN"),
		CookiePath:   os.Getenv("COOKIE_PATH"),
		CookieSecure: os.Getenv("COOKIE_SECURE") == "true",
	}

	// templates, pages
	ctx := &Context{
		tu:    tu,
		repos: repos,
	}
	ctx.Setup(tmplRoute)

	// main entry
	http.Handle("/", r)

	// metrics
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(":8080", nil)
}
