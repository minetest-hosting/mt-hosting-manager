package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (a *Api) CreateBackup(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	b := &types.Backup{}
	err := json.NewDecoder(r.Body).Decode(b)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	bs, err := a.repos.BackupSpaceRepo.GetByID(b.BackupSpaceID)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if bs == nil {
		SendError(w, 404, fmt.Errorf("backupspace not found: '%s'", b.MinetestServerID))
		return
	}
	if bs.UserID != c.UserID && c.Role != types.UserRoleAdmin {
		// not the owner and not admin
		SendError(w, 401, fmt.Errorf("user-id mismatch"))
		return
	}

	// validate given id's
	mtserver, err := a.repos.MinetestServerRepo.GetByID(b.MinetestServerID)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if mtserver == nil {
		SendError(w, 404, fmt.Errorf("server not found: '%s'", b.MinetestServerID))
		return
	}
	node, err := a.repos.UserNodeRepo.GetByID(mtserver.UserNodeID)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if node == nil {
		SendError(w, 404, fmt.Errorf("node not found: '%s'", mtserver.UserNodeID))
		return
	}
	if node.UserID != c.UserID && c.Role != types.UserRoleAdmin {
		SendError(w, 405, fmt.Errorf("invalid data"))
		return
	}

	b.State = types.BackupStateCreated
	b.Passphrase = core.RandStringRunes(64)
	b.ID = uuid.NewString()
	b.Size = 0
	b.Created = time.Now().Unix()

	err = a.repos.BackupRepo.Insert(b)
	Send(w, b, err)
}

// created -> progress
func (a *Api) MarkBackupProgress(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	backup, err := a.repos.BackupRepo.GetByID(id)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if backup == nil {
		SendError(w, 404, fmt.Errorf("backup not found: '%s'", id))
		return
	}
	if backup.State != types.BackupStateCreated {
		SendError(w, http.StatusConflict, fmt.Errorf("state invalid"))
		return
	}

	backup.State = types.BackupStateProgress
	err = a.repos.BackupRepo.Update(backup)
	Send(w, backup, err)
}

// progress -> complete
func (a *Api) CompleteBackup(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	backup, err := a.repos.BackupRepo.GetByID(id)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if backup == nil {
		SendError(w, 404, fmt.Errorf("backup not found: '%s'", id))
		return
	}
	if backup.State != types.BackupStateProgress {
		SendError(w, http.StatusConflict, fmt.Errorf("state invalid"))
		return
	}

	size, err := a.core.GetBackupSize(backup)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	backup.State = types.BackupStateComplete
	backup.Size = size
	err = a.repos.BackupRepo.Update(backup)

	Send(w, backup, err)
}

// progress -> error
func (a *Api) MarkBackupError(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	backup, err := a.repos.BackupRepo.GetByID(id)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if backup == nil {
		SendError(w, 404, fmt.Errorf("backup not found: '%s'", id))
		return
	}
	if backup.State != types.BackupStateProgress {
		SendError(w, http.StatusConflict, fmt.Errorf("state invalid"))
		return
	}

	backup.State = types.BackupStateError
	//TODO: notify someone?
	err = a.repos.BackupRepo.Update(backup)
	Send(w, backup, err)
}

func (a *Api) GetBackups(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	bs, err := a.repos.BackupSpaceRepo.GetByID(vars["id"])
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if bs == nil {
		SendError(w, 404, fmt.Errorf("backup_space not found"))
		return
	}
	if bs.UserID != c.UserID && c.Role != types.UserRoleAdmin {
		SendError(w, 401, fmt.Errorf("unauthorized"))
		return
	}

	list, err := a.repos.BackupRepo.GetByBackupSpaceID(bs.ID)
	Send(w, list, err)
}

func (a *Api) RemoveBackup(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	b, err := a.repos.BackupRepo.GetByID(vars["id"])
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if b == nil {
		SendError(w, 404, fmt.Errorf("backup not found"))
		return
	}
	bs, err := a.repos.BackupSpaceRepo.GetByID(b.BackupSpaceID)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if bs == nil {
		SendError(w, 404, fmt.Errorf("backup_space not found"))
		return
	}
	if bs != nil && bs.UserID != c.UserID && c.Role != types.UserRoleAdmin {
		SendError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	err = a.core.RemoveBackup(b)
	Send(w, true, err)
}
