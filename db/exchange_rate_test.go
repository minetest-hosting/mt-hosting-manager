package db_test

import (
	"mt-hosting-manager/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExchangeRateRepository(t *testing.T) {
	repos := SetupRepos(t)

	assert.NoError(t, repos.ExchangeRateRepo.Insert(&types.ExchangeRate{
		Currency:    "EUR",
		Rate:        "1.0",
		DisplayName: "Euro",
		Digits:      2,
	}))

	list, err := repos.ExchangeRateRepo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "EUR", list[0].Currency)
	assert.Equal(t, 2, list[0].Digits)
	assert.Equal(t, "Euro", list[0].DisplayName)
	assert.Equal(t, "1.0", list[0].Rate)

	r := list[0]
	r.Rate = "1.1"
	assert.NoError(t, repos.ExchangeRateRepo.Update(r))

	list, err = repos.ExchangeRateRepo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "EUR", list[0].Currency)
	assert.Equal(t, 2, list[0].Digits)
	assert.Equal(t, "Euro", list[0].DisplayName)
	assert.Equal(t, "1.1", list[0].Rate)
}
