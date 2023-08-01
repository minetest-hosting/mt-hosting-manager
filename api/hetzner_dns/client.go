package hetzner_dns

type HetznerDNSClient struct {
	Key string
}

func New(key string) *HetznerDNSClient {
	return &HetznerDNSClient{
		Key: key,
	}
}
