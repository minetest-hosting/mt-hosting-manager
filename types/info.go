package types

type Info struct {
	BaseURL             string `json:"baseurl"`
	HostingDomainSuffix string `json:"hostingdomain_suffix"`
	Stage               string `json:"stage"`
	CoinbaseEnabled     bool   `json:"coinbase_enabled"`
	WalleeEnabled       bool   `json:"wallee_enabled"`
	GithubLogin         string `json:"github_login"`
	DiscordLogin        string `json:"discord_login"`
	MesehubLogin        string `json:"mesehub_login"`
	CDBLogin            string `json:"cdb_login"`
	MaxBalance          int    `json:"max_balance"` //max balance in eurocents
}

func NewInfo(cfg *Config) *Info {
	return &Info{
		BaseURL:             cfg.BaseURL,
		HostingDomainSuffix: cfg.HostingDomainSuffix,
		Stage:               cfg.Stage,
		CoinbaseEnabled:     cfg.CoinbaseEnabled,
		WalleeEnabled:       cfg.WalleeEnabled,
		GithubLogin:         cfg.GithubOauthConfig.LoginURL,
		DiscordLogin:        cfg.DiscordOauthConfig.LoginURL,
		MesehubLogin:        cfg.MesehubOauthConfig.LoginURL,
		CDBLogin:            cfg.CDBOauthConfig.LoginURL,
		MaxBalance:          cfg.MaxBalance,
	}
}
