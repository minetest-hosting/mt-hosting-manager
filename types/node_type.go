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

type NodeType struct {
	ID                      string       `json:"id"`
	State                   string       `json:"state"`
	OrderID                 int          `json:"order_id"`
	Provider                ProviderType `json:"provider"`
	ServerType              string       `json:"server_type"`
	Name                    string       `json:"name"`
	Description             string       `json:"description"`
	MonthlyCost             string       `json:"monthly_cost"`
	MaxMonths               int          `json:"max_months"`
	MaxRecommendedInstances int          `json:"max_recommended_instances"`
	MaxInstances            int          `json:"max_instances"`
}

func (m *NodeType) Columns(action string) []string {
	return []string{
		"id",
		"state",
		"order_id",
		"provider",
		"server_type",
		"name",
		"description",
		"monthly_cost",
		"max_months",
		"max_recommended_instances",
		"max_instances",
	}
}

func (m *NodeType) Table() string {
	return "node_type"
}

func (m *NodeType) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.State, &m.OrderID, &m.Provider, &m.ServerType, &m.Name, &m.Description, &m.MonthlyCost, &m.MaxMonths, &m.MaxRecommendedInstances, &m.MaxInstances)
}

func (m *NodeType) Values(action string) []any {
	return []any{m.ID, m.State, m.OrderID, m.Provider, m.ServerType, m.Name, m.Description, m.MonthlyCost, m.MaxMonths, m.MaxRecommendedInstances, m.MaxInstances}
}
