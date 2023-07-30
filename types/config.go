package types

import "os"

type Config struct {
	BaseURL             string
	HostingDomainSuffix string
}

func NewConfig() *Config {
	return &Config{
		BaseURL:             os.Getenv("BASEURL"),
		HostingDomainSuffix: os.Getenv("HOSTING_DOMAIN_SUFFIX"),
	}
}
