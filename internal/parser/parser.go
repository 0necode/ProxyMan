package parser

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ProxyNode 代理节点配置
type ProxyNode struct {
	Name       string            `json:"name"`
	Type       string            `json:"type"` // vmess, vless, trojan, ss, ss2022, anytls
	Server     string            `json:"server"`
	Port       int               `json:"port"`
	UUID       string            `json:"uuid,omitempty"`
	Password   string            `json:"password,omitempty"`
	Cipher     string            `json:"cipher,omitempty"`
	SNI        string            `json:"sni,omitempty"`
	Network    string            `json:"network,omitempty"`
	WSPATH     string            `json:"ws-path,omitempty"`
	WSHost     string            `json:"ws-host,omitempty"`
	TLS        bool              `json:"tls"`
	Flow       string            `json:"flow,omitempty"`
	Alpn       []string          `json:"alpn,omitempty"`
	Fingerprint string            `json:"fingerprint,omitempty"`
	Params     map[string]string `json:"params,omitempty"`
}

// Subscription 订阅信息
type Subscription struct {
	Proxies []ProxyNode `json:"proxies"`
	Count   int         `json:"count"`
}

// DetectProtocol 检测协议类型
func DetectProtocol(input string) string {
	input = strings.TrimSpace(input)

	// 检测订阅链接
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		return "subscription"
	}

	// 检测协议链接
	if strings.HasPrefix(input, "vmess://") {
		return "vmess"
	}
	if strings.HasPrefix(input, "vless://") {
		return "vless"
	}
	if strings.HasPrefix(input, "trojan://") {
		return "trojan"
	}
	if strings.HasPrefix(input, "ss://") {
		return "shadowsocks"
	}
	if strings.HasPrefix(input, "ssr://") {
		return "shadowsocks-r"
	}
	if strings.HasPrefix(input, "hysteria://") || strings.HasPrefix(input, "hy2://") {
		return "hysteria"
	}
	if strings.HasPrefix(input, "tuic://") {
		return "tuic"
	}
	if strings.HasPrefix(input, "anytls://") {
		return "anytls"
	}

	// 检测 JSON 配置
	if strings.HasPrefix(input, "{") {
		return "json-config"
	}

	// 检测 YAML 配置
	if strings.HasPrefix(input, "mixed-port:") || strings.HasPrefix(input, "proxies:") {
		return "yaml-config"
	}

	return "unknown"
}

// ParseProxyLink 解析代理链接
func ParseProxyLink(link string) (*ProxyNode, error) {
	link = strings.TrimSpace(link)

	switch {
	case strings.HasPrefix(link, "vmess://"):
		return parseVMessLink(link)
	case strings.HasPrefix(link, "vless://"):
		return parseVLessLink(link)
	case strings.HasPrefix(link, "trojan://"):
		return parseTrojanLink(link)
	case strings.HasPrefix(link, "ss://"):
		return parseSSLink(link)
	case strings.HasPrefix(link, "hysteria://") || strings.HasPrefix(link, "hy2://"):
		return parseHysteriaLink(link)
	case strings.HasPrefix(link, "tuic://"):
		return parseTUICLink(link)
	case strings.HasPrefix(link, "anytls://"):
		return parseAnyTLSLink(link)
	default:
		return nil, fmt.Errorf("unsupported protocol link: %s", link[:20])
	}
}

// parseVMessLink 解析 vmess:// 链接
func parseVMessLink(link string) (*ProxyNode, error) {
	// vmess://base64(json)
	data := strings.TrimPrefix(link, "vmess://")
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("decode vmess link: %w", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(decoded, &config); err != nil {
		return nil, fmt.Errorf("parse vmess json: %w", err)
	}

	node := &ProxyNode{
		Type: "vmess",
	}

	if v, ok := config["add"].(string); ok {
		node.Server = v
	}
	if v, ok := config["port"].(string); ok {
		fmt.Sscanf(v, "%d", &node.Port)
	}
	if v, ok := config["id"].(string); ok {
		node.UUID = v
	}
	if v, ok := config["net"].(string); ok {
		node.Network = v
	}
	if v, ok := config["path"].(string); ok {
		node.WSPATH = v
	}
	if v, ok := config["host"].(string); ok {
		node.WSHost = v
	}
	if v, ok := config["tls"].(string); ok {
		node.TLS = v == "tls"
	}
	if v, ok := config["sni"].(string); ok {
		node.SNI = v
	}
	if v, ok := config["ps"].(string); ok {
		node.Name = v
	}

	if node.Name == "" {
		node.Name = fmt.Sprintf("vmess-%s:%d", node.Server, node.Port)
	}

	return node, nil
}

// parseVLessLink 解析 vless:// 链接
func parseVLessLink(link string) (*ProxyNode, error) {
	// vless://uuid@server:port?params#name
	link = strings.TrimPrefix(link, "vless://")

	// 分离名称
	parts := strings.Split(link, "#")
	if len(parts) == 2 {
		link = parts[0]
	}

	// 分离参数
	paramParts := strings.SplitN(link, "?", 2)
	mainPart := paramParts[0]

	// 解析 uuid@server:port
	atParts := strings.Split(mainPart, "@")
	if len(atParts) != 2 {
		return nil, fmt.Errorf("invalid vless link format")
	}

	uuid := atParts[0]
	serverPort := atParts[1]
	spParts := strings.SplitN(serverPort, ":", 2)
	if len(spParts) != 2 {
		return nil, fmt.Errorf("invalid server:port format")
	}

	node := &ProxyNode{
		Type:   "vless",
		UUID:   uuid,
		Server: spParts[0],
		TLS:    true,
	}

	fmt.Sscanf(spParts[1], "%d", &node.Port)

	// 解析参数
	if len(paramParts) == 2 {
		params, _ := url.ParseQuery(paramParts[1])
		if v := params.Get("security"); v != "" {
			node.TLS = v == "tls"
		}
		if v := params.Get("sni"); v != "" {
			node.SNI = v
		}
		if v := params.Get("type"); v != "" {
			node.Network = v
		}
		if v := params.Get("path"); v != "" {
			node.WSPATH = v
		}
		if v := params.Get("host"); v != "" {
			node.WSHost = v
		}
		if v := params.Get("flow"); v != "" {
			node.Flow = v
		}
		if v := params.Get("fp"); v != "" {
			node.Fingerprint = v
		}
	}

	// 提取名称
	if len(parts) == 2 {
		node.Name, _ = url.QueryUnescape(parts[1])
	}

	if node.Name == "" {
		node.Name = fmt.Sprintf("vless-%s:%d", node.Server, node.Port)
	}

	return node, nil
}

// parseTrojanLink 解析 trojan:// 链接
func parseTrojanLink(link string) (*ProxyNode, error) {
	// trojan://password@server:port?params#name
	link = strings.TrimPrefix(link, "trojan://")

	parts := strings.Split(link, "#")
	if len(parts) == 2 {
		link = parts[0]
	}

	paramParts := strings.SplitN(link, "?", 2)
	mainPart := paramParts[0]

	atParts := strings.Split(mainPart, "@")
	if len(atParts) != 2 {
		return nil, fmt.Errorf("invalid trojan link format")
	}

	password := atParts[0]
	serverPort := atParts[1]
	spParts := strings.SplitN(serverPort, ":", 2)
	if len(spParts) != 2 {
		return nil, fmt.Errorf("invalid server:port format")
	}

	node := &ProxyNode{
		Type:     "trojan",
		Password: password,
		Server:   spParts[0],
		TLS:      true,
	}

	fmt.Sscanf(spParts[1], "%d", &node.Port)

	if len(paramParts) == 2 {
		params, _ := url.ParseQuery(paramParts[1])
		if v := params.Get("sni"); v != "" {
			node.SNI = v
		}
		if v := params.Get("type"); v != "" {
			node.Network = v
		}
		if v := params.Get("path"); v != "" {
			node.WSPATH = v
		}
		if v := params.Get("host"); v != "" {
			node.WSHost = v
		}
	}

	if len(parts) == 2 {
		node.Name, _ = url.QueryUnescape(parts[1])
	}

	if node.Name == "" {
		node.Name = fmt.Sprintf("trojan-%s:%d", node.Server, node.Port)
	}

	return node, nil
}

// parseSSLink 解析 ss:// 链接
func parseSSLink(link string) (*ProxyNode, error) {
	// ss://base64(method:password)@server:port#name
	link = strings.TrimPrefix(link, "ss://")

	parts := strings.Split(link, "#")
	if len(parts) == 2 {
		link = parts[0]
	}

	atParts := strings.Split(link, "@")
	if len(atParts) != 2 {
		return nil, fmt.Errorf("invalid ss link format")
	}

	decoded, err := base64.StdEncoding.DecodeString(atParts[0])
	if err != nil {
		return nil, fmt.Errorf("decode ss link: %w", err)
	}

	cipherParts := strings.SplitN(string(decoded), ":", 2)
	if len(cipherParts) != 2 {
		return nil, fmt.Errorf("invalid cipher:password format")
	}

	serverPort := atParts[1]
	spParts := strings.SplitN(serverPort, ":", 2)
	if len(spParts) != 2 {
		return nil, fmt.Errorf("invalid server:port format")
	}

	node := &ProxyNode{
		Type:     "shadowsocks",
		Cipher:   cipherParts[0],
		Password: cipherParts[1],
		Server:   spParts[0],
	}

	fmt.Sscanf(spParts[1], "%d", &node.Port)

	if len(parts) == 2 {
		node.Name, _ = url.QueryUnescape(parts[1])
	}

	if node.Name == "" {
		node.Name = fmt.Sprintf("ss-%s:%d", node.Server, node.Port)
	}

	return node, nil
}

// parseHysteriaLink 解析 hysteria:// 链接
func parseHysteriaLink(link string) (*ProxyNode, error) {
	link = strings.TrimPrefix(link, "hysteria://")
	link = strings.TrimPrefix(link, "hy2://")

	parts := strings.Split(link, "#")
	if len(parts) == 2 {
		link = parts[0]
	}

	spParts := strings.SplitN(link, ":", 2)
	if len(spParts) != 2 {
		return nil, fmt.Errorf("invalid hysteria server:port format")
	}

	node := &ProxyNode{
		Type:   "hysteria",
		Server: spParts[0],
	}

	fmt.Sscanf(spParts[1], "%d", &node.Port)

	if len(parts) == 2 {
		node.Name, _ = url.QueryUnescape(parts[1])
	}

	if node.Name == "" {
		node.Name = fmt.Sprintf("hysteria-%s:%d", node.Server, node.Port)
	}

	return node, nil
}

// parseTUICLink 解析 tuic:// 链接
func parseTUICLink(link string) (*ProxyNode, error) {
	link = strings.TrimPrefix(link, "tuic://")

	parts := strings.Split(link, "#")
	if len(parts) == 2 {
		link = parts[0]
	}

	spParts := strings.SplitN(link, ":", 2)
	if len(spParts) != 2 {
		return nil, fmt.Errorf("invalid tuic server:port format")
	}

	node := &ProxyNode{
		Type:   "tuic",
		Server: spParts[0],
	}

	fmt.Sscanf(spParts[1], "%d", &node.Port)

	if len(parts) == 2 {
		node.Name, _ = url.QueryUnescape(parts[1])
	}

	if node.Name == "" {
		node.Name = fmt.Sprintf("tuic-%s:%d", node.Server, node.Port)
	}

	return node, nil
}

// parseAnyTLSLink 解析 anytls:// 链接
func parseAnyTLSLink(link string) (*ProxyNode, error) {
	link = strings.TrimPrefix(link, "anytls://")

	parts := strings.Split(link, "#")
	if len(parts) == 2 {
		link = parts[0]
	}

	spParts := strings.SplitN(link, ":", 2)
	if len(spParts) != 2 {
		return nil, fmt.Errorf("invalid anytls server:port format")
	}

	node := &ProxyNode{
		Type:   "anytls",
		Server: spParts[0],
		TLS:    true,
	}

	fmt.Sscanf(spParts[1], "%d", &node.Port)

	if len(parts) == 2 {
		node.Name, _ = url.QueryUnescape(parts[1])
	}

	if node.Name == "" {
		node.Name = fmt.Sprintf("anytls-%s:%d", node.Server, node.Port)
	}

	return node, nil
}

// ParseSubscription 解析订阅链接
func ParseSubscription(subURL string) (*Subscription, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(subURL)
	if err != nil {
		return nil, fmt.Errorf("fetch subscription: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("subscription returned status %d", resp.StatusCode)
	}

	// 读取响应
	buf := make([]byte, 0, 1024*1024) // 1MB buffer
	tmp := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(tmp)
		buf = append(buf, tmp[:n]...)
		if err != nil {
			break
		}
	}

	// 尝试 base64 解码
	decoded, err := base64.StdEncoding.DecodeString(string(buf))
	if err == nil {
		buf = decoded
	}

	// 按行分割
	lines := strings.Split(string(buf), "\n")

	sub := &Subscription{
		Proxies: make([]ProxyNode, 0),
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 尝试解析为代理链接
		protocol := DetectProtocol(line)
		if protocol != "unknown" && protocol != "subscription" {
			node, err := ParseProxyLink(line)
			if err == nil {
				sub.Proxies = append(sub.Proxies, *node)
			}
		}
	}

	sub.Count = len(sub.Proxies)
	return sub, nil
}

// FormatProxyLink 生成代理链接
func FormatProxyLink(node *ProxyNode) (string, error) {
	switch node.Type {
	case "vmess":
		return formatVMessLink(node)
	case "vless":
		return formatVLessLink(node)
	case "trojan":
		return formatTrojanLink(node)
	case "shadowsocks":
		return formatSSLink(node)
	default:
		return "", fmt.Errorf("unsupported protocol: %s", node.Type)
	}
}

func formatVMessLink(node *ProxyNode) (string, error) {
	config := map[string]interface{}{
		"v":    "2",
		"ps":   node.Name,
		"add":  node.Server,
		"port": fmt.Sprintf("%d", node.Port),
		"id":   node.UUID,
		"aid":  "0",
		"net":  node.Network,
		"type": "none",
		"host": node.WSHost,
		"path": node.WSPATH,
		"tls":  "",
		"sni":  node.SNI,
	}

	if node.TLS {
		config["tls"] = "tls"
	}

	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	return "vmess://" + base64.StdEncoding.EncodeToString(jsonBytes), nil
}

func formatVLessLink(node *ProxyNode) (string, error) {
	params := url.Values{}
	params.Set("security", "tls")
	params.Set("type", node.Network)
	if node.SNI != "" {
		params.Set("sni", node.SNI)
	}
	if node.WSPATH != "" {
		params.Set("path", node.WSPATH)
	}
	if node.Flow != "" {
		params.Set("flow", node.Flow)
	}

	return fmt.Sprintf("vless://%s@%s:%d?%s#%s",
		node.UUID, node.Server, node.Port,
		params.Encode(), url.QueryEscape(node.Name)), nil
}

func formatTrojanLink(node *ProxyNode) (string, error) {
	params := url.Values{}
	params.Set("security", "tls")
	if node.SNI != "" {
		params.Set("sni", node.SNI)
	}

	return fmt.Sprintf("trojan://%s@%s:%d?%s#%s",
		node.Password, node.Server, node.Port,
		params.Encode(), url.QueryEscape(node.Name)), nil
}

func formatSSLink(node *ProxyNode) (string, error) {
	encoded := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s", node.Cipher, node.Password)),
	)

	return fmt.Sprintf("ss://%s@%s:%d#%s",
		encoded, node.Server, node.Port,
		url.QueryEscape(node.Name)), nil
}
