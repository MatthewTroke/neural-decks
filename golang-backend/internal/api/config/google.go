package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewGoogleOAuthConfig(redirectURI string, clientID string, clientSecret string) oauth2.Config {
	return oauth2.Config{
		RedirectURL:  redirectURI,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}
}
