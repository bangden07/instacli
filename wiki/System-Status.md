# System Status

View detailed system information and installed tools.

## 🔍 Overview

The System Status view shows:
- Operating system details
- Architecture
- Package manager detected
- Installed tools with versions

## 📊 Accessing System Status

1. Navigate to **"System Status"** in main menu
2. Press `Enter`

## 🖥️ Information Displayed

### Operating System

```
Operating System
  OS: linux
  Distribution: Ubuntu 22.04 LTS
  Arch: amd64
  Kernel: 5.15.0-91-generic
```

### Package Manager

```
Package Manager
  Primary: apt
  Available: apt, snap, brew
```

### Installed Tools

```
Detected Tools
  ✅ Docker         v24.0.7
  ✅ Docker Compose v2.23.0
  ✅ Node.js        v20.10.0
  ✅ npm            v10.2.3
  ✅ Go             v1.21.5
  ✅ Python         v3.11.2
  ✅ PHP            v8.2.12
  ✅ Nginx          v1.24.0
  ✅ MySQL          v8.0.35
  ✅ Redis          v7.2.3
  ❌ MongoDB        (not installed)
  ❌ PostgreSQL     (not installed)
```

## 🔄 Refresh

The system status is refreshed each time you enter the view.

## 🛠️ How Detection Works

InstaCli uses various methods to detect tools:

### Command Check
```go
// Check if command exists
_, err := exec.LookPath("docker")
```

### Version Extraction
```go
// Run version command
output, _ := exec.Command("docker", "--version").Output()
// Parse: "Docker version 24.0.7, build ..."
```

### Service Status
```go
// Check systemd service
exec.Command("systemctl", "is-active", "docker").Run()
```

## 📋 Detected Tools List

| Tool | Detection Command |
| ---- | ----------------- |
| Docker | `docker --version` |
| Docker Compose | `docker compose version` |
| Node.js | `node --version` |
| npm | `npm --version` |
| Go | `go version` |
| Python | `python3 --version` |
| PHP | `php --version` |
| Nginx | `nginx -v` |
| Apache | `apache2 -v` |
| MySQL | `mysql --version` |
| PostgreSQL | `psql --version` |
| MongoDB | `mongod --version` |
| Redis | `redis-server --version` |
| Git | `git --version` |

---

**Next:** [[Project Structure]] →
