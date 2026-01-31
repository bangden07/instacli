package installers

import "strings"

// LaravelKitInstaller handles Laravel development environment setup
type LaravelKitInstaller struct {
	BaseInstaller
}

func NewLaravelKitInstaller() *LaravelKitInstaller {
	return &LaravelKitInstaller{
		BaseInstaller: NewBaseInstaller(
			"Laravel Kit",
			"Complete Laravel development environment",
			CategoryFramework,
			"🔴",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (l *LaravelKitInstaller) Dependencies() []string {
	return []string{"php", "composer", "node", "npm"}
}

func (l *LaravelKitInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMBrew, PMChoco}
}

func (l *LaravelKitInstaller) Install(executor Executor) error {
	script := l.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager())
	return executor.RunWithProgress(script, func(line string) {})
}

func (l *LaravelKitInstaller) Uninstall(executor Executor) error {
	script := l.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager())
	_, err := executor.Run(script)
	return err
}

func (l *LaravelKitInstaller) IsInstalled(executor Executor) (bool, error) {
	output, err := executor.Run("composer global show laravel/installer 2>/dev/null")
	if err != nil {
		return false, nil
	}
	return strings.Contains(output, "laravel/installer"), nil
}

func (l *LaravelKitInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	if os == OSLinux && pm == PMApt {
		return `#!/bin/bash
set -e

echo "🔴 Installing Laravel Kit..."

# Install PHP and extensions
sudo apt update
sudo apt install -y software-properties-common
sudo add-apt-repository -y ppa:ondrej/php
sudo apt update
sudo apt install -y php8.3 php8.3-cli php8.3-common php8.3-curl php8.3-mbstring \
    php8.3-mysql php8.3-xml php8.3-zip php8.3-gd php8.3-bcmath php8.3-intl \
    php8.3-sqlite3 php8.3-redis unzip

# Install Composer
echo "Installing Composer..."
curl -sS https://getcomposer.org/installer | php
sudo mv composer.phar /usr/local/bin/composer
composer --version

# Install Laravel Installer
echo "Installing Laravel Installer..."
composer global require laravel/installer

# Add Composer bin to PATH
if ! grep -q '.composer/vendor/bin' ~/.bashrc; then
    echo 'export PATH="$HOME/.config/composer/vendor/bin:$PATH"' >> ~/.bashrc
fi

# Install Node.js (for frontend assets)
echo "Installing Node.js..."
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt install -y nodejs

echo ""
echo "✅ Laravel Kit installed successfully!"
echo ""
echo "📦 Create new project: laravel new my-app"
echo "🚀 Run development: php artisan serve"
echo "🎨 Install frontend: npm install && npm run dev"
`
	}

	if os == OSMacOS {
		return `#!/bin/bash
set -e

echo "🔴 Installing Laravel Kit on macOS..."

# Install PHP via Homebrew
brew install php

# Install Composer
brew install composer

# Install Laravel Installer
composer global require laravel/installer

# Install Node.js
brew install node

echo ""
echo "✅ Laravel Kit installed!"
echo "📦 Create new project: laravel new my-app"
`
	}

	return `# Windows - install via Chocolatey
choco install php composer nodejs -y
composer global require laravel/installer
Write-Host "✅ Laravel Kit installed!" -ForegroundColor Green`
}

func (l *LaravelKitInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `composer global remove laravel/installer`
}

var _ Installer = (*LaravelKitInstaller)(nil)

// NextJSKitInstaller handles Next.js development environment setup
type NextJSKitInstaller struct {
	BaseInstaller
}

func NewNextJSKitInstaller() *NextJSKitInstaller {
	return &NextJSKitInstaller{
		BaseInstaller: NewBaseInstaller(
			"Next.js Kit",
			"Complete Next.js development environment",
			CategoryFramework,
			"▲",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (n *NextJSKitInstaller) Dependencies() []string {
	return []string{"node", "npm"}
}

func (n *NextJSKitInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMBrew, PMChoco}
}

func (n *NextJSKitInstaller) Install(executor Executor) error {
	script := n.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager())
	return executor.RunWithProgress(script, func(line string) {})
}

func (n *NextJSKitInstaller) Uninstall(executor Executor) error {
	script := n.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager())
	_, err := executor.Run(script)
	return err
}

func (n *NextJSKitInstaller) IsInstalled(executor Executor) (bool, error) {
	output, err := executor.Run("npm list -g create-next-app 2>/dev/null")
	if err != nil {
		// Check if npx is available (we don't need global install)
		output2, err2 := executor.Run("npx --version")
		if err2 != nil {
			return false, nil
		}
		return output2 != "", nil
	}
	return strings.Contains(output, "create-next-app"), nil
}

func (n *NextJSKitInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e

echo "▲ Installing Next.js Kit..."

# Install NVM and Node.js
export NVM_DIR="$HOME/.nvm"
if [ ! -d "$NVM_DIR" ]; then
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash
fi
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
nvm install --lts
nvm use --lts

# Install pnpm (recommended for Next.js)
npm install -g pnpm

# Install useful global packages
npm install -g typescript ts-node

echo ""
echo "✅ Next.js Kit installed successfully!"
echo ""
echo "📦 Create new project (recommended):"
echo "   pnpm create next-app my-app"
echo "   # or"
echo "   npx create-next-app@latest my-app"
echo ""
echo "🚀 Development: pnpm dev"
echo "🏗️  Build: pnpm build"
echo "🌐 Start production: pnpm start"
`
}

func (n *NextJSKitInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `npm uninstall -g pnpm typescript ts-node`
}

var _ Installer = (*NextJSKitInstaller)(nil)
