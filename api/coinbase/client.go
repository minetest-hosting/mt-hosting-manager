package coinbase

import (
	"net/http"
)

type CoinbaseClient struct {
	Key    string
	client http.Client
}

func New(key string) *CoinbaseClient {
	return &CoinbaseClient{
		Key:    key,
		client: http.Client{},
	}
}
