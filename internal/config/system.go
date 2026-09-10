package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SystemProxy manages system-wide proxy settings
type SystemProxy struct {
	cfg *Config
}

// Enable enables system-wide proxy
func (sp *SystemProxy) Enable() error {
	envFile := filepath.Join(sp.cfg.WorkDir, "proxy-env")
	content := fmt.Sprintf(`export HTTP_PROXY=http://127.0.0.1:7890
export HTTPS_PROXY=http://127.0.0.1:7890
export http_proxy=http://127.0.0.1:7890
export https_proxy=http://127.0.0.1:7890
export NO_PROXY=localhost,127.0.0.1
export no_proxy=localhost,127.0.0.1
`)

	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("write env file: %w", err)
	}

	if os.Geteuid() == 0 {
		sp.updateEtcEnvironment()
	}

	return nil
}

// Disable disables system-wide proxy
func (sp *SystemProxy) Disable() error {
	envFile := filepath.Join(sp.cfg.WorkDir, "proxy-env")
	if _, err := os.Stat(envFile); err == nil {
		return os.Remove(envFile)
	}
	return nil
}

// Status returns whether system proxy is enabled
func (sp *SystemProxy) Status() (bool, error) {
	envFile := filepath.Join(sp.cfg.WorkDir, "proxy-env")
	_, err := os.Stat(envFile)
	return err == nil, nil
}

func (sp *SystemProxy) updateEtcEnvironment() {
	data, err := os.ReadFile("/etc/environment")
	if err != nil {
		return
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	var newLines []string
	skip := false
	for _, line := range lines {
		upper := strings.ToUpper(strings.TrimSpace(line))
		if strings.HasPrefix(upper, "HTTP_PROXY=") ||
			strings.HasPrefix(upper, "HTTPS_PROXY=") ||
			strings.HasPrefix(upper, "NO_PROXY=") {
			skip = true
			continue
		}
		if skip && strings.HasPrefix(line, "export ") {
			skip = false
			continue
		}
		skip = false
		newLines = append(newLines, line)
	}

	newLines = append(newLines, "HTTP_PROXY=http://127.0.0.1:7890")
	newLines = append(newLines, "HTTPS_PROXY=http://127.0.0.1:7890")
	newLines = append(newLines, "NO_PROXY=localhost,127.0.0.1")

	os.WriteFile("/etc/environment", []byte(strings.Join(newLines, "\n")), 0644)
}
