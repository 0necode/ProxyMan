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
	version = "1.0.0"
)

// rootCmd is the base command
var rootCmd = &cobra.Command{
	Use:   "proxy-cli",
	Short: "A unified proxy CLI for Linux — V2Ray & Clash protocols",
	Long: `proxy-cli is a unified proxy CLI that manages V2Ray and Clash protocols
on Linux, enabling local system access to the internet.

Supports: V2Ray (vmess, vless, trojan, ss, socks, http, shadowsocks)
          Clash (vmess, vless, trojan, ss, http, https, socks5)`,
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.proxy-cli.yaml)")
	rootCmd.Flags().BoolP("verbose", "v", false, "verbose output")

	// Add subcommands
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(systemProxyCmd)
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

		// Search config in home directory with name ".proxy-cli" (without extension)
		viper.AddConfigPath(filepath.Join(home, ".config", "proxy-cli"))
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".proxy-cli")
	}

	viper.SetEnvPrefix("PROXYCLI")
	viper.BindEnv("workdir", "PROXYCLI_WORKDIR")
	viper.BindEnv("loglevel", "PROXYCLI_LOGLEVEL")
	viper.BindEnv("ipv6", "PROXYCLI_IPV6")
	viper.BindEnv("systemproxy", "PROXYCLI_SYSTEMPROXY")
	viper.BindEnv("v2ray_url", "PROXYCLI_V2RAY_URL")
	viper.BindEnv("clash_url", "PROXYCLI_CLASH_URL")

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
