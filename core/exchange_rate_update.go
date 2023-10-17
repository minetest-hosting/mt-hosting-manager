package core

import (
	"fmt"
	"mt-hosting-manager/api/coinbase"
)

var enabledCurrencies = map[string]bool{
	"EUR": true,
	"USD": true,
	"BTC": true,
	"ETH": true,
}

func (c *Core) UpdateExchangeRates() error {
	r, err := c.cbc.GetRates(coinbase.CURRENCY_EUR)
	if err != nil {
		return fmt.Errorf("GetRates error: %v", err)
	}

	cu, err := c.cbc.GetCurrencies()
	if err != nil {
		return fmt.Errorf("GetCurrencies error: %v", err)
	}

	cc, err := c.cbc.GetCryptoCurrencies()
	if err != nil {
		return fmt.Errorf("GetCryptoCurrencies error: %v", err)
	}

	list, err := c.repos.ExchangeRateRepo.GetAll()
	if err != nil {
		return fmt.Errorf("db get error: %v", err)
	}

	updatedCurrencies := map[string]bool{}
	needs_fetch := false

	for _, e := range list {
		updatedCurrencies[e.Currency] = true
		r.Data.Rates[e.Currency] = e.Rate
		err = c.repos.ExchangeRateRepo.Update(e)
	}

	return nil
}
