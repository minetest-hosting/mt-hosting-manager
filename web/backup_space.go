package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (a *Api) CreateBackupSpace(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	bs := &types.BackupSpace{}
	bs.UserID = c.UserID
	bs.Created = time.Now().Unix()
	bs.ValidUntil = time.Now().Add(time.Hour * 24).Unix()
	bs.ID = uuid.NewString()
	if bs.RetentionDays < 7 {
		bs.RetentionDays = 7
	}

	err := json.NewDecoder(r.Body).Decode(bs)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	err = a.repos.BackupSpaceRepo.Insert(bs)
	Send(w, bs, err)
}

func (a *Api) UpdateBackupSpace(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	updated_bs := &types.BackupSpace{}
	err := json.NewDecoder(r.Body).Decode(updated_bs)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	bs, err := a.repos.BackupSpaceRepo.GetByID(updated_bs.ID)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if bs == nil {
		SendError(w, 500, fmt.Errorf("not found"))
		return
	}
	if bs.UserID != c.UserID && c.Role != types.UserRoleAdmin {
		SendError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	bs.Name = updated_bs.Name
	bs.RetentionDays = updated_bs.RetentionDays
	if bs.RetentionDays < 7 {
		bs.RetentionDays = 7
	}

	err = a.repos.BackupSpaceRepo.Update(bs)
	Send(w, bs, err)
}

func (a *Api) GetBackupSpace(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	bs, err := a.repos.BackupSpaceRepo.GetByID(vars["id"])
	if err != nil && bs != nil && bs.UserID != c.UserID && c.Role != types.UserRoleAdmin {
		SendError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	Send(w, bs, err)
}

func (a *Api) GetBackupSpaces(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	list, err := a.repos.BackupSpaceRepo.GetByUserID(c.UserID)
	Send(w, list, err)
}

func (a *Api) RemoveBackupSpace(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	bs, err := a.repos.BackupSpaceRepo.GetByID(vars["id"])
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if bs != nil && bs.UserID != c.UserID && c.Role != types.UserRoleAdmin {
		SendError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	//TODO: remove all backups and data

	err = a.repos.BackupSpaceRepo.Delete(vars["id"])
	Send(w, true, err)
}
