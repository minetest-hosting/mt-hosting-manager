package wallee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func CreateTransaction(userID, key string, tx *TransactionRequest) (*TransactionResponse, error) {
	ts := time.Now().Unix()
	path := "/api/transaction/createTransactionCredentials"
	method := http.MethodPost

	mac, err := CreateMac(userID, key, method, path, ts)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, fmt.Sprintf("https://app-wallee.com%s", path), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-mac-version", "1")
	req.Header.Set("x-mac-userid", userID)
	req.Header.Set("x-mac-timestamp", fmt.Sprintf("%d", ts))
	req.Header.Set("x-mac-value", mac)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	txr := &TransactionResponse{}
	err = json.NewDecoder(resp.Body).Decode(txr)
	return txr, err
}

func GetTransaction(userID, key, id string) (*TransactionResponse, error) {
	return nil, nil
}

func CreatePaymentPageURL(userID, key, spaceID, transactionID string) (string, error) {
	return "", nil
}
