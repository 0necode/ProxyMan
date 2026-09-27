# Contributing to proxyman

Thank you for your interest in contributing to proxyman! This document provides guidelines and information for contributors.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Environment](#development-environment)
- [How to Contribute](#how-to-contribute)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)
- [Community](#community)

## 📜 Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to [nzl9100@gmail.com](mailto:nzl9100@gmail.com).

## 🚀 Getting Started

### Prerequisites

- Go 1.26 or higher
- Git
- Make (optional, for build commands)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/geek0ne/proxyman.git
   cd proxyman
   ```

3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/originalusername/proxyman.git
   ```

## 🛠️ Development Environment

### Setup

```bash
# Install dependencies
go mod tidy

# Build the project
make build

# Run tests
make test

# Format code
make fmt

# Run linter
make lint
```

### Project Structure

```
proxyman/
├── main.go                    # Entry point
├── go.mod                     # Go module definition
├── Makefile                   # Build automation
├── README.md                  # Project documentation
├── LICENSE                    # MIT License
├── CHANGELOG.md              # Version history
├── CONTRIBUTING.md           # This file
├── .gitignore               # Git ignore rules
└── internal/
    ├── cmd/
    │   ├── proxyman.go      # CLI entry point
    │   └── commands.go       # Command implementations
    ├── config/
    │   ├── config.go         # Configuration management
    │   └── system.go         # System proxy management
    ├── engine/
    │   ├── manager.go        # Engine manager
    │   ├── download.go       # GitHub download + SHA256
    │   ├── templates.go      # Config templates
    │   └── engine_test.go    # Unit tests
    └── rules/
        └── matcher.go        # Rule engine
```

## 🤝 How to Contribute

### Reporting Bugs

- Use the GitHub Issues tracker
- Include a clear title and description
- Provide steps to reproduce the issue
- Include your environment details (OS, Go version, etc.)

### Suggesting Enhancements

- Use the GitHub Issues tracker with the "enhancement" label
- Describe the feature you'd like to see
- Explain why this feature would be useful
- Provide examples if possible

### Contributing Code

1. **Find an issue** to work on or create one
2. **Fork** the repository
3. **Create a branch** for your feature/fix
4. **Make your changes**
5. **Test** your changes
6. **Commit** with a clear message
7. **Push** to your fork
8. **Create a Pull Request**

## 📝 Pull Request Process

### Before Submitting

- [ ] Code follows the project's coding standards
- [ ] Tests pass (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] No lint errors (`make lint`)
- [ ] Documentation is updated if needed
- [ ] CHANGELOG.md is updated

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Enhancement
- [ ] Documentation update
- [ ] Refactoring
- [ ] Other

## Testing
- [ ] Unit tests added/updated
- [ ] Manual testing performed
- [ ] Tested on Linux amd64

## Checklist
- [ ] Code follows project style
- [ ] Self-reviewed the code
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] No new warnings
- [ ] Tests pass
```

## 🎨 Coding Standards

### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` to format code
- Use `go vet` to check for errors
- Write meaningful variable and function names
- Add comments for exported functions
- Handle errors properly

### Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests after the first line

Example:
```
Add support for Trojan protocol

- Implement Trojan config template for Xray
- Add Trojan support to Mihomo config
- Update documentation with Trojan examples

Fixes #123
```

### Code Review

- All submissions require review before merging
- We may suggest changes, improvements, or alternatives
- Maintainers may request changes before approving
- Once approved, a maintainer will merge the PR

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run specific test
go test -v ./internal/engine/...

# Run with coverage
go test -cover ./...

# Run with race detector
go test -race ./...
```

### Writing Tests

- Test files should be in the same package as the code
- Use table-driven tests for multiple test cases
- Test both success and failure scenarios
- Mock external dependencies when needed

Example:
```go
func TestWriteXrayConfig(t *testing.T) {
    tests := []struct {
        name     string
        logLevel string
        protocol string
        wantErr  bool
    }{
        {"valid vmess", "info", "vmess", false},
        {"valid vless", "info", "vless", false},
        {"invalid protocol", "info", "invalid", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            dir := t.TempDir()
            _, err := WriteXrayConfig(dir, tt.logLevel, tt.protocol, nil)
            if (err != nil) != tt.wantErr {
                t.Errorf("WriteXrayConfig() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## 📚 Documentation

### Code Documentation

- Add comments for exported functions, types, and constants
- Use godoc format for documentation
- Include examples in documentation

### User Documentation

- Update README.md for new features
- Add usage examples
- Update CHANGELOG.md

### API Documentation

- Document all public APIs
- Include parameter descriptions
- Provide usage examples

## 🌐 Community

### Getting Help

- Open an issue on GitHub
- Join our discussions
- Read the documentation

### Staying Updated

- Watch the repository for updates
- Follow the release notes
- Check the CHANGELOG.md

## 🙏 Recognition

Contributors will be recognized in:
- README.md contributors section
- Release notes
- Special thanks in documentation

## 📞 Contact

- Email: nzl9100@gmail.com
- GitHub Issues: [proxyman Issues](https://github.com/geek0ne/proxyman/issues)

---

Thank you for contributing to proxyman! 🎉
