# Project Structure

Understanding the InstaCli codebase.

## 📁 Directory Layout

```
instacli/
├── cmd/
│   └── instacli/
│       └── main.go           # Entry point
├── internal/
│   ├── executor/
│   │   ├── executor.go       # Command execution interface
│   │   ├── system.go         # System detection
│   │   └── reposetup.go      # Clone & Setup logic
│   ├── installers/
│   │   ├── base.go           # Base types and interfaces
│   │   ├── registry.go       # Installer registry
│   │   ├── webserver.go      # Nginx, Apache, LAMP, LEMP
│   │   ├── runtime.go        # Node.js, Go, Python, PHP
│   │   ├── database.go       # MySQL, PostgreSQL, MongoDB, Redis
│   │   ├── container.go      # Docker, Docker Compose
│   │   ├── framework.go      # Laravel, Next.js kits
│   │   ├── automation.go     # N8N, Coolify, PM2
│   │   ├── security.go       # UFW, Certbot, Fail2ban
│   │   ├── monitoring.go     # Prometheus, Grafana, Netdata
│   │   ├── infrastructure.go # NPM, Traefik, MinIO
│   │   ├── cicd.go           # Jenkins, GitLab, GitHub runners
│   │   └── additional.go     # Pi-hole, WordPress, Ghost, Restic
│   ├── tui/
│   │   ├── app.go            # Main TUI application
│   │   └── styles.go         # UI styles and colors
│   └── version/
│       └── version.go        # Version information
├── scripts/
│   └── generated/            # Generated install scripts
├── wiki/                     # GitHub Wiki pages
├── .gitignore
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

## 🧩 Core Components

### Entry Point (`cmd/instacli/main.go`)

```go
func main() {
    app := tui.NewApp()
    p := tea.NewProgram(app, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### Installer Interface (`internal/installers/base.go`)

```go
type Installer interface {
    Name() string
    Description() string
    Category() Category
    Icon() string
    SupportedOS() []OS
    GenerateInstallScript(os OS, pm PackageManager) string
    GenerateUninstallScript(os OS, pm PackageManager) string
    Install(exec executor.Executor) error
    Uninstall(exec executor.Executor) error
    IsInstalled() bool
}
```

### TUI Application (`internal/tui/app.go`)

The TUI is built with [Bubble Tea](https://github.com/charmbracelet/bubbletea):

```go
type App struct {
    currentView       View
    targetMode        TargetMode
    categoryItems     []CategoryItem
    installerItems    []InstallerItem
    cursor            int
    selectedCategory  string
    selectedInstaller Installer
    registry          *Registry
    // ... more fields
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (a *App) View() string
```

### System Detection (`internal/executor/system.go`)

```go
type SystemCheck struct {
    OS             string
    Arch           string
    PackageManager PackageManager
    InstalledTools map[string]string
}

func (s *SystemCheck) DetectAll()
func (s *SystemCheck) DetectTool(name string) (version string, installed bool)
```

## 🔄 Data Flow

```
┌─────────────┐     ┌──────────────┐     ┌────────────┐
│   main.go   │────▶│   TUI App    │────▶│  Registry  │
└─────────────┘     └──────────────┘     └────────────┘
                           │                    │
                           ▼                    ▼
                    ┌──────────────┐     ┌────────────┐
                    │ SystemCheck  │     │ Installers │
                    └──────────────┘     └────────────┘
                                               │
                                               ▼
                                        ┌────────────┐
                                        │  Scripts   │
                                        └────────────┘
```

## 📦 Dependencies

| Package | Purpose |
| ------- | ------- |
| `github.com/charmbracelet/bubbletea` | TUI framework |
| `github.com/charmbracelet/lipgloss` | Styling |
| `github.com/charmbracelet/bubbles` | TUI components |

## 🎨 Styling (`internal/tui/styles.go`)

```go
var (
    Primary   = lipgloss.Color("#7C3AED")
    Secondary = lipgloss.Color("#10B981")
    Accent    = lipgloss.Color("#F59E0B")
    Muted     = lipgloss.Color("#6B7280")
    
    TitleStyle = lipgloss.NewStyle().
        Foreground(Primary).
        Bold(true)
)
```

---

**Next:** [[Adding Installers]] →
