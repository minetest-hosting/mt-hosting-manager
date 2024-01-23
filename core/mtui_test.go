package core_test

import (
	"mt-hosting-manager/core"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMTUIDigest(t *testing.T) {
	d, err := core.GetMTUIDigest()
	assert.NoError(t, err)
	assert.NotNil(t, d)
	assert.True(t, d != "")
	assert.True(t, strings.HasPrefix(d, "sha256:"))
}
