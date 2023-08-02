package hetzner_dns

import "net/http"

type HetznerDNSClient struct {
	Key    string
	zoneID string
	client http.Client
}

func New(key, zoneID string) *HetznerDNSClient {
	return &HetznerDNSClient{
		Key:    key,
		zoneID: zoneID,
		client: http.Client{},
	}
}
