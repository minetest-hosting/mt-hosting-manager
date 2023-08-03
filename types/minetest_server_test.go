package types_test

import (
	"mt-hosting-manager/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerNameValid(t *testing.T) {
	assert.True(t, types.ValidServerName.MatchString("abc"))
	assert.True(t, types.ValidServerName.MatchString("123"))
	assert.False(t, types.ValidServerName.MatchString("abc_"))
	assert.False(t, types.ValidServerName.MatchString("abc-"))
	assert.False(t, types.ValidServerName.MatchString(""))
}
