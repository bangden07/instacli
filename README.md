# 🚀 InstaCli

<p align="center">
  <img src="https://img.shields.io/badge/version-1.2.5-blue.svg" alt="Version">
  <img src="https://img.shields.io/badge/go-%3E%3D1.21-00ADD8.svg" alt="Go Version">
  <img src="https://img.shields.io/badge/license-MIT-green.svg" alt="License">
  <img src="https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg" alt="Platform">
</p>

<p align="center">
  <strong>Universal Server Installation Tool with Beautiful TUI</strong>
</p>

<p align="center">
  A powerful terminal-based installer tool for setting up development environments, servers, and applications with just a few keystrokes.
</p>

---

## ✨ Features

### 🎯 28+ One-Click Installers

| Category | Tools |
| -------- | ----- |
| **🌐 Web Server** | Nginx, Apache, LAMP Stack, LEMP Stack |
| **⚡ Runtime** | Node.js, Go, Python, PHP |
| **🐳 Containers** | Docker, Docker Compose |
| **🗄️ Databases** | MySQL, PostgreSQL, MongoDB, Redis |
| **🔧 Frameworks** | Laravel Kit, Next.js Kit |
| **🤖 Automation** | N8N, Coolify, PM2 |
| **🛡️ Security** | UFW Firewall, Certbot SSL, Fail2ban |
| **📊 Monitoring** | Prometheus, Grafana, Netdata |
| **🔀 Infrastructure** | Nginx Proxy Manager, Traefik, MinIO |
| **🔐 VPN** | WireGuard |
| **🚀 CI/CD** | Jenkins, GitLab Runner, GitHub Actions Runner |
| **🕳️ DNS** | Pi-hole |
| **📝 CMS** | WordPress, Ghost |
| **💾 Backup** | Restic |

### 📥 Clone & Setup (NEW in v1.2.0)

Automatically clone and setup any Git repository with dependency detection:

- **📦 Node.js** - Detects npm, pnpm, yarn, bun
- **🐍 Python** - Detects pip, pipenv, venv
- **🐹 Go** - Detects go.mod
- **🐘 PHP** - Detects composer.json, Laravel
- **💎 Ruby** - Detects Gemfile
- **🦀 Rust** - Detects Cargo.toml
- **🐳 Docker** - Detects docker-compose.yml

### 🎨 Beautiful TUI

- RGB animated borders
- Keyboard-driven navigation
- Real-time system status
- SSH remote installation support

---

## 📦 Installation

### Quick Install (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/bangden07/instacli/main/install.sh | bash
```

### Build from Source

```bash
# Clone repository
git clone https://github.com/bangden07/instacli.git
cd instacli

# Build
go build -o instacli ./cmd/instacli

# Run
./instacli
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/bangden07/instacli/releases).

---

## 🚀 Usage

### CLI Commands

```bash
# Run the TUI
instacli

# Check current version
instacli --version

# Check for updates
instacli --check-update

# Auto-update to latest version
instacli --update
```

### Keyboard Shortcuts

| Key | Action |
| --- | ------ |
| `↑` `↓` / `j` `k` | Navigate menu |
| `Enter` | Select item |
| `Tab` | Toggle Local/SSH mode |
| `Esc` | Go back |
| `i` | Install selected tool |
| `g` | Generate install script |
| `?` | Show help |
| `q` | Quit |

### Clone & Setup

1. Navigate to **"Clone & Setup"** menu
2. Paste your Git repository URL
3. Press `Enter` to clone and auto-setup
4. The tool will:
   - Clone the repository
   - Detect project type
   - Install runtime if needed
   - Install all dependencies

---

## 📊 System Requirements

- **OS**: Linux (Debian/Ubuntu, RHEL/CentOS), macOS, Windows (limited)
- **Go**: 1.21+ (for building from source)
- **Terminal**: 256-color support recommended

---

## 🔧 Configuration

### SSH Remote Installation

1. Go to **Settings** in the TUI
2. Press `e` to edit SSH configuration
3. Enter your server details:
   - Host: `192.168.1.100`
   - Port: `22`
   - User: `root`
   - Password: `****`
4. Press `Tab` to switch to SSH mode
5. Install tools remotely!

---

## 📁 Project Structure

```
instacli/
├── cmd/instacli/          # Main entry point
├── internal/
│   ├── executor/          # System check & repo setup
│   ├── installers/        # Installer definitions
│   │   ├── base.go        # Base installer interface
│   │   ├── registry.go    # Installer registry
│   │   ├── webserver.go   # Web server installers
│   │   ├── runtime.go     # Runtime installers
│   │   ├── database.go    # Database installers
│   │   ├── monitoring.go  # Monitoring installers
│   │   ├── infrastructure.go
│   │   ├── cicd.go        # CI/CD installers
│   │   └── additional.go  # Additional installers
│   ├── tui/               # Terminal UI
│   └── version/           # Version info
└── scripts/generated/     # Generated install scripts
```

---

## 🔄 Auto-Update

InstaCli can update itself! No need to run the install script again.

```bash
# Check if update is available
instacli --check-update

# Update to latest version
instacli --update
```

---

## 🏷️ Version History

| Version | Date | Changes |
| ------- | ---- | ------- |
| v1.2.5 | 2026-01-31 | Auto-update feature |
| v1.2.4 | 2026-01-31 | Fix SSH remote installation |
| v1.2.0 | 2026-01-31 | Clone & Setup feature |
| v1.1.0 | 2026-01-31 | 13 new installers (28 total) |
| v1.0.0 | 2026-01-31 | Initial release (15 installers) |

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👤 Author

**bangden07**

- GitHub: [@bangden07](https://github.com/bangden07)

---

<p align="center">
  Made with ❤️ and Go
</p>
