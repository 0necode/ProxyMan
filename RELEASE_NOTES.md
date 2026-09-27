# proxyman v1.0.0 Release Notes

## 🎉 首个稳定版本发布

**proxyman** 是一个一站式 Linux 网络代理 CLI 工具，统一管理 Xray (V2Ray) 和 Mihomo (Clash Meta)。

---

## ✨ 核心特性

### 🔌 全协议支持

#### V2Ray/Xray 协议 (12+ 种)
- **VMess** - V2Ray 主力协议
- **VLess** - 新一代轻量协议
- **VLess+Reality** - 最新抗审查协议 🔥
- **VLess+gRPC** - 高性能传输 🔥
- **VLess+XTLS** - 高性能传输
- **VLESS+WS+TLS** - WebSocket 传输
- **Trojan** - 伪装 HTTPS 协议
- **Trojan+gRPC** - Trojan + gRPC
- **Trojan+WS** - Trojan + WebSocket
- **Shadowsocks** - 经典加密代理
- **Shadowsocks-2022** - 新一代 SS 协议 🔥
- **AnyTLS** - 基于 TLS 的新型协议 🔥

#### Clash/Mihomo 协议 (10+ 种)
- **VMess** - V2Ray 主力协议
- **VLess** - 新一代轻量协议
- **Trojan** - 伪装 HTTPS 协议
- **Shadowsocks** - 经典加密代理
- **HTTP/HTTPS** - HTTP 代理
- **SOCKS5** - SOCKS5 代理
- **Tuic** - 基于 QUIC 的协议 🔥
- **Hysteria** - 高性能 UDP 协议 🔥
- **Hysteria2** - Hysteria 升级版 🔥
- **AnyTLS** - 基于 TLS 的新型协议 🔥

### 🔗 协议自动检测
- 自动识别 `vmess://` 链接
- 自动识别 `vless://` 链接
- 自动识别 `trojan://` 链接
- 自动识别 `ss://` 链接
- 自动识别 `hysteria://` / `hy2://` 链接
- 自动识别 `tuic://` 链接
- 自动识别 `anytls://` 链接

### 📡 订阅地址解析
- 支持机场订阅链接
- 自动解析所有节点
- 批量导入代理配置

### 🚀 多样化安装
- **一键安装脚本** - `curl | bash`
- **Git Clone** - `git clone + make install`
- **直接下载** - `wget` 二进制文件

### 🛡️ 智能分流
- **全代理模式** - 所有流量走代理
- **国内直连模式** - 国内流量直连，国际流量走代理
- **GEOIP/CN 规则** - 自动识别国内流量
- **GEOSITE 规则** - 基于域名的路由

### 🔒 安全特性
- **SHA256 校验** - 所有下载的二进制都经过哈希校验
- **官方源下载** - 从 GitHub 官方 release 下载
- **镜像加速** - 支持 ghfast.top 镜像加速
- **配置隔离** - 每个引擎独立配置目录

---

## 📦 安装方式

### 方式一：一键安装脚本
```bash
curl -fsSL https://raw.githubusercontent.com/geek0ne/proxyman/main/install.sh | bash
```

### 方式二：Git Clone
```bash
git clone https://github.com/geek0ne/proxyman.git
cd proxyman
make build
sudo make install
```

### 方式三：直接下载
```bash
# 下载
wget https://github.com/geek0ne/proxyman/releases/download/v1.0.0/proxyman-1.0.0-linux-amd64.tar.gz

# 解压
tar -xzf proxyman-1.0.0-linux-amd64.tar.gz

# 安装
chmod +x proxyman-linux-amd64
sudo cp proxyman-linux-amd64 /usr/local/bin/proxyman
```

---

## 🚀 快速开始

### 1. 安装引擎
```bash
proxyman install v2ray
proxyman install clash
```

### 2. 导入机场订阅
```bash
proxyman import "https://your-airport.com/sub?token=xxx"
```

### 3. 导入代理链接
```bash
proxyman import "vmess://eyJ2IjoiMiIs..."
proxyman import "vless://uuid@server:port?security=tls#name"
```

### 4. 配置分流模式
```bash
proxyman split clash china-split   # 国内直连模式
proxyman split clash full          # 全代理模式
```

### 5. 启动代理
```bash
proxyman start v2ray ~/.config/proxyman/engines/v2ray/config.json
proxyman start clash ~/.config/proxyman/engines/clash/config.yaml
```

### 6. 开启系统代理
```bash
proxyman system-proxy enable
```

---

## 📊 命令清单

| 命令 | 说明 |
|------|------|
| `proxyman install <v2ray\|clash>` | 安装引擎 |
| `proxyman uninstall <v2ray\|clash>` | 卸载引擎 |
| `proxyman start <engine> <config>` | 启动引擎 |
| `proxyman stop <engine>` | 停止引擎 |
| `proxyman status` | 查看引擎状态 |
| `proxyman config show` | 查看配置 |
| `proxyman config set <key> <value>` | 设置配置 |
| `proxyman config edit <engine>` | 编辑配置文件 |
| `proxyman split <engine> <mode>` | 分流模式 |
| `proxyman test <engine>` | 测试连接 |
| `proxyman import <link>` | 导入代理链接 |
| `proxyman system-proxy enable/disable/status` | 系统代理 |

---

## 📁 文件结构

```
proxyman/
├── README.md                  # 项目文档
├── LICENSE                    # MIT 许可证
├── CHANGELOG.md              # 版本历史
├── CONTRIBUTING.md           # 贡献指南
├── Makefile                  # 构建自动化
├── install.sh                # 一键安装脚本
├── go.mod                    # Go 模块定义
├── go.sum                    # 依赖锁定
├── main.go                   # 入口文件
└── internal/
    ├── cmd/                  # CLI 命令
    ├── config/               # 配置管理
    ├── engine/               # 引擎管理
    ├── parser/               # 协议解析
    └── rules/                # 规则引擎
```

---

## 🛡️ 安全说明

- 所有下载的二进制都经过 SHA256 校验
- 从 GitHub 官方 release 下载
- 支持 ghfast.top 镜像加速
- 配置文件隔离存储

---

## 🐛 已知问题

- 无

---

## 📞 反馈

- Issues: [GitHub Issues](https://github.com/geek0ne/proxyman/issues)
- Email: nzl9100@gmail.com

---

## 🙏 致谢

- [Xray-core](https://github.com/XTLS/Xray-core) - V2Ray 核心
- [Mihomo](https://github.com/MetaCubeX/mihomo) - Clash Meta 核心
- [GeoIP](https://github.com/Loyalsoldier/geoip) - IP 地理位置数据库
- [GeoSite](https://github.com/v2fly/domain-list-community) - 域名分类数据库

---

**如果这个项目对你有帮助，请给个 ⭐ Star 支持一下！**
