package config

import (
	"cardgame/bootstrap"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	GoogleLoginConfig oauth2.Config
}

var AppConfig Config

func NewGoogleOAuthConfig(env *bootstrap.Env) oauth2.Config {
	fmt.Println("GOOGLE OAUTH REDIRECT URI", env.GoogleOAuthRedirectURI)
	fmt.Println(env.GoogleOAuthClientID)
	fmt.Println(env.GoogleOAuthClientSecret)

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
