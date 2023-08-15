package core_test

import (
	"embed"
	"mt-hosting-manager/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/template.txt
var Files embed.FS

func TestTemplateFile(t *testing.T) {
	data, err := core.TemplateFile(Files, "testdata/template.txt", map[string]string{"Stuff": "xy"})
	assert.NoError(t, err)
	assert.Equal(t, "my-xy", string(data))
}
