package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (a *Api) GetMTServers(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	list, err := a.repos.MinetestServerRepo.GetByUserID(c.UserID)
	Send(w, list, err)
}

func (a *Api) CreateMTServer(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	create_mtserver := &types.MinetestServer{}
	err := json.NewDecoder(r.Body).Decode(create_mtserver)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	node, status, err := a.CheckedGetUserNode(create_mtserver.UserNodeID, c)
	if err != nil {
		SendError(w, status, err)
		return
	}

	other_servers, err := a.repos.MinetestServerRepo.GetByNodeID(node.ID)
	if err != nil {
		SendError(w, 500, fmt.Errorf("servers fetch error: %v", err))
		return
	}

	if create_mtserver.Port < 1000 || create_mtserver.Port > 65000 {
		SendError(w, 500, fmt.Errorf("invalid port: %d", create_mtserver.Port))
		return
	}

	for _, s := range other_servers {
		if s.Port == create_mtserver.Port {
			SendError(w, 500, fmt.Errorf("port already in use by: %s", s.ID))
			return
		}
	}

	mtserver := &types.MinetestServer{
		ID:         uuid.NewString(),
		UserNodeID: node.ID,
		Name:       create_mtserver.Name,
		DNSName:    create_mtserver.DNSName,
		Port:       create_mtserver.Port,
		Created:    time.Now().Unix(),
		UIVersion:  "latest",
		State:      types.MinetestServerStateCreated,
	}
	err = a.repos.MinetestServerRepo.Insert(mtserver)
	if err != nil {
		SendError(w, 500, fmt.Errorf("server insert error: %v", err))
		return
	}

	job := worker.SetupServerJob(node, mtserver)
	err = a.repos.JobRepo.Insert(job)
	if err != nil {
		SendError(w, 500, fmt.Errorf("job insert error: %v", err))
		return
	}
}

func (a *Api) DeleteMTServer(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	node, mtserver, status, err := a.CheckedGetMTServer(id, c)
	if err != nil {
		SendError(w, status, err)
		return
	}

	job := worker.RemoveServerJob(node, mtserver)
	err = a.repos.JobRepo.Insert(job)
	Send(w, true, err)
}

func (a *Api) UpdateMTServer(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	node, mtserver, status, err := a.CheckedGetMTServer(id, c)
	if err != nil {
		SendError(w, status, err)
		return
	}

	updated_mtserver := &types.MinetestServer{}
	err = json.NewDecoder(r.Body).Decode(updated_mtserver)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	// update allowed fields
	mtserver.Name = updated_mtserver.Name

	err = a.repos.MinetestServerRepo.Update(mtserver)
	Send(w, node, err)
}