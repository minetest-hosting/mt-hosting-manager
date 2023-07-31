package wallee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (c *WalleeClient) CreateTransaction(tx *TransactionRequest) (*TransactionResponse, error) {
	ts := time.Now().Unix()
	path := fmt.Sprintf("/api/transaction/create?spaceId=%s", c.SpaceID)
	method := http.MethodPost

	mac, err := CreateMac(c.UserID, c.Key, method, path, ts)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://app-wallee.com%s", path)
	fmt.Printf("Url: '%s', Data: %s\n", url, data)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-mac-version", "1")
	req.Header.Set("x-mac-userid", c.UserID)
	req.Header.Set("x-mac-timestamp", fmt.Sprintf("%d", ts))
	req.Header.Set("x-mac-value", mac)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	resp_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Response: %s\n", resp_bytes)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("api-response status: %d", resp.StatusCode)
	}

	txr := &TransactionResponse{}
	err = json.Unmarshal(resp_bytes, txr)
	return txr, err
}

func SearchTransaction(userID, key string, filter TransactionSearchFilter) ([]*TransactionResponse, error) {
	return nil, nil
}

func CreatePaymentPageURL(userID, key, spaceID, transactionID string) (string, error) {
	return "", nil
}
