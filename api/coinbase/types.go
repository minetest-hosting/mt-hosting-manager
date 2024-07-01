package coinbase

import "time"

const (
	CURRENCY_EUR = "EUR"
	CURRENCY_ETH = "ETH"
	CURRENCY_BTC = "BTC"
)

type RatesData struct {
	Currency string            `json:"currency"`
	Rates    map[string]string `json:"rates"`
}

type RatesResponse struct {
	Data *RatesData `json:"data"`
}

type CryptoCurrencyData struct {
	AssetID      string `json:"asset_id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Color        string `json:"color"`
	Exponent     int    `json:"exponent"`
	AddressRegex string `json:"address_regex"`
}

type CryptoCurrencies struct {
	Data []*CryptoCurrencyData `json:"data"`
}

type CurrencyData struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	MinSize string `json:"min_size"`
}

type Currencies struct {
	Data []*CurrencyData `json:"data"`
}

type PricingType string

const (
	PricingTypeNoPrice PricingType = "no_price"
	PricingTypeFixed   PricingType = "fixed_price"
)

type LocalPrice struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type CreateChargeRequest struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	PricingType PricingType `json:"pricing_type"`
	LocalPrice  *LocalPrice `json:"local_price"`
	RedirectURL string      `json:"redirect_url"`
	CancelURL   string      `json:"cancel_url"`
}

type ChargeAddresses struct {
	Bitcoin  string `json:"bitcoin"`
	Ethereum string `json:"ethereum"`
}

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusConfirmed PaymentStatus = "CONFIRMED"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
)

type PaymentValue struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

const (
	PaymentValueLocal  = "local"
	PaymentValueCrypto = "crypto"
)

type Payment struct {
	TransactionID string        `json:"transaction_id"`
	Status        PaymentStatus `json:"status"`
	Value         map[string]*PaymentValue
}

type ChargeData struct {
	HostedURL string           `json:"hosted_url"`
	Addresses *ChargeAddresses `json:"addresses"`
	Code      string           `json:"code"`
	CreatedAt time.Time        `json:"created_at"`
	ExpiredAt time.Time        `json:"expires_at"`
	ID        string           `json:"id"`
	Payments  []*Payment       `json:"payments"`
}

type Charge struct {
	Data *ChargeData `json:"data"`
}
