package core

import (
	"io"

	openssl "github.com/Luzifer/go-openssl/v4"
)

func StreamBackup(passphrase string, src io.Reader, dst io.Writer) (int64, error) {
	r := openssl.NewReader(src, passphrase, openssl.PBKDF2SHA256)
	return io.Copy(dst, r)
}
