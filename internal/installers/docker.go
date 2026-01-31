package installers

import (
	"strings"
)

// DockerInstaller handles Docker and Docker Compose installation
type DockerInstaller struct {
	BaseInstaller
}

func NewDockerInstaller() *DockerInstaller {
	return &DockerInstaller{
		BaseInstaller: NewBaseInstaller(
			"Docker",
			"Docker Engine & Docker Compose",
			CategoryContainer,
			"🐳",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (d *DockerInstaller) Dependencies() []string {
	return []string{"curl", "gnupg"}
}

func (d *DockerInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum, PMDnf, PMBrew, PMChoco}
}

func (d *DockerInstaller) Install(executor Executor) error {
	script := d.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager())
	return executor.RunWithProgress(script, func(line string) {
		// Progress callback
	})
}

func (d *DockerInstaller) Uninstall(executor Executor) error {
	script := d.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager())
	_, err := executor.Run(script)
	return err
}

func (d *DockerInstaller) IsInstalled(executor Executor) (bool, error) {
	output, err := executor.Run("docker --version")
	if err != nil {
		return false, nil
	}
	return strings.Contains(output, "Docker version"), nil
}

func (d *DockerInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	var script strings.Builder

	switch os {
	case OSLinux:
		switch pm {
		case PMApt:
			script.WriteString(`#!/bin/bash
set -e

echo "🐳 Installing Docker on Ubuntu/Debian..."

# Remove old versions
sudo apt-get remove -y docker docker-engine docker.io containerd runc 2>/dev/null || true

# Install dependencies
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg

# Add Docker's official GPG key
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

# Set up repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Add current user to docker group
sudo usermod -aG docker $USER

echo "✅ Docker installed successfully!"
echo "⚠️  Please log out and back in for group changes to take effect."
`)
		case PMYum, PMDnf:
			script.WriteString(`#!/bin/bash
set -e

echo "🐳 Installing Docker on RHEL/CentOS..."

# Remove old versions
sudo yum remove -y docker docker-client docker-client-latest docker-common docker-latest docker-latest-logrotate docker-logrotate docker-engine 2>/dev/null || true

# Install dependencies
sudo yum install -y yum-utils

# Add Docker repository
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo

# Install Docker Engine
sudo yum install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Start and enable Docker
sudo systemctl start docker
sudo systemctl enable docker

# Add current user to docker group
sudo usermod -aG docker $USER

echo "✅ Docker installed successfully!"
`)
		}

	case OSMacOS:
		script.WriteString(`#!/bin/bash
set -e

echo "🐳 Installing Docker on macOS..."

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    echo "Installing Homebrew first..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi

# Install Docker Desktop
brew install --cask docker

echo "✅ Docker Desktop installed!"
echo "⚠️  Please launch Docker Desktop from Applications."
`)

	case OSWindows:
		script.WriteString(`# PowerShell script for Windows
Write-Host "🐳 Installing Docker on Windows..." -ForegroundColor Cyan

# Check if running as admin
if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Host "Please run as Administrator!" -ForegroundColor Red
    exit 1
}

# Install via Chocolatey if available
if (Get-Command choco -ErrorAction SilentlyContinue) {
    choco install docker-desktop -y
} else {
    Write-Host "Please install Chocolatey first or download Docker Desktop manually from:"
    Write-Host "https://www.docker.com/products/docker-desktop/"
}

Write-Host "✅ Docker Desktop installed!" -ForegroundColor Green
`)
	}

	return script.String()
}

func (d *DockerInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	switch os {
	case OSLinux:
		if pm == PMApt {
			return `sudo apt-get purge -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
sudo rm -rf /var/lib/docker
sudo rm -rf /var/lib/containerd`
		}
		return `sudo yum remove -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin`
	case OSMacOS:
		return `brew uninstall --cask docker`
	case OSWindows:
		return `choco uninstall docker-desktop -y`
	default:
		return ""
	}
}

// Ensure DockerInstaller implements Installer
var _ Installer = (*DockerInstaller)(nil)
