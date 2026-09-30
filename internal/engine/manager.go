package engine

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Geek0ne/ProxyMan/internal/config"
)

// EngineState represents the state of a proxy engine
type EngineState struct {
	Engine  string `json:"engine"`
	State   string `json:"state"` // installed, running, not_installed
	Version string `json:"version,omitempty"`
	PID     int    `json:"pid,omitempty"`
}

// Manager manages proxy engines
type Manager struct {
	cfg   *Config
	mu    sync.RWMutex
	procs map[string]*exec.Cmd
}

// Config holds engine manager configuration
type Config struct {
	WorkDir  string
	LogLevel string
}

// NewManager creates a new engine manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		cfg: &Config{
			WorkDir:  cfg.WorkDir,
			LogLevel: cfg.LogLevel,
		},
		procs: make(map[string]*exec.Cmd),
	}
}

// XrayRelease defines the XTLS Xray-core GitHub release parameters.
var XrayRelease = EngineRelease{
	Name:        "Xray-core",
	BinaryName:  "xray",
	ArchiveName: "Xray-linux-64.zip",
	Owner:       "XTLS",
	Repo:        "Xray-core",
	IsZip:       true,
}

// MihomoRelease defines the MetaCubeX Mihomo (Clash Meta) release parameters.
var MihomoRelease = EngineRelease{
	Name:        "Mihomo (Clash Meta)",
	BinaryName:  "mihomo",
	ArchiveName: "mihomo-linux-amd64",
	Owner:       "MetaCubeX",
	Repo:        "mihomo",
	IsZip:       false,
	PostExtract: func(dir, binaryPath string) error {
		// Mihomo releases may name the binary differently; check common names
		candidates := []string{
			filepath.Join(dir, "mihomo"),
			filepath.Join(dir, "clash"),
			filepath.Join(dir, "mihomo-linux-amd64"),
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				if c != binaryPath {
					os.Rename(c, binaryPath)
				}
				return nil
			}
		}
		return fmt.Errorf("mihomo binary not found after extraction")
	},
}

// Install installs a proxy engine
func (m *Manager) Install(engine string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	engine = strings.ToLower(engine)
	switch engine {
	case "v2ray", "xray":
		return m.installXray()
	case "clash", "mihomo":
		return m.installMihomo()
	default:
		return fmt.Errorf("unsupported engine: %s (supported: v2ray/xray, clash/mihomo)", engine)
	}
}

// Uninstall removes a proxy engine
func (m *Manager) Uninstall(engine string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	engine = strings.ToLower(engine)
	// Normalize aliases
	switch engine {
	case "xray":
		engine = "v2ray"
	case "mihomo":
		engine = "clash"
	}

	engineDir := filepath.Join(m.cfg.WorkDir, "engines", engine)

	if _, err := os.Stat(engineDir); os.IsNotExist(err) {
		fmt.Printf("  %s is not installed\n", engine)
		return nil
	}

	if err := os.RemoveAll(engineDir); err != nil {
		return fmt.Errorf("remove engine dir: %w", err)
	}

	fmt.Printf("  ✓ %s uninstalled\n", engine)
	return nil
}

// Start starts a proxy engine with a config file
func (m *Manager) Start(engine, cfgFile string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	engine = strings.ToLower(engine)
	engineDir := filepath.Join(m.cfg.WorkDir, "engines", engine)

	if _, err := os.Stat(engineDir); os.IsNotExist(err) {
		return fmt.Errorf("%s is not installed. Run 'ProxyMan install %s' first", engine, engine)
	}

	var binary string
	switch engine {
	case "v2ray", "xray":
		binary = filepath.Join(engineDir, "xray")
	case "clash", "mihomo":
		binary = filepath.Join(engineDir, "mihomo")
	}

	if _, err := os.Stat(binary); os.IsNotExist(err) {
		return fmt.Errorf("engine binary not found: %s", binary)
	}

	// Stop existing process if running
	if cmd, ok := m.procs[engine]; ok && cmd.Process != nil {
		cmd.Process.Kill()
	}

	// Determine run command based on engine type
	var cmdArgs []string
	switch engine {
	case "v2ray", "xray":
		cmdArgs = []string{"run", "-config", cfgFile}
	case "clash", "mihomo":
		cmdArgs = []string{"-d", engineDir, "-f", cfgFile}
	}

	cmd := exec.Command(binary, cmdArgs...)
	cmd.Dir = engineDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start engine: %w", err)
	}

	m.procs[engine] = cmd
	fmt.Printf("  ✓ %s started (PID: %d)\n", engine, cmd.Process.Pid)
	return nil
}

// Stop stops a running proxy engine
func (m *Manager) Stop(engine string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	engine = strings.ToLower(engine)
	cmd, ok := m.procs[engine]
	if !ok || cmd.Process == nil {
		fmt.Printf("  %s is not running\n", engine)
		return nil
	}

	if err := cmd.Process.Kill(); err != nil {
		return fmt.Errorf("kill process: %w", err)
	}

	delete(m.procs, engine)
	fmt.Printf("  ✓ %s stopped\n", engine)
	return nil
}

// Status returns the status of all engines
func (m *Manager) Status() ([]EngineState, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var statuses []EngineState

	for _, engine := range []string{"v2ray", "clash"} {
		engineDir := filepath.Join(m.cfg.WorkDir, "engines", engine)
		state := EngineState{Engine: engine}

		if _, err := os.Stat(engineDir); os.IsNotExist(err) {
			state.State = "not_installed"
		} else {
			state.State = "installed"
			// Read version if available
			versionFile := filepath.Join(engineDir, "version")
			if data, err := os.ReadFile(versionFile); err == nil {
				lines := strings.SplitN(string(data), "\n", 2)
				if len(lines) > 0 {
					state.Version = strings.TrimSpace(lines[0])
				}
			}
			if cmd, ok := m.procs[engine]; ok && cmd.Process != nil {
				state.State = "running"
				state.PID = cmd.Process.Pid
			}
		}

		statuses = append(statuses, state)
	}

	return statuses, nil
}

func (m *Manager) installXray() error {
	engineDir := filepath.Join(m.cfg.WorkDir, "engines", "v2ray")

	fmt.Println("Installing Xray-core (XTLS)...")
	fmt.Println("  Source: github.com/XTLS/Xray-core")

	binaryPath, hash, err := DownloadEngine(XrayRelease, engineDir)
	if err != nil {
		return fmt.Errorf("download xray: %w", err)
	}

	fmt.Printf("  ✓ Binary installed: %s\n", binaryPath)
	fmt.Printf("  ✓ SHA256: %s\n", hash)

	// Generate default config
	configPath, err := WriteXrayConfig(engineDir, m.cfg.LogLevel, "vmess", nil)
	if err != nil {
		return fmt.Errorf("generate xray config: %w", err)
	}
	fmt.Printf("  ✓ Config: %s\n", configPath)

	// Generate systemd service
	servicePath, err := WriteSystemdService(engineDir, "xray", "config.json", "Xray-core proxy engine")
	if err != nil {
		return fmt.Errorf("generate service: %w", err)
	}
	fmt.Printf("  ✓ Systemd service: %s\n", servicePath)

	return nil
}

func (m *Manager) installMihomo() error {
	engineDir := filepath.Join(m.cfg.WorkDir, "engines", "clash")

	fmt.Println("Installing Mihomo (Clash Meta)...")
	fmt.Println("  Source: github.com/MetaCubeX/mihomo")

	binaryPath, hash, err := DownloadEngine(MihomoRelease, engineDir)
	if err != nil {
		return fmt.Errorf("download mihomo: %w", err)
	}

	fmt.Printf("  ✓ Binary installed: %s\n", binaryPath)
	fmt.Printf("  ✓ SHA256: %s\n", hash)

	// Generate default config
	configPath, err := WriteMihomoConfig(engineDir, m.cfg.LogLevel, "full")
	if err != nil {
		return fmt.Errorf("generate mihomo config: %w", err)
	}
	fmt.Printf("  ✓ Config: %s\n", configPath)

	// Generate systemd service
	servicePath, err := WriteSystemdService(engineDir, "mihomo", "config.yaml", "Mihomo (Clash Meta) proxy engine")
	if err != nil {
		return fmt.Errorf("generate service: %w", err)
	}
	fmt.Printf("  ✓ Systemd service: %s\n", servicePath)

	return nil
}

// GetConfigPath returns the config file path for an engine
func (m *Manager) GetConfigPath(engine string) string {
	engine = strings.ToLower(engine)
	switch engine {
	case "v2ray", "xray":
		return filepath.Join(m.cfg.WorkDir, "engines", "v2ray", "config.json")
	case "clash", "mihomo":
		return filepath.Join(m.cfg.WorkDir, "engines", "clash", "config.yaml")
	}
	return ""
}
