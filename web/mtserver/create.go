package mtserver

import (
	"fmt"
	"mt-hosting-manager/types"
	"mt-hosting-manager/web/components"
	"mt-hosting-manager/worker"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type NodeInfo struct {
	*types.UserNode
	Selected bool
}

type CreateServerModel struct {
	Nodes      []*NodeInfo
	NodeID     string
	Name       string
	NameErr    string
	Port       string
	PortErr    string
	DNSName    string
	DNSNameErr string
	Breadcrumb *components.Breadcrumb
}

func (ctx *Context) Create(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	nodes, err := ctx.repos.UserNodeRepo.GetByUserIDAndState(c.UserID, types.UserNodeStateRunning)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, fmt.Errorf("get nodes error: %v", err))
		return
	}

	nodeid := r.URL.Query().Get("node_id")
	if nodeid == "" {
		nodeid = r.FormValue("nodeid")
	}

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
		Name:    r.FormValue("Name"),
		DNSName: r.FormValue("DNSName"),
		Port:    r.FormValue("Port"),
		Breadcrumb: &components.Breadcrumb{
			Entries: []*components.BreadcrumbEntry{
				components.HomeBreadcrumb,
				components.ServersBreadcrumb, {
					Name:   "Create server",
					FAIcon: "plus",
				},
			},
		},
	}

	if r.Method == http.MethodPost {
		node, err := ctx.repos.UserNodeRepo.GetByID(m.NodeID)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, fmt.Errorf("usernode getbyid error: %v", err))
			return
		}
		if node == nil {
			ctx.tu.RenderError(w, r, 404, fmt.Errorf("usernode not found: %s", nodeid))
			return
		}

		// check for valid port number
		port_num, err := strconv.ParseInt(m.Port, 10, 64)
		if err != nil || port_num < 1000 || port_num > 65535 {
			m.PortErr = "port-number-port-err"
		}
		if m.Name == "" {
			m.NameErr = "server-name-empty-err"
		}

		// duplicate port number
		other_servers, err := ctx.repos.MinetestServerRepo.GetByNodeID(node.ID)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, fmt.Errorf("other_servers getby-nodeid error: %v", err))
			return
		}

		for _, os := range other_servers {
			if os.Port == int(port_num) {
				m.PortErr = "port-number-already-in-use"
			}
		}

		// check for valid dns name
		m.DNSName = r.FormValue("DNSName")
		if !types.ValidServerName.Match([]byte(m.DNSName)) {
			m.DNSNameErr = "invalid-server-name"
		}
		if strings.HasPrefix(m.DNSName, "node-") {
			m.DNSNameErr = "invalid-server-name"
		}

		// duplicate dns name
		existing_server, err := ctx.repos.MinetestServerRepo.GetByName(m.DNSName)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, fmt.Errorf("server getbyname error: %v", err))
			return
		}
		if existing_server != nil {
			m.DNSNameErr = "duplicate-server-name"
		}

		if m.DNSNameErr == "" && m.NameErr == "" && m.PortErr == "" {
			// valid name, create server
			server := &types.MinetestServer{
				ID:         uuid.NewString(),
				UserNodeID: node.ID,
				Name:       m.Name,
				DNSName:    m.DNSName,
				Port:       int(port_num),
				Created:    time.Now().Unix(),
				UIVersion:  "latest",
				State:      types.MinetestServerStateCreated,
			}
			err = ctx.repos.MinetestServerRepo.Insert(server)
			if err != nil {
				ctx.tu.RenderError(w, r, 500, fmt.Errorf("server insert error: %v", err))
				return
			}

			job := worker.SetupServerJob(node, server)
			err = ctx.repos.JobRepo.Insert(job)
			if err != nil {
				ctx.tu.RenderError(w, r, 500, fmt.Errorf("job insert error: %v", err))
				return
			}

			http.Redirect(w, r, fmt.Sprintf("%s/mtserver/%s", ctx.tu.BaseURL, server.ID), http.StatusSeeOther)
			return
		}

	}

	// default values
	if m.Port == "" {
		m.Port = "30000"
	}

	ctx.tu.ExecuteTemplate(w, r, "mtserver/create.html", m)
}
