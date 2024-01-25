package web

import (
	"bytes"
	"fmt"
	"io"
	"mt-hosting-manager/types"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (a *Api) CompleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	job_id := vars["job_id"]
	job, err := a.repos.JobRepo.GetByID(job_id)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if job == nil {
		SendError(w, 404, fmt.Errorf("job %s not found", job_id))
		return
	}
	if job.State != types.JobStateRunning {
		SendError(w, http.StatusConflict, fmt.Errorf("job state invalid: '%s'", job.State))
		return
	}
	job.State = types.JobStateDoneSuccess
	job.Finished = time.Now().Unix()
	err = a.repos.JobRepo.Update(job)
	Send(w, true, err)
}

func (a *Api) MarkJobFailure(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	job_id := vars["job_id"]
	job, err := a.repos.JobRepo.GetByID(job_id)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if job == nil {
		SendError(w, 404, fmt.Errorf("job %s not found", job_id))
		return
	}
	if job.State != types.JobStateRunning {
		SendError(w, http.StatusConflict, fmt.Errorf("job state invalid: '%s'", job.State))
		return
	}
	buf := bytes.NewBuffer([]byte{})
	size, err := io.Copy(buf, r.Body)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if size > 0 {
		// set error message
		job.Message = buf.String()
	}

	job.State = types.JobStateDoneFailure
	job.Finished = time.Now().Unix()
	err = a.repos.JobRepo.Update(job)
	Send(w, true, err)
}

func (a *Api) GetJobs(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	list, err := a.repos.JobRepo.GetAll()
	Send(w, list, err)
}

func (a *Api) RetryJob(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	job, err := a.repos.JobRepo.GetByID(id)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if job == nil {
		SendError(w, 404, fmt.Errorf("job %s not found", id))
		return
	}

	job.State = types.JobStateCreated
	err = a.repos.JobRepo.Update(job)
	Send(w, job, err)
}

func (a *Api) DeleteJob(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	err := a.repos.JobRepo.Delete(vars["id"])
	Send(w, true, err)
}
