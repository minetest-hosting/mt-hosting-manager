package types

type UserType string

const (
	UserTypeGithub  UserType = "GITHUB"
	UserTypeDiscord UserType = "DISCORD"
)

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Mail       string `json:"mail"`
	Created    int64  `json:"created"`
	ExternalID string `json:"external_id"`
	Type       string `json:"type"`
}

func (m *User) Columns(action string) []string {
	return []string{"id", "name", "mail", "created", "external_id", "type"}
}

func (m *User) Table() string {
	return "user"
}

func (m *User) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.Name, &m.Mail, &m.Created, &m.ExternalID, &m.Type)
}

func (m *User) Values(action string) []any {
	return []any{m.ID, m.Name, m.Mail, m.Created, m.ExternalID, m.Type}
}
