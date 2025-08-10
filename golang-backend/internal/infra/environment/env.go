package environment

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                   string `mapstructure:"APP_ENV"`
	GoogleOAuthClientID      string `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleOAuthRedirectURI   string `mapstructure:"GOOGLE_OAUTH_REDIRECT_URI"`
	GoogleOAuthClientSecret  string `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
	DiscordOAuthClientID     string `mapstructure:"DISCORD_OAUTH_CLIENT_ID"`
	DiscordOAuthRedirectURI  string `mapstructure:"DISCORD_OAUTH_REDIRECT_URI"`
	DiscordOAuthClientSecret string `mapstructure:"DISCORD_OAUTH_CLIENT_SECRET"`
	DatabaseDSN              string `mapstructure:"DATABASE_DSN"`
	TemporalHostPort         string `mapstructure:"TEMPORAL_HOST_PORT"`
	TemporalNamespace        string `mapstructure:"TEMPORAL_NAMESPACE"`
	ChatGPTAPIKey            string `mapstructure:"CHATGPT_API_KEY"`
	JWTVerifySecret          string `mapstructure:"JWT_VERIFY_SECRET"`
	RedisHost                string `mapstructure:"REDIS_HOST"`
	RedisPort                string `mapstructure:"REDIS_PORT"`
	RedisPassword            string `mapstructure:"REDIS_PASSWORD"`
	RedisDB                  int    `mapstructure:"REDIS_DB"`
	LocalDevBypass           bool   `mapstructure:"LOCAL_DEV_BYPASS"`
}

func NewEnv() *Env {
	env := Env{}

	// Try to read from .env file first
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic("No .env file found")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	// Debug: Print environment variables
	log.Println("✅ Environment Loaded")

	return &env
}
