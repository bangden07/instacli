# Containers

Container runtimes and orchestration tools.

## 🐳 Available Installers

### Docker

**Industry-standard container runtime**

| Property | Value |
| -------- | ----- |
| Supported OS | Linux, macOS, Windows |
| Default Socket | `/var/run/docker.sock` |

**Installation (Linux):**
```bash
# Quick install
curl -fsSL https://get.docker.com | sh

# Add user to docker group
sudo usermod -aG docker $USER

# Start service
sudo systemctl enable docker
sudo systemctl start docker

# Verify
docker --version
docker run hello-world
```

**Installation (macOS):**
```bash
# Install Docker Desktop
brew install --cask docker

# Or download from docker.com
```

**Basic Commands:**
```bash
# Run container
docker run -d -p 80:80 nginx

# List containers
docker ps

# Stop container
docker stop <container_id>

# Remove container
docker rm <container_id>

# List images
docker images

# Pull image
docker pull nginx:latest

# Build image
docker build -t myapp .
```

---

### Docker Compose

**Multi-container Docker applications**

| Property | Value |
| -------- | ----- |
| Supported OS | Linux, macOS, Windows |
| Config File | `docker-compose.yml` |

**Installation:**
```bash
# Docker Compose V2 (plugin)
sudo apt install docker-compose-plugin

# Verify
docker compose version
```

**Example docker-compose.yml:**
```yaml
version: '3.8'

services:
  web:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./html:/usr/share/nginx/html

  db:
    image: mysql:8
    environment:
      MYSQL_ROOT_PASSWORD: secret
      MYSQL_DATABASE: myapp
    volumes:
      - db_data:/var/lib/mysql

  redis:
    image: redis:alpine
    ports:
      - "6379:6379"

volumes:
  db_data:
```

**Commands:**
```bash
# Start services
docker compose up -d

# Stop services
docker compose down

# View logs
docker compose logs -f

# Rebuild
docker compose up -d --build

# Scale service
docker compose up -d --scale web=3
```

---

## 📦 Dockerfile Example

```dockerfile
FROM node:20-alpine

WORKDIR /app

COPY package*.json ./
RUN npm ci --only=production

COPY . .

EXPOSE 3000

CMD ["node", "server.js"]
```

## 🔧 Docker Configuration

### Daemon Configuration
`/etc/docker/daemon.json`:
```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  },
  "storage-driver": "overlay2"
}
```

### Limit Resources
```bash
docker run -d \
  --memory="512m" \
  --cpus="1.0" \
  nginx
```

## 🔒 Security

1. **Don't run as root** in containers
2. **Use official images**
3. **Scan for vulnerabilities**: `docker scan myimage`
4. **Limit capabilities**: `--cap-drop=ALL`
5. **Read-only filesystem**: `--read-only`

## 💡 Tips

### Clean Up
```bash
# Remove unused containers, networks, images
docker system prune -a

# Remove unused volumes
docker volume prune
```

### Networking
```bash
# Create network
docker network create mynetwork

# Run with network
docker run -d --network mynetwork nginx
```

---

← [[Databases]] | [[Monitoring]] →
