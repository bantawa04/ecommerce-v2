package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadConfig loads configuration from environment variables and .env files
func LoadConfig() (*Config, error) {
	// Configure Viper to read from .env file
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set default values
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_READ_TIMEOUT", "10s")
	viper.SetDefault("SERVER_WRITE_TIMEOUT", "10s")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_NAME", "beautyessentials")
	viper.SetDefault("DB_SSL_MODE", "allow")

	// Enable environment variables
	viper.AutomaticEnv()

	// Try to read .env file (optional)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
		// Continue with environment variables and defaults
	}

	// Parse the configuration
	var config Config
	config.Server = ServerConfig{
		Port:         viper.GetString("SERVER_PORT"),
		ReadTimeout:  viper.GetDuration("SERVER_READ_TIMEOUT"),
		WriteTimeout: viper.GetDuration("SERVER_WRITE_TIMEOUT"),
	}

	config.Database = DatabaseConfig{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		DBName:   viper.GetString("DB_NAME"),
		SSLMode:  viper.GetString("DB_SSL_MODE"),
	}

	return &config, nil
}
