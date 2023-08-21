package types

type Info struct {
	BaseURL             string
	HostingDomainSuffix string
	Stage               string
	GithubClientID      string
}

func NewInfo(cfg *Config) *Info {
	return &Info{
		BaseURL:             cfg.BaseURL,
		HostingDomainSuffix: cfg.HostingDomainSuffix,
		Stage:               cfg.Stage,
		GithubClientID:      cfg.GithubOauthConfig.ClientID,
	}
}
