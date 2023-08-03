package types

import (
	"time"
)

type UserNodeState string

const (
	UserNodeStateCreated      UserNodeState = "CREATED"
	UserNodeStateProvisioning UserNodeState = "PROVISIONING"
	UserNodeStateRunning      UserNodeState = "RUNNING"
	UserNodeStateStopped      UserNodeState = "STOPPED"
	UserNodeStateRemoving     UserNodeState = "REMOVING"
)

// Created -> Provisioning -> Running <-> Stopped
//                                     -> Removing

const ExpirationWarnThreshold = time.Hour * 24 * 14

type UserNode struct {
	ID          string        `json:"id"`
	UserID      string        `json:"user_id"`
	NodeTypeID  string        `json:"node_type_id"`
	ExternalID  string        `json:"external_id"`
	Created     int64         `json:"created"`
	Expires     int64         `json:"expires"`
	State       UserNodeState `json:"state"`
	Name        string        `json:"name"`
	IPv4        string        `json:"ipv4"`
	IPv6        string        `json:"ipv6"`
	LoadPercent int           `json:"load_percent"`
	DiskSize    int64         `json:"disk_size"`
	DiskUsed    int64         `json:"disk_used"`
	MemorySize  int64         `json:"memory_size"`
	MemoryUsed  int64         `json:"memory_used"`
}

func (m *UserNode) Columns(action string) []string {
	return []string{
		"id",
		"user_id",
		"node_type_id",
		"external_id",
		"created",
		"expires",
		"state",
		"name",
		"ipv4",
		"ipv6",
		"load_percent",
		"disk_size",
		"disk_used",
		"memory_size",
		"memory_used",
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
		&m.Expires,
		&m.State,
		&m.Name,
		&m.IPv4,
		&m.IPv6,
		&m.LoadPercent,
		&m.DiskSize,
		&m.DiskUsed,
		&m.MemorySize,
		&m.MemoryUsed,
	)
}

func (m *UserNode) Values(action string) []any {
	return []any{
		m.ID,
		m.UserID,
		m.NodeTypeID,
		m.ExternalID,
		m.Created,
		m.Expires,
		m.State,
		m.Name,
		m.IPv4,
		m.IPv6,
		m.LoadPercent,
		m.DiskSize,
		m.DiskUsed,
		m.MemorySize,
		m.MemoryUsed,
	}
}

func (m *UserNode) ExpirationWarning() bool {
	return time.Unix(m.Expires, 0).Add(-ExpirationWarnThreshold).Before(time.Now())
}
