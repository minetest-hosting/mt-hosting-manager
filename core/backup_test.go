package core_test

import (
	"bytes"
	"mt-hosting-manager/core"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamBackup(t *testing.T) {
	// echo hello world | openssl enc -aes-256-cbc -pbkdf2 -md sha256 -pass pass:enter > my.text.enc
	f, err := os.Open("testdata/my.text.enc")
	assert.NoError(t, err)
	assert.NotNil(t, f)

	buf := bytes.NewBuffer([]byte{})

	count, err := core.DecryptBackup("enter", f, buf)
	assert.NoError(t, err)
	assert.True(t, count > 0)

	assert.Equal(t, "hello world\n", buf.String())
}
