package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Github struct {
	cfg *oauth2.Config
}

func NewGithub(clientID, clientSecret, redirectURL string) *Github {
	return &Github{
		cfg: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Endpoint:     github.Endpoint,
			Scopes:       []string{"read:user", "user:email"},
		},
	}
}

func (g *Github) Name() string {
	return "github"
}

func (g *Github) AuthCodeURL(state string) string {
	return g.cfg.AuthCodeURL(state)
}

func (g *Github) Identity(ctx context.Context, code string) (Identity, error) {
	tok, err := g.cfg.Exchange(ctx, code)
	if err != nil {
		return Identity{}, fmt.Errorf("exchange code :%w", err)
	}
	client := g.cfg.Client(ctx, tok)

	var profile struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
	}

	if err := getJSON(ctx, client, "https://api.github.com/user", &profile); err != nil {
		return Identity{}, fmt.Errorf("fetch user :%w", err)
	}

	var email []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := getJSON(ctx, client, "https://api.github.com/user/emails", &email); err != nil {
		return Identity{}, fmt.Errorf("fetch email :%w", err)
	}

	id := Identity{
		ProviderUID: strconv.FormatInt(profile.ID, 10),
		Username:    profile.Login,
	}

	for _, e := range email {
		if e.Primary && e.Verified {
			id.Email = e.Email
			id.EmailVerified = true
			break
		}
	}

	return id, nil
}

func getJSON(ctx context.Context, c *http.Client, url string, dst any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(dst)
}
