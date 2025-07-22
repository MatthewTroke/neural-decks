package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                  string `mapstructure:"APP_ENV"`
	GoogleOAuthClientID     string `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleOAuthRedirectURI  string `mapstructure:"GOOGLE_OAUTH_REDIRECT_URI"`
	GoogleOAuthClientSecret string `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
	DatabaseDSN             string `mapstructure:"DATABASE_DSN"`
	TemporalHostPort        string `mapstructure:"TEMPORAL_HOST_PORT"`
	TemporalNamespace       string `mapstructure:"TEMPORAL_NAMESPACE"`
	ChatGPTAPIKey           string `mapstructure:"CHATGPT_API_KEY"`
	JWTVerifySecret         string `mapstructure:"JWT_VERIFY_SECRET"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
