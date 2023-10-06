package types

import (
	"fmt"
	"regexp"
)

type MinetestServerState string

const (
	MinetestServerStateCreated        MinetestServerState = "CREATED"
	MinetestServerStateProvisioning   MinetestServerState = "PROVISIONING"
	MinetestServerStateRunning        MinetestServerState = "RUNNING"
	MinetestServerStateRemoving       MinetestServerState = "REMOVING"
	MinetestServerStateDecommissioned MinetestServerState = "DECOMMISSIONED"
)

var ValidServerName = regexp.MustCompile(`^[a-z|A-Z|0-9]+$`)

var ValidPlayernameRegex = regexp.MustCompile(`^[a-zA-Z0-9\-_]*$`)

func ValidateUsername(username string) error {
	if len(username) == 0 {
		return fmt.Errorf("playername empty")
	}
	if len(username) > 20 {
		return fmt.Errorf("playername too long")
	}
	if !ValidPlayernameRegex.Match([]byte(username)) {
		return fmt.Errorf("playername can only contain chars a to z, A to Z, 0 to 9 and -, _")
	}
	return nil
}

func MinetestServerProvider() *MinetestServer { return &MinetestServer{} }

type MinetestServer struct {
	ID                 string              `json:"id"`
	UserNodeID         string              `json:"user_node_id"`
	Name               string              `json:"name"`
	DNSName            string              `json:"dns_name"`
	Admin              string              `json:"admin"`
	ExternalCNAMEDNSID string              `json:"external_cname_dns_id"`
	CustomDNS          string              `json:"custom_dns_name"`
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
		"admin",
		"external_cname_dns_id",
		"custom_dns_name",
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
	return r(&m.ID, &m.UserNodeID, &m.Name, &m.DNSName, &m.Admin, &m.ExternalCNAMEDNSID, &m.CustomDNS, &m.Port, &m.UIVersion, &m.JWTKey, &m.Created, &m.State)
}

func (m *MinetestServer) Values(action string) []any {
	return []any{m.ID, m.UserNodeID, m.Name, m.DNSName, m.Admin, m.ExternalCNAMEDNSID, m.CustomDNS, m.Port, m.UIVersion, m.JWTKey, m.Created, m.State}
}
