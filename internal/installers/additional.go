package installers

// ============================================================
// Pi-hole DNS Installer
// ============================================================
type PiholeInstaller struct {
	BaseInstaller
}

func NewPiholeInstaller() *PiholeInstaller {
	return &PiholeInstaller{
		BaseInstaller: BaseInstaller{
			name:        "Pi-hole",
			description: "Network-wide ad blocking DNS",
			category:    CategoryDNS,
			icon:        "🕳️",
			supportedOS: []OS{OSLinux},
		},
	}
}

func (i *PiholeInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum}
}
func (i *PiholeInstaller) Dependencies() []string { return []string{"curl"} }
func (i *PiholeInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *PiholeInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *PiholeInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("pihole -v")
	return err == nil, nil
}

func (i *PiholeInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "🕳️ Installing Pi-hole..."

# Run official installer
curl -sSL https://install.pi-hole.net | bash

echo "✅ Pi-hole installed!"
echo "🌐 Admin: http://localhost/admin"
echo "🔧 Change password: pihole -a -p"`
}

func (i *PiholeInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
pihole uninstall`
}

// ============================================================
// WordPress Installer
// ============================================================
type WordPressInstaller struct {
	BaseInstaller
}

func NewWordPressInstaller() *WordPressInstaller {
	return &WordPressInstaller{
		BaseInstaller: BaseInstaller{
			name:        "WordPress",
			description: "Popular CMS & blogging platform",
			category:    CategoryCMS,
			icon:        "📝",
			supportedOS: []OS{OSLinux},
		},
	}
}

func (i *WordPressInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum}
}
func (i *WordPressInstaller) Dependencies() []string { return []string{"docker"} }
func (i *WordPressInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *WordPressInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *WordPressInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("docker ps | grep wordpress")
	return err == nil, nil
}

func (i *WordPressInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "📝 Installing WordPress..."

sudo mkdir -p /opt/wordpress
cd /opt/wordpress

cat <<EOF > docker-compose.yml
version: '3.8'
services:
  db:
    image: mysql:8.0
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: wordpress_root
      MYSQL_DATABASE: wordpress
      MYSQL_USER: wordpress
      MYSQL_PASSWORD: wordpress
    volumes:
      - db_data:/var/lib/mysql

  wordpress:
    image: wordpress:latest
    restart: unless-stopped
    depends_on:
      - db
    ports:
      - "8080:80"
    environment:
      WORDPRESS_DB_HOST: db
      WORDPRESS_DB_USER: wordpress
      WORDPRESS_DB_PASSWORD: wordpress
      WORDPRESS_DB_NAME: wordpress
    volumes:
      - wp_data:/var/www/html

volumes:
  db_data:
  wp_data:
EOF

docker compose up -d

echo "✅ WordPress installed!"
echo "🌐 Access: http://localhost:8080"`
}

func (i *WordPressInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
cd /opt/wordpress && docker compose down -v
sudo rm -rf /opt/wordpress`
}

// ============================================================
// Ghost CMS Installer
// ============================================================
type GhostInstaller struct {
	BaseInstaller
}

func NewGhostInstaller() *GhostInstaller {
	return &GhostInstaller{
		BaseInstaller: BaseInstaller{
			name:        "Ghost",
			description: "Modern publishing platform",
			category:    CategoryCMS,
			icon:        "👻",
			supportedOS: []OS{OSLinux},
		},
	}
}

func (i *GhostInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum}
}
func (i *GhostInstaller) Dependencies() []string { return []string{"docker"} }
func (i *GhostInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *GhostInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *GhostInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("docker ps | grep ghost")
	return err == nil, nil
}

func (i *GhostInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "👻 Installing Ghost CMS..."

sudo mkdir -p /opt/ghost
cd /opt/ghost

cat <<EOF > docker-compose.yml
version: '3.8'
services:
  ghost:
    image: ghost:latest
    restart: unless-stopped
    ports:
      - "2368:2368"
    environment:
      url: http://localhost:2368
    volumes:
      - ghost_data:/var/lib/ghost/content

volumes:
  ghost_data:
EOF

docker compose up -d

echo "✅ Ghost CMS installed!"
echo "🌐 Access: http://localhost:2368"
echo "📝 Admin: http://localhost:2368/ghost"`
}

func (i *GhostInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
cd /opt/ghost && docker compose down -v
sudo rm -rf /opt/ghost`
}

// ============================================================
// Restic Backup Installer
// ============================================================
type ResticInstaller struct {
	BaseInstaller
}

func NewResticInstaller() *ResticInstaller {
	return &ResticInstaller{
		BaseInstaller: BaseInstaller{
			name:        "Restic",
			description: "Fast, secure backup program",
			category:    CategoryBackup,
			icon:        "💾",
			supportedOS: []OS{OSLinux, OSMacOS},
		},
	}
}

func (i *ResticInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum, PMBrew}
}
func (i *ResticInstaller) Dependencies() []string { return []string{} }
func (i *ResticInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *ResticInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *ResticInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("restic version")
	return err == nil, nil
}

func (i *ResticInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "💾 Installing Restic..."

if [ -f /etc/debian_version ]; then
    sudo apt-get update
    sudo apt-get install -y restic
elif [ -f /etc/redhat-release ]; then
    sudo yum install -y epel-release
    sudo yum install -y restic
elif command -v brew &> /dev/null; then
    brew install restic
fi

echo "✅ Restic installed!"
echo ""
echo "📝 Quick start:"
echo "   restic init --repo /path/to/backup"
echo "   restic backup /path/to/data --repo /path/to/backup"
restic version`
}

func (i *ResticInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
sudo apt-get remove -y restic || sudo yum remove -y restic || brew uninstall restic`
}
