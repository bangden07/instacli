# Security

Firewall, SSL certificates, and intrusion prevention tools.

## 🛡️ Available Installers

### UFW Firewall

**Uncomplicated Firewall for Linux**

| Property | Value |
| -------- | ----- |
| Supported OS | Linux (Debian/Ubuntu) |
| Backend | iptables |

**Installation:**
```bash
sudo apt update
sudo apt install -y ufw
```

**Basic Usage:**
```bash
# Enable firewall
sudo ufw enable

# Check status
sudo ufw status verbose

# Allow SSH (important!)
sudo ufw allow ssh
# or
sudo ufw allow 22

# Allow HTTP/HTTPS
sudo ufw allow 80
sudo ufw allow 443

# Allow specific port
sudo ufw allow 3000

# Allow from specific IP
sudo ufw allow from 192.168.1.100

# Deny port
sudo ufw deny 3306

# Delete rule
sudo ufw delete allow 3000

# Reset all rules
sudo ufw reset
```

**Common Rules:**
```bash
# Web server
sudo ufw allow 'Nginx Full'
sudo ufw allow 'Apache Full'

# Database (only from specific IP)
sudo ufw allow from 192.168.1.0/24 to any port 3306

# SSH from specific IP only
sudo ufw allow from 203.0.113.0/24 to any port 22
```

---

### Certbot SSL

**Free SSL certificates from Let's Encrypt**

| Property | Value |
| -------- | ----- |
| Supported OS | Linux |
| Certificate Validity | 90 days (auto-renew) |

**Installation:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y certbot

# For Nginx
sudo apt install -y python3-certbot-nginx

# For Apache
sudo apt install -y python3-certbot-apache
```

**Get Certificate (Nginx):**
```bash
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com
```

**Get Certificate (Apache):**
```bash
sudo certbot --apache -d yourdomain.com -d www.yourdomain.com
```

**Standalone (no web server):**
```bash
sudo certbot certonly --standalone -d yourdomain.com
```

**Auto-Renewal:**
```bash
# Test renewal
sudo certbot renew --dry-run

# Renewal is automatic via systemd timer
sudo systemctl status certbot.timer
```

**Certificate Location:**
```
/etc/letsencrypt/live/yourdomain.com/
├── cert.pem       # Certificate
├── chain.pem      # Intermediate chain
├── fullchain.pem  # Cert + chain
└── privkey.pem    # Private key
```

---

### Fail2ban

**Intrusion prevention system**

| Property | Value |
| -------- | ----- |
| Supported OS | Linux |
| Config Location | `/etc/fail2ban/` |

**Installation:**
```bash
sudo apt update
sudo apt install -y fail2ban
```

**Configuration:**
```bash
# Copy default config
sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local

# Edit
sudo nano /etc/fail2ban/jail.local
```

**Basic Configuration:**
```ini
# /etc/fail2ban/jail.local
[DEFAULT]
bantime = 1h
findtime = 10m
maxretry = 5
backend = systemd

[sshd]
enabled = true
port = ssh
filter = sshd
logpath = /var/log/auth.log
maxretry = 3
bantime = 24h
```

**Common Jails:**
```ini
# Nginx auth
[nginx-http-auth]
enabled = true
filter = nginx-http-auth
port = http,https
logpath = /var/log/nginx/error.log

# Nginx bad bots
[nginx-botsearch]
enabled = true
filter = nginx-botsearch
port = http,https
logpath = /var/log/nginx/access.log
maxretry = 2
```

**Commands:**
```bash
# Start/restart
sudo systemctl restart fail2ban

# Check status
sudo fail2ban-client status

# Check specific jail
sudo fail2ban-client status sshd

# Unban IP
sudo fail2ban-client set sshd unbanip 192.168.1.100

# Ban IP manually
sudo fail2ban-client set sshd banip 192.168.1.100
```

---

## 🔒 Security Best Practices

### 1. SSH Hardening
```bash
# /etc/ssh/sshd_config
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
Port 2222  # Change default port
```

### 2. Firewall Rules Order
```bash
# 1. Allow SSH first!
sudo ufw allow ssh
# 2. Enable firewall
sudo ufw enable
# 3. Add other rules
```

### 3. Regular Updates
```bash
# Enable automatic updates
sudo apt install unattended-upgrades
sudo dpkg-reconfigure unattended-upgrades
```

### 4. Monitoring
- Check `/var/log/auth.log` for SSH attempts
- Check fail2ban logs: `/var/log/fail2ban.log`
- Check firewall logs: `sudo ufw status verbose`

---

← [[CI CD]] | [[CMS]] →
