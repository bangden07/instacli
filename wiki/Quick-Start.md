# Quick Start Guide

Get up and running with InstaCli in 5 minutes!

## 🚀 Launch InstaCli

```bash
./instacli
```

You'll see the beautiful TUI interface:

```
╭──────────────────────────────────────────────────╮
│      ___           _        ___  _  _            │
│     |_ _| _ _   __| |_ __ _/ __|| ||_|           │
│      | | | ' \ (_-<  _/ _` \__ \| || |           │
│     |___||_||_|/__/\__\__,_|___/|_||_|           │
│                                                   │
│          🚀 Universal Installer Tool v1.2        │
╰──────────────────────────────────────────────────╯
```

## 🎯 Basic Navigation

| Key | Action |
| --- | ------ |
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` | Select |
| `Esc` | Go back |
| `q` | Quit |

## 📦 Install Your First Tool

### Example: Install Docker

1. **Select "Containers"** category
2. **Select "Docker"** from the list
3. **Press `i`** to install

InstaCli will:
- Detect your OS and package manager
- Generate the appropriate install script
- Execute the installation

## 🔍 Check System Status

1. **Select "System Status"** from main menu
2. View detected tools and versions:

```
🔍 System Status

Operating System
  OS: linux (Ubuntu 22.04)
  Arch: amd64
  Package Manager: apt

Detected Tools
  ✅ Docker v24.0.5
  ✅ Node.js v20.10.0
  ✅ Go v1.21.5
  ❌ Python (not installed)
```

## 📥 Clone & Setup a Repository

1. **Select "Clone & Setup"** from main menu
2. **Paste** your Git URL:
   ```
   https://github.com/user/nextjs-app.git
   ```
3. **Press Enter** to clone and setup

InstaCli will automatically:
- Clone the repository
- Detect the project type (Node.js, Python, Go, etc.)
- Install dependencies

## 🌐 Remote Installation (SSH)

1. **Go to Settings** → Configure SSH
2. **Enter server details**:
   - Host: `192.168.1.100`
   - Port: `22`
   - User: `root`
3. **Press Tab** to switch to SSH mode
4. **Install any tool** - it will run on the remote server!

## 📝 Generate Scripts Only

If you prefer to review scripts before running:

1. Select a tool
2. Press `g` to generate script
3. Find the script in `./scripts/generated/`
4. Review and run manually:
   ```bash
   chmod +x ./scripts/generated/docker.sh
   sudo ./scripts/generated/docker.sh
   ```

---

**Next:** [[Keyboard Shortcuts]] →
