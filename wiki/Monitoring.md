# Monitoring

System monitoring, metrics collection, and visualization tools.

## 📊 Available Installers

### Prometheus

**Metrics collection and alerting system**

| Property | Value |
| -------- | ----- |
| Default Port | 9090 |
| Supported OS | Linux, macOS |
| Config File | `/etc/prometheus/prometheus.yml` |

**Installation:**
```bash
# Download
wget https://github.com/prometheus/prometheus/releases/download/v2.48.0/prometheus-2.48.0.linux-amd64.tar.gz

# Extract
tar xvf prometheus-*.tar.gz
cd prometheus-*

# Create directories
sudo mkdir -p /etc/prometheus /var/lib/prometheus

# Copy files
sudo cp prometheus promtool /usr/local/bin/
sudo cp prometheus.yml /etc/prometheus/
```

**Configuration:**
```yaml
# /etc/prometheus/prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'node'
    static_configs:
      - targets: ['localhost:9100']
```

**Systemd Service:**
```ini
[Unit]
Description=Prometheus
After=network.target

[Service]
User=prometheus
ExecStart=/usr/local/bin/prometheus \
  --config.file=/etc/prometheus/prometheus.yml \
  --storage.tsdb.path=/var/lib/prometheus

[Install]
WantedBy=multi-user.target
```

---

### Grafana

**Analytics and visualization platform**

| Property | Value |
| -------- | ----- |
| Default Port | 3000 |
| Supported OS | Linux, macOS |
| Default Login | admin / admin |

**Installation (Ubuntu):**
```bash
# Add repository
sudo apt-get install -y apt-transport-https software-properties-common
wget -q -O - https://packages.grafana.com/gpg.key | sudo apt-key add -
echo "deb https://packages.grafana.com/oss/deb stable main" | sudo tee /etc/apt/sources.list.d/grafana.list

# Install
sudo apt-get update
sudo apt-get install -y grafana

# Start
sudo systemctl enable grafana-server
sudo systemctl start grafana-server
```

**Access:**
- URL: `http://localhost:3000`
- Username: `admin`
- Password: `admin` (change on first login)

**Add Prometheus Data Source:**
1. Go to Configuration → Data Sources
2. Add Prometheus
3. URL: `http://localhost:9090`
4. Save & Test

---

### Netdata

**Real-time performance monitoring**

| Property | Value |
| -------- | ----- |
| Default Port | 19999 |
| Supported OS | Linux, macOS |
| Features | Real-time, low overhead |

**Installation:**
```bash
# One-line install
bash <(curl -Ss https://my-netdata.io/kickstart.sh)

# Or via package manager
sudo apt install netdata
```

**Access:**
- URL: `http://localhost:19999`

**Features:**
- CPU, Memory, Disk, Network monitoring
- Per-application metrics
- Container monitoring
- 1-second granularity
- Web dashboard included

**Configuration:**
```bash
# Edit config
sudo nano /etc/netdata/netdata.conf

# Restart
sudo systemctl restart netdata
```

---

## 📈 Monitoring Stack

A common setup is the **Prometheus + Grafana** stack:

```
┌─────────────┐     ┌─────────────┐
│   Targets   │────▶│ Prometheus  │
│ (exporters) │     │  (metrics)  │
└─────────────┘     └──────┬──────┘
                           │
                           ▼
                    ┌─────────────┐
                    │   Grafana   │
                    │ (dashboard) │
                    └─────────────┘
```

## 🔧 Common Exporters

| Exporter | Port | Purpose |
| -------- | ---- | ------- |
| Node Exporter | 9100 | Linux metrics |
| MySQL Exporter | 9104 | MySQL metrics |
| PostgreSQL Exporter | 9187 | PostgreSQL metrics |
| Redis Exporter | 9121 | Redis metrics |
| Nginx Exporter | 9113 | Nginx metrics |

**Install Node Exporter:**
```bash
wget https://github.com/prometheus/node_exporter/releases/download/v1.7.0/node_exporter-1.7.0.linux-amd64.tar.gz
tar xvf node_exporter-*.tar.gz
sudo cp node_exporter-*/node_exporter /usr/local/bin/
```

## 🚨 Alerting

Prometheus alerting with Alertmanager:
```yaml
# alert.rules.yml
groups:
  - name: example
    rules:
      - alert: HighCPU
        expr: 100 - (avg by(instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High CPU usage"
```

---

← [[Containers]] | [[Infrastructure]] →
