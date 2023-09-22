package types

type UserType string

const (
	UserTypeGithub  UserType = "GITHUB"
	UserTypeDiscord UserType = "DISCORD"
)

type UserRole string

const (
	UserRoleUser  UserRole = "USER"
	UserRoleAdmin UserRole = "ADMIN"
)

type UserState string

const (
	UserStateActive UserState = "ACTIVE"
)

func UserProvider() *User { return &User{} }

type User struct {
	ID         string    `json:"id"`
	State      UserState `json:"state"`
	Name       string    `json:"name"`
	Mail       string    `json:"mail"`
	Created    int64     `json:"created"`
	Balance    int64     `json:"balance"`
	ExternalID string    `json:"external_id"`
	Type       UserType  `json:"type"`
	Role       UserRole  `json:"role"`
}

func (m *User) Columns(action string) []string {
	return []string{"id", "state", "name", "mail", "created", "balance", "external_id", "type", "role"}
}

func (m *User) Table() string {
	return "user"
}

func (m *User) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.State, &m.Name, &m.Mail, &m.Created, &m.Balance, &m.ExternalID, &m.Type, &m.Role)
}

func (m *User) Values(action string) []any {
	return []any{m.ID, m.State, m.Name, m.Mail, m.Created, m.Balance, m.ExternalID, m.Type, m.Role}
}
