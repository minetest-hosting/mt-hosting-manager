package types

type ProviderType string

const (
	ProviderHetzner ProviderType = "HETZNER"
)

type Node struct {
	ID                      string       `json:"id"`
	Deprecated              bool         `json:"deprecated"`
	OrderID                 int          `json:"order_id"`
	Provider                ProviderType `json:"provider"`
	ServerType              string       `json:"server_type"`
	Name                    string       `json:"name"`
	Description             string       `json:"description"`
	CostPerHour             int64        `json:"cost_per_hour"`
	MaxRecommendedInstances int          `json:"max_recommended_instances"`
	MaxInstances            int          `json:"max_instances"`
}

func (m *Node) Columns(action string) []string {
	return []string{
		"id",
		"deprecated",
		"order_id",
		"provider",
		"server_type",
		"name",
		"description",
		"cost_per_hour",
		"max_recommended_instances",
		"max_instances",
	}
}

func (m *Node) Table() string {
	return "node"
}

func (m *Node) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.Deprecated, &m.OrderID, &m.Provider, &m.ServerType, &m.Name, &m.Description, &m.CostPerHour, &m.MaxRecommendedInstances, &m.MaxInstances)
}

func (m *Node) Values(action string) []any {
	return []any{m.ID, m.Deprecated, m.OrderID, m.Provider, m.ServerType, m.Name, m.Description, m.CostPerHour, m.MaxRecommendedInstances, m.MaxInstances}
}
