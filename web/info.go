package web

import (
	"mt-hosting-manager/types"
	"net/http"
)

func (a *Api) GetInfo(w http.ResponseWriter, r *http.Request) {
	SendJson(w, types.NewInfo(a.cfg))
}
