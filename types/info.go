package types

type Info struct {
	BaseURL             string `json:"baseurl"`
	HostingDomainSuffix string `json:"hostingdomain_suffix"`
	Stage               string `json:"stage"`
	GithubClientID      string `json:"github_client_id"`
	DiscordClientID     string `json:"discord_client_id"`
	MesehubClientID     string `json:"mesehub_client_id"`
	MaxBalance          int    `json:"max_balance"` //max balance in eurocents
}

func NewInfo(cfg *Config) *Info {
	return &Info{
		BaseURL:             cfg.BaseURL,
		HostingDomainSuffix: cfg.HostingDomainSuffix,
		Stage:               cfg.Stage,
		GithubClientID:      cfg.GithubOauthConfig.ClientID,
		DiscordClientID:     cfg.DiscordOauthConfig.ClientID,
		MesehubClientID:     cfg.MesehubOauthConfig.ClientID,
		MaxBalance:          cfg.MaxBalance,
	}
}
