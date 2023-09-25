package types

import (
	"os"
)

type OAuthConfig struct {
	ClientID string
	Secret   string
}

type Config struct {
	BaseURL             string
	HostingDomainSuffix string
	Stage               string
	Webdev              bool
	EnableWorker        bool
	JWTKey              string
	LogStreamKey        string
	LogStreamDir        string
	CookieName          string
	CookieDomain        string
	CookiePath          string
	CookieSecure        bool
	HetznerCloudKey     string
	HetznerApiKey       string
	HetznerApiZoneID    string
	GithubOauthConfig   *OAuthConfig
	MailHost            string
	MailAddress         string
	MailPassword        string
	MaxBalance          int //max balance in eurocents
}

func NewConfig() *Config {
	return &Config{
		BaseURL:             os.Getenv("BASEURL"),
		HostingDomainSuffix: os.Getenv("HOSTING_DOMAIN_SUFFIX"),
		Stage:               os.Getenv("STAGE"),
		Webdev:              os.Getenv("WEBDEV") == "true",
		EnableWorker:        os.Getenv("ENABLE_WORKER") == "true",
		JWTKey:              os.Getenv("JWT_KEY"),
		LogStreamKey:        os.Getenv("LOG_STREAM_KEY"),
		LogStreamDir:        os.Getenv("LOG_STREAM_DIR"),
		CookieName:          "mt-hosting-manager",
		CookieDomain:        os.Getenv("COOKIE_DOMAIN"),
		CookiePath:          os.Getenv("COOKIE_PATH"),
		CookieSecure:        os.Getenv("COOKIE_SECURE") == "true",
		HetznerCloudKey:     os.Getenv("HETZNER_CLOUD_KEY"),
		HetznerApiKey:       os.Getenv("HETZNER_API_KEY"),
		HetznerApiZoneID:    os.Getenv("HETZNER_API_ZONEID"),
		GithubOauthConfig: &OAuthConfig{
			ClientID: os.Getenv("GITHUB_CLIENTID"),
			Secret:   os.Getenv("GITHUB_SECRET"),
		},
		MailHost:     os.Getenv("MAIL_HOST"),
		MailAddress:  os.Getenv("MAIL_ADDRESS"),
		MailPassword: os.Getenv("MAIL_PASSWORD"),
		MaxBalance:   100 * 100,
	}
}
