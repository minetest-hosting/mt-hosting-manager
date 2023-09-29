package coinbase

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *CoinbaseClient) GetCurrencies() (*Currencies, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.coinbase.com/v2/currencies", nil)
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

	cu := &Currencies{}
	err = json.NewDecoder(resp.Body).Decode(cu)

	return cu, err
}

func (c *CoinbaseClient) GetCryptoCurrencies() (*CryptoCurrencies, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.coinbase.com/v2/currencies/crypto", nil)
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

	cu := &CryptoCurrencies{}
	err = json.NewDecoder(resp.Body).Decode(cu)

	return cu, err
}
