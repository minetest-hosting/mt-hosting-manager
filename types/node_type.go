package types

type ProviderType string

const (
	ProviderHetzner ProviderType = "HETZNER"
)

type NodeTypeState string

const (
	NodeTypeStateInactive   NodeTypeState = "INACTIVE"
	NodeTypeStateActive     NodeTypeState = "ACTIVE"
	NodeTypeStateDeprecated NodeTypeState = "DEPRECATED"
)

func NodeTypeProvider() *NodeType { return &NodeType{} }

type NodeType struct {
	ID                      string        `json:"id"`
	State                   NodeTypeState `json:"state"`
	OrderID                 int           `json:"order_id"`
	Provider                ProviderType  `json:"provider"`
	ServerType              string        `json:"server_type"`
	Locations               string        `json:"locations"`
	Name                    string        `json:"name"`
	Description             string        `json:"description"`
	CpuCount                int           `json:"cpu_count"`
	RamGB                   int           `json:"ram_gb"`
	DiskGB                  int           `json:"disk_gb"`
	Dedicated               bool          `json:"dedicated"`
	DailyCost               int64         `json:"daily_cost"`
	MaxRecommendedInstances int           `json:"max_recommended_instances"`
	MaxInstances            int           `json:"max_instances"`
}

func (m *NodeType) Columns(action string) []string {
	return []string{
		"id",
		"state",
		"order_id",
		"provider",
		"server_type",
		"locations",
		"name",
		"description",
		"cpu_count",
		"ram_gb",
		"disk_gb",
		"dedicated",
		"daily_cost",
		"max_recommended_instances",
		"max_instances",
	}
}

func (m *NodeType) Table() string {
	return "node_type"
}

func (m *NodeType) Scan(action string, r func(dest ...any) error) error {
	return r(
		&m.ID,
		&m.State,
		&m.OrderID,
		&m.Provider,
		&m.ServerType,
		&m.Locations,
		&m.Name,
		&m.Description,
		&m.CpuCount,
		&m.RamGB,
		&m.DiskGB,
		&m.Dedicated,
		&m.DailyCost,
		&m.MaxRecommendedInstances,
		&m.MaxInstances,
	)
}

func (m *NodeType) Values(action string) []any {
	return []any{
		m.ID,
		m.State,
		m.OrderID,
		m.Provider,
		m.ServerType,
		m.Locations,
		m.Name,
		m.Description,
		m.CpuCount,
		m.RamGB,
		m.DiskGB,
		m.Dedicated,
		m.DailyCost,
		m.MaxRecommendedInstances,
		m.MaxInstances,
	}
}
