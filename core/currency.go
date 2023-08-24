package core

import (
	"fmt"
	"mt-hosting-manager/types"

	"github.com/bojanz/currency"
)

func ConvertCurrency(src *currency.Amount, src_xr *types.ExchangeRate, target_xr *types.ExchangeRate) (*currency.Amount, error) {

	amount_eur := *src
	if src.CurrencyCode() != "EUR" {
		// convert to "base-unit" euro first
		one_a, err := currency.NewAmountFromInt64(1, src_xr.Currency)
		if err != nil {
			return nil, fmt.Errorf("could not parse 'one' amount in currency '%s': %v", src_xr.Currency, err)
		}

		src_xr_a, err := one_a.Div(src_xr.Rate)
		if err != nil {
			return nil, fmt.Errorf("could not invert rate '%s': %v", src_xr.Rate, err)
		}

		amount_eur, err = src.Convert("EUR", src_xr_a.String())
		if err != nil {
			return nil, fmt.Errorf("could not convert to 'EUR': %v", err)
		}
	}

	target_amount, err := amount_eur.Convert(target_xr.Currency, target_xr.Rate)
	if err != nil {
		return nil, fmt.Errorf("could not convert to currency '%s': %v", target_xr.Currency, err)
	}

	return &target_amount, err
}
