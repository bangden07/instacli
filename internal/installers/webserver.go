package installers

import "strings"

// LAMPInstaller handles LAMP stack installation
type LAMPInstaller struct {
	BaseInstaller
	phpVersion string
}

func NewLAMPInstaller() *LAMPInstaller {
	return &LAMPInstaller{
		BaseInstaller: NewBaseInstaller(
			"LAMP Stack",
			"Linux, Apache, MySQL, PHP",
			CategoryWebServer,
			"🌐",
			[]OS{OSLinux},
		),
		phpVersion: "8.3",
	}
}

func (l *LAMPInstaller) Dependencies() []string {
	return []string{}
}

func (l *LAMPInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum, PMDnf}
}

func (l *LAMPInstaller) Install(executor Executor) error {
	script := l.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager())
	return executor.RunWithProgress(script, func(line string) {})
}

func (l *LAMPInstaller) Uninstall(executor Executor) error {
	script := l.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager())
	_, err := executor.Run(script)
	return err
}

func (l *LAMPInstaller) IsInstalled(executor Executor) (bool, error) {
	output, err := executor.Run("apache2 -v && php -v && mysql --version")
	if err != nil {
		return false, nil
	}
	return strings.Contains(output, "Apache") && strings.Contains(output, "PHP"), nil
}

func (l *LAMPInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	if pm == PMApt {
		return `#!/bin/bash
set -e

echo "🌐 Installing LAMP Stack..."

# Update system
sudo apt update && sudo apt upgrade -y

# Install Apache
echo "Installing Apache..."
sudo apt install -y apache2
sudo systemctl start apache2
sudo systemctl enable apache2

# Install MySQL
echo "Installing MySQL..."
sudo apt install -y mysql-server
sudo systemctl start mysql
sudo systemctl enable mysql

# Secure MySQL (optional - requires interaction)
# sudo mysql_secure_installation

# Add PHP repository for latest version
echo "Installing PHP 8.3..."
sudo apt install -y software-properties-common
sudo add-apt-repository -y ppa:ondrej/php
sudo apt update

# Install PHP and common extensions
sudo apt install -y php8.3 php8.3-fpm php8.3-mysql php8.3-curl php8.3-gd \
    php8.3-mbstring php8.3-xml php8.3-zip php8.3-bcmath php8.3-intl \
    libapache2-mod-php8.3

# Enable Apache modules
sudo a2enmod php8.3
sudo a2enmod rewrite

# Restart services
sudo systemctl restart apache2

# Set permissions
sudo chown -R $USER:www-data /var/www/html
sudo chmod -R 755 /var/www/html

echo ""
echo "✅ LAMP Stack installed successfully!"
echo ""
echo "📁 Web root: /var/www/html"
echo "🌐 Apache: http://localhost"
echo "🐬 MySQL: mysql -u root -p"
echo "🐘 PHP version: $(php -v | head -n 1)"
`
	}

	// CentOS/RHEL
	return `#!/bin/bash
set -e

echo "🌐 Installing LAMP Stack on RHEL/CentOS..."

# Install Apache
sudo yum install -y httpd
sudo systemctl start httpd
sudo systemctl enable httpd

# Install MySQL
sudo yum install -y mysql-server
sudo systemctl start mysqld
sudo systemctl enable mysqld

# Install PHP
sudo yum install -y php php-mysqlnd php-fpm php-gd php-mbstring php-xml php-curl

# Restart Apache
sudo systemctl restart httpd

echo "✅ LAMP Stack installed!"
`
}

func (l *LAMPInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	if pm == PMApt {
		return `sudo apt purge -y apache2 mysql-server php8.3* libapache2-mod-php8.3
sudo apt autoremove -y`
	}
	return `sudo yum remove -y httpd mysql-server php*`
}

var _ Installer = (*LAMPInstaller)(nil)

// LEMPInstaller handles LEMP stack installation
type LEMPInstaller struct {
	BaseInstaller
	phpVersion string
}

func NewLEMPInstaller() *LEMPInstaller {
	return &LEMPInstaller{
		BaseInstaller: NewBaseInstaller(
			"LEMP Stack",
			"Linux, Nginx, MySQL, PHP-FPM",
			CategoryWebServer,
			"🌐",
			[]OS{OSLinux},
		),
		phpVersion: "8.3",
	}
}

func (l *LEMPInstaller) Dependencies() []string {
	return []string{}
}

func (l *LEMPInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum, PMDnf}
}

func (l *LEMPInstaller) Install(executor Executor) error {
	script := l.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager())
	return executor.RunWithProgress(script, func(line string) {})
}

func (l *LEMPInstaller) Uninstall(executor Executor) error {
	script := l.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager())
	_, err := executor.Run(script)
	return err
}

func (l *LEMPInstaller) IsInstalled(executor Executor) (bool, error) {
	output, err := executor.Run("nginx -v && php -v && mysql --version")
	if err != nil {
		return false, nil
	}
	return strings.Contains(output, "nginx") && strings.Contains(output, "PHP"), nil
}

func (l *LEMPInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	if pm == PMApt {
		return `#!/bin/bash
set -e

echo "🌐 Installing LEMP Stack..."

# Update system
sudo apt update && sudo apt upgrade -y

# Install Nginx
echo "Installing Nginx..."
sudo apt install -y nginx
sudo systemctl start nginx
sudo systemctl enable nginx

# Install MySQL
echo "Installing MySQL..."
sudo apt install -y mysql-server
sudo systemctl start mysql
sudo systemctl enable mysql

# Add PHP repository
echo "Installing PHP 8.3 FPM..."
sudo apt install -y software-properties-common
sudo add-apt-repository -y ppa:ondrej/php
sudo apt update

# Install PHP-FPM
sudo apt install -y php8.3-fpm php8.3-mysql php8.3-curl php8.3-gd \
    php8.3-mbstring php8.3-xml php8.3-zip php8.3-bcmath php8.3-intl

# Configure PHP-FPM
sudo systemctl start php8.3-fpm
sudo systemctl enable php8.3-fpm

# Create sample Nginx config for PHP
cat << 'NGINX' | sudo tee /etc/nginx/sites-available/default
server {
    listen 80 default_server;
    listen [::]:80 default_server;

    root /var/www/html;
    index index.php index.html index.htm;

    server_name _;

    location / {
        try_files $uri $uri/ =404;
    }

    location ~ \.php$ {
        include snippets/fastcgi-php.conf;
        fastcgi_pass unix:/var/run/php/php8.3-fpm.sock;
    }

    location ~ /\.ht {
        deny all;
    }
}
NGINX

# Restart services
sudo systemctl restart nginx
sudo systemctl restart php8.3-fpm

echo ""
echo "✅ LEMP Stack installed successfully!"
echo ""
echo "📁 Web root: /var/www/html"
echo "🌐 Nginx: http://localhost"
echo "🐬 MySQL: mysql -u root -p"
echo "🐘 PHP-FPM version: $(php -v | head -n 1)"
`
	}

	return `#!/bin/bash
echo "Installing LEMP on CentOS..."
sudo yum install -y nginx mysql-server php php-fpm php-mysqlnd
sudo systemctl start nginx mysqld php-fpm
sudo systemctl enable nginx mysqld php-fpm
echo "✅ LEMP installed!"`
}

func (l *LEMPInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	if pm == PMApt {
		return `sudo apt purge -y nginx mysql-server php8.3*
sudo apt autoremove -y`
	}
	return `sudo yum remove -y nginx mysql-server php*`
}

var _ Installer = (*LEMPInstaller)(nil)
