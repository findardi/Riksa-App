package config

type OAuthProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// LoadOAuth reads {prefix}_CLIENT_ID / _CLIENT_SECRET / _REDIRECT_URL.
func LoadOAuth(prefix string) OAuthProviderConfig {
	return OAuthProviderConfig{
		ClientID:     GetEnv(prefix+"_CLIENT_ID", ""),
		ClientSecret: GetEnv(prefix+"_CLIENT_SECRET", ""),
		RedirectURL:  GetEnv(prefix+"_REDIRECT_URL", ""),
	}
}
