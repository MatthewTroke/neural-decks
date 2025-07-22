package config

import (
	"cardgame/bootstrap"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	GoogleLoginConfig oauth2.Config
}

var AppConfig Config

func NewGoogleOAuthConfig(env *bootstrap.Env) oauth2.Config {
	AppConfig.GoogleLoginConfig = oauth2.Config{
		RedirectURL:  env.GoogleOAuthRedirectURI,
		ClientID:     env.GoogleOAuthClientID,
		ClientSecret: env.GoogleOAuthClientSecret,
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	return AppConfig.GoogleLoginConfig
}
