package types

type UserSetting struct {
	UserID string `json:"user_id" gorm:"primarykey;column:user_id"`
	Key    string `json:"key" gorm:"column:key"`
	Value  string `json:"value" gorm:"column:value"`
}

func (*UserSetting) TableName() string {
	return "user_setting"
}
