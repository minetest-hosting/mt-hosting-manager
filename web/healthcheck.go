package web

import (
	"mt-hosting-manager/types"
	"net/http"
)

func (a *Api) Healthcheck(w http.ResponseWriter, r *http.Request) {
	if !a.running.Load() {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("shutting down"))
		return
	}

	_, err := a.repos.JobRepo.GetByState(types.JobStateCreated)
	Send(w, true, err)
}
