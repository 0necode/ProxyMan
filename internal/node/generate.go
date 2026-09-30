package node

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// SplitMode selects the routing strategy written into a generated config
type SplitMode string

const (
	// SplitFull routes all traffic through the proxy group
	SplitFull SplitMode = "full"
	// SplitChina routes domestic traffic direct and the rest via the proxy
	SplitChina SplitMode = "china-split"
)

// Generator writes node sets into engine configuration files
type Generator struct {
	EngineDir string
	LogLevel  string
	Mode      SplitMode
	GroupName string
	// Preview renders to memory instead of writing to disk.
	Preview bool
}

// NewGenerator builds a generator for an engine directory
func NewGenerator(engineDir, logLevel string, mode SplitMode) *Generator {
	if mode == "" {
		mode = SplitFull
	}
	return &Generator{
		EngineDir: engineDir,
		LogLevel:  logLevel,
		Mode:      mode,
		GroupName: "ProxyMan",
	}
}

// ConfigPath returns the config file the generator will write
func (g *Generator) ConfigPath() string {
	return filepath.Join(g.EngineDir, "nodes-config.json")
}

// WriteXray renders nodes into a complete Xray config and returns its path
func (g *Generator) WriteXray(nodes []*Node) (string, error) {
	outbounds := []any{map[string]any{
		"tag":      "direct",
		"protocol": "freedom",
		"settings": map[string]any{},
	}}
	for _, n := range nodes {
		outbounds = append(outbounds, XrayOutbound(n))
	}
	outbounds = append(outbounds, map[string]any{
		"tag":      "block",
		"protocol": "blackhole",
		"settings": map[string]any{},
	})

	rules := []any{
		map[string]any{"type": "field", "ip": []string{"geoip:private"}, "outboundTag": "direct"},
	}
	if g.Mode == SplitChina {
		rules = append(rules,
			map[string]any{"type": "field", "ip": []string{"geoip:cn"}, "outboundTag": "direct"},
			map[string]any{"type": "field", "domain": []string{"geosite:cn"}, "outboundTag": "direct"},
		)
	}
	if len(nodes) > 0 {
		rules = append(rules, map[string]any{
			"type":        "field",
			"network":     "tcp,udp",
			"outboundTag": outboundTag(nodes[0]),
		})
	}

	cfg := map[string]any{
		"log":       map[string]any{"loglevel": g.LogLevel},
		"inbounds":  xrayInbounds(),
		"outbounds": outbounds,
		"routing":   map[string]any{"domainStrategy": "IPIfNonMatch", "rules": rules},
	}
	return g.writeJSON(cfg, "config.json")
}

// WriteClash renders nodes into a Mihomo/Clash config and returns its path
func (g *Generator) WriteClash(nodes []*Node) (string, error) {
	proxies := make([]any, 0, len(nodes))
	names := make([]string, 0, len(nodes))
	for _, n := range nodes {
		proxies = append(proxies, ClashProxy(n))
		names = append(names, n.Name)
	}

	group := make([]any, 0, len(names)+1)
	if g.Mode == SplitChina && len(names) > 0 {
		for _, n := range names {
			group = append(group, map[string]any{
				"type":     "url-test",
				"name":     "Auto-" + n,
				"proxies":  []any{n},
				"url":      "https://www.gstatic.com/generate_204",
				"interval": 300,
			})
		}
	}
	group = append(group, map[string]any{
		"type":    "select",
		"name":    g.GroupName,
		"proxies": append(names, "DIRECT"),
	})

	rules := []any{
		"GEOIP,private,DIRECT,no-resolve",
	}
	if g.Mode == SplitChina {
		rules = append(rules, "GEOSITE,cn,DIRECT", "GEOIP,cn,DIRECT,no-resolve")
	}
	rules = append(rules, "MATCH,"+g.GroupName)

	cfg := map[string]any{
		"mixed-port": 7890,
		"allow-lan":  false,
		"mode":       "rule",
		"log-level":  g.LogLevel,
		"dns": map[string]any{
			"enable":        true,
			"enhanced-mode": "fake-ip",
			"nameserver":    []string{"223.5.5.5", "119.29.29.29"},
			"fallback":      []string{"https://dns.google/dns-query"},
		},
		"proxies":      proxies,
		"proxy-groups": group,
		"rules":        rules,
	}
	return g.writeYAML(cfg, "config.yaml")
}

func xrayInbounds() []any {
	return []any{
		map[string]any{
			"tag":      "http-in",
			"port":     10809,
			"protocol": "http",
			"listen":   "127.0.0.1",
			"settings": map[string]any{"allowTransparent": false},
		},
		map[string]any{
			"tag":      "socks-in",
			"port":     10808,
			"protocol": "socks",
			"listen":   "127.0.0.1",
			"settings": map[string]any{"udp": true, "auth": "noauth"},
		},
	}
}

func outboundTag(n *Node) string {
	return fmt.Sprintf("proxy-%s", sanitise(n.Name))
}

// RenderXray returns the generated Xray config as a string
func (g *Generator) RenderXray(nodes []*Node) (string, string, error) {
	outbounds := []any{map[string]any{
		"tag":      "direct",
		"protocol": "freedom",
		"settings": map[string]any{},
	}}
	for _, n := range nodes {
		outbounds = append(outbounds, XrayOutbound(n))
	}
	outbounds = append(outbounds, map[string]any{
		"tag":      "block",
		"protocol": "blackhole",
		"settings": map[string]any{},
	})

	rules := []any{
		map[string]any{"type": "field", "ip": []string{"geoip:private"}, "outboundTag": "direct"},
	}
	if g.Mode == SplitChina {
		rules = append(rules,
			map[string]any{"type": "field", "ip": []string{"geoip:cn"}, "outboundTag": "direct"},
			map[string]any{"type": "field", "domain": []string{"geosite:cn"}, "outboundTag": "direct"},
		)
	}
	if len(nodes) > 0 {
		rules = append(rules, map[string]any{
			"type":        "field",
			"network":     "tcp,udp",
			"outboundTag": outboundTag(nodes[0]),
		})
	}

	cfg := map[string]any{
		"log":       map[string]any{"loglevel": g.LogLevel},
		"inbounds":  xrayInbounds(),
		"outbounds": outbounds,
		"routing":   map[string]any{"domainStrategy": "IPIfNonMatch", "rules": rules},
	}
	path := filepath.Join(g.EngineDir, "config.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", "", fmt.Errorf("marshal config: %w", err)
	}
	return path, string(data), nil
}

// RenderClash returns the generated Clash config as a string
func (g *Generator) RenderClash(nodes []*Node) (string, string, error) {
	proxies := make([]any, 0, len(nodes))
	names := make([]string, 0, len(nodes))
	for _, n := range nodes {
		proxies = append(proxies, ClashProxy(n))
		names = append(names, n.Name)
	}

	group := make([]any, 0, len(names)+1)
	if g.Mode == SplitChina && len(names) > 0 {
		for _, n := range names {
			group = append(group, map[string]any{
				"type":     "url-test",
				"name":     "Auto-" + n,
				"proxies":  []any{n},
				"url":      "https://www.gstatic.com/generate_204",
				"interval": 300,
			})
		}
	}
	group = append(group, map[string]any{
		"type":    "select",
		"name":    g.GroupName,
		"proxies": append(names, "DIRECT"),
	})

	rules := []any{"GEOIP,private,DIRECT,no-resolve"}
	if g.Mode == SplitChina {
		rules = append(rules, "GEOSITE,cn,DIRECT", "GEOIP,cn,DIRECT,no-resolve")
	}
	rules = append(rules, "MATCH,"+g.GroupName)

	cfg := map[string]any{
		"mixed-port": 7890,
		"allow-lan":  false,
		"mode":       "rule",
		"log-level":  g.LogLevel,
		"dns": map[string]any{
			"enable":        true,
			"enhanced-mode": "fake-ip",
			"nameserver":    []string{"223.5.5.5", "119.29.29.29"},
			"fallback":      []string{"https://dns.google/dns-query"},
		},
		"proxies":      proxies,
		"proxy-groups": group,
		"rules":        rules,
	}
	path := filepath.Join(g.EngineDir, "config.yaml")
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return "", "", fmt.Errorf("marshal config: %w", err)
	}
	var sb strings.Builder
	sb.WriteString("# Generated by ProxyMan — edit via 'ProxyMan apply'\n")
	sb.Write(data)
	return path, sb.String(), nil
}

func (g *Generator) writeJSON(cfg map[string]any, name string) (string, error) {
	path := filepath.Join(g.EngineDir, name)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return "", fmt.Errorf("write config: %w", err)
	}
	return path, nil
}

func (g *Generator) writeYAML(cfg map[string]any, name string) (string, error) {
	path := filepath.Join(g.EngineDir, name)
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("marshal config: %w", err)
	}
	var sb strings.Builder
	sb.WriteString("# Generated by ProxyMan — edit via 'ProxyMan import --apply'\n")
	sb.Write(data)
	if err := os.WriteFile(path, []byte(sb.String()), 0o600); err != nil {
		return "", fmt.Errorf("write config: %w", err)
	}
	return path, nil
}
