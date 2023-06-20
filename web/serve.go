package web

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"mt-hosting-manager/tmpl"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"

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

func baseurl() string {
	u := os.Getenv("BASEURL")
	if u == "" {
		u = "http://127.0.0.1:8080"
	}
	return u
}

func Serve() error {

	key := "mykey"

	r := mux.NewRouter()
	r.Use(prometheusMiddleware)
	r.Use(loggingMiddleware)

	tmplRoute := r.NewRoute().Subrouter()
	tmplRoute.Use(csrf.Protect([]byte(key)))

	tu := &tmpl.TemplateUtil{
		Files: Files,
		AddFuncs: func(funcs template.FuncMap, r *http.Request) {
			funcs["BaseURL"] = baseurl
			funcs["prettysize"] = prettysize
			funcs["formattime"] = formattime
			funcs["CSRFField"] = func() template.HTML { return csrf.TemplateField(r) }
		},
		JWTKey:       key,
		CookieName:   "mt-hosting-manager",
		CookieDomain: "127.0.0.1",
		CookiePath:   "/",
		CookieSecure: false,
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
