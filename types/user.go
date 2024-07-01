package types

type UserType string

const (
	UserTypeGithub  UserType = "GITHUB"
	UserTypeDiscord UserType = "DISCORD"
	UserTypeMesehub UserType = "MESEHUB"
	UserTypeLocal   UserType = "LOCAL"
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
	Hash       string    `json:"hash"`
	Created    int64     `json:"created"`
	Balance    int64     `json:"balance"`
	ExternalID string    `json:"external_id"`
	Currency   string    `json:"currency"`
	Type       UserType  `json:"type"`
	Role       UserRole  `json:"role"`
}

func (u *User) RemoveSensitiveFields() {
	u.Hash = ""
	u.ExternalID = ""
}

type UserSearch struct {
	UserID   *string `json:"user_id"`
	NameLike *string `json:"name_like"`
	Limit    *int    `json:"limit"`
}

func (m *User) Columns(action string) []string {
	return []string{
		"id",
		"state",
		"name",
		"hash",
		"created",
		"balance",
		"external_id",
		"currency",
		"type",
		"role",
	}
}

func (m *User) Table() string {
	return "public.user"
}

func (m *User) Scan(action string, r func(dest ...any) error) error {
	return r(
		&m.ID,
		&m.State,
		&m.Name,
		&m.Hash,
		&m.Created,
		&m.Balance,
		&m.ExternalID,
		&m.Currency,
		&m.Type,
		&m.Role,
	)
}

func (m *User) Values(action string) []any {
	return []any{
		m.ID,
		m.State,
		m.Name,
		m.Hash,
		m.Created,
		m.Balance,
		m.ExternalID,
		m.Currency,
		m.Type,
		m.Role,
	}
}
