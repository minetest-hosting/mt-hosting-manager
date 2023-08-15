package mtserver

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type NodeInfo struct {
	*types.UserNode
	Selected bool
}

type CreateServerModel struct {
	Nodes   []*NodeInfo
	NodeID  string
	NameErr string
	Name    string
	DNSName string
}

func (ctx *Context) Create(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	nodes, err := ctx.repos.UserNodeRepo.GetByUserID(c.UserID)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("get nodes error: %v", err))
		return
	}

	nodeid := r.URL.Query().Get("node_id")

	nodeinfos := make([]*NodeInfo, len(nodes))
	for i, node := range nodes {
		nodeinfos[i] = &NodeInfo{
			UserNode: node,
			Selected: node.ID == nodeid,
		}
	}

	m := &CreateServerModel{
		Nodes:   nodeinfos,
		NodeID:  nodeid,
		Name:    r.FormValue("name"),
		DNSName: r.FormValue("DNSName"),
	}

	if r.Method == http.MethodPost {
		m.DNSName = r.FormValue("DNSName")
		if !types.ValidServerName.Match([]byte(m.DNSName)) {
			m.NameErr = "invalid-server-name"
		}
		if strings.HasPrefix(m.DNSName, "node-") {
			m.NameErr = "invalid-server-name"
		}

		existing_server, err := ctx.repos.MinetestServerRepo.GetByName(m.DNSName)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, fmt.Errorf("server getbyname error: %v", err))
			return
		}
		if existing_server != nil {
			m.NameErr = "duplicate-server-name"
		}

		if m.NameErr == "" {
			// valid name, create server
			node, err := ctx.repos.UserNodeRepo.GetByID(nodeid)
			if err != nil {
				ctx.tu.RenderError(w, r, 500, fmt.Errorf("usernode getbyid error: %v", err))
				return
			}
			if node == nil {
				ctx.tu.RenderError(w, r, 404, fmt.Errorf("usernode not found: %s", nodeid))
				return
			}

			server := &types.MinetestServer{
				ID:         uuid.NewString(),
				UserNodeID: node.ID,
				Name:       m.Name,
				DNSName:    m.DNSName,
				Created:    time.Now().Unix(),
				State:      types.MinetestServerStateCreated,
			}
			err = ctx.repos.MinetestServerRepo.Insert(server)
			if err != nil {
				ctx.tu.RenderError(w, r, 500, fmt.Errorf("server insert error: %v", err))
				return
			}

			job := &types.Job{
				ID:               uuid.NewString(),
				Type:             types.JobTypeServerSetup,
				State:            types.JobStateCreated,
				UserNodeID:       &node.ID,
				MinetestServerID: &server.ID,
			}
			err = ctx.repos.JobRepo.Insert(job)
			if err != nil {
				ctx.tu.RenderError(w, r, 500, fmt.Errorf("job insert error: %v", err))
				return
			}

			http.Redirect(w, r, fmt.Sprintf("%s/mtserver/%s", ctx.tu.BaseURL, server.ID), http.StatusSeeOther)
			return
		}

	}

	ctx.tu.ExecuteTemplate(w, r, "mtserver/create.html", m)
}
