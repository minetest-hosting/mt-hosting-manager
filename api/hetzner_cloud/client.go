package hetzner_cloud

import "net/http"

type HetznerCloudClient struct {
	Key    string
	client http.Client
}

func New(key string) *HetznerCloudClient {
	return &HetznerCloudClient{
		Key:    key,
		client: http.Client{},
	}
}
