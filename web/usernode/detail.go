package usernode

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"

	"github.com/gorilla/mux"
)

type DetailModel struct {
	UserNode      *types.UserNode
	LatestJob     *types.Job
	DiskPercent   int
	DiskGBUsed    float64
	DiskGBTotal   float64
	DiskWarn      bool
	DiskDanger    bool
	MemoryPercent int
	MemoryGBUsed  float64
	MemoryGBTotal float64
	MemoryWarn    bool
	MemoryDanger  bool
	AliasUpdated  bool
}

// view details
func (ctx *Context) Detail(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	node, err := ctx.repos.UserNodeRepo.GetByID(id)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	job, err := ctx.repos.JobRepo.GetLatestByUserNodeID(node.ID)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	bytes_in_gb := 1024.0 * 1024.0 * 1024.0
	m := &DetailModel{
		UserNode:      node,
		LatestJob:     job,
		DiskPercent:   int(float64(node.DiskUsed) / float64(node.DiskSize) * 100),
		DiskGBUsed:    float64(node.DiskUsed) / bytes_in_gb,
		DiskGBTotal:   float64(node.DiskSize) / bytes_in_gb,
		MemoryPercent: int(float64(node.MemoryUsed) / float64(node.MemorySize) * 100),
		MemoryGBUsed:  float64(node.MemoryUsed) / bytes_in_gb,
		MemoryGBTotal: float64(node.MemorySize) / bytes_in_gb,
	}

	if m.DiskPercent > 90 {
		m.DiskDanger = true
	} else if m.DiskPercent > 75 {
		m.DiskWarn = true
	}

	if m.MemoryPercent > 90 {
		m.MemoryDanger = true
	} else if m.MemoryPercent > 75 {
		m.MemoryWarn = true
	}

	if r.Method == http.MethodPost {

		switch r.FormValue("action") {
		case "set-alias":
			node.Alias = r.FormValue("alias")
			m.AliasUpdated = true
		case "start":
			err = ctx.hcc.PowerOnServer(node.ExternalID)
			node.State = types.UserNodeStateRunning
		case "stop":
			err = ctx.hcc.PowerOffServer(node.ExternalID)
			node.State = types.UserNodeStateStopped
		}

		if err != nil {
			ctx.tu.RenderError(w, r, 500, fmt.Errorf("start/stop error: %v", err))
			return
		}
		err = ctx.repos.UserNodeRepo.Update(node)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, fmt.Errorf("node-update error: %v", err))
			return
		}
	}

	ctx.tu.ExecuteTemplate(w, r, "usernode/detail.html", m)
}
