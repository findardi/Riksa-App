package oauth

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Google struct {
	cfg *oauth2.Config
}

func NewGoogle(clientID, clientSecret, redirectURL string) *Google {
	return &Google{
		cfg: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Endpoint:     google.Endpoint,
			Scopes:       []string{"openid", "email", "profile"},
		},
	}
}

func (g *Google) Name() string {
	return "google"
}

func (g *Google) AuthCodeURL(state string) string {
	return g.cfg.AuthCodeURL(state)
}

// Identity exchanges the code, then reads the OIDC userinfo endpoint
// (email_verified is provided directly, so no extra call is needed).
func (g *Google) Identity(ctx context.Context, code string) (Identity, error) {
	tok, err := g.cfg.Exchange(ctx, code)
	if err != nil {
		return Identity{}, fmt.Errorf("exchange code :%w", err)
	}
	client := g.cfg.Client(ctx, tok)

	var profile struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
	}

	if err := getJSON(ctx, client, "https://www.googleapis.com/oauth2/v3/userinfo", &profile); err != nil {
		return Identity{}, fmt.Errorf("fetch userinfo :%w", err)
	}

	return Identity{
		ProviderUID:   profile.Sub,
		Email:         profile.Email,
		EmailVerified: profile.EmailVerified,
		Username:      profile.Name,
	}, nil
}
