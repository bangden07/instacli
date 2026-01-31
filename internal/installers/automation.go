package installers

import "strings"

// N8NInstaller handles N8N workflow automation installation
type N8NInstaller struct {
	BaseInstaller
}

func NewN8NInstaller() *N8NInstaller {
	return &N8NInstaller{
		BaseInstaller: NewBaseInstaller(
			"N8N",
			"Workflow automation tool (self-hosted)",
			CategoryAutomation,
			"🤖",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (n *N8NInstaller) Dependencies() []string {
	return []string{"docker", "docker-compose"}
}

func (n *N8NInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum, PMBrew, PMChoco}
}

func (n *N8NInstaller) Install(executor Executor) error {
	script := n.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager())
	return executor.RunWithProgress(script, func(line string) {})
}

func (n *N8NInstaller) Uninstall(executor Executor) error {
	script := n.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager())
	_, err := executor.Run(script)
	return err
}

func (n *N8NInstaller) IsInstalled(executor Executor) (bool, error) {
	output, err := executor.Run("docker ps | grep n8n")
	if err != nil {
		return false, nil
	}
	return strings.Contains(output, "n8n"), nil
}

func (n *N8NInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e

echo "🤖 Installing N8N..."

# Create directory
mkdir -p ~/n8n
cd ~/n8n

# Create docker-compose.yml
cat << 'COMPOSE' > docker-compose.yml
version: '3.8'

services:
  n8n:
    image: n8nio/n8n:latest
    container_name: n8n
    restart: unless-stopped
    ports:
      - "5678:5678"
    environment:
      - N8N_BASIC_AUTH_ACTIVE=true
      - N8N_BASIC_AUTH_USER=admin
      - N8N_BASIC_AUTH_PASSWORD=changeme
      - N8N_HOST=localhost
      - N8N_PORT=5678
      - N8N_PROTOCOL=http
      - GENERIC_TIMEZONE=Asia/Jakarta
    volumes:
      - n8n_data:/home/node/.n8n

volumes:
  n8n_data:
COMPOSE

# Start N8N
docker compose up -d

echo ""
echo "✅ N8N installed successfully!"
echo ""
echo "🌐 Access: http://localhost:5678"
echo "👤 Username: admin"
echo "🔑 Password: changeme (change this!)"
echo ""
echo "📁 Data location: ~/n8n"
echo "🛑 Stop: cd ~/n8n && docker compose down"
`
}

func (n *N8NInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `cd ~/n8n && docker compose down -v
rm -rf ~/n8n`
}

var _ Installer = (*N8NInstaller)(nil)

// GolangInstaller handles Go programming language installation
type GolangInstaller struct {
	BaseInstaller
	version string
}

func NewGolangInstaller() *GolangInstaller {
	return &GolangInstaller{
		BaseInstaller: NewBaseInstaller(
			"Golang",
			"Go programming language",
			CategoryRuntime,
			"🐹",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
		version: "1.23.5",
	}
}

func (g *GolangInstaller) Dependencies() []string {
	return []string{"curl", "tar"}
}

func (g *GolangInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum, PMBrew, PMChoco}
}

func (g *GolangInstaller) Install(executor Executor) error {
	script := g.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager())
	return executor.RunWithProgress(script, func(line string) {})
}

func (g *GolangInstaller) Uninstall(executor Executor) error {
	script := g.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager())
	_, err := executor.Run(script)
	return err
}

func (g *GolangInstaller) IsInstalled(executor Executor) (bool, error) {
	output, err := executor.Run("go version")
	if err != nil {
		return false, nil
	}
	return strings.Contains(output, "go version"), nil
}

func (g *GolangInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	switch os {
	case OSLinux:
		return `#!/bin/bash
set -e

echo "🐹 Installing Go..."

VERSION="1.23.5"
ARCH=$(uname -m)

if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "aarch64" ]; then
    ARCH="arm64"
fi

# Remove old installation
sudo rm -rf /usr/local/go

# Download and install
curl -LO "https://go.dev/dl/go${VERSION}.linux-${ARCH}.tar.gz"
sudo tar -C /usr/local -xzf "go${VERSION}.linux-${ARCH}.tar.gz"
rm "go${VERSION}.linux-${ARCH}.tar.gz"

# Add to PATH
if ! grep -q '/usr/local/go/bin' ~/.bashrc; then
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
fi

if [ -f ~/.zshrc ] && ! grep -q '/usr/local/go/bin' ~/.zshrc; then
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
    echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc
fi

export PATH=$PATH:/usr/local/go/bin

echo ""
echo "✅ Go installed successfully!"
go version
`
	case OSMacOS:
		return `#!/bin/bash
brew install go
echo "✅ Go installed!"
go version`
	case OSWindows:
		return `choco install golang -y
Write-Host "✅ Go installed!" -ForegroundColor Green`
	default:
		return ""
	}
}

func (g *GolangInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	switch os {
	case OSLinux:
		return `sudo rm -rf /usr/local/go
sed -i '/\/usr\/local\/go/d' ~/.bashrc ~/.zshrc 2>/dev/null || true`
	case OSMacOS:
		return `brew uninstall go`
	case OSWindows:
		return `choco uninstall golang -y`
	default:
		return ""
	}
}

var _ Installer = (*GolangInstaller)(nil)
