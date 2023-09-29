package coinbase

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *CoinbaseClient) GetRates(currency string) (*RatesResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.coinbase.com/v2/exchange-rates?currency=%s", currency), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response-code: %d", resp.StatusCode)
	}

	rr := &RatesResponse{}
	err = json.NewDecoder(resp.Body).Decode(rr)

	return rr, err
}
