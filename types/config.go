package types

import (
	"os"
	"strings"
)

type OAuthConfig struct {
	ClientID string
	Secret   string
}

type Config struct {
	BaseURL             string
	HostingDomainSuffix string
	ReservedPrefixes    []string
	Stage               string
	SignupWhitelist     []string
	InitialBalance      string
	Webdev              bool
	EnableWorker        bool
	JWTKey              string
	CoinbaseKey         string
	CoinbaseEnabled     bool
	LogStreamKey        string
	LogStreamDir        string
	CookieName          string
	CookiePath          string
	CookieSecure        bool
	HetznerCloudKey     string
	HetznerApiKey       string
	HetznerApiZoneID    string
	WalleeUserID        string
	WalleeSpaceID       string
	WalleeKey           string
	WalleeEnabled       bool
	GithubOauthConfig   *OAuthConfig
	DiscordOauthConfig  *OAuthConfig
	MesehubOauthConfig  *OAuthConfig
	MailHost            string
	MailAddress         string
	MailPassword        string
	MaxBalance          int //max balance in eurocents
	S3Endpoint          string
	S3Bucket            string
	S3KeyID             string
	S3AccessKey         string
}

func NewConfig() *Config {
	return &Config{
		BaseURL:             os.Getenv("BASEURL"),
		HostingDomainSuffix: os.Getenv("HOSTING_DOMAIN_SUFFIX"),
		ReservedPrefixes:    strings.Split(os.Getenv("RESERVED_PREFIXES"), ","),
		Stage:               os.Getenv("STAGE"),
		SignupWhitelist:     strings.Split(os.Getenv("SIGNUP_WHITELIST"), ","),
		InitialBalance:      os.Getenv("INITIAL_BALANCE"),
		Webdev:              os.Getenv("WEBDEV") == "true",
		EnableWorker:        os.Getenv("ENABLE_WORKER") == "true",
		JWTKey:              os.Getenv("JWT_KEY"),
		CoinbaseKey:         os.Getenv("COINBASE_KEY"),
		CoinbaseEnabled:     os.Getenv("COINBASE_ENABLED") == "true",
		LogStreamKey:        os.Getenv("LOG_STREAM_KEY"),
		LogStreamDir:        os.Getenv("LOG_STREAM_DIR"),
		CookieName:          "mt-hosting-manager",
		CookiePath:          os.Getenv("COOKIE_PATH"),
		CookieSecure:        os.Getenv("COOKIE_SECURE") == "true",
		HetznerCloudKey:     os.Getenv("HETZNER_CLOUD_KEY"),
		HetznerApiKey:       os.Getenv("HETZNER_API_KEY"),
		HetznerApiZoneID:    os.Getenv("HETZNER_API_ZONEID"),
		WalleeUserID:        os.Getenv("WALLEE_USERID"),
		WalleeSpaceID:       os.Getenv("WALLEE_SPACEID"),
		WalleeKey:           os.Getenv("WALLEE_KEY"),
		WalleeEnabled:       os.Getenv("WALLEE_ENABLED") == "true",
		GithubOauthConfig: &OAuthConfig{
			ClientID: os.Getenv("GITHUB_CLIENTID"),
			Secret:   os.Getenv("GITHUB_SECRET"),
		},
		DiscordOauthConfig: &OAuthConfig{
			ClientID: os.Getenv("DISCORD_CLIENTID"),
			Secret:   os.Getenv("DISCORD_SECRET"),
		},
		MesehubOauthConfig: &OAuthConfig{
			ClientID: os.Getenv("MESEHUB_CLIENTID"),
			Secret:   os.Getenv("MESEHUB_SECRET"),
		},
		MailHost:     os.Getenv("MAIL_HOST"),
		MailAddress:  os.Getenv("MAIL_ADDRESS"),
		MailPassword: os.Getenv("MAIL_PASSWORD"),
		MaxBalance:   100 * 100,
		S3Endpoint:   os.Getenv("S3_ENDPOINT"),
		S3Bucket:     os.Getenv("S3_BUCKET"),
		S3KeyID:      os.Getenv("S3_KEY_ID"),
		S3AccessKey:  os.Getenv("S3_ACCESS_KEY"),
	}
}
