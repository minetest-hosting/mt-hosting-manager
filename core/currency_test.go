package core_test

import (
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"testing"

	"github.com/bojanz/currency"
	"github.com/stretchr/testify/assert"
)

func TestCurrency(t *testing.T) {
	a, err := currency.NewAmount("10", "EUR")
	assert.NoError(t, err)
	assert.NotNil(t, a)

	// exchange rate EUR -> USD
	usd, err := a.Convert("USD", "1.211")
	assert.NoError(t, err)
	assert.NotNil(t, usd)

	expected, err := currency.NewAmount("12.11", "USD")
	assert.NoError(t, err)
	assert.NotNil(t, expected)

	assert.True(t, expected.Equal(usd))
}

func TestConvertUSDtoCHF(t *testing.T) {
	a, err := currency.NewAmount("10", "USD")
	assert.NoError(t, err)
	assert.NotNil(t, a)

	usd_xr := &types.ExchangeRate{Currency: "USD", Rate: "1.2"}
	chf_xr := &types.ExchangeRate{Currency: "CHF", Rate: "1.4"}

	chf, err := core.ConvertCurrency(&a, usd_xr, chf_xr)
	assert.NoError(t, err)
	assert.NotNil(t, chf)

	expected, err := currency.NewAmount("11.67", "CHF")
	assert.NoError(t, err)
	assert.True(t, expected.Equal(*chf))
}
