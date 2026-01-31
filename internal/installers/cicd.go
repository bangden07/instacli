package installers

// ============================================================
// Jenkins Installer
// ============================================================
type JenkinsInstaller struct {
	BaseInstaller
}

func NewJenkinsInstaller() *JenkinsInstaller {
	return &JenkinsInstaller{
		BaseInstaller: BaseInstaller{
			name:        "Jenkins",
			description: "Automation server for CI/CD",
			category:    CategoryCICD,
			icon:        "🔧",
			supportedOS: []OS{OSLinux},
		},
	}
}

func (i *JenkinsInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum}
}
func (i *JenkinsInstaller) Dependencies() []string { return []string{"java"} }
func (i *JenkinsInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *JenkinsInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *JenkinsInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("systemctl status jenkins")
	return err == nil, nil
}

func (i *JenkinsInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "🔧 Installing Jenkins..."

if [ -f /etc/debian_version ]; then
    # Install Java
    sudo apt-get update
    sudo apt-get install -y fontconfig openjdk-17-jre
    
    # Add Jenkins repo
    sudo wget -O /usr/share/keyrings/jenkins-keyring.asc https://pkg.jenkins.io/debian-stable/jenkins.io-2023.key
    echo "deb [signed-by=/usr/share/keyrings/jenkins-keyring.asc] https://pkg.jenkins.io/debian-stable binary/" | sudo tee /etc/apt/sources.list.d/jenkins.list > /dev/null
    
    sudo apt-get update
    sudo apt-get install -y jenkins
    
elif [ -f /etc/redhat-release ]; then
    sudo yum install -y java-17-openjdk
    sudo wget -O /etc/yum.repos.d/jenkins.repo https://pkg.jenkins.io/redhat-stable/jenkins.repo
    sudo rpm --import https://pkg.jenkins.io/redhat-stable/jenkins.io-2023.key
    sudo yum install -y jenkins
fi

sudo systemctl enable jenkins
sudo systemctl start jenkins

echo "✅ Jenkins installed!"
echo "🌐 Access: http://localhost:8080"
echo "🔑 Initial password: sudo cat /var/lib/jenkins/secrets/initialAdminPassword"`
}

func (i *JenkinsInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
sudo systemctl stop jenkins
sudo apt-get remove -y jenkins || sudo yum remove -y jenkins`
}

// ============================================================
// GitLab Runner Installer
// ============================================================
type GitLabRunnerInstaller struct {
	BaseInstaller
}

func NewGitLabRunnerInstaller() *GitLabRunnerInstaller {
	return &GitLabRunnerInstaller{
		BaseInstaller: BaseInstaller{
			name:        "GitLab Runner",
			description: "CI/CD runner for GitLab",
			category:    CategoryCICD,
			icon:        "🦊",
			supportedOS: []OS{OSLinux},
		},
	}
}

func (i *GitLabRunnerInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum}
}
func (i *GitLabRunnerInstaller) Dependencies() []string { return []string{} }
func (i *GitLabRunnerInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *GitLabRunnerInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *GitLabRunnerInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("gitlab-runner --version")
	return err == nil, nil
}

func (i *GitLabRunnerInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "🦊 Installing GitLab Runner..."

# Add repository
curl -L "https://packages.gitlab.com/install/repositories/runner/gitlab-runner/script.deb.sh" | sudo bash

if [ -f /etc/debian_version ]; then
    sudo apt-get install -y gitlab-runner
elif [ -f /etc/redhat-release ]; then
    curl -L "https://packages.gitlab.com/install/repositories/runner/gitlab-runner/script.rpm.sh" | sudo bash
    sudo yum install -y gitlab-runner
fi

echo "✅ GitLab Runner installed!"
echo "🔧 Register: sudo gitlab-runner register"
gitlab-runner --version`
}

func (i *GitLabRunnerInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
sudo gitlab-runner stop
sudo apt-get remove -y gitlab-runner || sudo yum remove -y gitlab-runner`
}

// ============================================================
// GitHub Actions Runner Installer
// ============================================================
type GitHubActionsRunnerInstaller struct {
	BaseInstaller
}

func NewGitHubActionsRunnerInstaller() *GitHubActionsRunnerInstaller {
	return &GitHubActionsRunnerInstaller{
		BaseInstaller: BaseInstaller{
			name:        "GitHub Actions Runner",
			description: "Self-hosted runner for GitHub Actions",
			category:    CategoryCICD,
			icon:        "🐙",
			supportedOS: []OS{OSLinux},
		},
	}
}

func (i *GitHubActionsRunnerInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum}
}
func (i *GitHubActionsRunnerInstaller) Dependencies() []string { return []string{} }
func (i *GitHubActionsRunnerInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *GitHubActionsRunnerInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *GitHubActionsRunnerInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("ls /opt/actions-runner")
	return err == nil, nil
}

func (i *GitHubActionsRunnerInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "🐙 Installing GitHub Actions Runner..."

# Create directory
sudo mkdir -p /opt/actions-runner
cd /opt/actions-runner

# Get latest version
RUNNER_VERSION=$(curl -s https://api.github.com/repos/actions/runner/releases/latest | grep tag_name | cut -d '"' -f 4 | tr -d 'v')

# Download runner
curl -o actions-runner-linux-x64.tar.gz -L https://github.com/actions/runner/releases/download/v${RUNNER_VERSION}/actions-runner-linux-x64-${RUNNER_VERSION}.tar.gz

# Extract
tar xzf actions-runner-linux-x64.tar.gz

echo "✅ GitHub Actions Runner downloaded!"
echo ""
echo "📝 To configure, run:"
echo "   cd /opt/actions-runner"
echo "   ./config.sh --url https://github.com/YOUR_ORG/YOUR_REPO --token YOUR_TOKEN"
echo "   sudo ./svc.sh install"
echo "   sudo ./svc.sh start"`
}

func (i *GitHubActionsRunnerInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
cd /opt/actions-runner
sudo ./svc.sh stop 2>/dev/null || true
sudo ./svc.sh uninstall 2>/dev/null || true
sudo rm -rf /opt/actions-runner`
}
