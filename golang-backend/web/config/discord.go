package config

import (
	"cardgame/bootstrap"

	"golang.org/x/oauth2"
)

func NewDiscordOAuthConfig(env *bootstrap.Env) oauth2.Config {
	return oauth2.Config{
		RedirectURL:  env.DiscordOAuthRedirectURI,
		ClientID:     env.DiscordOAuthClientID,
		ClientSecret: env.DiscordOAuthClientSecret,
		Scopes:       []string{"identify", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
	}
}
