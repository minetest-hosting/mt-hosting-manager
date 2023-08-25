package types

type Info struct {
	BaseURL             string `json:"baseurl"`
	HostingDomainSuffix string `json:"hostingdomain_suffix"`
	Stage               string `json:"stage"`
	GithubClientID      string `json:"github_client_id"`
}

func NewInfo(cfg *Config) *Info {
	return &Info{
		BaseURL:             cfg.BaseURL,
		HostingDomainSuffix: cfg.HostingDomainSuffix,
		Stage:               cfg.Stage,
		GithubClientID:      cfg.GithubOauthConfig.ClientID,
	}
}
