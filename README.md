# 🚀 InstaCli

<p align="center">
  <img src="https://img.shields.io/badge/version-1.3.0-00D9FF.svg" alt="Version">
  <img src="https://img.shields.io/badge/go-%3E%3D1.21-00ADD8.svg" alt="Go Version">
  <img src="https://img.shields.io/badge/license-MIT-00FF88.svg" alt="License">
  <img src="https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-FF00E5.svg" alt="Platform">
</p>

<p align="center">
  <strong>⚡ Universal Server Installation Tool with Premium TUI</strong>
</p>

<p align="center">
  A powerful terminal-based installer tool for setting up development environments, servers, AI tools, and MCP servers with just a few keystrokes.
</p>

---

## ✨ Features

### 🎯 42+ One-Click Installers

| Category | Tools |
| -------- | ----- |
| **⚡ Runtime** | Node.js, Go, Python, PHP |
| **🗃️ Database** | MySQL, PostgreSQL, MongoDB, Redis |
| **🌐 Web Server** | Nginx, Apache, LAMP Stack, LEMP Stack |
| **🐳 Containers** | Docker, Docker Compose |
| **🔧 Dev Tools** | Laravel Kit, Next.js Kit, PM2 |
| **🤖 AI CLI** | Claude CLI, Gemini CLI, Codex CLI, Aider, Kilo Code, Continue |
| **📦 MCP Servers** | Context7, Playwright, GitHub, Filesystem, PostgreSQL, Brave Search, Memory, Sequential Thinking |
| **🛡️ Security** | UFW Firewall, Certbot SSL, Fail2ban |
| **📊 Monitoring** | Prometheus, Grafana, Netdata |
| **🔀 Infrastructure** | Nginx Proxy Manager, Traefik, MinIO |
| **🚀 CI/CD** | Jenkins, GitLab Runner, GitHub Actions Runner |
| **🤖 Automation** | N8N, Coolify |

### 🤖 AI CLI Tools (NEW in v1.3.0)

Install your favorite AI coding assistants:

| Tool | Description |
|------|-------------|
| **Claude CLI** | Anthropic's Claude Code CLI |
| **Gemini CLI** | Google's Gemini AI CLI |
| **Codex CLI** | OpenAI's Codex CLI |
| **Aider** | AI pair programming assistant |
| **Kilo Code** | VS Code AI extension |
| **Continue** | Open-source AI assistant |
| **OpenCode CLI** | Open-source AI coding CLI |


### 📦 MCP Servers (NEW in v1.3.0)

Pre-configured Model Context Protocol servers:

| MCP | Purpose |
|-----|---------|
| **Context7** | Documentation lookup |
| **Playwright** | Browser automation |
| **GitHub** | Repository integration |
| **Filesystem** | File system access |
| **PostgreSQL** | Database access |
| **Brave Search** | Web search |
| **Memory** | Persistent memory |
| **Sequential Thinking** | Step-by-step reasoning |

### 📥 Clone & Setup

Automatically clone and setup any Git repository:

- **📦 Node.js** - npm, pnpm, yarn, bun
- **🐍 Python** - pip, pipenv, venv
- **🐹 Go** - go.mod
- **🐘 PHP** - composer, Laravel
- **💎 Ruby** - Gemfile
- **🦀 Rust** - Cargo.toml
- **🐳 Docker** - docker-compose.yml

### 🎨 Premium TUI

- Cyberpunk neon theme
- RGB animated borders
- Real-time system status
- SSH remote installation support
- Auto-update feature

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

---

## 🔄 Auto-Update

InstaCli updates itself automatically:

```bash
# Check for updates
instacli --check-update

# Update to latest version
instacli --update
```

---

## 🏷️ Version History

| Version | Date | Changes |
| ------- | ---- | ------- |
| v1.3.0 | 2026-01-31 | AI CLI + MCP Server installers, Premium UI redesign |
| v1.2.9 | 2026-01-31 | UI performance improvements |
| v1.2.5 | 2026-01-31 | Auto-update feature |
| v1.2.0 | 2026-01-31 | Clone & Setup feature |
| v1.1.0 | 2026-01-31 | 13 new installers |
| v1.0.0 | 2026-01-31 | Initial release |

---

## 🤝 Contributing

See [Contributing Guide](wiki/Contributing.md) for detailed instructions on:

- Adding new installers
- Adding new features
- Fixing bugs
- Improving UI

Quick start:

```bash
# Fork and clone
git clone https://github.com/YOUR-USERNAME/instacli.git
cd instacli

# Build and test
go build -o instacli ./cmd/instacli
./instacli

# Create branch
git checkout -b feature/my-feature

# Submit PR
```

---

## 📄 License

MIT License - see [LICENSE](LICENSE) file.

---

## 👤 Author

**bangden07** - [@bangden07](https://github.com/bangden07)

---

<p align="center">
  Made with ❤️ and Go
</p>
