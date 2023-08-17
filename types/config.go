package types

import "os"

type Config struct {
	BaseURL             string
	HostingDomainSuffix string
	Stage               string
	MockOrchestration   bool
}

func NewConfig() *Config {
	return &Config{
		BaseURL:             os.Getenv("BASEURL"),
		HostingDomainSuffix: os.Getenv("HOSTING_DOMAIN_SUFFIX"),
		Stage:               os.Getenv("STAGE"),
		MockOrchestration:   os.Getenv("MOCK_ORCHESTRATION") == "true",
	}
}
