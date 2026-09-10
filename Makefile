.PHONY: build clean test install uninstall

# 项目信息
APP_NAME := proxy-cli
VERSION := 1.0.0
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell go version | cut -d " " -f 3)

# 编译参数
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# 默认目标
all: build

# 编译
build:
	@echo "Building $(APP_NAME)..."
	go build $(LDFLAGS) -o $(APP_NAME) .
	@echo "Build complete: $(APP_NAME)"

# 清理
clean:
	@echo "Cleaning..."
	rm -f $(APP_NAME)
	rm -rf dist/
	@echo "Clean complete"

# 测试
test:
	@echo "Running tests..."
	go test -v ./...
	@echo "Tests complete"

# 安装到系统路径
install: build
	@echo "Installing $(APP_NAME) to /usr/local/bin..."
	sudo cp $(APP_NAME) /usr/local/bin/
	@echo "Install complete"

# 卸载
uninstall:
	@echo "Uninstalling $(APP_NAME)..."
	sudo rm -f /usr/local/bin/$(APP_NAME)
	@echo "Uninstall complete"

# 格式化代码
fmt:
	@echo "Formatting code..."
	gofmt -w .
	@echo "Format complete"

# 静态检查
lint:
	@echo "Running go vet..."
	go vet ./...
	@echo "Lint complete"

# 生成版本信息
version:
	@echo "$(APP_NAME) $(VERSION)"
	@echo "Build time: $(BUILD_TIME)"
	@echo "Go version: $(GO_VERSION)"

# 交叉编译
build-linux:
	@echo "Building for Linux amd64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(APP_NAME)-linux-amd64 .
	@echo "Build complete: $(APP_NAME)-linux-amd64"

# 发布准备
release: clean build-linux
	@echo "Creating release archive..."
	tar -czf dist/$(APP_NAME)-$(VERSION)-linux-amd64.tar.gz $(APP_NAME)-linux-amd64 README.md LICENSE
	@echo "Release archive created: dist/$(APP_NAME)-$(VERSION)-linux-amd64.tar.gz"

# 帮助
help:
	@echo "Available commands:"
	@echo "  make build      - Build the application"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make test       - Run tests"
	@echo "  make install    - Install to /usr/local/bin"
	@echo "  make uninstall  - Uninstall from /usr/local/bin"
	@echo "  make fmt        - Format code"
	@echo "  make lint       - Run go vet"
	@echo "  make version    - Show version info"
	@echo "  make build-linux - Cross-compile for Linux"
	@echo "  make release    - Create release archive"
	@echo "  make help       - Show this help"
