package oauth

import "context"

type Identity struct {
	ProviderUID   string
	Email         string
	EmailVerified bool
	Username      string
}

type Provider interface {
	Name() string
	AuthCodeURL(state string) string
	Identity(ctx context.Context, code string) (Identity, error)
}
