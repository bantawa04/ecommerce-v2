package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	// Server config
	ServerPort         string        `mapstructure:"SERVER_PORT"`
	ServerReadTimeout  time.Duration `mapstructure:"SERVER_READ_TIMEOUT"`
	ServerWriteTimeout time.Duration `mapstructure:"SERVER_WRITE_TIMEOUT"`
	
	// Database config
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE"`
	
	// ImageKit config
	ImageKitPublicKey   string `mapstructure:"IMAGEKIT_PUBLIC_KEY"`
	ImageKitPrivateKey  string `mapstructure:"IMAGEKIT_PRIVATE_KEY"`
	ImageKitURLEndpoint string `mapstructure:"IMAGEKIT_URL_ENDPOINT"`
}

// ServerConfig returns the server configuration
func (c *Config) Server() ServerConfig {
	return ServerConfig{
		Port:         c.ServerPort,
		ReadTimeout:  c.ServerReadTimeout,
		WriteTimeout: c.ServerWriteTimeout,
	}
}

// Database returns the database configuration
func (c *Config) Database() DatabaseConfig {
	return DatabaseConfig{
		Host:     c.DBHost,
		Port:     c.DBPort,
		User:     c.DBUser,
		Password: c.DBPassword,
		DBName:   c.DBName,
		SSLMode:  c.DBSSLMode,
	}
}

// ImageKit returns the ImageKit configuration
func (c *Config) ImageKit() ImageKitConfig {
	return ImageKitConfig{
		PublicKey:   c.ImageKitPublicKey,
		PrivateKey:  c.ImageKitPrivateKey,
		URLEndpoint: c.ImageKitURLEndpoint,
	}
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

// ImageKitConfig holds ImageKit-related configuration
type ImageKitConfig struct {
	PublicKey   string
	PrivateKey  string
	URLEndpoint string
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
	viper.SetDefault("IMAGEKIT_PUBLIC_KEY", "")
	viper.SetDefault("IMAGEKIT_PRIVATE_KEY", "")
	viper.SetDefault("IMAGEKIT_URL_ENDPOINT", "")

	// Enable environment variables
	viper.AutomaticEnv()

	// Try to read .env file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
		// Continue with environment variables and defaults
	}

	// Create config instance and unmarshal into it
	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
