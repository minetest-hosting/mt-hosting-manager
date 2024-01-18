package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) ResolveGeoIP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	result := a.core.GeoIP.Resolve(vars["ip"])
	if result == nil {
		SendError(w, 404, fmt.Errorf("could not resolve ip"))
		return
	}
	Send(w, result, nil)
}
