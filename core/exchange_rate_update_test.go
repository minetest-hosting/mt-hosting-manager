package core_test

import (
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExchangeRateUpdate(t *testing.T) {
	repos := SetupRepos(t)
	cfg := types.NewConfig()
	c := core.New(repos, cfg)

	assert.NotNil(t, c)

	assert.NoError(t, c.UpdateExchangeRates())
}
