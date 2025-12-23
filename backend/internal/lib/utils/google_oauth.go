package utils

import (
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuth struct {
	config *oauth2.Config
}

func NewGoogleOAuth(
	clientID string,
	clientSecret string,
	redirectURL string,
) *GoogleOAuth {
	return &GoogleOAuth{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"openid",
				"email",
				"profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (g *GoogleOAuth) AuthURL(state string) string {
	return g.config.AuthCodeURL(state)
}

func (g *GoogleOAuth) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.config.Exchange(ctx, code)
}
