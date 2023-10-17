package core

import (
	"fmt"
	"mt-hosting-manager/api/coinbase"
	"time"
)

func (c *Core) UpdateExchangeRates() error {
	r, err := c.cbc.GetRates(coinbase.CURRENCY_EUR)
	if err != nil {
		return fmt.Errorf("GetRates error: %v", err)
	}

	list, err := c.repos.ExchangeRateRepo.GetAll()
	if err != nil {
		return fmt.Errorf("db get error: %v", err)
	}

	for _, e := range list {
		e.Rate = r.Data.Rates[e.Currency]
		if e.Rate == "" {
			return fmt.Errorf("currency update failed for %s:", e.Currency)
		}

		e.Updated = time.Now().Unix()
		err = c.repos.ExchangeRateRepo.Update(e)
		if err != nil {
			return err
		}
	}

	return nil
}
