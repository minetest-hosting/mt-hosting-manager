package web

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (a *Api) GetNodes(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	if c.Role == types.UserRoleAdmin {
		list, err := a.repos.UserNodeRepo.GetAll()
		Send(w, list, err)
	} else {
		list, err := a.repos.UserNodeRepo.GetByUserID(c.UserID)
		Send(w, list, err)
	}
}

func (a *Api) GetNode(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	node, _, err := a.CheckedGetUserNode(id, c)
	Send(w, node, err)
}

func (a *Api) GetNodeServers(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	node, _, err := a.CheckedGetUserNode(id, c)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	servers, err := a.repos.MinetestServerRepo.GetByNodeID(node.ID)
	Send(w, servers, err)
}

func (a *Api) GetLatestNodeJob(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	node, _, err := a.CheckedGetUserNode(id, c)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	job, err := a.repos.JobRepo.GetLatestByUserNodeID(node.ID)
	Send(w, job, err)
}

func (a *Api) GetNodeStats(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	node, _, err := a.CheckedGetUserNode(id, c)
	metrics := &core.NodeExporterMetrics{}

	if node.State != types.UserNodeStateRunning {
		SendError(w, 500, fmt.Errorf("node not in running state"))
		return
	}

	if a.cfg.EnableDummyWorker {
		metrics.DiskSize = 1000 * 1000 * 1000 * 10
		metrics.DiskUsed = 1000 * 1000 * 1000 * 2.5
		metrics.MemorySize = 1000 * 1000 * 1000 * 2
		metrics.MemoryUsed = 1000 * 1000 * 1000 * 0.2
		metrics.LoadPercent = rand.Intn(20)

	} else {
		metrics, err = core.FetchMetrics(node.IPv4)
	}

	Send(w, metrics, err)
}

func (a *Api) UpdateNode(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	node, status, err := a.CheckedGetUserNode(id, c)
	if err != nil {
		SendError(w, status, err)
		return
	}

	updated_node := &types.UserNode{}
	err = json.NewDecoder(r.Body).Decode(updated_node)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	// update allowed fields
	node.Alias = updated_node.Alias

	err = a.repos.UserNodeRepo.Update(node)
	Send(w, node, err)
}

func (a *Api) DeleteNode(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	node, status, err := a.CheckedGetUserNode(id, c)
	if err != nil {
		SendError(w, status, err)
		return
	}

	serverlist, err := a.repos.MinetestServerRepo.GetByNodeID(node.ID)
	if err != nil {
		SendError(w, 500, fmt.Errorf("server fetch error: %v", err))
		return
	}
	if len(serverlist) > 0 {
		SendError(w, 500, fmt.Errorf("node still contains servers"))
		return
	}

	job := worker.RemoveNodeJob(node)
	err = a.repos.JobRepo.Insert(job)
	if err != nil {
		SendError(w, 500, fmt.Errorf("job insert error: %v", err))
		return
	}

	a.core.AddAuditLog(&types.AuditLog{
		Type:       types.AuditLogNodeRemoved,
		UserID:     c.UserID,
		UserNodeID: &node.ID,
	})

	Send(w, true, nil)
}

func (a *Api) CreateNode(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	create_node := &types.UserNode{}
	err := json.NewDecoder(r.Body).Decode(create_node)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	nt, err := a.repos.NodeTypeRepo.GetByID(create_node.NodeTypeID)
	if err != nil {
		SendError(w, 500, fmt.Errorf("nodetype fetch error: %v", err))
		return
	}
	if nt == nil {
		SendError(w, 404, fmt.Errorf("nodetype not found: %s", create_node.NodeTypeID))
		return
	}
	if nt.State != types.NodeTypeStateActive {
		SendError(w, 405, fmt.Errorf("nodetype in wrong state: %s", nt.State))
		return
	}

	user, err := a.repos.UserRepo.GetByID(c.UserID)
	if err != nil {
		SendError(w, 500, fmt.Errorf("user fetch error: %v", err))
		return
	}
	if user == nil {
		SendError(w, 404, fmt.Errorf("user not found: %s", c.UserID))
		return
	}

	if user.Balance < nt.DailyCost {
		SendError(w, 405, fmt.Errorf("remaining balance is less than the daily cost of the node-type"))
		return
	}

	randstr := core.RandStringRunes(7)

	node := &types.UserNode{
		ID:                uuid.NewString(),
		UserID:            c.UserID,
		NodeTypeID:        create_node.NodeTypeID,
		Created:           time.Now().Unix(),
		LastCollectedTime: time.Now().Unix(),
		State:             types.UserNodeStateCreated,
		Name:              fmt.Sprintf("node-%s-%s", a.cfg.Stage, randstr),
		Alias:             create_node.Alias,
	}
	err = a.repos.UserNodeRepo.Insert(node)
	if err != nil {
		SendError(w, 500, fmt.Errorf("node insert error: %v", err))
		return
	}

	job := worker.SetupNodeJob(node)
	err = a.repos.JobRepo.Insert(job)
	if err != nil {
		SendError(w, 500, fmt.Errorf("job insert error: %v", err))
		return
	}

	err = a.core.SubtractBalance(c.UserID, nt.DailyCost)
	if err != nil {
		SendError(w, 500, fmt.Errorf("balance update error: %v", err))
		return
	}

	a.core.AddAuditLog(&types.AuditLog{
		Type:       types.AuditLogNodeCreated,
		UserID:     c.UserID,
		UserNodeID: &node.ID,
	})

	Send(w, node, nil)
}
