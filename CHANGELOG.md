# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-01-01

### Added

- **Core Features**
  - Xray-core (V2Ray) support with VMess, VLess, Trojan, Shadowsocks protocols
  - Mihomo (Clash Meta) support with VMess, VLess, Trojan, Shadowsocks, HTTP/HTTPS, SOCKS5, Tuic, Hysteria2 protocols
  - Unified CLI for managing both engines
  - Automatic engine download from GitHub releases
  - SHA256 verification for all downloads
  - GitHub mirror acceleration (ghfast.top)

- **Traffic Splitting**
  - Full proxy mode (all traffic through proxy)
  - China-split mode (domestic direct, international through proxy)
  - GEOIP/CN rules for automatic domestic traffic detection
  - GEOSITE rules for domain-based routing

- **Configuration Management**
  - Auto-generated config templates for both engines
  - Xray JSON config with inbound/outbound/routing
  - Mihomo YAML config with DNS, proxies, rules
  - Systemd service file generation

- **System Integration**
  - System proxy configuration (environment variables)
  - Systemd service support
  - Config file editing via $EDITOR

- **Commands**
  - `install` - Install engine binaries
  - `uninstall` - Remove engine binaries
  - `start` - Start engine with config
  - `stop` - Stop running engine
  - `status` - Show engine status
  - `config show/set/edit` - Configuration management
  - `split` - Traffic splitting configuration
  - `test` - Test engine configuration
  - `system-proxy` - System proxy management

- **Testing**
  - Unit tests for config generation
  - Unit tests for service file generation
  - Unit tests for protocol-specific configs

### Security

- SHA256 verification for all downloaded binaries
- Official GitHub release sources only
- Mirror acceleration for faster downloads in China
- Config file isolation per engine

## [Unreleased]

### Planned

- Support for more proxy protocols
- GUI configuration tool
- Multi-user support
- Traffic statistics
- Rule subscription support
- Docker container support

---

## Version History

- **v1.0.0** - Initial release with full protocol support
- **v0.9.0** - Beta testing
- **v0.8.0** - Alpha testing
- **v0.7.0** - Core development
- **v0.6.0** - Feature implementation
- **v0.5.0** - Basic framework
- **v0.4.0** - Design phase
- **v0.3.0** - Architecture planning
- **v0.2.0** - Requirements gathering
- **v0.1.0** - Project inception
