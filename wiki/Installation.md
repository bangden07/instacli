# Installation Guide

## 📦 Installation Methods

### Method 1: Quick Install Script (Recommended)

For Linux and macOS:

```bash
curl -fsSL https://raw.githubusercontent.com/bangden07/instacli/main/install.sh | bash
```

This will:
1. Download the latest release
2. Install to `/usr/local/bin/instacli`
3. Make it available system-wide

### Method 2: Build from Source

**Prerequisites:**
- Go 1.21 or higher
- Git

**Steps:**

```bash
# Clone the repository
git clone https://github.com/bangden07/instacli.git
cd instacli

# Build the binary
go build -o instacli ./cmd/instacli

# Optional: Move to PATH
sudo mv instacli /usr/local/bin/
```

### Method 3: Download Binary

1. Go to [Releases](https://github.com/bangden07/instacli/releases)
2. Download the binary for your OS:
   - `instacli-linux-amd64` for Linux
   - `instacli-darwin-amd64` for macOS (Intel)
   - `instacli-darwin-arm64` for macOS (Apple Silicon)
   - `instacli-windows-amd64.exe` for Windows
3. Make it executable (Linux/macOS):
   ```bash
   chmod +x instacli-*
   ```

## ✅ Verify Installation

```bash
instacli --version
# Output: InstaCli v1.2.5
```

## 🔄 Auto-Update

InstaCli can update itself! No need to run the install script again.

```bash
# Check if update is available
instacli --check-update

# Update to latest version
instacli --update
```

## 🔧 System Requirements

| Requirement | Minimum |
| ----------- | ------- |
| OS | Linux, macOS, Windows |
| Terminal | 256-color support |
| Go | 1.21+ (for building) |
| Disk | 10MB free space |

## 🐧 Supported Linux Distributions

- Ubuntu 18.04+
- Debian 10+
- CentOS 7+
- RHEL 7+
- Fedora 30+
- Arch Linux
- Alpine Linux

## 🍎 macOS Support

- macOS 10.15 (Catalina) or later
- Intel and Apple Silicon supported

## 🪟 Windows Support

- Windows 10/11
- Windows Terminal recommended
- Limited installer support (mainly generates scripts)

---

**Next:** [[Quick Start]] →

