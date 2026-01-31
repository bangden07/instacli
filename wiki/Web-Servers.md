# Web Servers

Install and configure web servers for hosting your applications.

## 🌐 Available Installers

### Nginx

**High-performance web server and reverse proxy**

| Property | Value |
| -------- | ----- |
| Category | Web Server |
| Supported OS | Linux, macOS |
| Default Port | 80, 443 |

**Features:**
- High concurrency handling
- Low memory footprint
- Reverse proxy support
- Load balancing
- SSL/TLS termination

**Installation:**
```bash
# Ubuntu/Debian
sudo apt update && sudo apt install -y nginx

# CentOS/RHEL
sudo yum install -y nginx

# macOS
brew install nginx
```

**Post-install:**
```bash
sudo systemctl enable nginx
sudo systemctl start nginx
```

---

### Apache

**Popular, feature-rich web server**

| Property | Value |
| -------- | ----- |
| Category | Web Server |
| Supported OS | Linux, macOS |
| Default Port | 80, 443 |

**Features:**
- .htaccess support
- mod_rewrite for URL rewriting
- Virtual hosts
- Extensive module ecosystem

**Installation:**
```bash
# Ubuntu/Debian
sudo apt update && sudo apt install -y apache2

# CentOS/RHEL
sudo yum install -y httpd

# macOS
brew install httpd
```

---

### LAMP Stack

**Linux, Apache, MySQL, PHP - Complete web development stack**

| Component | Description |
| --------- | ----------- |
| **L**inux | Operating system |
| **A**pache | Web server |
| **M**ySQL | Database |
| **P**HP | Programming language |

**Installation (Ubuntu):**
```bash
sudo apt update
sudo apt install -y apache2 mysql-server php libapache2-mod-php php-mysql
sudo systemctl enable apache2 mysql
sudo systemctl start apache2 mysql
```

**Verify:**
```bash
# Create test file
echo "<?php phpinfo(); ?>" | sudo tee /var/www/html/info.php

# Open http://localhost/info.php
```

---

### LEMP Stack

**Linux, Nginx, MySQL, PHP - High-performance stack**

| Component | Description |
| --------- | ----------- |
| **L**inux | Operating system |
| **E**ngine-X (Nginx) | Web server |
| **M**ySQL/MariaDB | Database |
| **P**HP | Programming language |

**Installation (Ubuntu):**
```bash
sudo apt update
sudo apt install -y nginx mysql-server php-fpm php-mysql
sudo systemctl enable nginx mysql php8.1-fpm
sudo systemctl start nginx mysql php8.1-fpm
```

**Configure Nginx for PHP:**
```nginx
server {
    listen 80;
    root /var/www/html;
    index index.php index.html;
    
    location ~ \.php$ {
        include snippets/fastcgi-php.conf;
        fastcgi_pass unix:/var/run/php/php8.1-fpm.sock;
    }
}
```

---

## 🔧 Configuration Files

| Server | Config Location |
| ------ | --------------- |
| Nginx | `/etc/nginx/nginx.conf` |
| Apache | `/etc/apache2/apache2.conf` |
| PHP | `/etc/php/8.1/fpm/php.ini` |

## 🔒 SSL with Certbot

After installing a web server, secure it with SSL:
```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d yourdomain.com
```

---

← [[Installers]] | [[Runtimes]] →
