package config

import "os"

// EnvConfig holds all environment configuration
type EnvConfig struct {
	// Database
	DatabaseURL      string
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// JWT
	JWTVerifySecret string

	// OAuth
	GoogleOAuthClientID      string
	GoogleOAuthClientSecret  string
	GoogleOAuthRedirectURI   string
	DiscordOAuthClientID     string
	DiscordOAuthClientSecret string
	DiscordOAuthRedirectURI  string

	// App
	AppEnv            string
	LocalDevBypass    bool
	LocalDevUserID    string
	LocalDevUserName  string
	LocalDevUserEmail string
}

// NewEnvConfig creates a new environment configuration
func NewEnvConfig() *EnvConfig {
	return &EnvConfig{
		// Initialize with default values or load from environment
		DatabaseURL:      getEnvOrDefault("DATABASE_URL", ""),
		DatabaseHost:     getEnvOrDefault("DATABASE_HOST", "localhost"),
		DatabasePort:     getEnvOrDefault("DATABASE_PORT", "5432"),
		DatabaseName:     getEnvOrDefault("DATABASE_NAME", "cardgame"),
		DatabaseUser:     getEnvOrDefault("DATABASE_USER", "postgres"),
		DatabasePassword: getEnvOrDefault("DATABASE_PASSWORD", ""),

		RedisHost:     getEnvOrDefault("REDIS_HOST", "localhost"),
		RedisPort:     getEnvOrDefault("REDIS_PORT", "6379"),
		RedisPassword: getEnvOrDefault("REDIS_PASSWORD", ""),
		RedisDB:       0,

		JWTVerifySecret: getEnvOrDefault("JWT_VERIFY_SECRET", "your-secret-key"),

		GoogleOAuthClientID:     getEnvOrDefault("GOOGLE_OAUTH_CLIENT_ID", ""),
		GoogleOAuthClientSecret: getEnvOrDefault("GOOGLE_OAUTH_CLIENT_SECRET", ""),
		GoogleOAuthRedirectURI:  getEnvOrDefault("GOOGLE_OAUTH_REDIRECT_URI", ""),

		DiscordOAuthClientID:     getEnvOrDefault("DISCORD_OAUTH_CLIENT_ID", ""),
		DiscordOAuthClientSecret: getEnvOrDefault("DISCORD_OAUTH_CLIENT_SECRET", ""),
		DiscordOAuthRedirectURI:  getEnvOrDefault("DISCORD_OAUTH_REDIRECT_URI", ""),

		AppEnv:            getEnvOrDefault("APP_ENV", "development"),
		LocalDevBypass:    getEnvOrDefault("LOCAL_DEV_BYPASS", "false") == "true",
		LocalDevUserID:    getEnvOrDefault("LOCAL_DEV_USER_ID", "dev-user"),
		LocalDevUserName:  getEnvOrDefault("LOCAL_DEV_USER_NAME", "Dev User"),
		LocalDevUserEmail: getEnvOrDefault("LOCAL_DEV_USER_EMAIL", "dev@example.com"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
