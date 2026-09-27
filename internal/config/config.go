package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the proxyman configuration
type Config struct {
	WorkDir     string `mapstructure:"workdir"`
	LogLevel    string `mapstructure:"loglevel"`
	IPv6        bool   `mapstructure:"ipv6"`
	SystemProxy bool   `mapstructure:"systemproxy"`
	V2RayURL    string `mapstructure:"v2ray_url"`
	ClashURL    string `mapstructure:"clash_url"`
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		WorkDir:     filepath.Join(home, ".config", "proxyman"),
		LogLevel:    "info",
		IPv6:        false,
		SystemProxy: false,
	}
}

// LoadConfig loads configuration from file or environment
func LoadConfig() *Config {
	cfg := DefaultConfig()

	// Ensure work directory exists
	os.MkdirAll(cfg.WorkDir, 0755)

	// Try to load from viper
	if viper.IsSet("workdir") {
		cfg.WorkDir = viper.GetString("workdir")
	}
	if viper.IsSet("loglevel") {
		cfg.LogLevel = viper.GetString("loglevel")
	}
	if viper.IsSet("ipv6") {
		cfg.IPv6 = viper.GetBool("ipv6")
	}
	if viper.IsSet("systemproxy") {
		cfg.SystemProxy = viper.GetBool("systemproxy")
	}
	if viper.IsSet("v2ray_url") {
		cfg.V2RayURL = viper.GetString("v2ray_url")
	}
	if viper.IsSet("clash_url") {
		cfg.ClashURL = viper.GetString("clash_url")
	}

	return cfg
}

// SaveConfig saves configuration to disk
func SaveConfig(cfg *Config) error {
	viper.Set("workdir", cfg.WorkDir)
	viper.Set("loglevel", cfg.LogLevel)
	viper.Set("ipv6", cfg.IPv6)
	viper.Set("systemproxy", cfg.SystemProxy)
	viper.Set("v2ray_url", cfg.V2RayURL)
	viper.Set("clash_url", cfg.ClashURL)

	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config", "proxyman")
	os.MkdirAll(configPath, 0755)

	if err := viper.WriteConfig(); err != nil {
		// If config file doesn't exist, create it
		if err := viper.WriteConfigAs(filepath.Join(configPath, ".proxyman.yaml")); err != nil {
			return fmt.Errorf("write config: %w", err)
		}
	}
	return nil
}

// NewSystemProxy creates a system proxy manager
func NewSystemProxy(cfg *Config) *SystemProxy {
	return &SystemProxy{cfg: cfg}
}
