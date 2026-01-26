package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	App      AppConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port        string `mapstructure:"port"`
	Environment string `mapstructure:"environment"`
	Host        string `mapstructure:"host"`
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Neo4jURI      string `mapstructure:"neo4j_uri"`
	Neo4jUsername string `mapstructure:"neo4j_username"`
	Neo4jPassword string `mapstructure:"neo4j_password"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret          string `mapstructure:"jwt_secret"`
	GoogleClientID     string `mapstructure:"google_client_id"`
	GoogleClientSecret string `mapstructure:"google_client_secret"`
	AppleClientID      string `mapstructure:"apple_client_id"`
	AppleClientSecret  string `mapstructure:"apple_client_secret"`
	SessionSecret      string `mapstructure:"session_secret"`
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	LogLevel    string `mapstructure:"log_level"`
	FrontendURL string `mapstructure:"frontend_url"`
}

// Load reads configuration from environment variables and config files
func Load() (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Try to load .env file manually first (search in multiple locations)
	envPaths := []string{".env", "../.env", "../../.env"}
	for _, envPath := range envPaths {
		if err := loadEnvFile(envPath); err == nil {
			break
		}
	}

	// Read from environment variables with GRGN_STACK prefix
	// Map environment variables to nested config structure
	v.SetEnvPrefix("GRGN_STACK")
	v.AutomaticEnv()

	// Explicitly bind environment variables to config keys
	v.BindEnv("server.port", "GRGN_STACK_SERVER_PORT")
	v.BindEnv("server.environment", "GRGN_STACK_SERVER_ENVIRONMENT")
	v.BindEnv("server.host", "GRGN_STACK_SERVER_HOST")

	v.BindEnv("database.neo4j_uri", "GRGN_STACK_DATABASE_NEO4J_URI")
	v.BindEnv("database.neo4j_username", "GRGN_STACK_DATABASE_NEO4J_USERNAME")
	v.BindEnv("database.neo4j_password", "GRGN_STACK_DATABASE_NEO4J_PASSWORD")

	v.BindEnv("auth.jwt_secret", "GRGN_STACK_AUTH_JWT_SECRET")
	v.BindEnv("auth.google_client_id", "GRGN_STACK_AUTH_GOOGLE_CLIENT_ID")
	v.BindEnv("auth.google_client_secret", "GRGN_STACK_AUTH_GOOGLE_CLIENT_SECRET")
	v.BindEnv("auth.apple_client_id", "GRGN_STACK_AUTH_APPLE_CLIENT_ID")
	v.BindEnv("auth.apple_client_secret", "GRGN_STACK_AUTH_APPLE_CLIENT_SECRET")
	v.BindEnv("auth.session_secret", "GRGN_STACK_AUTH_SESSION_SECRET")

	v.BindEnv("app.name", "GRGN_STACK_APP_NAME")
	v.BindEnv("app.version", "GRGN_STACK_APP_VERSION")
	v.BindEnv("app.log_level", "GRGN_STACK_APP_LOG_LEVEL")
	v.BindEnv("app.frontend_url", "GRGN_STACK_APP_FRONTEND_URL")

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &config, nil
}

// loadEnvFile loads environment variables from a .env file
func loadEnvFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove surrounding quotes if present
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		// Only set if not already set (env vars take precedence)
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.environment", "development")
	v.SetDefault("server.host", "0.0.0.0")

	// Database defaults
	v.SetDefault("database.neo4j_uri", "bolt://localhost:7687")
	v.SetDefault("database.neo4j_username", "neo4j")
	v.SetDefault("database.neo4j_password", "password")

	// App defaults
	v.SetDefault("app.name", "GRGN Stack")
	v.SetDefault("app.version", "0.1.0")
	v.SetDefault("app.log_level", "info")
	v.SetDefault("app.frontend_url", "http://localhost:5173")
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// IsStaging returns true if running in staging mode
func (c *Config) IsStaging() bool {
	return c.Server.Environment == "staging"
}
