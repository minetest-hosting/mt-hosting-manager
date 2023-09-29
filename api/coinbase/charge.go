package coinbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *CoinbaseClient) CreateCharge(charge *CreateChargeRequest) (*Charge, error) {
	data, err := json.Marshal(charge)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.commerce.coinbase.com/charges/", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-CC-Api-Key", c.Key)
	req.Header.Set("X-CC-Version", "2018-03-22")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response-code: %d", resp.StatusCode)
	}
	cr := &Charge{}
	err = json.NewDecoder(resp.Body).Decode(cr)

	return cr, err
}

func (c *CoinbaseClient) GetCharge(code string) (*Charge, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.commerce.coinbase.com/charges/%s", code), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-CC-Version", "2018-03-22")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response-code: %d", resp.StatusCode)
	}

	cr := &Charge{}
	err = json.NewDecoder(resp.Body).Decode(cr)

	return cr, err
}
