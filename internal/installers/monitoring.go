package installers

// ============================================================
// Prometheus Installer
// ============================================================
type PrometheusInstaller struct {
	BaseInstaller
}

func NewPrometheusInstaller() *PrometheusInstaller {
	return &PrometheusInstaller{
		BaseInstaller: BaseInstaller{
			name:        "Prometheus",
			description: "Metrics monitoring & alerting system",
			category:    CategoryMonitoring,
			icon:        "📊",
			supportedOS: []OS{OSLinux, OSMacOS},
		},
	}
}

func (i *PrometheusInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum, PMBrew}
}
func (i *PrometheusInstaller) Dependencies() []string { return []string{"wget"} }
func (i *PrometheusInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *PrometheusInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *PrometheusInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("prometheus --version")
	return err == nil, nil
}

func (i *PrometheusInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "📊 Installing Prometheus..."

# Create prometheus user
sudo useradd --no-create-home --shell /bin/false prometheus 2>/dev/null || true

# Create directories
sudo mkdir -p /etc/prometheus /var/lib/prometheus
sudo chown prometheus:prometheus /etc/prometheus /var/lib/prometheus

# Download latest Prometheus
PROM_VERSION=$(curl -s https://api.github.com/repos/prometheus/prometheus/releases/latest | grep tag_name | cut -d '"' -f 4 | tr -d 'v')
cd /tmp
wget -q https://github.com/prometheus/prometheus/releases/download/v${PROM_VERSION}/prometheus-${PROM_VERSION}.linux-amd64.tar.gz
tar xzf prometheus-${PROM_VERSION}.linux-amd64.tar.gz
cd prometheus-${PROM_VERSION}.linux-amd64

# Install binaries
sudo cp prometheus promtool /usr/local/bin/
sudo chown prometheus:prometheus /usr/local/bin/prometheus /usr/local/bin/promtool

# Install config files
sudo cp -r consoles console_libraries prometheus.yml /etc/prometheus/
sudo chown -R prometheus:prometheus /etc/prometheus

# Create systemd service
cat <<EOF | sudo tee /etc/systemd/system/prometheus.service
[Unit]
Description=Prometheus Monitoring
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \
    --config.file /etc/prometheus/prometheus.yml \
    --storage.tsdb.path /var/lib/prometheus/ \
    --web.console.templates=/etc/prometheus/consoles \
    --web.console.libraries=/etc/prometheus/console_libraries

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable prometheus
sudo systemctl start prometheus

echo "✅ Prometheus installed!"
echo "🌐 Access: http://localhost:9090"
prometheus --version`
}

func (i *PrometheusInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
sudo systemctl stop prometheus
sudo systemctl disable prometheus
sudo rm -f /etc/systemd/system/prometheus.service
sudo rm -rf /usr/local/bin/prometheus /usr/local/bin/promtool
sudo rm -rf /etc/prometheus /var/lib/prometheus
sudo userdel prometheus 2>/dev/null || true`
}

// ============================================================
// Grafana Installer
// ============================================================
type GrafanaInstaller struct {
	BaseInstaller
}

func NewGrafanaInstaller() *GrafanaInstaller {
	return &GrafanaInstaller{
		BaseInstaller: BaseInstaller{
			name:        "Grafana",
			description: "Analytics & visualization platform",
			category:    CategoryMonitoring,
			icon:        "📈",
			supportedOS: []OS{OSLinux, OSMacOS},
		},
	}
}

func (i *GrafanaInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum, PMBrew}
}
func (i *GrafanaInstaller) Dependencies() []string { return []string{} }
func (i *GrafanaInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *GrafanaInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *GrafanaInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("grafana-server -v")
	return err == nil, nil
}

func (i *GrafanaInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "📈 Installing Grafana..."

if [ -f /etc/debian_version ]; then
    # Install dependencies
    sudo apt-get update
    sudo apt-get install -y apt-transport-https software-properties-common wget
    
    # Add Grafana GPG key
    sudo mkdir -p /etc/apt/keyrings/
    wget -q -O - https://apt.grafana.com/gpg.key | gpg --dearmor | sudo tee /etc/apt/keyrings/grafana.gpg > /dev/null
    
    # Add repository
    echo "deb [signed-by=/etc/apt/keyrings/grafana.gpg] https://apt.grafana.com stable main" | sudo tee /etc/apt/sources.list.d/grafana.list
    
    # Install Grafana
    sudo apt-get update
    sudo apt-get install -y grafana
    
elif [ -f /etc/redhat-release ]; then
    cat <<EOF | sudo tee /etc/yum.repos.d/grafana.repo
[grafana]
name=grafana
baseurl=https://rpm.grafana.com
repo_gpgcheck=1
enabled=1
gpgcheck=1
gpgkey=https://rpm.grafana.com/gpg.key
sslverify=1
sslcacert=/etc/pki/tls/certs/ca-bundle.crt
EOF
    sudo yum install -y grafana
    
elif command -v brew &> /dev/null; then
    brew install grafana
fi

sudo systemctl daemon-reload
sudo systemctl enable grafana-server
sudo systemctl start grafana-server

echo "✅ Grafana installed!"
echo "🌐 Access: http://localhost:3000"
echo "📝 Default login: admin / admin"`
}

func (i *GrafanaInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
sudo systemctl stop grafana-server
sudo apt-get remove -y grafana || sudo yum remove -y grafana || brew uninstall grafana`
}

// ============================================================
// Netdata Installer
// ============================================================
type NetdataInstaller struct {
	BaseInstaller
}

func NewNetdataInstaller() *NetdataInstaller {
	return &NetdataInstaller{
		BaseInstaller: BaseInstaller{
			name:        "Netdata",
			description: "Real-time performance monitoring",
			category:    CategoryMonitoring,
			icon:        "⚡",
			supportedOS: []OS{OSLinux, OSMacOS},
		},
	}
}

func (i *NetdataInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum, PMBrew}
}
func (i *NetdataInstaller) Dependencies() []string { return []string{"curl"} }
func (i *NetdataInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *NetdataInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *NetdataInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("netdata -v")
	return err == nil, nil
}

func (i *NetdataInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "⚡ Installing Netdata..."

# One-liner install from official script
curl https://get.netdata.cloud/kickstart.sh > /tmp/netdata-kickstart.sh && sh /tmp/netdata-kickstart.sh --stable-channel --disable-telemetry

echo "✅ Netdata installed!"
echo "🌐 Access: http://localhost:19999"`
}

func (i *NetdataInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
sudo /usr/libexec/netdata/netdata-uninstaller.sh --yes --env /etc/netdata/.environment`
}
