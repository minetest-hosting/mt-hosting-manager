package tmpl

import (
	"bytes"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/sirupsen/logrus"
)

type TemplateUtil struct {
	Files        fs.FS
	AddFuncs     func(funcs template.FuncMap, r *http.Request)
	JWTKey       string
	CookieName   string
	CookieDomain string
	CookiePath   string
	CookieSecure bool
}

func (tu *TemplateUtil) CreateTemplate(pagename string, r *http.Request) (*template.Template, error) {
	funcs := template.FuncMap{
		"Claims": func() (any, error) { return tu.GetClaims(r) },
	}
	tu.AddFuncs(funcs, r)
	return template.New("").Funcs(funcs).ParseFS(tu.Files, "components/*.html", pagename)
}

func (tu *TemplateUtil) StaticPage(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := tu.CreateTemplate(name, r)
		if err != nil {
			tu.RenderError(w, r, 500, err)
			return
		}
		t.ExecuteTemplate(w, "layout", nil)
	}
}

func (tu *TemplateUtil) ExecuteTemplate(w http.ResponseWriter, r *http.Request, name string, data any) {
	t, err := tu.CreateTemplate(name, r)
	if err != nil {
		tu.RenderError(w, r, 500, err)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	err = t.ExecuteTemplate(buf, "layout", data)
	if err != nil {
		tu.RenderError(w, r, 500, err)
	} else {
		w.Write(buf.Bytes())
	}
}

func (tu *TemplateUtil) RenderError(w http.ResponseWriter, r *http.Request, code int, err error) {
	logrus.WithFields(logrus.Fields{
		"error": err,
		"code":  code,
	}).Error()
	w.WriteHeader(code)
	t, terr := tu.CreateTemplate("error.html", r)
	if terr != nil {
		panic(terr)
	}
	t.ExecuteTemplate(w, "error", map[string]any{
		"Error": err,
		"Code":  code,
	})
}
