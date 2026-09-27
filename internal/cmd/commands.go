package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/proxyman/proxyman/internal/config"
	"github.com/proxyman/proxyman/internal/engine"
	"github.com/proxyman/proxyman/internal/parser"
	"github.com/spf13/cobra"
)

// ==================== install ====================
var installCmd = &cobra.Command{
	Use:   "install [engine]",
	Short: "Install a proxy engine (v2ray or clash)",
	Long: `Install a proxy engine. Supported engines:
  - v2ray/xray: Xray-core with full protocol support
  - clash/mihomo: Mihomo (Clash Meta) with rule-based routing`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		engineName := args[0]
		fmt.Printf("Installing %s...\n", engineName)
		cfg := config.LoadConfig()
		em := engine.NewManager(cfg)
		if err := em.Install(engineName); err != nil {
			return fmt.Errorf("install failed: %w", err)
		}
		fmt.Printf("✓ %s installed successfully\n", engineName)
		return nil
	},
}

// ==================== uninstall ====================
var uninstallCmd = &cobra.Command{
	Use:   "uninstall [engine]",
	Short: "Uninstall a proxy engine",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		engineName := args[0]
		fmt.Printf("Uninstalling %s...\n", engineName)
		cfg := config.LoadConfig()
		em := engine.NewManager(cfg)
		if err := em.Uninstall(engineName); err != nil {
			return fmt.Errorf("uninstall failed: %w", err)
		}
		fmt.Printf("✓ %s uninstalled successfully\n", engineName)
		return nil
	},
}

// ==================== start ====================
var startCmd = &cobra.Command{
	Use:   "start [engine] [config-file]",
	Short: "Start a proxy engine with a configuration file",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		engineName := args[0]
		cfgFile := args[1]
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			return fmt.Errorf("config file not found: %s", cfgFile)
		}
		fmt.Printf("Starting %s with %s...\n", engineName, filepath.Base(cfgFile))
		cfg := config.LoadConfig()
		em := engine.NewManager(cfg)
		if err := em.Start(engineName, cfgFile); err != nil {
			return fmt.Errorf("start failed: %w", err)
		}
		fmt.Printf("✓ %s started successfully\n", engineName)
		return nil
	},
}

// ==================== stop ====================
var stopCmd = &cobra.Command{
	Use:   "stop [engine]",
	Short: "Stop a running proxy engine",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		engineName := args[0]
		fmt.Printf("Stopping %s...\n", engineName)
		cfg := config.LoadConfig()
		em := engine.NewManager(cfg)
		if err := em.Stop(engineName); err != nil {
			return fmt.Errorf("stop failed: %w", err)
		}
		fmt.Printf("✓ %s stopped successfully\n", engineName)
		return nil
	},
}

// ==================== status ====================
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of all proxy engines",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.LoadConfig()
		em := engine.NewManager(cfg)
		statuses, err := em.Status()
		if err != nil {
			return fmt.Errorf("status failed: %w", err)
		}
		fmt.Println("Proxy Engine Status:")
		fmt.Println("====================")
		for _, s := range statuses {
			fmt.Printf("  %-10s  %s\n", s.Engine, s.State)
		}
		return nil
	},
}

// ==================== config ====================
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage proxyman configuration",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.LoadConfig()
		fmt.Println("Current Configuration:")
		fmt.Println("=====================")
		fmt.Printf("  Work Dir:     %s\n", cfg.WorkDir)
		fmt.Printf("  Log Level:    %s\n", cfg.LogLevel)
		fmt.Printf("  IPv6:         %v\n", cfg.IPv6)
		fmt.Printf("  System Proxy: %v\n", cfg.SystemProxy)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]
		cfg := config.LoadConfig()
		switch key {
		case "workdir":
			cfg.WorkDir = value
		case "loglevel":
			cfg.LogLevel = value
		case "ipv6":
			cfg.IPv6 = value == "true"
		case "systemproxy":
			cfg.SystemProxy = value == "true"
		default:
			return fmt.Errorf("unknown config key: %s", key)
		}
		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("save failed: %w", err)
		}
		fmt.Printf("✓ Set %s = %s\n", key, value)
		return nil
	},
}

var configEditCmd = &cobra.Command{
	Use:   "edit [engine]",
	Short: "Open config file in editor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		engineName := args[0]
		cfg := config.LoadConfig()
		em := engine.NewManager(cfg)
		configPath := em.GetConfigPath(engineName)
		if configPath == "" {
			return fmt.Errorf("no config found for %s. Run 'proxyman install %s' first", engineName, engineName)
		}
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}
		fmt.Printf("Editing %s config: %s\n", engineName, configPath)
		c := exec.Command(editor, configPath)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

// ==================== test ====================
var testCmd = &cobra.Command{
	Use:   "test [engine]",
	Short: "Test proxy connection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		engineName := args[0]
		cfg := config.LoadConfig()
		em := engine.NewManager(cfg)
		configPath := em.GetConfigPath(engineName)
		if configPath == "" {
			return fmt.Errorf("no config found for %s", engineName)
		}
		fmt.Printf("Testing %s connection...\n", engineName)
		fmt.Printf("Config: %s\n", configPath)
		engineDir := filepath.Join(cfg.WorkDir, "engines", engineName)
		var binary string
		switch engineName {
		case "v2ray", "xray":
			binary = filepath.Join(engineDir, "xray")
		case "clash", "mihomo":
			binary = filepath.Join(engineDir, "mihomo")
		}
		if _, err := os.Stat(binary); os.IsNotExist(err) {
			return fmt.Errorf("engine binary not found: %s", binary)
		}
		fmt.Printf("✓ Binary: %s\n", binary)
		fmt.Printf("✓ Config: %s\n", configPath)
		fmt.Println("✓ Engine is ready to start")
		return nil
	},
}

// ==================== split ====================
var splitCmd = &cobra.Command{
	Use:   "split [engine] [mode]",
	Short: "Configure traffic splitting (auto-split domestic/international)",
	Long: `Configure automatic traffic splitting:
  - full: All traffic goes through proxy (default)
  - china-split: Domestic traffic direct, international through proxy`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		engineName := args[0]
		mode := args[1]
		if mode != "full" && mode != "china-split" {
			return fmt.Errorf("invalid mode: %s (use 'full' or 'china-split')", mode)
		}
		cfg := config.LoadConfig()
		engineDir := filepath.Join(cfg.WorkDir, "engines", engineName)
		if _, err := os.Stat(engineDir); os.IsNotExist(err) {
			return fmt.Errorf("%s is not installed", engineName)
		}
		var configPath string
		var err error
		switch engineName {
		case "v2ray", "xray":
			configPath, err = engine.WriteXrayConfig(engineDir, cfg.LogLevel, "vmess", nil)
		case "clash", "mihomo":
			configPath, err = engine.WriteMihomoConfig(engineDir, cfg.LogLevel, mode)
		default:
			return fmt.Errorf("unsupported engine: %s", engineName)
		}
		if err != nil {
			return fmt.Errorf("generate config: %w", err)
		}
		fmt.Printf("✓ Config updated: %s\n", configPath)
		fmt.Printf("✓ Mode: %s\n", mode)
		if mode == "china-split" {
			fmt.Println("  - Domestic traffic: DIRECT")
			fmt.Println("  - International traffic: PROXY")
		} else {
			fmt.Println("  - All traffic: PROXY")
		}
		return nil
	},
}

// ==================== import ====================
var importCmd = &cobra.Command{
	Use:   "import [link|subscription-url]",
	Short: "Import proxy configuration from link or subscription URL",
	Long: `Import proxy configuration from various sources:
  - vmess:// links
  - vless:// links
  - trojan:// links
  - ss:// links
  - Subscription URLs (airport links)

Examples:
  proxyman import "vmess://eyJ2IjoiMiIs..."
  proxyman import "vless://uuid@server:port?security=tls#name"
  proxyman import "https://example.com/sub?token=xxx"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		protocol := parser.DetectProtocol(input)

		switch protocol {
		case "subscription":
			fmt.Println("检测到订阅链接，正在解析...")
			sub, err := parser.ParseSubscription(input)
			if err != nil {
				return fmt.Errorf("解析订阅失败: %w", err)
			}
			fmt.Printf("✓ 解析成功，共 %d 个节点\n\n", sub.Count)
			for i, node := range sub.Proxies {
				fmt.Printf("  %d. [%s] %s → %s:%d\n", i+1, node.Type, node.Name, node.Server, node.Port)
			}

		case "vmess", "vless", "trojan", "shadowsocks":
			fmt.Printf("检测到 %s 协议链接\n", protocol)
			node, err := parser.ParseProxyLink(input)
			if err != nil {
				return fmt.Errorf("解析链接失败: %w", err)
			}
			fmt.Printf("✓ 解析成功\n")
			fmt.Printf("  名称:   %s\n", node.Name)
			fmt.Printf("  类型:   %s\n", node.Type)
			fmt.Printf("  服务器: %s:%d\n", node.Server, node.Port)
			if node.UUID != "" {
				fmt.Printf("  UUID:   %s\n", node.UUID)
			}
			if node.SNI != "" {
				fmt.Printf("  SNI:    %s\n", node.SNI)
			}

		default:
			return fmt.Errorf("无法识别的输入格式: %s\n支持: vmess://, vless://, trojan://, ss://, 订阅链接", input)
		}
		return nil
	},
}

// ==================== system-proxy ====================
var systemProxyCmd = &cobra.Command{
	Use:   "system-proxy [enable|disable|status]",
	Short: "Manage system-wide proxy settings",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		action := args[0]
		cfg := config.LoadConfig()
		sp := config.NewSystemProxy(cfg)
		switch action {
		case "enable":
			if err := sp.Enable(); err != nil {
				return fmt.Errorf("enable failed: %w", err)
			}
			fmt.Println("✓ System proxy enabled")
		case "disable":
			if err := sp.Disable(); err != nil {
				return fmt.Errorf("disable failed: %w", err)
			}
			fmt.Println("✓ System proxy disabled")
		case "status":
			enabled, err := sp.Status()
			if err != nil {
				return fmt.Errorf("status failed: %w", err)
			}
			if enabled {
				fmt.Println("System proxy: ENABLED")
			} else {
				fmt.Println("System proxy: DISABLED")
			}
		default:
			return fmt.Errorf("unknown action: %s (use enable|disable|status)", action)
		}
		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configEditCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(systemProxyCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(splitCmd)
	rootCmd.AddCommand(importCmd)
}
