package core_test

import (
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
