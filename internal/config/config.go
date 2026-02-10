package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the Muster CLI configuration
type Config struct {
	User         *User         `mapstructure:"user"`
	Organization *Organization `mapstructure:"organization"`
	Auth         *Auth         `mapstructure:"auth"`
	API          *API          `mapstructure:"api"`
}

// User represents the authenticated user
type User struct {
	ID    int    `mapstructure:"id"`
	Email string `mapstructure:"email"`
	Name  string `mapstructure:"name"`
	Role  string `mapstructure:"role"`
}

// Organization represents the user's organization
type Organization struct {
	ID     int    `mapstructure:"id"`
	Name   string `mapstructure:"name"`
	Status string `mapstructure:"status"`
}

// Auth represents authentication data
type Auth struct {
	Token string `mapstructure:"token"`
}

// API represents API configuration
type API struct {
	BaseURL string `mapstructure:"base_url"`
}

var (
	configDir  string
	configFile string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("failed to get user home directory: %v", err))
	}

	configDir = filepath.Join(homeDir, ".muster")
	configFile = filepath.Join(configDir, "config.yaml")
}

// GetConfigDir returns the config directory path
func GetConfigDir() string {
	return configDir
}

// GetConfigFile returns the config file path
func GetConfigFile() string {
	return configFile
}

// Load reads the configuration from ~/.muster/config.yaml
func Load() (*Config, error) {
	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Config file not found, return default config
		return &Config{
			API: &API{
				BaseURL: "http://localhost:3000",
			},
		}, nil
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	// Set defaults
	viper.SetDefault("api.base_url", "http://localhost:3000")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// Save writes the configuration to ~/.muster/config.yaml
func Save(cfg *Config) error {
	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	viper.Set("user", cfg.User)
	viper.Set("organization", cfg.Organization)
	viper.Set("auth", cfg.Auth)
	viper.Set("api", cfg.API)

	if err := viper.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// IsAuthenticated checks if the user is authenticated
func (c *Config) IsAuthenticated() bool {
	return c.Auth != nil && c.Auth.Token != ""
}

// Clear removes all authentication data
func (c *Config) Clear() {
	c.User = nil
	c.Organization = nil
	c.Auth = nil
}

// SetAuthData sets the authentication data from login/register response
func (c *Config) SetAuthData(token string, user *User, org *Organization) {
	c.Auth = &Auth{Token: token}
	c.User = user
	c.Organization = org
}
