package installers

// ============================================================
// UFW Firewall Installer
// ============================================================
type UFWInstaller struct {
	BaseInstaller
}

func NewUFWInstaller() *UFWInstaller {
	return &UFWInstaller{
		BaseInstaller: BaseInstaller{
			name:        "UFW Firewall",
			description: "Uncomplicated Firewall for Linux",
			category:    CategorySecurity,
			icon:        "🛡️",
			supportedOS: []OS{OSLinux},
		},
	}
}

func (i *UFWInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum}
}
func (i *UFWInstaller) Dependencies() []string { return []string{} }
func (i *UFWInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *UFWInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *UFWInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("ufw version")
	return err == nil, nil
}

func (i *UFWInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "📦 Installing UFW Firewall..."
if [ -f /etc/debian_version ]; then
    sudo apt-get update && sudo apt-get install -y ufw
elif [ -f /etc/redhat-release ]; then
    sudo yum install -y epel-release && sudo yum install -y ufw
fi
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
echo "y" | sudo ufw enable
echo "✅ UFW installed and configured!"
sudo ufw status`
}

func (i *UFWInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
sudo ufw disable && sudo apt-get remove -y ufw`
}

// ============================================================
// Certbot SSL Installer
// ============================================================
type CertbotInstaller struct {
	BaseInstaller
}

func NewCertbotInstaller() *CertbotInstaller {
	return &CertbotInstaller{
		BaseInstaller: BaseInstaller{
			name:        "Certbot SSL",
			description: "Let's Encrypt SSL Certificate Manager",
			category:    CategorySecurity,
			icon:        "🔒",
			supportedOS: []OS{OSLinux, OSMacOS},
		},
	}
}

func (i *CertbotInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum, PMBrew}
}
func (i *CertbotInstaller) Dependencies() []string { return []string{} }
func (i *CertbotInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *CertbotInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *CertbotInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("certbot --version")
	return err == nil, nil
}

func (i *CertbotInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "📦 Installing Certbot..."
if [ -f /etc/debian_version ]; then
    sudo apt-get update && sudo apt-get install -y certbot
    command -v nginx && sudo apt-get install -y python3-certbot-nginx
    command -v apache2 && sudo apt-get install -y python3-certbot-apache
elif [ -f /etc/redhat-release ]; then
    sudo yum install -y epel-release certbot
elif command -v brew &> /dev/null; then
    brew install certbot
fi
echo "✅ Certbot installed!"
certbot --version
echo "Usage: sudo certbot --nginx -d yourdomain.com"`
}

func (i *CertbotInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
sudo apt-get remove -y certbot || brew uninstall certbot`
}

// ============================================================
// Fail2ban Installer
// ============================================================
type Fail2banInstaller struct {
	BaseInstaller
}

func NewFail2banInstaller() *Fail2banInstaller {
	return &Fail2banInstaller{
		BaseInstaller: BaseInstaller{
			name:        "Fail2ban",
			description: "Intrusion Prevention System",
			category:    CategorySecurity,
			icon:        "🚫",
			supportedOS: []OS{OSLinux},
		},
	}
}

func (i *Fail2banInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum}
}
func (i *Fail2banInstaller) Dependencies() []string { return []string{} }
func (i *Fail2banInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *Fail2banInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *Fail2banInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("fail2ban-client --version")
	return err == nil, nil
}

func (i *Fail2banInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "📦 Installing Fail2ban..."
if [ -f /etc/debian_version ]; then
    sudo apt-get update && sudo apt-get install -y fail2ban
elif [ -f /etc/redhat-release ]; then
    sudo yum install -y epel-release fail2ban
fi
sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local 2>/dev/null || true
sudo systemctl enable fail2ban && sudo systemctl start fail2ban
echo "✅ Fail2ban installed!"
sudo fail2ban-client status`
}

func (i *Fail2banInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
sudo systemctl stop fail2ban && sudo apt-get remove -y fail2ban`
}
