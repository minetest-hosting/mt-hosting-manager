package wallee

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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

func (c *WalleeClient) request(path, method string, req_obj, resp_obj any) error {
	ts := time.Now().Unix()

	mac, err := CreateMac(c.UserID, c.Key, method, path, ts)
	if err != nil {
		return err
	}

	data, err := json.Marshal(req_obj)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://app-wallee.com%s", path)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("x-mac-version", "1")
	req.Header.Set("x-mac-userid", c.UserID)
	req.Header.Set("x-mac-timestamp", fmt.Sprintf("%d", ts))
	req.Header.Set("x-mac-value", mac)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	resp_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("api-response status: %d", resp.StatusCode)
	}

	err = json.Unmarshal(resp_bytes, resp_obj)
	return err
}
