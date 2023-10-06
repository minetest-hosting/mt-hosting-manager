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

func (a *Api) GetBackups(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	list, err := a.repos.BackupRepo.GetByUserID(c.UserID)
	Send(w, list, err)
}

func (a *Api) CreateBackup(w http.ResponseWriter, r *http.Request) {
	b := &types.Backup{}
	err := json.NewDecoder(r.Body).Decode(b)
	if err != nil {
		SendError(w, 500, err)
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
	if node.UserID != b.UserID {
		SendError(w, 405, fmt.Errorf("invalid data"))
		return
	}

	b.State = types.BackupStateCreated
	b.ID = uuid.NewString()
	b.Size = 0
	b.Created = time.Now().Unix()

	err = a.repos.BackupRepo.Insert(b)
	Send(w, b, err)
}

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

	backup.State = types.BackupStateComplete
	//TODO: get final size
	err = a.repos.BackupRepo.Update(backup)
	Send(w, backup, err)
}

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

	backup.State = types.BackupStateError
	//TODO: notify someone?
	err = a.repos.BackupRepo.Update(backup)
	Send(w, backup, err)
}
