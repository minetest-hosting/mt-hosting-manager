package web

import (
	"mt-hosting-manager/types"
	"net/http"
)

func (a *Api) Healthcheck(w http.ResponseWriter, r *http.Request) {
	_, err := a.repos.JobRepo.GetByState(types.JobStateCreated)
	Send(w, true, err)
}
