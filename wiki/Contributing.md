# Contributing to InstaCli

Thank you for your interest in contributing! 🎉

## 🚀 Quick Start

1. Fork the repository
2. Clone your fork
3. Create a branch
4. Make changes
5. Submit a PR

## 📋 Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/instacli.git
cd instacli

# Install dependencies
go mod download

# Build
go build -o instacli ./cmd/instacli

# Run
./instacli
```

## 🔧 Types of Contributions

### 1. Add New Installer

See [[Adding Installers]] for detailed guide.

Quick steps:
1. Create installer in `internal/installers/`
2. Register in `registry.go`
3. Test locally
4. Submit PR

### 2. Bug Fixes

1. Open an issue describing the bug
2. Reference the issue in your PR
3. Include test case if possible

### 3. Documentation

- Wiki pages in `wiki/` folder
- README improvements
- Code comments

### 4. UI Improvements

- New views or screens
- Better styling
- Accessibility improvements

## 📝 Code Style

### Go Formatting

```bash
# Format code
go fmt ./...

# Lint
golangci-lint run
```

### Naming Conventions

- **Files**: `lowercase_with_underscores.go`
- **Types**: `PascalCase`
- **Functions**: `PascalCase` (exported), `camelCase` (private)
- **Variables**: `camelCase`
- **Constants**: `PascalCase`

### Comments

```go
// MyFunction does X and returns Y.
// It handles Z edge cases.
func MyFunction() error {
    // ...
}
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## 📬 Pull Request Process

### 1. Branch Naming

```
feature/add-redis-installer
fix/ssh-connection-timeout
docs/update-installation-guide
```

### 2. Commit Messages

```
feat: add Redis installer

- Add Redis installer with apt/yum/brew support
- Add systemd service configuration
- Add cluster mode option

Closes #123
```

Prefixes:
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation
- `style:` Formatting
- `refactor:` Code restructure
- `test:` Tests
- `chore:` Maintenance

### 3. PR Description

```markdown
## Description
Brief description of changes.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation

## Testing
- [ ] Tested on Linux
- [ ] Tested on macOS

## Screenshots (if UI changes)
```

### 4. Review Process

1. Automated checks must pass
2. At least one maintainer review
3. Address feedback
4. Merge!

## 🐛 Reporting Issues

### Bug Reports

Include:
- InstaCli version
- OS and version
- Steps to reproduce
- Expected vs actual behavior
- Error messages/logs

### Feature Requests

Include:
- Use case description
- Proposed solution
- Alternatives considered

## 💬 Communication

- GitHub Issues: Bug reports, features
- GitHub Discussions: Questions, ideas
- Pull Requests: Code contributions

## 📜 License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing! 🙏
