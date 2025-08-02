package bootstrap

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
	LocalDevUserID           string `mapstructure:"LOCAL_DEV_USER_ID"`
	LocalDevUserName         string `mapstructure:"LOCAL_DEV_USER_NAME"`
	LocalDevUserEmail        string `mapstructure:"LOCAL_DEV_USER_EMAIL"`
}

func NewEnv() *Env {
	env := Env{}

	// Try to read from .env file first
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("No .env file found, using environment variables")
		// If no .env file, use environment variables
		viper.AutomaticEnv()
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	// Set defaults if not provided
	if env.RedisHost == "" {
		env.RedisHost = "localhost"
	}
	if env.RedisPort == "" {
		env.RedisPort = "6379"
	}

	// Debug: Print environment variables
	log.Printf("🔍 Environment - RedisHost: %s, RedisPort: %s", env.RedisHost, env.RedisPort)
	if env.DatabaseDSN == "" {
		env.DatabaseDSN = "postgresql://postgres:password@localhost:5432/neural_decks"
	}
	if env.JWTVerifySecret == "" {
		env.JWTVerifySecret = "your-secret-key-here"
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	// Set defaults for local development bypass
	if env.LocalDevUserID == "" {
		env.LocalDevUserID = "local-dev-user"
	}
	if env.LocalDevUserName == "" {
		env.LocalDevUserName = "Local Developer"
	}
	if env.LocalDevUserEmail == "" {
		env.LocalDevUserEmail = "dev@localhost"
	}

	return &env
}
