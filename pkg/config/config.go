package config

import (
	"fmt"
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

	// Read from environment variables
	v.SetEnvPrefix("GRGN_STACK")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Try to read from config file
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("..")
	v.AddConfigPath("../../")

	// Read config file if it exists (optional)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; ignore error since we can use env vars
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &config, nil
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
