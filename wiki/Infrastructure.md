# Infrastructure

Reverse proxies, load balancers, and storage solutions.

## 🔀 Available Installers

### Nginx Proxy Manager

**Easy-to-use reverse proxy with SSL management**

| Property | Value |
| -------- | ----- |
| Web UI Port | 81 |
| HTTP Port | 80 |
| HTTPS Port | 443 |
| Default Login | admin@example.com / changeme |

**Installation (Docker):**
```yaml
# docker-compose.yml
version: '3.8'

services:
  nginx-proxy-manager:
    image: 'jc21/nginx-proxy-manager:latest'
    restart: unless-stopped
    ports:
      - '80:80'
      - '81:81'
      - '443:443'
    volumes:
      - ./data:/data
      - ./letsencrypt:/etc/letsencrypt
```

```bash
docker compose up -d
```

**Features:**
- Beautiful web UI
- Free SSL with Let's Encrypt
- Access lists and authentication
- Redirect and stream hosts
- Custom Nginx configurations

**Setup Steps:**
1. Access `http://your-ip:81`
2. Login: `admin@example.com` / `changeme`
3. Change password
4. Add Proxy Hosts

---

### Traefik

**Cloud-native reverse proxy and load balancer**

| Property | Value |
| -------- | ----- |
| HTTP Port | 80 |
| HTTPS Port | 443 |
| Dashboard Port | 8080 |

**Installation (Docker):**
```yaml
# docker-compose.yml
version: '3.8'

services:
  traefik:
    image: traefik:v2.10
    command:
      - "--api.dashboard=true"
      - "--providers.docker=true"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./traefik.yml:/etc/traefik/traefik.yml

  whoami:
    image: traefik/whoami
    labels:
      - "traefik.http.routers.whoami.rule=Host(`whoami.localhost`)"
```

**Features:**
- Automatic service discovery
- Docker, Kubernetes integration
- Automatic SSL with Let's Encrypt
- Middleware (auth, rate limiting)
- Metrics and tracing

**Configuration File:**
```yaml
# traefik.yml
api:
  dashboard: true
  insecure: true

entryPoints:
  web:
    address: ":80"
  websecure:
    address: ":443"

providers:
  docker:
    exposedByDefault: false

certificatesResolvers:
  letsencrypt:
    acme:
      email: your@email.com
      storage: /letsencrypt/acme.json
      httpChallenge:
        entryPoint: web
```

---

### MinIO

**S3-compatible object storage**

| Property | Value |
| -------- | ----- |
| API Port | 9000 |
| Console Port | 9001 |
| Supported OS | Linux, macOS |

**Installation (Docker):**
```yaml
# docker-compose.yml
version: '3.8'

services:
  minio:
    image: minio/minio
    command: server /data --console-address ":9001"
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - minio_data:/data

volumes:
  minio_data:
```

**Binary Installation:**
```bash
wget https://dl.min.io/server/minio/release/linux-amd64/minio
chmod +x minio
sudo mv minio /usr/local/bin/

# Run
export MINIO_ROOT_USER=admin
export MINIO_ROOT_PASSWORD=password
minio server /data --console-address ":9001"
```

**Access:**
- Console: `http://localhost:9001`
- API: `http://localhost:9000`

**Use with AWS CLI:**
```bash
# Configure
aws configure --profile minio
# Access Key: minioadmin
# Secret Key: minioadmin
# Region: us-east-1

# Create bucket
aws --endpoint-url http://localhost:9000 --profile minio s3 mb s3://mybucket

# Upload file
aws --endpoint-url http://localhost:9000 --profile minio s3 cp file.txt s3://mybucket/
```

---

## 🔄 Comparison

| Tool | Best For | UI |
| ---- | -------- | -- |
| Nginx Proxy Manager | Simple setups, beginners | ✅ Web UI |
| Traefik | Docker/K8s, microservices | ✅ Dashboard |
| MinIO | S3-compatible storage | ✅ Console |

## 🔒 SSL Certificates

All these tools support automatic SSL via Let's Encrypt:

- **NPM**: Built-in, click to enable
- **Traefik**: Configure ACME resolver
- **MinIO**: Use reverse proxy for SSL

---

← [[Monitoring]] | [[CI CD]] →
