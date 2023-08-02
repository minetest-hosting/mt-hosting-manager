package web

import (
	"mt-hosting-manager/types"
	"net/http"
)

type JobModel struct {
	Jobs []*types.Job
}

func (ctx *Context) Jobs(w http.ResponseWriter, r *http.Request, c *types.Claims) {

	if r.Method == http.MethodPost {
		id := r.FormValue("job_id")
		job, err := ctx.repos.JobRepo.GetByID(id)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}

		switch r.FormValue("action") {
		case "retry":
			job.State = types.JobStateCreated
			job.Message = ""
			err = ctx.repos.JobRepo.Update(job)

		case "delete":
			err = ctx.repos.JobRepo.Delete(id)
		}

		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}
	}

	list, err := ctx.repos.JobRepo.GetAll()
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	model := &JobModel{
		Jobs: list,
	}

	ctx.tu.ExecuteTemplate(w, r, "jobs.html", model)
}
