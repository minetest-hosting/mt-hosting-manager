package usernode

import (
	"fmt"
	"math"
	"mt-hosting-manager/types"
	"net/http"

	"github.com/gorilla/mux"
)

type DetailModel struct {
	UserNode      *types.UserNode
	LatestJob     *types.Job
	Transactions  []*types.PaymentTransaction
	Servers       []*types.MinetestServer
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
	if node == nil {
		ctx.tu.RenderError(w, r, 404, fmt.Errorf("node not found: %s", id))
		return
	}

	tx_list, err := ctx.repos.PaymentTransactionRepo.GetByNodeID(node.ID)
	if node == nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("get transactions error: %s", err))
		return
	}

	servers, err := ctx.repos.MinetestServerRepo.GetByNodeID(node.ID)
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
		Servers:       servers,
		Transactions:  tx_list,
		DiskPercent:   int(float64(node.DiskUsed) / float64(node.DiskSize) * 100),
		DiskGBUsed:    float64(node.DiskUsed) / bytes_in_gb,
		DiskGBTotal:   float64(node.DiskSize) / bytes_in_gb,
		MemoryPercent: int(float64(node.MemoryUsed) / float64(node.MemorySize) * 100),
		MemoryGBUsed:  math.Max(0, float64(node.MemoryUsed)/bytes_in_gb),
		MemoryGBTotal: math.Max(0, float64(node.MemorySize)/bytes_in_gb),
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
