package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (a *Api) GetNodes(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	list, err := a.repos.UserNodeRepo.GetByUserID(c.UserID)
	Send(w, list, err)
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

	randstr := core.RandStringRunes(7)

	node := &types.UserNode{
		ID:         uuid.NewString(),
		UserID:     c.UserID,
		NodeTypeID: create_node.NodeTypeID,
		Created:    time.Now().Unix(),
		State:      types.UserNodeStateCreated,
		Name:       fmt.Sprintf("node-%s-%s", os.Getenv("STAGE"), randstr),
		Alias:      create_node.Alias,
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

	Send(w, node, nil)
}
