package oauth

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type OAuthConfig struct {
	ClientID string
	Secret   string
}

type OauthUserInfo struct {
	Name       string
	Email      string
	ExternalID string
}

type OauthImplementation interface {
	RequestAccessToken(code, baseurl string, cfg *OAuthConfig) (string, error)
	RequestUserInfo(access_token string, cfg *OAuthConfig) (*OauthUserInfo, error)
}
