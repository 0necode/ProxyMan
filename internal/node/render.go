package node

import (
	"fmt"
	"strings"
)

// ClashProxy renders a node as a Mihomo/Clash proxy mapping
func ClashProxy(n *Node) map[string]any {
	p := map[string]any{
		"name":   n.Name,
		"type":   clashType(n.Type),
		"server": n.Server,
		"port":   n.Port,
	}
	if n.UUID != "" {
		p["uuid"] = n.UUID
	}
	if n.Password != "" {
		p["password"] = n.Password
	}
	if n.Cipher != "" {
		p["cipher"] = n.Cipher
	}
	if n.SNI != "" {
		p["sni"] = n.SNI
	}
	if n.TLS {
		p["tls"] = true
	}
	if n.Network != "" {
		p["network"] = n.Network
	}
	if n.Network == "ws" {
		opts := map[string]any{}
		if n.WSPath != "" {
			opts["path"] = n.WSPath
		}
		if n.WSHost != "" {
			opts["headers"] = map[string]any{"Host": n.WSHost}
		}
		if len(opts) > 0 {
			p["ws-opts"] = opts
		}
	}
	if n.Flow != "" {
		p["flow"] = n.Flow
	}
	if len(n.Alpn) > 0 {
		p["alpn"] = n.Alpn
	}
	for k, v := range n.Params {
		p[k] = v
	}
	return p
}

func clashType(t string) string {
	switch strings.ToLower(t) {
	case "shadowsocks", "ss":
		return "ss"
	case "shadowsocks-r", "ssr":
		return "ssr"
	case "hysteria", "hysteria2", "hy2":
		return "hysteria2"
	case "tuic":
		return "tuic"
	default:
		return strings.ToLower(t)
	}
}

// XrayOutbound renders a node as an Xray outbound mapping
func XrayOutbound(n *Node) map[string]any {
	settings := map[string]any{}
	switch strings.ToLower(n.Type) {
	case "vmess":
		users := []any{map[string]any{"id": n.UUID, "alterId": 0, "security": "auto"}}
		if n.Cipher != "" && n.Cipher != "auto" {
			users[0].(map[string]any)["security"] = n.Cipher
		}
		settings["vnext"] = []any{map[string]any{
			"address": n.Server,
			"port":    n.Port,
			"users":   users,
		}}
	case "vless":
		user := map[string]any{"id": n.UUID, "encryption": "none"}
		if n.Flow != "" {
			user["flow"] = n.Flow
		}
		settings["vnext"] = []any{map[string]any{
			"address": n.Server,
			"port":    n.Port,
			"users":   []any{user},
		}}
	case "trojan":
		settings["servers"] = []any{map[string]any{
			"address":  n.Server,
			"port":     n.Port,
			"password": n.Password,
		}}
	case "shadowsocks", "ss":
		m := map[string]any{
			"address":  n.Server,
			"port":     n.Port,
			"password": n.Password,
		}
		if n.Cipher != "" {
			m["method"] = n.Cipher
		}
		settings["servers"] = []any{m}
	default:
		settings["address"] = n.Server
		settings["port"] = n.Port
	}

	out := map[string]any{
		"tag":      fmt.Sprintf("proxy-%s", sanitise(n.Name)),
		"protocol": strings.ToLower(n.Type),
		"settings": settings,
	}

	stream := map[string]any{}
	if n.Network != "" {
		stream["network"] = n.Network
	}
	if n.Network == "ws" {
		ws := map[string]any{}
		if n.WSPath != "" {
			ws["path"] = n.WSPath
		}
		if n.WSHost != "" {
			ws["headers"] = map[string]any{"Host": n.WSHost}
		}
		stream["wsSettings"] = ws
	}
	if n.SNI != "" {
		stream["security"] = "tls"
		tls := map[string]any{"serverName": n.SNI}
		if len(n.Alpn) > 0 {
			tls["alpn"] = n.Alpn
		}
		stream["tlsSettings"] = tls
	} else if n.TLS {
		stream["security"] = "tls"
	}
	if len(stream) > 0 {
		out["streamSettings"] = stream
	}
	return out
}

func sanitise(name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(name) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			b.WriteRune('-')
		}
	}
	s := strings.Trim(b.String(), "-")
	if s == "" {
		s = "node"
	}
	return s
}
