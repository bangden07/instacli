package installers

// ============================================================
// Nginx Proxy Manager Installer
// ============================================================
type NginxProxyManagerInstaller struct {
	BaseInstaller
}

func NewNginxProxyManagerInstaller() *NginxProxyManagerInstaller {
	return &NginxProxyManagerInstaller{
		BaseInstaller: BaseInstaller{
			name:        "Nginx Proxy Manager",
			description: "Reverse proxy with SSL & visual UI",
			category:    CategoryInfrastructure,
			icon:        "🔀",
			supportedOS: []OS{OSLinux},
		},
	}
}

func (i *NginxProxyManagerInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum}
}
func (i *NginxProxyManagerInstaller) Dependencies() []string { return []string{"docker"} }
func (i *NginxProxyManagerInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *NginxProxyManagerInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *NginxProxyManagerInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("docker ps | grep nginx-proxy-manager")
	return err == nil, nil
}

func (i *NginxProxyManagerInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "🔀 Installing Nginx Proxy Manager..."

# Create directory
sudo mkdir -p /opt/nginx-proxy-manager
cd /opt/nginx-proxy-manager

# Create docker-compose.yml
cat <<EOF > docker-compose.yml
version: '3.8'
services:
  app:
    image: 'jc21/nginx-proxy-manager:latest'
    restart: unless-stopped
    ports:
      - '80:80'
      - '81:81'
      - '443:443'
    volumes:
      - ./data:/data
      - ./letsencrypt:/etc/letsencrypt
EOF

# Start container
docker compose up -d

echo "✅ Nginx Proxy Manager installed!"
echo "🌐 Admin UI: http://localhost:81"
echo "📝 Default login: admin@example.com / changeme"`
}

func (i *NginxProxyManagerInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
cd /opt/nginx-proxy-manager && docker compose down -v
sudo rm -rf /opt/nginx-proxy-manager`
}

// ============================================================
// Traefik Installer
// ============================================================
type TraefikInstaller struct {
	BaseInstaller
}

func NewTraefikInstaller() *TraefikInstaller {
	return &TraefikInstaller{
		BaseInstaller: BaseInstaller{
			name:        "Traefik",
			description: "Cloud-native reverse proxy & load balancer",
			category:    CategoryInfrastructure,
			icon:        "🚦",
			supportedOS: []OS{OSLinux},
		},
	}
}

func (i *TraefikInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum}
}
func (i *TraefikInstaller) Dependencies() []string { return []string{"docker"} }
func (i *TraefikInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *TraefikInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *TraefikInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("docker ps | grep traefik")
	return err == nil, nil
}

func (i *TraefikInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "🚦 Installing Traefik..."

sudo mkdir -p /opt/traefik
cd /opt/traefik

# Create traefik config
cat <<EOF > traefik.yml
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
    endpoint: "unix:///var/run/docker.sock"
    exposedByDefault: false
EOF

# Create docker-compose
cat <<EOF > docker-compose.yml
version: '3.8'
services:
  traefik:
    image: traefik:latest
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./traefik.yml:/traefik.yml:ro
EOF

docker compose up -d

echo "✅ Traefik installed!"
echo "🌐 Dashboard: http://localhost:8080"`
}

func (i *TraefikInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
cd /opt/traefik && docker compose down -v
sudo rm -rf /opt/traefik`
}

// ============================================================
// WireGuard VPN Installer
// ============================================================
type WireGuardInstaller struct {
	BaseInstaller
}

func NewWireGuardInstaller() *WireGuardInstaller {
	return &WireGuardInstaller{
		BaseInstaller: BaseInstaller{
			name:        "WireGuard VPN",
			description: "Fast, modern VPN tunnel",
			category:    CategoryVPN,
			icon:        "🔐",
			supportedOS: []OS{OSLinux},
		},
	}
}

func (i *WireGuardInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum}
}
func (i *WireGuardInstaller) Dependencies() []string { return []string{} }
func (i *WireGuardInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *WireGuardInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *WireGuardInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("wg --version")
	return err == nil, nil
}

func (i *WireGuardInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "🔐 Installing WireGuard VPN..."

if [ -f /etc/debian_version ]; then
    sudo apt-get update
    sudo apt-get install -y wireguard wireguard-tools
elif [ -f /etc/redhat-release ]; then
    sudo yum install -y epel-release elrepo-release
    sudo yum install -y kmod-wireguard wireguard-tools
fi

# Generate server keys
sudo mkdir -p /etc/wireguard
cd /etc/wireguard
wg genkey | sudo tee privatekey | wg pubkey | sudo tee publickey

# Create basic config
PRIVATE_KEY=$(sudo cat privatekey)
cat <<EOF | sudo tee /etc/wireguard/wg0.conf
[Interface]
PrivateKey = $PRIVATE_KEY
Address = 10.0.0.1/24
ListenPort = 51820
SaveConfig = true
PostUp = iptables -A FORWARD -i wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE
EOF

sudo chmod 600 /etc/wireguard/wg0.conf
sudo systemctl enable wg-quick@wg0

echo ""
echo "✅ WireGuard installed!"
echo "📝 Config: /etc/wireguard/wg0.conf"
echo "🔧 Start: sudo wg-quick up wg0"
wg --version`
}

func (i *WireGuardInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
sudo wg-quick down wg0 2>/dev/null || true
sudo apt-get remove -y wireguard wireguard-tools || sudo yum remove -y wireguard-tools
sudo rm -rf /etc/wireguard`
}

// ============================================================
// MinIO S3 Storage Installer
// ============================================================
type MinIOInstaller struct {
	BaseInstaller
}

func NewMinIOInstaller() *MinIOInstaller {
	return &MinIOInstaller{
		BaseInstaller: BaseInstaller{
			name:        "MinIO",
			description: "S3-compatible object storage",
			category:    CategoryInfrastructure,
			icon:        "📦",
			supportedOS: []OS{OSLinux},
		},
	}
}

func (i *MinIOInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum}
}
func (i *MinIOInstaller) Dependencies() []string { return []string{"wget"} }
func (i *MinIOInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *MinIOInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *MinIOInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("minio --version")
	return err == nil, nil
}

func (i *MinIOInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "📦 Installing MinIO..."

# Download MinIO
wget -q https://dl.min.io/server/minio/release/linux-amd64/minio -O /tmp/minio
sudo mv /tmp/minio /usr/local/bin/
sudo chmod +x /usr/local/bin/minio

# Create minio user and directories
sudo useradd -r minio-user 2>/dev/null || true
sudo mkdir -p /data/minio
sudo chown minio-user:minio-user /data/minio

# Create environment file
cat <<EOF | sudo tee /etc/default/minio
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin123
MINIO_VOLUMES="/data/minio"
MINIO_OPTS="--console-address :9001"
EOF

# Create systemd service
cat <<EOF | sudo tee /etc/systemd/system/minio.service
[Unit]
Description=MinIO
After=network.target

[Service]
User=minio-user
Group=minio-user
EnvironmentFile=/etc/default/minio
ExecStart=/usr/local/bin/minio server \$MINIO_OPTS \$MINIO_VOLUMES
Restart=always
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable minio
sudo systemctl start minio

echo "✅ MinIO installed!"
echo "🌐 Console: http://localhost:9001"
echo "📝 User: minioadmin / minioadmin123"`
}

func (i *MinIOInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
sudo systemctl stop minio
sudo systemctl disable minio
sudo rm -f /etc/systemd/system/minio.service
sudo rm -f /usr/local/bin/minio
sudo rm -rf /data/minio`
}
