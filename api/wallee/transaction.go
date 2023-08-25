package wallee

import (
	"fmt"
	"io"
	"net/http"
)

func (c *WalleeClient) CreateTransaction(tx *TransactionRequest) (*TransactionResponse, error) {
	path := fmt.Sprintf("/api/transaction/create?spaceId=%s", c.SpaceID)
	txr := &TransactionResponse{}
	err := c.jsonRequest(path, http.MethodPost, tx, &txr)
	return txr, err
}

func (c *WalleeClient) SearchTransaction(txreq *TransactionSearchRequest) ([]*TransactionResponse, error) {
	path := fmt.Sprintf("/api/transaction/search?spaceId=%s", c.SpaceID)
	txr := []*TransactionResponse{}
	err := c.jsonRequest(path, http.MethodPost, txreq, &txr)
	return txr, err
}

func (c *WalleeClient) CreatePaymentPageURL(transactionID int64) (string, error) {
	path := fmt.Sprintf("/api/transaction-payment-page/payment-page-url?spaceId=%s&id=%d", c.SpaceID, transactionID)

	req, err := c.request(path, http.MethodGet, nil)
	if err != nil {
		return "", err
	}

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
