package types

import "os"

type Config struct {
	BaseURL             string
	HostingDomainSuffix string
	Stage               string
	EnableWorker        bool
	EnableDummyWorker   bool
	JWTKey              string
	CookieName          string
	CookieDomain        string
	CookiePath          string
	CookieSecure        bool
}

func NewConfig() *Config {
	return &Config{
		BaseURL:             os.Getenv("BASEURL"),
		HostingDomainSuffix: os.Getenv("HOSTING_DOMAIN_SUFFIX"),
		Stage:               os.Getenv("STAGE"),
		EnableWorker:        os.Getenv("ENABLE_WORKER") == "true",
		EnableDummyWorker:   os.Getenv("ENABLE_DUMMY_WORKER") == "true",
		JWTKey:              os.Getenv("JWT_KEY"),
		CookieName:          "mt-hosting-manager",
		CookieDomain:        os.Getenv("COOKIE_DOMAIN"),
		CookiePath:          os.Getenv("COOKIE_PATH"),
		CookieSecure:        os.Getenv("COOKIE_SECURE") == "true",
	}
}
