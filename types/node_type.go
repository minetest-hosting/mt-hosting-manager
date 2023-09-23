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
	Location                string        `json:"location"`
	Name                    string        `json:"name"`
	Description             string        `json:"description"`
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
		"location",
		"name",
		"description",
		"daily_cost",
		"max_recommended_instances",
		"max_instances",
	}
}

func (m *NodeType) Table() string {
	return "node_type"
}

func (m *NodeType) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.State, &m.OrderID, &m.Provider, &m.ServerType, &m.Location, &m.Name, &m.Description, &m.DailyCost, &m.MaxRecommendedInstances, &m.MaxInstances)
}

func (m *NodeType) Values(action string) []any {
	return []any{m.ID, m.State, m.OrderID, m.Provider, m.ServerType, m.Location, m.Name, m.Description, m.DailyCost, m.MaxRecommendedInstances, m.MaxInstances}
}
