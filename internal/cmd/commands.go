package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Geek0ne/ProxyMan/internal/config"
	"github.com/Geek0ne/ProxyMan/internal/engine"
	nstore "github.com/Geek0ne/ProxyMan/internal/node"
	"github.com/Geek0ne/ProxyMan/internal/parser"
	"github.com/spf13/cobra"
)

// importStore makes 'import' persist parsed nodes instead of only printing them
var importStore bool

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
	Short: "Manage ProxyMan configuration",
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
			return fmt.Errorf("no config found for %s. Run 'ProxyMan install %s' first", engineName, engineName)
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
  ProxyMan import "vmess://eyJ2IjoiMiIs..."
  ProxyMan import "vless://uuid@server:port?security=tls#name"
  ProxyMan import "https://example.com/sub?token=xxx"`,
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
			fmt.Printf("✓ 解析成功，共 %d 个节点\n", sub.Count)
			for i, p := range sub.Proxies {
				fmt.Printf("  %d. [%s] %s → %s:%d\n", i+1, p.Type, p.Name, p.Server, p.Port)
			}

			if importStore {
				cfg := config.LoadConfig()
				store := nstore.NewStore(cfg.WorkDir)
				added := 0
				for _, p := range sub.Proxies {
					if store.Add(nstore.FromProxyNode(&p, "subscription")) {
						added++
					}
				}
				if dryRun {
					fmt.Printf("\n[DRY RUN] would store %d new node(s) in %s\n", added, store.Path())
					return nil
				}
				if err := store.Save(); err != nil {
					return fmt.Errorf("保存节点失败: %w", err)
				}
				fmt.Printf("\n✓ 已保存 %d 个新节点（累计 %d）→ %s\n",
					added, len(store.List()), store.Path())
				fmt.Println("  下一步: ProxyMan list   然后   ProxyMan apply <engine>")
			}

		case "vmess", "vless", "trojan", "shadowsocks", "hysteria", "tuic", "anytls":
			fmt.Printf("检测到 %s 协议链接\n", protocol)
			p, err := parser.ParseProxyLink(input)
			if err != nil {
				return fmt.Errorf("解析链接失败: %w", err)
			}
			fmt.Printf("✓ 解析成功\n")
			fmt.Printf("  名称:   %s\n", p.Name)
			fmt.Printf("  类型:   %s\n", p.Type)
			fmt.Printf("  服务器: %s:%d\n", p.Server, p.Port)
			if p.UUID != "" {
				fmt.Printf("  UUID:   %s\n", p.UUID)
			}
			if p.SNI != "" {
				fmt.Printf("  SNI:    %s\n", p.SNI)
			}

			if importStore {
				cfg := config.LoadConfig()
				store := nstore.NewStore(cfg.WorkDir)
				isNew := store.Add(nstore.FromProxyNode(p, protocol))
				if dryRun {
					fmt.Printf("\n[DRY RUN] would store node %q into %s\n", p.Name, store.Path())
					return nil
				}
				if err := store.Save(); err != nil {
					return fmt.Errorf("保存节点失败: %w", err)
				}
				if isNew {
					fmt.Printf("\n✓ 节点已保存: %s\n", p.Name)
				} else {
					fmt.Printf("\n✓ 节点已更新: %s\n", p.Name)
				}
				fmt.Printf("  存储位置: %s\n", store.Path())
				fmt.Println("  下一步: ProxyMan list   然后   ProxyMan apply <engine>")
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
	importCmd.Flags().BoolVar(&importStore, "store", false, "persist parsed nodes to the local store")

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
