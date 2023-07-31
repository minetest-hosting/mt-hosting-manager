package wallee

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
)

func CreateMac(userID, key, method, path string, ts int64) (string, error) {
	str := fmt.Sprintf("%d|%s|%d|%s|%s", 1, userID, ts, method, path)
	decoded_secret, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}

	hashfn := hmac.New(sha512.New, decoded_secret)
	hashfn.Write([]byte(str))
	hash := hashfn.Sum(nil)

	return base64.StdEncoding.EncodeToString(hash), nil
}
