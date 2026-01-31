# CMS & Blog

Content Management Systems and blogging platforms.

## 📝 Available Installers

### WordPress

**World's most popular CMS**

| Property | Value |
| -------- | ----- |
| Default Port | 80/443 |
| Requirements | PHP, MySQL/MariaDB |
| Admin Path | `/wp-admin` |

**Docker Installation (Recommended):**
```yaml
# docker-compose.yml
version: '3.8'

services:
  wordpress:
    image: wordpress:latest
    restart: always
    ports:
      - "8080:80"
    environment:
      WORDPRESS_DB_HOST: db
      WORDPRESS_DB_USER: wordpress
      WORDPRESS_DB_PASSWORD: wordpress
      WORDPRESS_DB_NAME: wordpress
    volumes:
      - wordpress_data:/var/www/html

  db:
    image: mysql:8.0
    restart: always
    environment:
      MYSQL_DATABASE: wordpress
      MYSQL_USER: wordpress
      MYSQL_PASSWORD: wordpress
      MYSQL_ROOT_PASSWORD: rootpassword
    volumes:
      - db_data:/var/lib/mysql

volumes:
  wordpress_data:
  db_data:
```

```bash
docker compose up -d
# Access: http://localhost:8080
```

**Manual Installation:**
```bash
# Install LAMP stack first
sudo apt install -y apache2 mysql-server php php-mysql libapache2-mod-php php-curl php-gd php-xml php-mbstring

# Download WordPress
cd /tmp
curl -O https://wordpress.org/latest.tar.gz
tar xzf latest.tar.gz

# Move to web root
sudo mv wordpress /var/www/html/wordpress

# Set permissions
sudo chown -R www-data:www-data /var/www/html/wordpress
sudo chmod -R 755 /var/www/html/wordpress

# Create database
sudo mysql -e "CREATE DATABASE wordpress; CREATE USER 'wpuser'@'localhost' IDENTIFIED BY 'password'; GRANT ALL ON wordpress.* TO 'wpuser'@'localhost';"
```

**WP-CLI (Command Line):**
```bash
# Install WP-CLI
curl -O https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar
chmod +x wp-cli.phar
sudo mv wp-cli.phar /usr/local/bin/wp

# Usage
wp core download
wp config create --dbname=wordpress --dbuser=wpuser --dbpass=password
wp core install --url=example.com --title="My Blog" --admin_user=admin --admin_password=password --admin_email=admin@example.com
```

---

### Ghost

**Modern publishing platform**

| Property | Value |
| -------- | ----- |
| Default Port | 2368 |
| Requirements | Node.js 18+ |
| Admin Path | `/ghost` |

**Docker Installation:**
```yaml
# docker-compose.yml
version: '3.8'

services:
  ghost:
    image: ghost:5-alpine
    restart: always
    ports:
      - "2368:2368"
    environment:
      url: http://localhost:2368
      database__client: mysql
      database__connection__host: db
      database__connection__user: ghost
      database__connection__password: ghost
      database__connection__database: ghost
    volumes:
      - ghost_content:/var/lib/ghost/content
    depends_on:
      - db

  db:
    image: mysql:8.0
    restart: always
    environment:
      MYSQL_DATABASE: ghost
      MYSQL_USER: ghost
      MYSQL_PASSWORD: ghost
      MYSQL_ROOT_PASSWORD: rootpassword
    volumes:
      - ghost_db:/var/lib/mysql

volumes:
  ghost_content:
  ghost_db:
```

**Manual Installation (Ghost-CLI):**
```bash
# Install Node.js 18
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

# Install Ghost-CLI
sudo npm install -g ghost-cli

# Create directory
sudo mkdir -p /var/www/ghost
sudo chown $USER:$USER /var/www/ghost
cd /var/www/ghost

# Install Ghost
ghost install

# Answer prompts:
# Blog URL: https://yourdomain.com
# MySQL hostname: localhost
# MySQL username: ghost
# MySQL password: (create in MySQL first)
# Ghost database name: ghost
# Set up Nginx: yes
# Set up SSL: yes
# Set up systemd: yes
```

**Ghost CLI Commands:**
```bash
cd /var/www/ghost

# Start/stop
ghost start
ghost stop
ghost restart

# Update
ghost update

# Check status
ghost status

# View logs
ghost log
```

---

## 📊 Comparison

| Feature | WordPress | Ghost |
| ------- | --------- | ----- |
| Ease of Use | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| Performance | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| Themes | 10,000+ | 100+ |
| Plugins | 60,000+ | Limited |
| SEO | ⭐⭐⭐⭐ (with plugins) | ⭐⭐⭐⭐⭐ (built-in) |
| Best For | Any website | Blogs, publications |

## 🔧 Recommended Plugins/Integrations

### WordPress
- **Yoast SEO** - SEO optimization
- **WP Super Cache** - Caching
- **Wordfence** - Security
- **Elementor** - Page builder

### Ghost
- Built-in SEO
- Built-in memberships
- Built-in newsletters
- Zapier integration

## 🔒 Security Tips

1. **Keep updated** - Both platforms regularly release security updates
2. **Strong passwords** - Use password manager
3. **Limit login attempts** - Use fail2ban or plugin
4. **Backup regularly** - Automate daily backups
5. **SSL certificate** - Always use HTTPS

---

← [[Security]] | [[Installers]] →
