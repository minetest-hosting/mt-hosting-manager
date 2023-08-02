package hetzner_cloud

type LocationType string

const (
	LocationNuernberg LocationType = "nbg1"
)

type PublicNet struct {
	EnableIPv4 bool `json:"enable_ipv4"`
	EnableIPv6 bool `json:"enable_ipv6"`
}

type CreateServerRequest struct {
	Image            string            `json:"image"`
	Labels           map[string]string `json:"labels"`
	Location         LocationType      `json:"location"`
	Name             string            `json:"name"`
	PublicNet        *PublicNet        `json:"public_net"`
	ServerType       string            `json:"server_type"`
	SSHKeys          []string          `json:"ssh_keys"`
	StartAfterCreate bool              `json:"start_after_create"`
}

type PublicNetEntry struct {
	IP string `json:"ip"`
}

type PublicNetInfo struct {
	IPv4 *PublicNetEntry `json:"ipv4"`
	IPv6 *PublicNetEntry `json:"ipv6"`
}

type ServerInfo struct {
	ID        int            `json:"id"`
	Status    string         `json:"status"`
	PublicNet *PublicNetInfo `json:"public_net"`
}

type CreateServerResponse struct {
	Server *ServerInfo `json:"server"`
}
