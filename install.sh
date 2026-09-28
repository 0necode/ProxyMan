#!/bin/bash
# ============================================
# ProxyMan 一键安装脚本
# 支持: Linux (amd64)
# ============================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 版本信息
VERSION="1.0.0"
REPO="yourusername/ProxyMan"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.config/ProxyMan"

# 打印信息
info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# 检查系统
check_system() {
    info "检查系统环境..."
    
    # 检查操作系统
    if [[ "$OSTYPE" != "linux-gnu"* ]]; then
        error "仅支持 Linux 系统"
    fi
    
    # 检查架构
    ARCH=$(uname -m)
    if [[ "$ARCH" != "x86_64" ]]; then
        error "仅支持 amd64 架构"
    fi
    
    # 检查必要工具
    for cmd in curl wget tar; do
        if ! command -v $cmd &> /dev/null; then
            warn "$cmd 未安装，尝试安装..."
            install_deps $cmd
        fi
    done
    
    success "系统检查通过"
}

# 安装依赖
install_deps() {
    local cmd=$1
    
    if command -v apt-get &> /dev/null; then
        apt-get update && apt-get install -y $cmd
    elif command -v yum &> /dev/null; then
        yum install -y $cmd
    elif command -v dnf &> /dev/null; then
        dnf install -y $cmd
    elif command -v pacman &> /dev/null; then
        pacman -Sy --noconfirm $cmd
    else
        error "无法自动安装 $cmd，请手动安装"
    fi
}

# 下载二进制
download_binary() {
    info "下载 ProxyMan v${VERSION}..."
    
    # 尝试 GitHub 镜像
    local urls=(
        "https://ghfast.top/https://github.com/${REPO}/releases/download/v${VERSION}/ProxyMan-linux-amd64.tar.gz"
        "https://ghproxy.net/https://github.com/${REPO}/releases/download/v${VERSION}/ProxyMan-linux-amd64.tar.gz"
        "https://github.com/${REPO}/releases/download/v${VERSION}/ProxyMan-linux-amd64.tar.gz"
    )
    
    local downloaded=false
    for url in "${urls[@]}"; do
        info "尝试下载: $url"
        if curl -fsSL "$url" -o /tmp/ProxyMan.tar.gz 2>/dev/null || \
           wget -q "$url" -O /tmp/ProxyMan.tar.gz 2>/dev/null; then
            downloaded=true
            break
        fi
    done
    
    if [ "$downloaded" = false ]; then
        error "下载失败，请检查网络连接"
    fi
    
    success "下载完成"
}

# 安装二进制
install_binary() {
    info "安装 ProxyMan..."
    
    # 创建临时目录
    mkdir -p /tmp/ProxyMan-install
    cd /tmp/ProxyMan-install
    
    # 解压
    tar -xzf /tmp/ProxyMan.tar.gz
    
    # 安装到系统路径
    if [ -w "$INSTALL_DIR" ]; then
        cp ProxyMan-linux-amd64 "$INSTALL_DIR/ProxyMan"
    else
        sudo cp ProxyMan-linux-amd64 "$INSTALL_DIR/ProxyMan"
    fi
    
    chmod +x "$INSTALL_DIR/ProxyMan"
    
    # 清理
    rm -rf /tmp/ProxyMan-install /tmp/ProxyMan.tar.gz
    
    success "安装完成"
}

# 创建配置目录
setup_config() {
    info "创建配置目录..."
    
    mkdir -p "$CONFIG_DIR/engines/v2ray"
    mkdir -p "$CONFIG_DIR/engines/clash"
    
    success "配置目录创建完成"
}

# 验证安装
verify_installation() {
    info "验证安装..."
    
    if command -v ProxyMan &> /dev/null; then
        local version=$(ProxyMan --version 2>/dev/null || echo "unknown")
        success "ProxyMan 已安装: $version"
    else
        error "安装验证失败"
    fi
}

# 打印使用说明
print_usage() {
    echo ""
    echo -e "${GREEN}============================================${NC}"
    echo -e "${GREEN}  ProxyMan 安装完成！${NC}"
    echo -e "${GREEN}============================================${NC}"
    echo ""
    echo "使用方法:"
    echo ""
    echo "1. 安装引擎:"
    echo "   ProxyMan install v2ray"
    echo "   ProxyMan install clash"
    echo ""
    echo "2. 配置代理:"
    echo "   ProxyMan config edit v2ray"
    echo "   ProxyMan config edit clash"
    echo ""
    echo "3. 启动代理:"
    echo "   ProxyMan start v2ray ~/.config/ProxyMan/engines/v2ray/config.json"
    echo "   ProxyMan start clash ~/.config/ProxyMan/engines/clash/config.yaml"
    echo ""
    echo "4. 开启系统代理:"
    echo "   ProxyMan system-proxy enable"
    echo ""
    echo "5. 分流模式:"
    echo "   ProxyMan split clash china-split"
    echo ""
    echo "6. 查看帮助:"
    echo "   ProxyMan --help"
    echo ""
    echo "文档: https://github.com/${REPO}"
    echo ""
}

# 主函数
main() {
    echo -e "${GREEN}============================================${NC}"
    echo -e "${GREEN}  ProxyMan 安装脚本 v${VERSION}${NC}"
    echo -e "${GREEN}============================================${NC}"
    echo ""
    
    check_system
    download_binary
    install_binary
    setup_config
    verify_installation
    print_usage
}

# 执行主函数
main "$@"
