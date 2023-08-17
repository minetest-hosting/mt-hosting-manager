package mtserver

import (
	"fmt"
	"mt-hosting-manager/types"
	"mt-hosting-manager/web/components"
	"mt-hosting-manager/worker"
	"net/http"

	"github.com/gorilla/mux"
)

type ServerDetailModel struct {
	Server      *types.MinetestServer
	Node        *types.UserNode
	Job         *types.Job
	Breadcrumb  *components.Breadcrumb
	IsDeploying bool
	DeployError bool
}

func (ctx *Context) Detail(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	server, err := ctx.repos.MinetestServerRepo.GetByID(id)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("get server error: %v", err))
		return
	}
	if server == nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("server not found: %s", id))
		return
	}

	node, err := ctx.repos.UserNodeRepo.GetByID(server.UserNodeID)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("get node error: %v", err))
		return
	}

	var job *types.Job

	if r.Method == http.MethodPost {
		switch r.FormValue("action") {
		case "update-deployment":
			job = worker.SetupServerJob(node, server)
			err = ctx.repos.JobRepo.Insert(job)
			if err != nil {
				ctx.tu.RenderError(w, r, 500, fmt.Errorf("could not schedule job: %s", err))
				return
			}
		}
	}

	if job == nil {
		// fetch job
		job, err = ctx.repos.JobRepo.GetLatestByMinetestServerID(server.ID)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, fmt.Errorf("get job error: %v", err))
			return
		}
	}

	m := &ServerDetailModel{
		Server: server,
		Node:   node,
		Breadcrumb: &components.Breadcrumb{
			Entries: []*components.BreadcrumbEntry{
				{
					Name: "Home",
					Link: "/",
				}, {
					Name: "Nodes",
					Link: "/nodes",
				}, {
					Name: fmt.Sprintf("Node '%s'", node.Alias),
					Link: fmt.Sprintf("/nodes/%s", node.ID),
				}, {
					Name:   fmt.Sprintf("Server '%s'", server.DNSName),
					Link:   fmt.Sprintf("/mtserver/%s", server.ID),
					Active: true,
				},
			},
		},
		Job:         job,
		IsDeploying: job != nil && job.Type == types.JobTypeServerSetup && (job.State == types.JobStateRunning || job.State == types.JobStateCreated),
		DeployError: job != nil && job.Type == types.JobTypeServerSetup && job.State == types.JobStateDoneFailure,
	}

	ctx.tu.ExecuteTemplate(w, r, "mtserver/detail.html", m)
}
