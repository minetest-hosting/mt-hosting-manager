package types_test

import (
	"mt-hosting-manager/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserNodeNameValid(t *testing.T) {
	assert.True(t, types.ValidUserNodeName.MatchString("abc"))
	assert.True(t, types.ValidUserNodeName.MatchString("123"))
	assert.False(t, types.ValidUserNodeName.MatchString("abc_"))
	assert.False(t, types.ValidUserNodeName.MatchString("abc-"))
	assert.False(t, types.ValidUserNodeName.MatchString(""))
}
