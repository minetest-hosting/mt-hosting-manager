package core_test

import (
	"bytes"
	"io"
	"mt-hosting-manager/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryption(t *testing.T) {
	b := []byte("hello world")
	key := core.RandStringRunes(16)

	// write encrypted
	encrypted := bytes.NewBuffer([]byte{})
	w, err := core.EncryptedWriter(key, encrypted)
	assert.NoError(t, err)
	_, err = io.Copy(w, bytes.NewReader(b))
	assert.NoError(t, err)

	// read encrypted
	r, err := core.EncryptedReader(key, encrypted)
	assert.NoError(t, err)
	b2 := bytes.NewBuffer([]byte{})
	_, err = io.Copy(b2, r)
	assert.NoError(t, err)

	assert.Equal(t, b, b2.Bytes())
}
