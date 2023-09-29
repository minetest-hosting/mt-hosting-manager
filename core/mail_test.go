package core_test

import (
	"mt-hosting-manager/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripMailPlusExtension(t *testing.T) {
	assert.Equal(t, "abc@xy.com", core.StripMailPlusExtension("abc+test@xy.com"))
	assert.Equal(t, "abc@xy.com", core.StripMailPlusExtension("abc+@xy.com"))
	assert.Equal(t, "abc@xy.com", core.StripMailPlusExtension("abc@xy.com"))
}
