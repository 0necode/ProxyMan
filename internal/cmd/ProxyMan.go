package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	version = "1.1.0"
)

// rootCmd is the base command
var rootCmd = &cobra.Command{
	Use:   "ProxyMan",
	Short: "A unified proxy CLI for Linux — V2Ray & Clash protocols",
	Long: `ProxyMan is a unified proxy CLI that manages V2Ray and Clash protocols
on Linux, enabling local system access to the internet.

Supports: V2Ray (vmess, vless, vless-reality, vless-grpc, trojan, trojan-grpc, ss, ss2022, anytls, socks, http)
          Clash (vmess, vless, trojan, ss, http, https, socks5, tuic, hysteria, anytls)`,
	Version: version,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.SetVersionTemplate("{{.Version}}\n")

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ProxyMan.yaml)")
	rootCmd.Flags().BoolP("verbose", "v", false, "verbose output")

	// Subcommands are registered in commands.go
}

// initConfig reads in config file and ENV variables
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".ProxyMan" (without extension)
		viper.AddConfigPath(filepath.Join(home, ".config", "ProxyMan"))
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ProxyMan")
	}

	viper.SetEnvPrefix("PROXYMAN")
	viper.BindEnv("workdir", "PROXYMAN_WORKDIR")
	viper.BindEnv("loglevel", "PROXYMAN_LOGLEVEL")
	viper.BindEnv("ipv6", "PROXYMAN_IPV6")
	viper.BindEnv("systemproxy", "PROXYMAN_SYSTEMPROXY")
	viper.BindEnv("v2ray_url", "PROXYMAN_V2RAY_URL")
	viper.BindEnv("clash_url", "PROXYMAN_CLASH_URL")

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
