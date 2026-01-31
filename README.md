# 🚀 InstaCli - Server Setup Automation Tool

A modern, interactive Terminal User Interface (TUI) for automating server setup and software installation. Built with Go and [Bubbletea](https://github.com/charmbracelet/bubbletea).

![InstaCli TUI](https://img.shields.io/badge/TUI-Modern%20Terminal%20UI-blueviolet)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![License](https://img.shields.io/badge/License-MIT-green)

---

## ✨ Features

### 🌈 Visual Experience
- **RGB Animated Border** - Rainbow color cycling border for modern aesthetic
- **Adaptive Header** - Compact or full ASCII logo based on terminal size
- **Dark Theme** - Eye-friendly colors optimized for terminal use

### 📦 15+ Pre-built Installers

| Category | Tools | Description |
|----------|-------|-------------|
| **🗄️ Databases** | MySQL, PostgreSQL, MongoDB, Redis | Latest stable versions via official repos |
| **🛡️ Security** | UFW, Certbot SSL, Fail2ban | Firewall, SSL certificates, intrusion prevention |
| **🐳 Containers** | Docker + Compose | Container runtime with buildx plugin |
| **⚡ Runtime** | Node.js (NVM), Golang | Programming language runtimes |
| **🌐 Web Server** | LAMP, LEMP | Full web server stacks |
| **🤖 Automation** | N8N | Workflow automation platform |
| **🔧 Frameworks** | Laravel Kit, Next.js Kit | Development framework setups |

### 🎯 Key Capabilities
- **Local Installation** - Install directly on your machine
- **SSH Remote Installation** - Deploy to remote servers
- **Script Generation** - Export installation scripts for manual use
- **System Status** - Check installed tools and environment

---

## 📥 Installation

### Quick Install (Linux/macOS)
```bash
curl -fsSL https://raw.githubusercontent.com/instacli/instacli/main/install.sh | bash
```

### Build from Source
```bash
# Clone repository
git clone https://github.com/instacli/instacli.git
cd instacli

# Build
go build -o instacli ./cmd/instacli

# Run
./instacli
```

### Windows
```powershell
# Download latest release or build from source
go build -o instacli.exe ./cmd/instacli
.\instacli.exe
```

---

## 🎮 Usage Guide

### Starting the Application
```bash
./instacli          # Linux/macOS
.\instacli.exe      # Windows
```

### Keyboard Navigation

| Key | Action | Context |
|-----|--------|---------|
| `↑` `↓` or `k` `j` | Navigate up/down | All views |
| `Enter` | Select item | Main menu, Categories |
| `Esc` / `Backspace` | Go back | All views |
| `Tab` | Toggle Local/SSH mode | Main menu |
| `g` | Generate install script | Installer view |
| `i` | Run installation | Installer view |
| `e` | Edit SSH settings | Settings view |
| `?` | Show help | Anywhere |
| `q` | Quit application | Main menu |

### Workflow Example

#### 1. Local Installation
```
1. Select a category (e.g., "Databases")
2. Choose an installer (e.g., "MySQL")
3. Press `i` to install locally
4. Follow the generated commands
```

#### 2. Remote SSH Installation
```
1. Press `Tab` to switch to SSH mode
2. Go to Settings → Edit SSH configuration
3. Enter: Host, Port (22), User, Password
4. Navigate to desired installer
5. Press `i` to generate SSH commands
```

#### 3. Generate Scripts for Later
```
1. Navigate to any installer
2. Press `g` to generate script
3. Scripts saved to: ./scripts/generated/{tool_name}.sh
4. Run manually: chmod +x script.sh && sudo ./script.sh
```

---

## 📋 Available Installers

### 🗄️ Databases

#### MySQL
- **Version**: Latest stable from official MySQL APT repository
- **Includes**: MySQL Server, MySQL Client
- **Service**: Enabled and started automatically

#### PostgreSQL
- **Version**: Latest stable from official PostgreSQL repository
- **Includes**: PostgreSQL Server, postgresql-contrib
- **Service**: Enabled and started automatically

#### MongoDB
- **Version**: MongoDB 7.0 (latest stable)
- **Repository**: Official MongoDB repo with GPG signing
- **Service**: mongod enabled and started

#### Redis
- **Version**: Latest stable from distribution repos
- **Includes**: Redis Server with default configuration
- **Service**: redis-server enabled and started

---

### 🛡️ Security

#### UFW Firewall
- **Description**: Uncomplicated Firewall for Linux
- **Default Rules**:
  - Deny all incoming
  - Allow all outgoing
  - Allow SSH (22), HTTP (80), HTTPS (443)

#### Certbot SSL
- **Description**: Let's Encrypt SSL certificate manager
- **Plugins**: Auto-detects Nginx/Apache and installs appropriate plugin
- **Usage**: `sudo certbot --nginx -d yourdomain.com`

#### Fail2ban
- **Description**: Intrusion prevention system
- **Default Config**: SSH protection with 5 max retries, 1 hour ban
- **Service**: Enabled and started with jail.local configuration

---

### 🐳 Containers

#### Docker
- **Components**: Docker CE, CLI, containerd, buildx, compose plugin
- **Repository**: Official Docker CE repository
- **Post-install**: User added to docker group

---

### ⚡ Runtime & Languages

#### Node.js (via NVM)
- **Version**: LTS (Long Term Support) - recommended for production
- **Manager**: NVM v0.40.1 for version management
- **Extras**: pnpm and yarn installed globally

#### Golang
- **Version**: Latest stable from official Go downloads
- **Path**: Configured in ~/.bashrc and ~/.zshrc

---

### 🌐 Web Server Stacks

#### LAMP Stack
- Apache2, MySQL, PHP (latest)
- phpMyAdmin for database management
- Virtual host configuration

#### LEMP Stack
- Nginx, MySQL, PHP-FPM
- Optimized for performance
- SSL-ready configuration

---

### 🤖 Automation

#### N8N
- **Description**: Workflow automation platform
- **Method**: Docker-based installation
- **Access**: http://localhost:5678

---

### 🔧 Frameworks

#### Laravel Kit
- PHP 8.2+, Composer, Laravel installer
- Required PHP extensions installed

#### Next.js Kit
- Node.js LTS, pnpm, create-next-app
- Ready for Next.js 14+ development

---

## ⚙️ Configuration

### SSH Settings
Configure SSH in the Settings menu:
- **Host**: Server IP or hostname
- **Port**: SSH port (default: 22)
- **User**: SSH username (default: root)
- **Password**: SSH password (masked input)

### Generated Scripts Location
All generated scripts are saved to:
```
./scripts/generated/
├── mysql.sh
├── postgresql.sh
├── docker.sh
└── ...
```

---

## 🔧 System Requirements

### Minimum
- Terminal with 256-color support
- 80x24 terminal size (recommended: 120x40)

### Target Systems
| OS | Package Managers |
|----|-----------------|
| Ubuntu/Debian | apt |
| RHEL/CentOS/Fedora | yum/dnf |
| Arch Linux | pacman |
| macOS | brew |
| Windows | choco/winget |

---

## 🛠️ Development

### Project Structure
```
instacli/
├── cmd/instacli/       # Main entry point
├── internal/
│   ├── tui/            # Terminal UI (app.go, styles.go)
│   ├── installers/     # All installer definitions
│   │   ├── base.go     # Base installer interface
│   │   ├── database.go # MySQL, PostgreSQL, MongoDB, Redis
│   │   ├── security.go # UFW, Certbot, Fail2ban
│   │   ├── docker.go   # Docker installer
│   │   ├── nodejs.go   # Node.js installer
│   │   └── ...
│   └── executor/       # Command execution, system checks
├── scripts/
│   ├── build.sh        # Build script
│   └── generated/      # Generated install scripts
└── README.md
```

### Building
```bash
# All platforms
./scripts/build.sh

# Windows only
go build -o instacli.exe ./cmd/instacli

# Linux only
GOOS=linux GOARCH=amd64 go build -o instacli-linux ./cmd/instacli
```

---

## 📝 License

MIT License - See [LICENSE](LICENSE) for details.

---

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📧 Support

- **Issues**: [GitHub Issues](https://github.com/instacli/instacli/issues)
- **Author**: InstaCli Team

---

Made with ❤️ using Go and Bubbletea
