package core

import (
	"fmt"
	"mt-hosting-manager/types"
)

type CreateServerResult struct {
	Valid             bool `json:"valid"`
	PortInvalid       bool `json:"port_invalid"`
	PortUsed          bool `json:"port_used"`
	AdminNameInvalid  bool `json:"admin_name_invalid"`
	ServerNameInvalid bool `json:"server_name_invalid"`
}

func (c *Core) ValidateCreateServer(server *types.MinetestServer, node *types.UserNode) (*CreateServerResult, error) {
	csr := &CreateServerResult{
		Valid: true,
	}

	other_servers, err := c.repos.MinetestServerRepo.GetByNodeID(node.ID)
	if err != nil {
		return nil, fmt.Errorf("servers fetch error: %v", err)
	}

	if server.Port < 1000 || server.Port > 65000 {
		csr.PortInvalid = true
		csr.Valid = false
	}

	err = types.ValidateUsername(server.Admin)
	if err != nil {
		csr.AdminNameInvalid = true
		csr.Valid = false
	}

	for _, s := range other_servers {
		if s.Port == server.Port {
			csr.PortUsed = true
			csr.Valid = false
			break
		}
	}

	if !types.ValidServerName.Match([]byte(server.DNSName)) {
		csr.ServerNameInvalid = true
		csr.Valid = false
	}

	return csr, nil
}
