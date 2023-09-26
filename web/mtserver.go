package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (a *Api) GetMTServers(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	if c.Role == types.UserRoleAdmin {
		list, err := a.repos.MinetestServerRepo.GetAll()
		Send(w, list, err)
	} else {
		list, err := a.repos.MinetestServerRepo.GetByUserID(c.UserID)
		Send(w, list, err)
	}
}

func (a *Api) GetMTServer(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, server, _, err := a.CheckedGetMTServer(id, c)
	Send(w, server, err)
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
		JWTKey:     core.RandSeq(16),
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

	a.core.AddAuditLog(&types.AuditLog{
		Type:             types.AuditLogServerCreated,
		UserID:           c.UserID,
		UserNodeID:       &node.ID,
		MinetestServerID: &mtserver.ID,
	})

	Send(w, mtserver, nil)
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

	a.core.AddAuditLog(&types.AuditLog{
		Type:             types.AuditLogServerRemoved,
		UserID:           c.UserID,
		UserNodeID:       &node.ID,
		MinetestServerID: &mtserver.ID,
	})

	Send(w, true, err)
}

func (a *Api) UpdateMTServer(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, mtserver, status, err := a.CheckedGetMTServer(id, c)
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
	Send(w, mtserver, err)
}

func (a *Api) SetupMTServer(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	node, mtserver, status, err := a.CheckedGetMTServer(id, c)
	if err != nil {
		SendError(w, status, err)
		return
	}

	// check for current execution
	latest_job, err := a.repos.JobRepo.GetLatestByMinetestServerID(mtserver.ID)
	if err != nil {
		SendError(w, 500, fmt.Errorf("could not fetch latest job: %v", err))
		return
	}
	if latest_job != nil && (latest_job.State == types.JobStateCreated || latest_job.State == types.JobStateRunning) {
		// already running or created
		SendError(w, 500, fmt.Errorf("job already scheduled: %s", latest_job.ID))
		return
	}

	job := worker.SetupServerJob(node, mtserver)
	err = a.repos.JobRepo.Insert(job)

	Send(w, job, err)
}

func (a *Api) GetLatestMTServerJob(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, mtserver, status, err := a.CheckedGetMTServer(id, c)
	if err != nil {
		SendError(w, status, err)
		return
	}

	job, err := a.repos.JobRepo.GetLatestByMinetestServerID(mtserver.ID)
	Send(w, job, err)
}
