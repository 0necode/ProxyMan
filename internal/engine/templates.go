package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// ==================== Xray 配置模板 ====================

// XrayVMessConfig 支持 VMess 协议
const XrayVMessConfig = `{
  "log": {"loglevel": "{{.LogLevel}}"},
  "inbounds": [
    {"tag":"http-in","port":10809,"protocol":"http","listen":"127.0.0.1","settings":{}},
    {"tag":"socks-in","port":10808,"protocol":"socks","listen":"127.0.0.1","settings":{"udp":true}},
    {"tag":"tproxy-in","port":12345,"protocol":"dokodemo-door","listen":"0.0.0.0","settings":{"network":"tcp,udp","followRedirect":true},"sniffing":{"enabled":true,"destOverride":["http","tls"]}}
  ],
  "outbounds": [
    {
      "tag":"proxy","protocol":"vmess",
      "settings":{"vnext":[{"address":"{{.ServerAddr}}","port":{{.ServerPort}},"users":[{"id":"{{.UUID}}","alterId":0,"security":"auto"}]}]},
      "streamSettings":{"network":"{{.Network}}","security":"tls","tlsSettings":{"serverName":"{{.SNI}}"}}
    },
    {"tag":"direct","protocol":"freedom","settings":{}},
    {"tag":"block","protocol":"blackhole","settings":{}}
  ],
  "routing":{"domainStrategy":"IPIfNonMatch","rules":[
    {"type":"field","ip":["geoip:cn","geoip:private"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:cn"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:geolocation-!cn"],"outboundTag":"proxy"}
  ]}
}`

// XrayVLessConfig 支持 VLess+TCP+TLS/XTLS/WS/gRPC
const XrayVLessConfig = `{
  "log": {"loglevel": "{{.LogLevel}}"},
  "inbounds": [
    {"tag":"http-in","port":10809,"protocol":"http","listen":"127.0.0.1","settings":{}},
    {"tag":"socks-in","port":10808,"protocol":"socks","listen":"127.0.0.1","settings":{"udp":true}},
    {"tag":"tproxy-in","port":12345,"protocol":"dokodemo-door","listen":"0.0.0.0","settings":{"network":"tcp,udp","followRedirect":true},"sniffing":{"enabled":true,"destOverride":["http","tls"]}}
  ],
  "outbounds": [
    {
      "tag":"proxy","protocol":"vless",
      "settings":{"vnext":[{"address":"{{.ServerAddr}}","port":{{.ServerPort}},"users":[{"id":"{{.UUID}}","encryption":"none","flow":"xtls-rprx-vision"}]}]},
      "streamSettings":{"network":"tcp","security":"tls","tlsSettings":{"serverName":"{{.SNI}}","fingerprint":"chrome"}}
    },
    {"tag":"direct","protocol":"freedom","settings":{}},
    {"tag":"block","protocol":"blackhole","settings":{}}
  ],
  "routing":{"domainStrategy":"IPIfNonMatch","rules":[
    {"type":"field","ip":["geoip:cn","geoip:private"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:cn"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:geolocation-!cn"],"outboundTag":"proxy"}
  ]}
}`

// XrayTrojanConfig 支持 Trojan 协议
const XrayTrojanConfig = `{
  "log": {"loglevel": "{{.LogLevel}}"},
  "inbounds": [
    {"tag":"http-in","port":10809,"protocol":"http","listen":"127.0.0.1","settings":{}},
    {"tag":"socks-in","port":10808,"protocol":"socks","listen":"127.0.0.1","settings":{"udp":true}},
    {"tag":"tproxy-in","port":12345,"protocol":"dokodemo-door","listen":"0.0.0.0","settings":{"network":"tcp,udp","followRedirect":true},"sniffing":{"enabled":true,"destOverride":["http","tls"]}}
  ],
  "outbounds": [
    {
      "tag":"proxy","protocol":"trojan",
      "settings":{"servers":[{"address":"{{.ServerAddr}}","port":{{.ServerPort}},"password":"{{.Password}}","level":"0"}]},
      "streamSettings":{"network":"tcp","security":"tls","tlsSettings":{"serverName":"{{.SNI}}"}}
    },
    {"tag":"direct","protocol":"freedom","settings":{}},
    {"tag":"block","protocol":"blackhole","settings":{}}
  ],
  "routing":{"domainStrategy":"IPIfNonMatch","rules":[
    {"type":"field","ip":["geoip:cn","geoip:private"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:cn"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:geolocation-!cn"],"outboundTag":"proxy"}
  ]}
}`

// XrayShadowsocksConfig 支持 Shadowsocks
const XrayShadowsocksConfig = `{
  "log": {"loglevel": "{{.LogLevel}}"},
  "inbounds": [
    {"tag":"http-in","port":10809,"protocol":"http","listen":"127.0.0.1","settings":{}},
    {"tag":"socks-in","port":10808,"protocol":"socks","listen":"127.0.0.1","settings":{"udp":true}},
    {"tag":"tproxy-in","port":12345,"protocol":"dokodemo-door","listen":"0.0.0.0","settings":{"network":"tcp,udp","followRedirect":true},"sniffing":{"enabled":true,"destOverride":["http","tls"]}}
  ],
  "outbounds": [
    {
      "tag":"proxy","protocol":"shadowsocks",
      "settings":{"servers":[{"address":"{{.ServerAddr}}","port":{{.ServerPort}},"method":"{{.Cipher}}","password":"{{.Password}}"}]}
    },
    {"tag":"direct","protocol":"freedom","settings":{}},
    {"tag":"block","protocol":"blackhole","settings":{}}
  ],
  "routing":{"domainStrategy":"IPIfNonMatch","rules":[
    {"type":"field","ip":["geoip:cn","geoip:private"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:cn"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:geolocation-!cn"],"outboundTag":"proxy"}
  ]}
}`

// ==================== Mihomo 配置模板 ====================

// MihomoFullConfig 支持全协议 + 自动分流
const MihomoFullConfig = `# ============================================
# proxyman Mihomo 全协议配置
# 支持: VMess, VLess, Trojan, Shadowsocks,
#       HTTP/HTTPS, SOCKS5, Tuic, Hysteria2
# ============================================

mixed-port: 7890
allow-lan: false
bind-address: "127.0.0.1"
mode: rule
log-level: {{.LogLevel}}
ipv6: false
find-process-mode: strict

# ==================== DNS 配置 ====================
dns:
  enable: true
  listen: 0.0.0.0:1053
  enhanced-mode: fake-ip
  fake-ip-range: 198.18.0.1/16
  fake-ip-filter:
    - "*.lan"
    - "*.local"
    - "*.localhost"
    - "dns.msftncsi.com"
    - "www.msftncsi.com"
    - "www.msftconnecttest.com"
  default-nameserver:
    - 223.5.5.5
    - 119.29.29.29
  nameserver:
    - https://doh.pub/dns-query
    - https://dns.alidns.com/dns-query
  fallback:
    - https://dns.google/dns-query
    - https://cloudflare-dns.com/dns-query
  fallback-filter:
    geoip: true
    geoip-code: CN
    ipcidr:
      - 240.0.0.0/4

# ==================== 代理节点 ====================
proxies:
  # --- VMess ---
  - name: "vmess-node"
    type: vmess
    server: YOUR_SERVER_ADDRESS
    port: 443
    uuid: YOUR_UUID
    alterId: 0
    cipher: auto
    tls: true
    network: ws
    ws-opts:
      path: /path
      headers:
        Host: your.domain.com

  # --- VLess ---
  # - name: "vless-node"
  #   type: vless
  #   server: YOUR_SERVER_ADDRESS
  #   port: 443
  #   uuid: YOUR_UUID
  #   flow: xtls-rprx-vision
  #   tls: true
  #   servername: your.domain.com
  #   client-fingerprint: chrome

  # --- Trojan ---
  # - name: "trojan-node"
  #   type: trojan
  #   server: YOUR_SERVER_ADDRESS
  #   port: 443
  #   password: YOUR_PASSWORD
  #   sni: your.domain.com

  # --- Shadowsocks ---
  # - name: "ss-node"
  #   type: ss
  #   server: YOUR_SERVER_ADDRESS
  #   port: 8388
  #   cipher: aes-256-gcm
  #   password: YOUR_PASSWORD

  # --- SOCKS5 ---
  # - name: "socks5-node"
  #   type: socks5
  #   server: YOUR_SERVER_ADDRESS
  #   port: 1080
  #   username: user
  #   password: pass

  # --- HTTP/HTTPS ---
  # - name: "http-node"
  #   type: http
  #   server: YOUR_SERVER_ADDRESS
  #   port: 8080
  #   username: user
  #   password: pass
  #   tls: true

  # --- Tuic ---
  # - name: "tuic-node"
  #   type: tuic
  #   server: YOUR_SERVER_ADDRESS
  #   port: 443
  #   uuid: YOUR_UUID
  #   password: YOUR_PASSWORD
  #   congestion-control: bbr
  #   tls:
  #     sni: your.domain.com
  #     skip-cert-verify: false

  # --- Hysteria2 ---
  # - name: "hy2-node"
  #   type: hysteria2
  #   server: YOUR_SERVER_ADDRESS
  #   port: 443
  #   password: YOUR_PASSWORD
  #   sni: your.domain.com
  #   up: 100 Mbps
  #   down: 100 Mbps

# ==================== 代理组 ====================
proxy-groups:
  # 自动选择最快的节点
  - name: "Auto"
    type: url-test
    proxies:
      - vmess-node
      - DIRECT
    url: http://www.gstatic.com/generate_204
    interval: 300
    tolerance: 50

  # 手动选择节点
  - name: "Proxy"
    type: select
    proxies:
      - Auto
      - vmess-node
      # - vless-node
      # - trojan-node
      # - ss-node
      # - socks5-node
      # - http-node
      # - tuic-node
      # - hy2-node
      - DIRECT

# ==================== 分流规则 ====================
rules:
  # --- 国内直连 ---
  - GEOSITE,cn,DIRECT
  - GEOIP,cn,DIRECT,no-resolve

  # --- 私有网络直连 ---
  - GEOIP,private,DIRECT,no-resolve
  - GEOSITE,private,DIRECT,no-resolve

  # --- 常用国内域名直连 ---
  - GEOSITE,apple-cn,DIRECT
  - GEOSITE,google-cn,DIRECT
  - GEOSITE,microsoft-cn,DIRECT
  - GEOSITE,steam@cn,DIRECT

  # --- 广告拦截 ---
  - GEOSITE,category-ads-all,REJECT

  # --- 流媒体走代理 ---
  - GEOSITE,netflix,Proxy
  - GEOSITE,youtube,Proxy
  - GEOSITE,google,Proxy
  - GEOSITE,telegram,Proxy
  - GEOSITE,twitter,Proxy
  - GEOSITE,facebook,Proxy
  - GEOSITE,github,Proxy

  # --- 国际网站走代理 ---
  - GEOSITE,geolocation-!cn,Proxy

  # --- 兜底规则 ---
  - MATCH,Proxy
`

// MihomoChinaSplitConfig 专注国内直连分流
const MihomoChinaSplitConfig = `# ============================================
# proxyman 国内直连分流配置
# 特点: 国内流量直连, 国际流量走代理
# ============================================

mixed-port: 7890
allow-lan: false
mode: rule
log-level: {{.LogLevel}}

dns:
  enable: true
  listen: 0.0.0.0:1053
  enhanced-mode: fake-ip
  fake-ip-range: 198.18.0.1/16
  nameserver:
    - 223.5.5.5
    - 119.29.29.29
  fallback:
    - https://dns.google/dns-query
  fallback-filter:
    geoip: true
    geoip-code: CN

proxies:
  - name: "proxy-node"
    type: vmess
    server: YOUR_SERVER_ADDRESS
    port: 443
    uuid: YOUR_UUID
    alterId: 0
    cipher: auto
    tls: true
    network: ws

proxy-groups:
  - name: "Proxy"
    type: select
    proxies:
      - proxy-node
      - DIRECT

rules:
  - GEOSITE,cn,DIRECT
  - GEOIP,cn,DIRECT,no-resolve
  - GEOIP,private,DIRECT,no-resolve
  - GEOSITE,category-ads-all,REJECT
  - MATCH,Proxy
`

// TemplateData holds variables for template rendering.
type TemplateData struct {
	LogLevel    string
	EngineDir   string
	BinaryPath  string
	ConfigPath  string
	Description string
}

// SystemdServiceTemplate is a systemd unit for engine services.
const SystemdServiceTemplate = `[Unit]
Description={{.Description}}
After=network.target

[Service]
Type=simple
WorkingDirectory={{.EngineDir}}
ExecStart={{.BinaryPath}} run -config {{.ConfigPath}}
Restart=on-failure
RestartSec=5
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
`

// ==================== 模板数据 ====================

// ProxyConfigData 代理配置参数
// ProxyConfigData 代理配置参数
// ProxyConfigData 代理配置参数
// ProxyConfigData 代理配置参数
type ProxyConfigData struct {
	Protocol    string
	ServerAddr  string
	ServerPort  int
	UUID        string
	Password    string
	Cipher      string
	SNI         string
	Network     string
	LogLevel    string
	PublicKey   string  // Reality 协议
	ShortId     string  // Reality 协议
	SpiderX     string  // Reality 协议
	GRPCService string  // gRPC 协议
	WSPATH      string  // WebSocket 协议
	Flow        string  // VLess XTLS flow
}

// WriteXrayConfig 根据协议生成 Xray 配置
func WriteXrayConfig(engineDir string, logLevel string, protocol string, data *ProxyConfigData) (string, error) {
	if data == nil {
		data = &ProxyConfigData{
			Protocol:   "vmess",
			ServerAddr: "YOUR_SERVER_ADDRESS",
			ServerPort: 443,
			UUID:       "YOUR_UUID_HERE",
			Network:    "ws",
			SNI:        "your.domain.com",
			LogLevel:   logLevel,
		}
	}
	data.LogLevel = logLevel

	var configTemplate string
	switch protocol {
	case "vmess":
		configTemplate = XrayVMessConfig
	case "vless":
		configTemplate = XrayVLessConfig
	case "trojan":
		configTemplate = XrayTrojanConfig
	case "shadowsocks", "ss":
		configTemplate = XrayShadowsocksConfig
	case "vless-reality", "reality":
		configTemplate = XrayVLessRealityConfig
	case "vless-grpc", "grpc":
		configTemplate = XrayVLessGRPCConfig
	case "trojan-grpc":
		configTemplate = XrayTrojanGRPCConfig
	case "ss2022", "shadowsocks-2022":
		configTemplate = XrayShadowsocks2022Config
	case "anytls":
		configTemplate = XrayAnyTLSConfig
	case "vless-xtls", "xtls":
		configTemplate = XrayVLessXTLSConfig
	case "trojan-ws", "trojan-websocket":
		configTemplate = XrayTrojanWSConfig
	case "vless-ws", "vless-websocket":
		configTemplate = XrayVLessWSConfig
	case "vmess-ws", "vmess-websocket":
		configTemplate = XrayVMessWSConfig
	default:
		configTemplate = XrayVMessConfig
	}

	tmpl, err := template.New("xray").Parse(configTemplate)
	if err != nil {
		return "", fmt.Errorf("parse xray template: %w", err)
	}

	configPath := filepath.Join(engineDir, "config.json")
	f, err := os.Create(configPath)
	if err != nil {
		return "", fmt.Errorf("create config: %w", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return configPath, nil
}

// WriteMihomoConfig 生成 Mihomo 配置（全协议 + 自动分流）
func WriteMihomoConfig(engineDir string, logLevel string, splitMode string) (string, error) {
	var configTemplate string
	switch splitMode {
	case "china-split":
		configTemplate = MihomoChinaSplitConfig
	default:
		configTemplate = MihomoFullConfig
	}

	tmpl, err := template.New("mihomo").Parse(configTemplate)
	if err != nil {
		return "", fmt.Errorf("parse mihomo template: %w", err)
	}

	data := TemplateData{LogLevel: logLevel}
	configPath := filepath.Join(engineDir, "config.yaml")
	f, err := os.Create(configPath)
	if err != nil {
		return "", fmt.Errorf("create config: %w", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return configPath, nil
}

// WriteSystemdService 生成 systemd 服务文件
func WriteSystemdService(engineDir, binaryName, configName, description string) (string, error) {
	data := TemplateData{
		EngineDir:   engineDir,
		BinaryPath:  filepath.Join(engineDir, binaryName),
		ConfigPath:  filepath.Join(engineDir, configName),
		Description: description,
	}

	tmpl, err := template.New("systemd").Parse(SystemdServiceTemplate)
	if err != nil {
		return "", fmt.Errorf("parse systemd template: %w", err)
	}

	serviceName := fmt.Sprintf("proxyman-%s.service", binaryName)
	servicePath := filepath.Join(engineDir, serviceName)
	f, err := os.Create(servicePath)
	if err != nil {
		return "", fmt.Errorf("create service file: %w", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return servicePath, nil
}

// WriteTProxyService 生成 TProxy 透明代理 systemd 服务
func WriteTProxyService(engineDir, binaryName, configName string) (string, error) {
	tproxyService := `[Unit]
Description=proxyman TProxy Transparent Proxy
After=network.target

[Service]
Type=simple
ExecStartPre=/sbin/iptables -t mangle -N PROXYMAN
ExecStartPre=/sbin/iptables -t mangle -A PROXYMAN -d 0.0.0.0/8 -j RETURN
ExecStartPre=/sbin/iptables -t mangle -A PROXYMAN -d 10.0.0.0/8 -j RETURN
ExecStartPre=/sbin/iptables -t mangle -A PROXYMAN -d 127.0.0.0/8 -j RETURN
ExecStartPre=/sbin/iptables -t mangle -A PROXYMAN -d 169.254.0.0/16 -j RETURN
ExecStartPre=/sbin/iptables -t mangle -A PROXYMAN -d 172.16.0.0/12 -j RETURN
ExecStartPre=/sbin/iptables -t mangle -A PROXYMAN -d 192.168.0.0/16 -j RETURN
ExecStartPre=/sbin/iptables -t mangle -A PROXYMAN -d 224.0.0.0/4 -j RETURN
ExecStartPre=/sbin/iptables -t mangle -A PROXYMAN -d 240.0.0.0/4 -j RETURN
ExecStartPre=/sbin/iptables -t mangle -A PREROUTING -p tcp -j PROXYMAN
ExecStartPre=/sbin/iptables -t mangle -A PREROUTING -p udp -j PROXYMAN
ExecStart=%s -d %s
ExecStopPost=/sbin/iptables -t mangle -D PREROUTING -p tcp -j PROXYMAN
ExecStopPost=/sbin/iptables -t mangle -D PREROUTING -p udp -j PROXYMAN
ExecStopPost=/sbin/iptables -t mangle -F PROXYMAN
ExecStopPost=/sbin/iptables -t mangle -X PROXYMAN
Restart=on-failure
RestartSec=5
CapabilityBoundingSet=CAP_NET_ADMIN CAP_NET_RAW
AmbientCapabilities=CAP_NET_ADMIN CAP_NET_RAW

[Install]
WantedBy=multi-user.target
`
	servicePath := filepath.Join(engineDir, "proxyman-tproxy.service")
	f, err := os.Create(servicePath)
	if err != nil {
		return "", fmt.Errorf("create tproxy service: %w", err)
	}
	defer f.Close()

	content := fmt.Sprintf(tproxyService,
		filepath.Join(engineDir, binaryName),
		filepath.Join(engineDir, configName),
	)
	if _, err := f.WriteString(content); err != nil {
		return "", fmt.Errorf("write tproxy service: %w", err)
	}
	return servicePath, nil
}

// ==================== 扩展协议模板 ====================

// XrayVLessRealityConfig 支持 VLess+Reality (最新抗审查协议)
const XrayVLessRealityConfig = `{
  "log": {"loglevel": "{{.LogLevel}}"},
  "inbounds": [
    {"tag":"http-in","port":10809,"protocol":"http","listen":"127.0.0.1","settings":{}},
    {"tag":"socks-in","port":10808,"protocol":"socks","listen":"127.0.0.1","settings":{"udp":true}}
  ],
  "outbounds": [
    {
      "tag":"proxy","protocol":"vless",
      "settings":{"vnext":[{"address":"{{.ServerAddr}}","port":{{.ServerPort}},"users":[{"id":"{{.UUID}}","encryption":"none","flow":"xtls-rprx-vision"}]}]},
      "streamSettings":{"network":"tcp","security":"reality","realitySettings":{"serverName":"{{.SNI}}","fingerprint":"chrome","publicKey":"{{.PublicKey}}","shortId":"{{.ShortId}}","spiderX":"{{.SpiderX}}"}}
    },
    {"tag":"direct","protocol":"freedom","settings":{}},
    {"tag":"block","protocol":"blackhole","settings":{}}
  ],
  "routing":{"domainStrategy":"IPIfNonMatch","rules":[
    {"type":"field","ip":["geoip:cn","geoip:private"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:cn"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:geolocation-!cn"],"outboundTag":"proxy"}
  ]}
}`

// XrayVLessGRPCConfig 支持 VLess+gRPC 高性能传输
const XrayVLessGRPCConfig = `{
  "log": {"loglevel": "{{.LogLevel}}"},
  "inbounds": [
    {"tag":"http-in","port":10809,"protocol":"http","listen":"127.0.0.1","settings":{}},
    {"tag":"socks-in","port":10808,"protocol":"socks","listen":"127.0.0.1","settings":{"udp":true}}
  ],
  "outbounds": [
    {
      "tag":"proxy","protocol":"vless",
      "settings":{"vnext":[{"address":"{{.ServerAddr}}","port":{{.ServerPort}},"users":[{"id":"{{.UUID}}","encryption":"none"}]}]},
      "streamSettings":{"network":"grpc","security":"tls","tlsSettings":{"serverName":"{{.SNI}}"},"grpcSettings":{"serviceName":"{{.GRPCService}}"}}
    },
    {"tag":"direct","protocol":"freedom","settings":{}},
    {"tag":"block","protocol":"blackhole","settings":{}}
  ],
  "routing":{"domainStrategy":"IPIfNonMatch","rules":[
    {"type":"field","ip":["geoip:cn","geoip:private"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:cn"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:geolocation-!cn"],"outboundTag":"proxy"}
  ]}
}`

// XrayTrojanGRPCConfig 支持 Trojan+gRPC 传输
const XrayTrojanGRPCConfig = `{
  "log": {"loglevel": "{{.LogLevel}}"},
  "inbounds": [
    {"tag":"http-in","port":10809,"protocol":"http","listen":"127.0.0.1","settings":{}},
    {"tag":"socks-in","port":10808,"protocol":"socks","listen":"127.0.0.1","settings":{"udp":true}}
  ],
  "outbounds": [
    {
      "tag":"proxy","protocol":"trojan",
      "settings":{"servers":[{"address":"{{.ServerAddr}}","port":{{.ServerPort}},"password":"{{.Password}}"}]},
      "streamSettings":{"network":"grpc","security":"tls","tlsSettings":{"serverName":"{{.SNI}}"},"grpcSettings":{"serviceName":"{{.GRPCService}}"}}
    },
    {"tag":"direct","protocol":"freedom","settings":{}},
    {"tag":"block","protocol":"blackhole","settings":{}}
  ],
  "routing":{"domainStrategy":"IPIfNonMatch","rules":[
    {"type":"field","ip":["geoip:cn","geoip:private"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:cn"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:geolocation-!cn"],"outboundTag":"proxy"}
  ]}
}`

// XrayShadowsocks2022Config 支持 Shadowsocks-2022 新一代协议
const XrayShadowsocks2022Config = `{
  "log": {"loglevel": "{{.LogLevel}}"},
  "inbounds": [
    {"tag":"http-in","port":10809,"protocol":"http","listen":"127.0.0.1","settings":{}},
    {"tag":"socks-in","port":10808,"protocol":"socks","listen":"127.0.0.1","settings":{"udp":true}}
  ],
  "outbounds": [
    {
      "tag":"proxy","protocol":"shadowsocks",
      "settings":{"servers":[{"address":"{{.ServerAddr}}","port":{{.ServerPort}},"method":"2022-blake3-aes-128-gcm","password":"{{.Password}}"}]}
    },
    {"tag":"direct","protocol":"freedom","settings":{}},
    {"tag":"block","protocol":"blackhole","settings":{}}
  ],
  "routing":{"domainStrategy":"IPIfNonMatch","rules":[
    {"type":"field","ip":["geoip:cn","geoip:private"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:cn"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:geolocation-!cn"],"outboundTag":"proxy"}
  ]}
}`

// XrayAnyTLSConfig 支持 AnyTLS 协议 (基于 TLS 的新型代理协议)
const XrayAnyTLSConfig = `{
  "log": {"loglevel": "{{.LogLevel}}"},
  "inbounds": [
    {"tag":"http-in","port":10809,"protocol":"http","listen":"127.0.0.1","settings":{}},
    {"tag":"socks-in","port":10808,"protocol":"socks","listen":"127.0.0.1","settings":{"udp":true}}
  ],
  "outbounds": [
    {
      "tag":"proxy","protocol":"anytls",
      "settings":{"servers":[{"address":"{{.ServerAddr}}","port":{{.ServerPort}},"password":"{{.Password}}","serverName":"{{.SNI}}"}]},
      "streamSettings":{"network":"tcp","security":"tls","tlsSettings":{"serverName":"{{.SNI}}","fingerprint":"chrome"}}
    },
    {"tag":"direct","protocol":"freedom","settings":{}},
    {"tag":"block","protocol":"blackhole","settings":{}}
  ],
  "routing":{"domainStrategy":"IPIfNonMatch","rules":[
    {"type":"field","ip":["geoip:cn","geoip:private"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:cn"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:geolocation-!cn"],"outboundTag":"proxy"}
  ]}
}`

// ProxyConfigData 扩展字段
// PublicKey 用于 Reality 协议
// ShortId 用于 Reality 协议
// SpiderX 用于 Reality 协议
// GRPCService 用于 gRPC 协议

// ==================== Mihomo 扩展协议 ====================

// MihomoAnyTLSConfig 支持 AnyTLS 协议
const MihomoAnyTLSConfig = `# ============================================
# proxyman Mihomo AnyTLS 配置
# ============================================

mixed-port: 7890
allow-lan: false
mode: rule
log-level: {{.LogLevel}}

dns:
  enable: true
  enhanced-mode: fake-ip
  nameserver:
    - 223.5.5.5
    - 119.29.29.29
  fallback:
    - https://dns.google/dns-query
  fallback-filter:
    geoip: true
    geoip-code: CN

proxies:
  - name: "anytls-node"
    type: anytls
    server: YOUR_SERVER_ADDRESS
    port: 443
    password: YOUR_PASSWORD
    sni: your.domain.com
    client-fingerprint: chrome

proxy-groups:
  - name: "Proxy"
    type: select
    proxies:
      - anytls-node
      - DIRECT

rules:
  - GEOSITE,cn,DIRECT
  - GEOIP,cn,DIRECT,no-resolve
  - GEOIP,private,DIRECT,no-resolve
  - MATCH,Proxy
`

// MihomoHysteriaConfig 支持 Hysteria 协议
const MihomoHysteriaConfig = `# ============================================
# proxyman Mihomo Hysteria 配置
# ============================================

mixed-port: 7890
allow-lan: false
mode: rule
log-level: {{.LogLevel}}

dns:
  enable: true
  enhanced-mode: fake-ip
  nameserver:
    - 223.5.5.5
    - 119.29.29.29
  fallback:
    - https://dns.google/dns-query
  fallback-filter:
    geoip: true
    geoip-code: CN

proxies:
  - name: "hysteria-node"
    type: hysteria
    server: YOUR_SERVER_ADDRESS
    port: 443
    ports: 443
    password: YOUR_PASSWORD
    sni: your.domain.com
    up: 100 Mbps
    down: 100 Mbps

proxy-groups:
  - name: "Proxy"
    type: select
    proxies:
      - hysteria-node
      - DIRECT

rules:
  - GEOSITE,cn,DIRECT
  - GEOIP,cn,DIRECT,no-resolve
  - GEOIP,private,DIRECT,no-resolve
  - MATCH,Proxy
`

// MihomoTuicConfig 支持 Tuic 协议
const MihomoTuicConfig = `# ============================================
# proxyman Mihomo Tuic 配置
# ============================================

mixed-port: 7890
allow-lan: false
mode: rule
log-level: {{.LogLevel}}

dns:
  enable: true
  enhanced-mode: fake-ip
  nameserver:
    - 223.5.5.5
    - 119.29.29.29
  fallback:
    - https://dns.google/dns-query
  fallback-filter:
    geoip: true
    geoip-code: CN

proxies:
  - name: "tuic-node"
    type: tuic
    server: YOUR_SERVER_ADDRESS
    port: 443
    uuid: YOUR_UUID
    password: YOUR_PASSWORD
    congestion-control: bbr
    tls:
      sni: your.domain.com
      skip-cert-verify: false
      fingerprint: chrome

proxy-groups:
  - name: "Proxy"
    type: select
    proxies:
      - tuic-node
      - DIRECT

rules:
  - GEOSITE,cn,DIRECT
  - GEOIP,cn,DIRECT,no-resolve
  - GEOIP,private,DIRECT,no-resolve
  - MATCH,Proxy
`

// ==================== 更多协议模板 ====================

// XrayVLessXTLSConfig 支持 VLess+XTLS 高性能传输
const XrayVLessXTLSConfig = `{
  "log": {"loglevel": "{{.LogLevel}}"},
  "inbounds": [
    {"tag":"http-in","port":10809,"protocol":"http","listen":"127.0.0.1","settings":{}},
    {"tag":"socks-in","port":10808,"protocol":"socks","listen":"127.0.0.1","settings":{"udp":true}}
  ],
  "outbounds": [
    {
      "tag":"proxy","protocol":"vless",
      "settings":{"vnext":[{"address":"{{.ServerAddr}}","port":{{.ServerPort}},"users":[{"id":"{{.UUID}}","encryption":"none","flow":"xtls-rprx-direct"}]}]},
      "streamSettings":{"network":"tcp","security":"xtls","xtlsSettings":{"serverName":"{{.SNI}}","fingerprint":"chrome"}}
    },
    {"tag":"direct","protocol":"freedom","settings":{}},
    {"tag":"block","protocol":"blackhole","settings":{}}
  ],
  "routing":{"domainStrategy":"IPIfNonMatch","rules":[
    {"type":"field","ip":["geoip:cn","geoip:private"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:cn"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:geolocation-!cn"],"outboundTag":"proxy"}
  ]}
}`

// XrayTrojanWSConfig 支持 Trojan+WebSocket 传输
const XrayTrojanWSConfig = `{
  "log": {"loglevel": "{{.LogLevel}}"},
  "inbounds": [
    {"tag":"http-in","port":10809,"protocol":"http","listen":"127.0.0.1","settings":{}},
    {"tag":"socks-in","port":10808,"protocol":"socks","listen":"127.0.0.1","settings":{"udp":true}}
  ],
  "outbounds": [
    {
      "tag":"proxy","protocol":"trojan",
      "settings":{"servers":[{"address":"{{.ServerAddr}}","port":{{.ServerPort}},"password":"{{.Password}}"}]},
      "streamSettings":{"network":"ws","security":"tls","tlsSettings":{"serverName":"{{.SNI}}"},"wsSettings":{"path":"{{.WSPATH}}","headers":{"Host":"{{.SNI}}"}}}
    },
    {"tag":"direct","protocol":"freedom","settings":{}},
    {"tag":"block","protocol":"blackhole","settings":{}}
  ],
  "routing":{"domainStrategy":"IPIfNonMatch","rules":[
    {"type":"field","ip":["geoip:cn","geoip:private"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:cn"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:geolocation-!cn"],"outboundTag":"proxy"}
  ]}
}`

// XrayVLessWSConfig 支持 VLess+WebSocket+TLS 传输
const XrayVLessWSConfig = `{
  "log": {"loglevel": "{{.LogLevel}}"},
  "inbounds": [
    {"tag":"http-in","port":10809,"protocol":"http","listen":"127.0.0.1","settings":{}},
    {"tag":"socks-in","port":10808,"protocol":"socks","listen":"127.0.0.1","settings":{"udp":true}}
  ],
  "outbounds": [
    {
      "tag":"proxy","protocol":"vless",
      "settings":{"vnext":[{"address":"{{.ServerAddr}}","port":{{.ServerPort}},"users":[{"id":"{{.UUID}}","encryption":"none"}]}]},
      "streamSettings":{"network":"ws","security":"tls","tlsSettings":{"serverName":"{{.SNI}}"},"wsSettings":{"path":"{{.WSPATH}}","headers":{"Host":"{{.SNI}}"}}}
    },
    {"tag":"direct","protocol":"freedom","settings":{}},
    {"tag":"block","protocol":"blackhole","settings":{}}
  ],
  "routing":{"domainStrategy":"IPIfNonMatch","rules":[
    {"type":"field","ip":["geoip:cn","geoip:private"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:cn"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:geolocation-!cn"],"outboundTag":"proxy"}
  ]}
}`

// XrayVMessWSConfig 支持 VMess+WebSocket+TLS 传输
const XrayVMessWSConfig = `{
  "log": {"loglevel": "{{.LogLevel}}"},
  "inbounds": [
    {"tag":"http-in","port":10809,"protocol":"http","listen":"127.0.0.1","settings":{}},
    {"tag":"socks-in","port":10808,"protocol":"socks","listen":"127.0.0.1","settings":{"udp":true}}
  ],
  "outbounds": [
    {
      "tag":"proxy","protocol":"vmess",
      "settings":{"vnext":[{"address":"{{.ServerAddr}}","port":{{.ServerPort}},"users":[{"id":"{{.UUID}}","alterId":0,"security":"auto"}]}]},
      "streamSettings":{"network":"ws","security":"tls","tlsSettings":{"serverName":"{{.SNI}}"},"wsSettings":{"path":"{{.WSPATH}}","headers":{"Host":"{{.SNI}}"}}}
    },
    {"tag":"direct","protocol":"freedom","settings":{}},
    {"tag":"block","protocol":"blackhole","settings":{}}
  ],
  "routing":{"domainStrategy":"IPIfNonMatch","rules":[
    {"type":"field","ip":["geoip:cn","geoip:private"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:cn"],"outboundTag":"direct"},
    {"type":"field","domain":["geosite:geolocation-!cn"],"outboundTag":"proxy"}
  ]}
}`
