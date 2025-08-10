package config

import (
	"golang.org/x/oauth2"
)

func NewDiscordOAuthConfig(redirectURI string, clientID string, clientSecret string) oauth2.Config {
	return oauth2.Config{
		RedirectURL:  redirectURI,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"identify", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
	}
}
