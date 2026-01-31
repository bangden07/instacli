# Databases

Install and configure database systems for your applications.

## 🗄️ Available Installers

### MySQL

**Popular open-source relational database**

| Property | Value |
| -------- | ----- |
| Default Port | 3306 |
| Supported OS | Linux, macOS |
| Type | Relational (SQL) |

**Installation:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y mysql-server

# CentOS/RHEL
sudo yum install -y mysql-server

# macOS
brew install mysql
```

**Secure Installation:**
```bash
sudo mysql_secure_installation
```

**Create Database & User:**
```sql
CREATE DATABASE myapp;
CREATE USER 'myuser'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON myapp.* TO 'myuser'@'localhost';
FLUSH PRIVILEGES;
```

---

### PostgreSQL

**Advanced open-source relational database**

| Property | Value |
| -------- | ----- |
| Default Port | 5432 |
| Supported OS | Linux, macOS |
| Type | Relational (SQL) |

**Installation:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y postgresql postgresql-contrib

# CentOS/RHEL
sudo yum install -y postgresql-server postgresql-contrib
sudo postgresql-setup --initdb

# macOS
brew install postgresql
```

**Setup:**
```bash
# Start service
sudo systemctl enable postgresql
sudo systemctl start postgresql

# Access psql
sudo -u postgres psql
```

**Create Database & User:**
```sql
CREATE USER myuser WITH PASSWORD 'password';
CREATE DATABASE myapp OWNER myuser;
GRANT ALL PRIVILEGES ON DATABASE myapp TO myuser;
```

---

### MongoDB

**NoSQL document database**

| Property | Value |
| -------- | ----- |
| Default Port | 27017 |
| Supported OS | Linux |
| Type | NoSQL (Document) |

**Installation (Ubuntu):**
```bash
# Import GPG key
curl -fsSL https://pgp.mongodb.com/server-7.0.asc | sudo gpg -o /usr/share/keyrings/mongodb-server-7.0.gpg --dearmor

# Add repository
echo "deb [ signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list

# Install
sudo apt update
sudo apt install -y mongodb-org

# Start
sudo systemctl enable mongod
sudo systemctl start mongod
```

**Connect:**
```bash
mongosh
```

---

### Redis

**In-memory data structure store**

| Property | Value |
| -------- | ----- |
| Default Port | 6379 |
| Supported OS | Linux, macOS |
| Type | Key-Value (Cache) |

**Installation:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y redis-server

# CentOS/RHEL
sudo yum install -y redis

# macOS
brew install redis
```

**Configuration:**
```bash
# Edit config
sudo nano /etc/redis/redis.conf

# Set password (uncomment and edit)
requirepass yourpassword

# Enable persistence
appendonly yes

# Restart
sudo systemctl restart redis
```

**Test:**
```bash
redis-cli
> PING
PONG
> SET key "value"
OK
> GET key
"value"
```

---

## 📊 Comparison

| Database | Type | Best For |
| -------- | ---- | -------- |
| MySQL | SQL | Web apps, WordPress |
| PostgreSQL | SQL | Complex queries, GIS |
| MongoDB | NoSQL | Flexible schemas, JSON |
| Redis | Cache | Sessions, caching, queues |

## 🔒 Security Tips

1. **Change default ports**
2. **Use strong passwords**
3. **Limit network access** (bind to localhost)
4. **Enable SSL/TLS**
5. **Regular backups**

## 💾 Backup Commands

```bash
# MySQL
mysqldump -u root -p database > backup.sql

# PostgreSQL
pg_dump database > backup.sql

# MongoDB
mongodump --db database --out /backup/

# Redis
redis-cli BGSAVE
```

---

← [[Runtimes]] | [[Containers]] →
