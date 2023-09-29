package web

import (
	"mt-hosting-manager/api/coinbase"
	"net/http"
	"time"
)

type CurrencyResponse struct {
	Rates            *coinbase.RatesData            `json:"rates"`
	Currencies       []*coinbase.CurrencyData       `json:"currencies"`
	CryptoCurrencies []*coinbase.CryptoCurrencyData `json:"crypto_currencies"`
}

var curr *CurrencyResponse
var curr_updated time.Time

func (a *Api) GetCurrencies(w http.ResponseWriter, r *http.Request) {
	if curr == nil || time.Since(curr_updated) > time.Hour {
		r, err := a.cbc.GetRates(coinbase.CURRENCY_EUR)
		if err != nil {
			SendError(w, 500, err)
			return
		}

		c, err := a.cbc.GetCurrencies()
		if err != nil {
			SendError(w, 500, err)
			return
		}

		cc, err := a.cbc.GetCryptoCurrencies()
		if err != nil {
			SendError(w, 500, err)
			return
		}

		curr = &CurrencyResponse{
			Rates:            r.Data,
			Currencies:       c.Data,
			CryptoCurrencies: cc.Data,
		}

		curr_updated = time.Now()
	}
	Send(w, curr, nil)
}
