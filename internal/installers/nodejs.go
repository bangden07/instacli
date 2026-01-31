package installers

import (
	"fmt"
	"strings"
)

// NodeJSInstaller handles Node.js installation via NVM
type NodeJSInstaller struct {
	BaseInstaller
	version string // "lts" or specific version like "20"
}

func NewNodeJSInstaller() *NodeJSInstaller {
	return &NodeJSInstaller{
		BaseInstaller: NewBaseInstaller(
			"Node.js",
			"Node.js runtime with NVM (version manager)",
			CategoryRuntime,
			"⬢",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
		version: "lts",
	}
}

func (n *NodeJSInstaller) WithVersion(version string) *NodeJSInstaller {
	n.version = version
	return n
}

func (n *NodeJSInstaller) Dependencies() []string {
	return []string{"curl", "git"}
}

func (n *NodeJSInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum, PMDnf, PMBrew, PMChoco}
}

func (n *NodeJSInstaller) Install(executor Executor) error {
	script := n.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager())
	return executor.RunWithProgress(script, func(line string) {
		// Progress callback
	})
}

func (n *NodeJSInstaller) Uninstall(executor Executor) error {
	script := n.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager())
	_, err := executor.Run(script)
	return err
}

func (n *NodeJSInstaller) IsInstalled(executor Executor) (bool, error) {
	output, err := executor.Run("node --version")
	if err != nil {
		return false, nil
	}
	return strings.HasPrefix(output, "v"), nil
}

func (n *NodeJSInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	var script strings.Builder

	switch os {
	case OSLinux, OSMacOS:
		script.WriteString(fmt.Sprintf(`#!/bin/bash
set -e

echo "⬢ Installing Node.js via NVM..."

# Install NVM
export NVM_DIR="$HOME/.nvm"
if [ ! -d "$NVM_DIR" ]; then
    echo "Installing NVM..."
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash
fi

# Load NVM
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"

# Install Node.js
echo "Installing Node.js %s..."
nvm install %s
nvm use %s
nvm alias default %s

# Verify installation
echo ""
echo "✅ Node.js installed successfully!"
node --version
npm --version

# Install some global packages
echo ""
echo "Installing useful global packages..."
npm install -g pnpm yarn

echo ""
echo "🎉 Setup complete!"
echo "Available commands: node, npm, npx, pnpm, yarn"
`, n.version, n.version, n.version, n.version))

	case OSWindows:
		script.WriteString(fmt.Sprintf(`# PowerShell script for Windows
Write-Host "⬢ Installing Node.js on Windows..." -ForegroundColor Green

# Option 1: Using Chocolatey
if (Get-Command choco -ErrorAction SilentlyContinue) {
    choco install nvm -y
    refreshenv
    nvm install %s
    nvm use %s
} else {
    # Option 2: Using winget
    if (Get-Command winget -ErrorAction SilentlyContinue) {
        winget install CoreyButler.NVMforWindows
        Write-Host "Please restart your terminal and run: nvm install %s"
    } else {
        Write-Host "Please install NVM for Windows manually from:"
        Write-Host "https://github.com/coreybutler/nvm-windows/releases"
    }
}

Write-Host "✅ Node.js installation initiated!" -ForegroundColor Green
`, n.version, n.version, n.version))
	}

	return script.String()
}

func (n *NodeJSInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	switch os {
	case OSLinux, OSMacOS:
		return `#!/bin/bash
# Remove NVM and Node.js
rm -rf "$HOME/.nvm"
# Remove NVM lines from shell config
sed -i '/NVM_DIR/d' ~/.bashrc ~/.zshrc 2>/dev/null || true
echo "Node.js and NVM removed. Please restart your terminal."`
	case OSWindows:
		return `# Uninstall via Chocolatey
choco uninstall nvm nodejs -y`
	default:
		return ""
	}
}

// Ensure NodeJSInstaller implements Installer
var _ Installer = (*NodeJSInstaller)(nil)
