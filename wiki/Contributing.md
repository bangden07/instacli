# Contributing to InstaCli

Thank you for your interest in contributing! 🎉

## 🚀 Quick Start

```bash
# 1. Fork the repository on GitHub
# 2. Clone your fork
git clone https://github.com/YOUR-USERNAME/instacli.git
cd instacli

# 3. Add upstream remote
git remote add upstream https://github.com/bangden07/instacli.git

# 4. Install dependencies
go mod download

# 5. Build
go build -o instacli.exe ./cmd/instacli  # Windows
go build -o instacli ./cmd/instacli      # Linux/macOS

# 6. Run
./instacli
```

## 📁 Project Structure

```
instacli/
├── cmd/
│   └── instacli/
│       └── main.go           # Entry point
├── internal/
│   ├── config/               # Configuration management
│   │   └── config.go
│   ├── executor/             # SSH and local execution
│   │   ├── executor.go
│   │   ├── local.go
│   │   └── ssh.go
│   ├── installers/           # All installer implementations
│   │   ├── base.go           # Base installer interface
│   │   ├── registry.go       # Installer registry
│   │   ├── runtime.go        # Node.js, Python, Go, etc.
│   │   ├── databases.go      # MySQL, PostgreSQL, Redis, etc.
│   │   ├── infrastructure.go # Nginx, Docker, etc.
│   │   ├── security.go       # Certbot, Fail2ban, UFW
│   │   ├── monitoring.go     # Prometheus, Grafana, etc.
│   │   ├── ai_cli.go         # Claude, Gemini, Codex, etc.
│   │   └── mcp.go            # MCP servers
│   ├── tui/                  # Terminal UI components
│   │   ├── app.go            # Main application logic
│   │   ├── styles.go         # UI styling (colors, borders)
│   │   └── views.go          # View rendering
│   ├── updater/              # Auto-update functionality
│   └── version/
│       └── version.go        # Version information
├── wiki/                     # GitHub Wiki documentation
├── README.md
└── go.mod
```

## 🔧 Adding a New Feature

### Adding a New Installer

**Files to modify:**
1. Create new installer file OR add to existing category file
2. `internal/installers/registry.go` - Register the installer

**Step 1: Create the Installer**

Create a new file or add to existing category in `internal/installers/`:

```go
// internal/installers/my_tool.go
package installers

import (
    "fmt"
    "strings"
)

// MyToolInstaller installs MyTool
type MyToolInstaller struct {
    BaseInstaller
}

// NewMyToolInstaller creates a new installer
func NewMyToolInstaller() *MyToolInstaller {
    return &MyToolInstaller{
        BaseInstaller: NewBaseInstaller(
            "MyTool",                              // Name
            "Description of what MyTool does",    // Description
            CategoryDevTools,                      // Category (see base.go)
            "🔧",                                  // Icon emoji
            []OS{OSLinux, OSMacOS, OSWindows},    // Supported OS
        ),
    }
}

// Dependencies returns required dependencies
func (i *MyToolInstaller) Dependencies() []string {
    return []string{"git"} // Optional: tools that must be installed first
}

// GenerateInstallScript generates the installation script
func (i *MyToolInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
    var script strings.Builder

    switch os {
    case OSLinux, OSMacOS:
        script.WriteString(`#!/bin/bash
set -e

echo "📦 Installing MyTool..."

# Check dependencies
if ! command -v git &> /dev/null; then
    echo "❌ Git is required"
    exit 1
fi

# Install based on package manager
`)
        if pm == PMApt {
            script.WriteString(`sudo apt-get update && sudo apt-get install -y mytool
`)
        } else if pm == PMBrew {
            script.WriteString(`brew install mytool
`)
        }
        script.WriteString(`
echo "✅ MyTool installed successfully!"
mytool --version
`)
    case OSWindows:
        script.WriteString(`# PowerShell
Write-Host "📦 Installing MyTool..." -ForegroundColor Cyan
choco install mytool -y
Write-Host "✅ MyTool installed!" -ForegroundColor Green
`)
    }

    return script.String()
}

// Required interface methods
func (i *MyToolInstaller) RequiredPackageManagers() []PackageManager { return nil }

func (i *MyToolInstaller) Install(executor Executor) error {
    return executor.RunWithProgress(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()), nil)
}

func (i *MyToolInstaller) Uninstall(executor Executor) error {
    _, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
    return err
}

func (i *MyToolInstaller) IsInstalled(executor Executor) (bool, error) {
    _, err := executor.Run("mytool --version")
    return err == nil, nil
}

func (i *MyToolInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
    switch os {
    case OSLinux:
        return "sudo apt-get remove -y mytool"
    case OSMacOS:
        return "brew uninstall mytool"
    case OSWindows:
        return "choco uninstall mytool -y"
    default:
        return ""
    }
}

// Ensure interface compliance
var _ Installer = (*MyToolInstaller)(nil)
```

**Step 2: Register the Installer**

Edit `internal/installers/registry.go`:

```go
func DefaultRegistry() *Registry {
    r := NewRegistry()

    // ... existing registrations ...

    // Add your new installer
    r.Register(NewMyToolInstaller())

    return r
}
```

**Step 3: Test**

```bash
go build -o instacli.exe ./cmd/instacli
./instacli
# Navigate to your category and verify installer appears
```

### Adding a New Category

**Files to modify:**
1. `internal/installers/base.go` - Add category constant
2. `internal/tui/app.go` - Add menu item and icon

**Step 1: Add Category Constant**

```go
// internal/installers/base.go
const (
    CategoryRuntime     Category = "Runtime"
    CategoryDatabase    Category = "Database"
    // ... existing categories ...
    CategoryMyCategory  Category = "My Category"  // ADD THIS
)
```

**Step 2: Update UI Menu**

In `internal/tui/app.go`, find the category menu rendering and add your category:

```go
// Look for category rendering code and add icon mapping
categoryIcons := map[installers.Category]string{
    installers.CategoryRuntime:    "⚡",
    installers.CategoryDatabase:   "🗃️",
    // ... add your category ...
    installers.CategoryMyCategory: "🆕",
}
```

### Modifying the UI

**Files to modify:**
- `internal/tui/styles.go` - Colors, styles, fonts
- `internal/tui/app.go` - Layout, rendering logic

**Example: Changing Colors**

```go
// internal/tui/styles.go
var (
    Primary   = lipgloss.Color("#00D9FF") // Change hex color
    Secondary = lipgloss.Color("#FF00E5")
)
```

**Example: Adding a New View**

1. Add view constant in `app.go`:
```go
const (
    ViewMain View = iota
    ViewCategory
    // ...
    ViewMyNewScreen  // ADD THIS
)
```

2. Add render function:
```go
func (a *App) renderMyNewScreenView(width int) string {
    var b strings.Builder
    b.WriteString("My New Screen Content\n")
    return b.String()
}
```

3. Add to switch statement in main render:
```go
switch a.currentView {
case ViewMyNewScreen:
    content = a.renderMyNewScreenView(contentWidth)
}
```

## 🐛 Fixing Bugs

### Bug Fix Workflow

1. **Create an issue** describing the bug
2. **Create a branch**:
   ```bash
   git checkout -b fix/issue-description
   ```
3. **Locate the problematic code**:
   - Use `grep` to find related code
   - Check related functions in affected files
4. **Make the fix**
5. **Test thoroughly**:
   ```bash
   go build -o instacli.exe ./cmd/instacli
   ./instacli
   ```
6. **Submit PR** referencing the issue

### Common Bug Areas

| Issue Type | Files to Check |
|------------|----------------|
| Version not updating | `internal/version/version.go`, `internal/tui/app.go` |
| SSH connection issues | `internal/executor/ssh.go` |
| Installer not working | `internal/installers/*.go` |
| UI rendering issues | `internal/tui/app.go`, `internal/tui/styles.go` |
| Config not saving | `internal/config/config.go` |

### Example: Fixing a Hardcoded Value

**Problem:** Version shows v1.0 instead of actual version

**Solution:**

1. Find the hardcoded value:
   ```bash
   grep -r "v1.0" internal/
   ```

2. Replace with dynamic version:
   ```go
   // Before
   subtitle := "Universal Installer Tool v1.0"

   // After
   import "github.com/instacli/instacli/internal/version"
   subtitle := fmt.Sprintf("Universal Installer Tool v%s", version.Version)
   ```

3. Update version when releasing:
   ```go
   // internal/version/version.go
   const Version = "1.3.1"  // Bump this
   ```

## 📝 Code Style

```bash
# Format all code
go fmt ./...

# Run linter
golangci-lint run

# Run tests
go test ./...
```

### Naming Conventions

| Type | Style | Example |
|------|-------|---------|
| Files | `snake_case.go` | `my_installer.go` |
| Exported types | `PascalCase` | `MyToolInstaller` |
| Private types | `camelCase` | `myHelper` |
| Constants | `PascalCase` | `CategoryDevTools` |

## 📬 Pull Request Process

### Branch Naming

```
feature/add-redis-installer
fix/ssh-connection-timeout
docs/update-installation-guide
ui/improve-header-styling
```

### Commit Message Format

```
<type>: <short description>

[Optional body explaining the changes]

[Optional: Closes #123]
```

**Types:**
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation
- `style:` Code formatting
- `refactor:` Code restructure
- `ui:` UI/UX changes
- `chore:` Maintenance

### Example PR

```markdown
## Description
Add Redis installer with cluster mode support

## Changes
- Created `internal/installers/redis.go`
- Added to registry in `registry.go`
- Supports apt, yum, and brew

## Type
- [x] New feature
- [ ] Bug fix

## Testing
- [x] Tested on Ubuntu 22.04
- [x] Tested on Windows 11

## Checklist
- [x] Code formatted with `go fmt`
- [x] Builds successfully
- [x] Tested manually
```

## 🧪 Testing

```bash
# Build and test
go build -o instacli.exe ./cmd/instacli && ./instacli

# Run unit tests
go test ./...

# Run with verbose output
go test -v ./...

# Check coverage
go test -cover ./...
```

## 📜 License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing! 🙏
