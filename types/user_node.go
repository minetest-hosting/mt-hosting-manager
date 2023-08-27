package types

type UserNodeState string

const (
	UserNodeStateCreated      UserNodeState = "CREATED"
	UserNodeStateProvisioning UserNodeState = "PROVISIONING"
	UserNodeStateRunning      UserNodeState = "RUNNING"
	UserNodeStateRemoving     UserNodeState = "REMOVING"
)

// Created -> Provisioning -> Running <-> Stopped
//                                     -> Removing

func UserNodeProvider() *UserNode { return &UserNode{} }

type UserNode struct {
	ID                string        `json:"id"`
	UserID            string        `json:"user_id"`
	NodeTypeID        string        `json:"node_type_id"`
	ExternalID        string        `json:"external_id"`
	Created           int64         `json:"created"`
	LastCollectedTime int64         `json:"last_collected_time"`
	State             UserNodeState `json:"state"`
	Name              string        `json:"name"`
	Alias             string        `json:"alias"`
	IPv4              string        `json:"ipv4"`
	IPv6              string        `json:"ipv6"`
	Fingerprint       string        `json:"fingerprint"`
}

func (m *UserNode) Columns(action string) []string {
	return []string{
		"id",
		"user_id",
		"node_type_id",
		"external_id",
		"created",
		"last_collected_time",
		"state",
		"name",
		"alias",
		"ipv4",
		"ipv6",
		"fingerprint",
	}
}

func (m *UserNode) Table() string {
	return "user_node"
}

func (m *UserNode) Scan(action string, r func(dest ...any) error) error {
	return r(
		&m.ID,
		&m.UserID,
		&m.NodeTypeID,
		&m.ExternalID,
		&m.Created,
		&m.LastCollectedTime,
		&m.State,
		&m.Name,
		&m.Alias,
		&m.IPv4,
		&m.IPv6,
		&m.Fingerprint,
	)
}

func (m *UserNode) Values(action string) []any {
	return []any{
		m.ID,
		m.UserID,
		m.NodeTypeID,
		m.ExternalID,
		m.Created,
		m.LastCollectedTime,
		m.State,
		m.Name,
		m.Alias,
		m.IPv4,
		m.IPv6,
		m.Fingerprint,
	}
}
