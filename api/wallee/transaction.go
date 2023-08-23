package wallee

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func (c *WalleeClient) CreateTransaction(tx *TransactionRequest) (*TransactionResponse, error) {
	path := fmt.Sprintf("/api/transaction/create?spaceId=%s", c.SpaceID)
	txr := &TransactionResponse{}
	err := c.request(path, http.MethodPost, tx, &txr)
	return txr, err
}

func (c *WalleeClient) SearchTransaction(txreq *TransactionSearchRequest) ([]*TransactionResponse, error) {
	path := fmt.Sprintf("/api/transaction/search?spaceId=%s", c.SpaceID)
	txr := []*TransactionResponse{}
	err := c.request(path, http.MethodPost, txreq, &txr)
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
