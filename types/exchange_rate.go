package types

type ExchangeRate struct {
	Currency    string `json:"currency"`
	Rate        string `json:"rate"`
	DisplayName string `json:"display_name"`
	Digits      int    `json:"digits"`
	Updated     int64  `json:"updated"`
}

func ExchangeRateProvider() *ExchangeRate { return &ExchangeRate{} }

func (m *ExchangeRate) Columns(action string) []string {
	return []string{"currency", "rate", "display_name", "digits", "updated"}
}

func (m *ExchangeRate) Table() string {
	return "exchange_rate"
}

func (m *ExchangeRate) Scan(action string, r func(dest ...any) error) error {
	return r(&m.Currency, &m.Rate, &m.DisplayName, &m.Digits, &m.Updated)
}

func (m *ExchangeRate) Values(action string) []any {
	return []any{m.Currency, m.Rate, m.DisplayName, m.Digits, m.Updated}
}
