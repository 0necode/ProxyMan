# proxy-cli 项目总结

## 📊 项目概况

**proxy-cli** 是一个一站式 Linux 网络代理 CLI 工具，统一管理 Xray (V2Ray) 和 Mihomo (Clash Meta)。

## 🎯 核心功能

### 1. 全协议支持

#### V2Ray/Xray 协议
- ✅ VMess
- ✅ VLess
- ✅ Trojan
- ✅ Shadowsocks
- ✅ SOCKS5
- ✅ HTTP

#### Clash/Mihomo 协议
- ✅ VMess
- ✅ VLess
- ✅ Trojan
- ✅ Shadowsocks
- ✅ HTTP/HTTPS
- ✅ SOCKS5
- ✅ Tuic
- ✅ Hysteria2

### 2. 智能分流
- 🇨🇳 **国内直连**：GEOIP/CN 规则自动直连
- 🌍 **国际代理**：海外流量自动走代理
- 🛡️ **广告拦截**：内置广告过滤规则

### 3. 一键管理
- 📦 自动下载官方二进制
- 🔐 SHA256 校验确保安全
- 🚀 ghfast.top 镜像加速
- ⚙️ systemd 服务集成
- 📝 配置文件模板自动生成

## 📁 项目结构

```
proxy-cli/
├── main.go                    # 入口文件
├── go.mod                     # Go 模块定义
├── go.sum                     # 依赖锁定
├── Makefile                   # 构建自动化
├── README.md                  # 项目文档
├── LICENSE                    # MIT 许可证
├── CHANGELOG.md              # 版本历史
├── CONTRIBUTING.md           # 贡献指南
├── .gitignore               # Git 忽略规则
└── internal/
    ├── cmd/
    │   ├── proxy-cli.go      # CLI 入口
    │   └── commands.go       # 命令实现
    ├── config/
    │   ├── config.go         # 配置管理
    │   └── system.go         # 系统代理管理
    ├── engine/
    │   ├── manager.go        # 引擎管理器
    │   ├── download.go       # GitHub 下载 + SHA256
    │   ├── templates.go      # 配置模板
    │   └── engine_test.go    # 单元测试
    └── rules/
        └── matcher.go        # 规则引擎
```

## 🚀 快速开始

### 安装

```bash
# 克隆仓库
git clone https://github.com/yourusername/proxy-cli.git
cd proxy-cli

# 编译
make build

# 安装到系统路径
make install
```

### 使用

```bash
# 安装引擎
proxy-cli install v2ray
proxy-cli install clash

# 配置代理服务器
proxy-cli config edit v2ray
proxy-cli config edit clash

# 启动代理
proxy-cli start v2ray ~/.config/proxy-cli/engines/v2ray/config.json
proxy-cli start clash ~/.config/proxy-cli/engines/clash/config.yaml

# 开启系统代理
proxy-cli system-proxy enable

# 分流模式
proxy-cli split clash china-split
proxy-cli split clash full
```

## 📊 命令清单

| 命令 | 说明 |
|------|------|
| `proxy-cli install <v2ray\|clash>` | 安装引擎 |
| `proxy-cli uninstall <v2ray\|clash>` | 卸载引擎 |
| `proxy-cli start <engine> <config>` | 启动引擎 |
| `proxy-cli stop <engine>` | 停止引擎 |
| `proxy-cli status` | 查看引擎状态 |
| `proxy-cli config show` | 查看配置 |
| `proxy-cli config set <key> <value>` | 设置配置 |
| `proxy-cli config edit <engine>` | 编辑配置 |
| `proxy-cli split <engine> <mode>` | 分流模式 |
| `proxy-cli test <engine>` | 测试连接 |
| `proxy-cli system-proxy enable/disable/status` | 系统代理 |

## 🛡️ 安全特性

- **SHA256 校验**：所有下载的二进制都经过哈希校验
- **官方源下载**：从 GitHub 官方 release 下载
- **镜像加速**：支持 ghfast.top 镜像加速
- **配置隔离**：每个引擎独立配置目录

## 📊 系统要求

- **操作系统**：Linux (amd64)
- **Go 版本**：1.26+ (编译时)
- **磁盘空间**：约 150MB (两个引擎)
- **依赖**：无 (静态编译)

## 📈 版本历史

### v1.0.0 (2024-01-01)
- ✅ 初始发布
- ✅ 全协议支持
- ✅ 智能分流
- ✅ 一键管理
- ✅ 单元测试

## 🤝 贡献

欢迎贡献代码！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解详情。

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [Xray-core](https://github.com/XTLS/Xray-core) - V2Ray 核心
- [Mihomo](https://github.com/MetaCubeX/mihomo) - Clash Meta 核心
- [GeoIP](https://github.com/Loyalsoldier/geoip) - IP 地理位置数据库
- [GeoSite](https://github.com/v2fly/domain-list-community) - 域名分类数据库

---

**项目状态：✅ 完成并可用**
