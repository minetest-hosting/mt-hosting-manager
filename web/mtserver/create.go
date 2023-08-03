package mtserver

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
	"strings"
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
		Nodes:  nodeinfos,
		NodeID: nodeid,
	}

	if r.Method == http.MethodPost {
		m.Name = r.FormValue("name")
		if !types.ValidServerName.Match([]byte(m.Name)) {
			m.NameErr = "invalid-server-name"
		}
		if strings.HasPrefix(m.Name, "node-") {
			m.NameErr = "invalid-server-name"
		}

		existing_server, err := ctx.repos.MinetestServerRepo.GetByName(m.Name)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, fmt.Errorf("server getbyname error: %v", err))
			return
		}
		if existing_server != nil {
			m.NameErr = "duplicate-server-name"
		}

		//TODO: create server stuff
	}

	ctx.tu.ExecuteTemplate(w, r, "mtserver/create.html", m)
}
