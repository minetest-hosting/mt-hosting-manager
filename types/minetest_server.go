package types

import "regexp"

type MinetestServerState string

const (
	MinetestServerStateCreated      MinetestServerState = "CREATED"
	MinetestServerStateProvisioning MinetestServerState = "PROVISIONING"
	MinetestServerStateRunning      MinetestServerState = "RUNNING"
	MinetestServerStateRemoving     MinetestServerState = "REMOVING"
)

var ValidServerName = regexp.MustCompile(`^[a-z|A-Z|0-9]+$`)

func MinetestServerProvider() *MinetestServer { return &MinetestServer{} }

type MinetestServer struct {
	ID                 string              `json:"id"`
	UserNodeID         string              `json:"user_node_id"`
	Name               string              `json:"name"`
	DNSName            string              `json:"dns_name"`
	ExternalCNAMEDNSID string              `json:"external_cname_dns_id"`
	Port               int                 `json:"port"`
	UIVersion          string              `json:"ui_version"`
	JWTKey             string              `json:"jwt_key"`
	Created            int64               `json:"created"`
	State              MinetestServerState `json:"state"`
}

func (m *MinetestServer) Columns(action string) []string {
	return []string{
		"id",
		"user_node_id",
		"name",
		"dns_name",
		"external_cname_dns_id",
		"port",
		"ui_version",
		"jwt_key",
		"created",
		"state",
	}
}

func (m *MinetestServer) Table() string {
	return "minetest_server"
}

func (m *MinetestServer) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.UserNodeID, &m.Name, &m.DNSName, &m.ExternalCNAMEDNSID, &m.Port, &m.UIVersion, &m.JWTKey, &m.Created, &m.State)
}

func (m *MinetestServer) Values(action string) []any {
	return []any{m.ID, m.UserNodeID, m.Name, m.DNSName, m.ExternalCNAMEDNSID, m.Port, m.UIVersion, m.JWTKey, m.Created, m.State}
}
