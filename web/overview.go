package web

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
	"slices"
	"strings"

	"github.com/gorilla/mux"
)

func (a *Api) GetOverviewData(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	user_id := vars["user_id"]

	if user_id != c.UserID && c.Role != types.UserRoleAdmin {
		// not an admin, can't look at other user's data
		SendError(w, 403, fmt.Errorf("unauthorized"))
		return
	}

	nodes, err := a.repos.UserNodeRepo.GetByUserID(user_id)
	if err != nil {
		SendError(w, 500, fmt.Errorf("node fetch error: %v", err))
		return
	}

	data := []*types.OverviewData{}
	for _, node := range nodes {
		if node.State == types.UserNodeStateDecommissioned {
			// skip decommissioned nodes
			continue
		}

		job, err := a.repos.JobRepo.GetLatestByUserNodeID(node.ID)
		if err != nil {
			SendError(w, 500, fmt.Errorf("node job fetch error: %v", err))
			return
		}

		od := &types.OverviewData{
			UserNode: node,
			Job:      job,
			Servers:  []*types.MinetestServerOverview{},
		}

		servers, err := a.repos.MinetestServerRepo.GetByNodeID(node.ID)
		if err != nil {
			SendError(w, 500, fmt.Errorf("servers fetch error: %v", err))
			return
		}
		for _, server := range servers {
			if server.State == types.MinetestServerStateDecommissioned {
				// skip decommissioned servers
				continue
			}

			server_job, err := a.repos.JobRepo.GetLatestByMinetestServerID(server.ID)
			if err != nil {
				SendError(w, 500, fmt.Errorf("server job fetch error: %v", err))
				return
			}

			od.Servers = append(od.Servers, &types.MinetestServerOverview{
				MinetestServer: server,
				Job:            server_job,
			})
		}

		slices.SortFunc(od.Servers, func(s1, s2 *types.MinetestServerOverview) int {
			return strings.Compare(s1.ID, s2.ID)
		})

		data = append(data, od)
	}

	slices.SortFunc(data, func(o1, o2 *types.OverviewData) int {
		return strings.Compare(o1.ID, o2.ID)
	})

	Send(w, data, nil)
}
