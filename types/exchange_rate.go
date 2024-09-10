package types

type ExchangeRate struct {
	Currency    string `json:"currency" gorm:"primarykey;column:currency"`
	Rate        string `json:"rate" gorm:"column:rate"`
	DisplayName string `json:"display_name" gorm:"column:display_name"`
	Digits      int    `json:"digits" gorm:"column:digits"`
	Updated     int64  `json:"updated" gorm:"column:updated"`
}

func (m *ExchangeRate) TableName() string {
	return "exchange_rate"
}
