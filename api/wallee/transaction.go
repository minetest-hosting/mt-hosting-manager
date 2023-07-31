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

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-mac-version", "1")
	req.Header.Set("x-mac-userid", c.UserID)
	req.Header.Set("x-mac-timestamp", fmt.Sprintf("%d", ts))
	req.Header.Set("x-mac-value", mac)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	resp_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("api-response status: %d", resp.StatusCode)
	}

	txr := &TransactionResponse{}
	err = json.Unmarshal(resp_bytes, txr)
	return txr, err
}

func (c *WalleeClient) SearchTransaction(txreq *TransactionSearchRequest) ([]*TransactionResponse, error) {
	ts := time.Now().Unix()
	path := fmt.Sprintf("/api/transaction/search?spaceId=%s", c.SpaceID)
	method := http.MethodPost

	mac, err := CreateMac(c.UserID, c.Key, method, path, ts)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(txreq)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://app-wallee.com%s", path)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-mac-version", "1")
	req.Header.Set("x-mac-userid", c.UserID)
	req.Header.Set("x-mac-timestamp", fmt.Sprintf("%d", ts))
	req.Header.Set("x-mac-value", mac)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	resp_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Response: %s\n", resp_bytes)
		return nil, fmt.Errorf("api-response status: %d", resp.StatusCode)
	}

	txr := []*TransactionResponse{}
	err = json.Unmarshal(resp_bytes, &txr)
	return txr, err
}

func (c *WalleeClient) CreatePaymentPageURL(transactionID int64) (string, error) {
	ts := time.Now().Unix()
	path := fmt.Sprintf("/api/transaction-payment-page/payment-page-url?spaceId=%s&id=%d", c.SpaceID, transactionID)
	method := http.MethodGet

	mac, err := CreateMac(c.UserID, c.Key, method, path, ts)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://app-wallee.com%s", path)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("x-mac-version", "1")
	req.Header.Set("x-mac-userid", c.UserID)
	req.Header.Set("x-mac-timestamp", fmt.Sprintf("%d", ts))
	req.Header.Set("x-mac-value", mac)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}

	resp_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(resp_bytes), nil
}
