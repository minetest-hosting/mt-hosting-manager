package web

import (
	"net/http"

	"github.com/dchest/captcha"
)

func (a *Api) CreateCaptcha(w http.ResponseWriter, r *http.Request) {
	c := captcha.New()
	w.Write([]byte(c))
}
