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
	ID         string    `json:"id" gorm:"primarykey;column:id"`
	State      UserState `json:"state" gorm:"column:state"`
	Name       string    `json:"name" gorm:"column:name"`
	Hash       string    `json:"hash" gorm:"column:hash"`
	Created    int64     `json:"created" gorm:"column:created"`
	LastLogin  int64     `json:"lastlogin" gorm:"column:lastlogin"`
	Balance    int64     `json:"balance" gorm:"column:balance"`
	ExternalID string    `json:"external_id" gorm:"column:external_id"`
	Currency   string    `json:"currency" gorm:"column:currency"`
	Type       UserType  `json:"type" gorm:"column:type"`
	Role       UserRole  `json:"role" gorm:"column:role"`
}

func (*User) TableName() string {
	return "public.user"
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
		"lastlogin",
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
		&m.LastLogin,
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
		m.LastLogin,
		m.Balance,
		m.ExternalID,
		m.Currency,
		m.Type,
		m.Role,
	}
}
