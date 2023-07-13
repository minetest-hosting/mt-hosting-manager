package web

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"time"

	"mt-hosting-manager/tmpl"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

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
	t := time.UnixMilli(ts)
	return t.Format(time.UnixDate)
}

func Serve() error {

	r := mux.NewRouter()
	r.Use(prometheusMiddleware)
	r.Use(loggingMiddleware)

	tmplRoute := r.NewRoute().Subrouter()
	tmplRoute.Use(csrf.Protect([]byte(os.Getenv("CSRF_KEY"))))

	var files fs.FS
	if os.Getenv("WEBDEV") == "true" {
		logrus.Warn("Webdev mode enabled")
		files = os.DirFS("web")
	} else {
		files = Files
	}

	tu := &tmpl.TemplateUtil{
		Files: files,
		AddFuncs: func(funcs template.FuncMap, r *http.Request) {
			funcs["BaseURL"] = func() string { return os.Getenv("BASEURL") }
			funcs["prettysize"] = prettysize
			funcs["formattime"] = formattime
			funcs["CSRFField"] = func() template.HTML { return csrf.TemplateField(r) }
		},
		JWTKey:       os.Getenv("JWT_KEY"),
		CookieName:   "mt-hosting-manager",
		CookieDomain: os.Getenv("COOKIE_DOMAIN"),
		CookiePath:   os.Getenv("COOKIE_PATH"),
		CookieSecure: os.Getenv("COOKIE_SECURE") == "true",
	}

	// templates, pages
	ctx := &Context{
		tu: tu,
	}
	ctx.Setup(tmplRoute)

	// main entry
	http.Handle("/", r)

	// metrics
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(":8080", nil)
}
