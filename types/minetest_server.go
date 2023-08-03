package types

import "regexp"

type MinetestServerState string

const (
	MinetestServerStateCreated      MinetestServerState = "CREATED"
	MinetestServerStateProvisioning MinetestServerState = "PROVISIONING"
	MinetestServerStateRunning      MinetestServerState = "RUNNING"
	MinetestServerStateStopped      MinetestServerState = "STOPPED"
	MinetestServerStateRemoving     MinetestServerState = "REMOVING"
)

var ValidServerName = regexp.MustCompile(`^[a-z|A-Z|0-9]+$`)

type MinetestServer struct {
	ID         string        `json:"id"`
	UserNodeID string        `json:"user_node_id"`
	Name       string        `json:"name"`
	Created    int64         `json:"created"`
	State      UserNodeState `json:"state"`
}

func (m *MinetestServer) Columns(action string) []string {
	return []string{
		"id",
		"user_node_id",
		"name",
		"created",
	}
}

func (m *MinetestServer) Table() string {
	return "minetest_server"
}

func (m *MinetestServer) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.UserNodeID, &m.Name, &m.Created)
}

func (m *MinetestServer) Values(action string) []any {
	return []any{m.ID, m.UserNodeID, m.Name, m.Created}
}
