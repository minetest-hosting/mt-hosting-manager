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
	ID                      string        `json:"id" gorm:"primarykey;column:id"`
	State                   NodeTypeState `json:"state" gorm:"column:state"`
	OrderID                 int           `json:"order_id" gorm:"column:order_id"`
	Provider                ProviderType  `json:"provider" gorm:"column:provider"`
	ServerType              string        `json:"server_type" gorm:"column:server_type"`
	Locations               string        `json:"locations" gorm:"column:locations"`
	Name                    string        `json:"name" gorm:"column:name"`
	Description             string        `json:"description" gorm:"column:description"`
	CpuCount                int           `json:"cpu_count" gorm:"column:cpu_count"`
	RamGB                   int           `json:"ram_gb" gorm:"column:ram_gb"`
	DiskGB                  int           `json:"disk_gb" gorm:"column:disk_gb"`
	Dedicated               bool          `json:"dedicated" gorm:"column:dedicated"`
	DailyCost               int64         `json:"daily_cost" gorm:"column:daily_cost"`
	MaxRecommendedInstances int           `json:"max_recommended_instances" gorm:"column:max_recommended_instances"`
	MaxInstances            int           `json:"max_instances" gorm:"column:max_instances"`
}

func (m *NodeType) TableName() string {
	return "node_type"
}
