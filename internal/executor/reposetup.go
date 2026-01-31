package executor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProjectType represents the detected project type
type ProjectType string

const (
	ProjectNodeJS  ProjectType = "nodejs"
	ProjectPython  ProjectType = "python"
	ProjectGo      ProjectType = "go"
	ProjectPHP     ProjectType = "php"
	ProjectRuby    ProjectType = "ruby"
	ProjectRust    ProjectType = "rust"
	ProjectDocker  ProjectType = "docker"
	ProjectNextJS  ProjectType = "nextjs"
	ProjectLaravel ProjectType = "laravel"
	ProjectDjango  ProjectType = "django"
	ProjectFlask   ProjectType = "flask"
	ProjectUnknown ProjectType = "unknown"
)

// ProjectInfo contains detected project information
type ProjectInfo struct {
	Type             ProjectType
	Framework        string
	RuntimeVersion   string
	PackageManager   string
	Dependencies     []string
	DevDeps          []string
	Scripts          map[string]string
	HasDocker        bool
	HasDockerCompose bool
}

// RepoSetup handles repository cloning and setup
type RepoSetup struct {
	URL       string
	LocalPath string
	Project   *ProjectInfo
}

// NewRepoSetup creates a new repository setup handler
func NewRepoSetup(url, localPath string) *RepoSetup {
	return &RepoSetup{
		URL:       url,
		LocalPath: localPath,
	}
}

// ParseRepoURL extracts repo name from URL
func ParseRepoURL(url string) (owner, repo string, err error) {
	// Handle various URL formats
	url = strings.TrimSuffix(url, ".git")
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "git@")

	// github.com/owner/repo or github.com:owner/repo
	url = strings.Replace(url, ":", "/", 1)
	parts := strings.Split(url, "/")

	if len(parts) < 3 {
		return "", "", fmt.Errorf("invalid repository URL")
	}

	return parts[1], parts[2], nil
}

// DetectProject analyzes the cloned repository
func (r *RepoSetup) DetectProject() (*ProjectInfo, error) {
	info := &ProjectInfo{
		Type:    ProjectUnknown,
		Scripts: make(map[string]string),
	}

	// Check for various project files
	checks := []struct {
		file     string
		projType ProjectType
		pm       string
	}{
		{"package.json", ProjectNodeJS, "npm"},
		{"pnpm-lock.yaml", ProjectNodeJS, "pnpm"},
		{"yarn.lock", ProjectNodeJS, "yarn"},
		{"bun.lockb", ProjectNodeJS, "bun"},
		{"requirements.txt", ProjectPython, "pip"},
		{"pyproject.toml", ProjectPython, "pip"},
		{"Pipfile", ProjectPython, "pipenv"},
		{"go.mod", ProjectGo, "go"},
		{"composer.json", ProjectPHP, "composer"},
		{"Gemfile", ProjectRuby, "bundler"},
		{"Cargo.toml", ProjectRust, "cargo"},
	}

	for _, check := range checks {
		if fileExists(filepath.Join(r.LocalPath, check.file)) {
			info.Type = check.projType
			info.PackageManager = check.pm
			break
		}
	}

	// Detect Docker
	info.HasDocker = fileExists(filepath.Join(r.LocalPath, "Dockerfile"))
	info.HasDockerCompose = fileExists(filepath.Join(r.LocalPath, "docker-compose.yml")) ||
		fileExists(filepath.Join(r.LocalPath, "docker-compose.yaml")) ||
		fileExists(filepath.Join(r.LocalPath, "compose.yml"))

	// Detect specific frameworks
	r.detectFramework(info)

	r.Project = info
	return info, nil
}

// detectFramework identifies specific frameworks
func (r *RepoSetup) detectFramework(info *ProjectInfo) {
	switch info.Type {
	case ProjectNodeJS:
		if fileExists(filepath.Join(r.LocalPath, "next.config.js")) ||
			fileExists(filepath.Join(r.LocalPath, "next.config.mjs")) ||
			fileExists(filepath.Join(r.LocalPath, "next.config.ts")) {
			info.Framework = "Next.js"
		} else if fileExists(filepath.Join(r.LocalPath, "nuxt.config.js")) ||
			fileExists(filepath.Join(r.LocalPath, "nuxt.config.ts")) {
			info.Framework = "Nuxt.js"
		} else if fileExists(filepath.Join(r.LocalPath, "vite.config.js")) ||
			fileExists(filepath.Join(r.LocalPath, "vite.config.ts")) {
			info.Framework = "Vite"
		} else if fileExists(filepath.Join(r.LocalPath, "angular.json")) {
			info.Framework = "Angular"
		} else if fileExists(filepath.Join(r.LocalPath, "svelte.config.js")) {
			info.Framework = "Svelte"
		}

	case ProjectPython:
		if fileExists(filepath.Join(r.LocalPath, "manage.py")) {
			info.Framework = "Django"
		} else if fileExists(filepath.Join(r.LocalPath, "app.py")) ||
			fileExists(filepath.Join(r.LocalPath, "wsgi.py")) {
			info.Framework = "Flask"
		} else if fileExists(filepath.Join(r.LocalPath, "main.py")) {
			// Check for FastAPI
			info.Framework = "Python App"
		}

	case ProjectPHP:
		if fileExists(filepath.Join(r.LocalPath, "artisan")) {
			info.Framework = "Laravel"
		} else if fileExists(filepath.Join(r.LocalPath, "bin", "console")) {
			info.Framework = "Symfony"
		}

	case ProjectGo:
		if fileExists(filepath.Join(r.LocalPath, "main.go")) {
			info.Framework = "Go App"
		}
	}
}

// GenerateSetupScript creates the installation script
func (r *RepoSetup) GenerateSetupScript() string {
	if r.Project == nil {
		return "echo 'No project detected'"
	}

	var script strings.Builder
	script.WriteString("#!/bin/bash\nset -e\n\n")
	script.WriteString(fmt.Sprintf("echo '📦 Setting up %s project...'\n\n", r.Project.Framework))

	// Check and install runtime if needed
	script.WriteString(r.generateRuntimeCheck())

	// Install dependencies based on project type
	script.WriteString(r.generateDependencyInstall())

	// Docker setup if available
	if r.Project.HasDockerCompose {
		script.WriteString("\n# Docker Compose found\n")
		script.WriteString("echo '🐳 Docker Compose detected'\n")
		script.WriteString("read -p 'Start with Docker Compose? [y/N] ' docker_choice\n")
		script.WriteString("if [[ $docker_choice =~ ^[Yy]$ ]]; then\n")
		script.WriteString("    docker compose up -d\n")
		script.WriteString("fi\n")
	}

	script.WriteString("\necho '✅ Setup complete!'\n")

	return script.String()
}

// generateRuntimeCheck creates runtime installation check
func (r *RepoSetup) generateRuntimeCheck() string {
	var script strings.Builder

	switch r.Project.Type {
	case ProjectNodeJS:
		script.WriteString("# Check Node.js\n")
		script.WriteString("if ! command -v node &> /dev/null; then\n")
		script.WriteString("    echo '⚠️ Node.js not found. Installing via NVM...'\n")
		script.WriteString("    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash\n")
		script.WriteString("    export NVM_DIR=\"$HOME/.nvm\"\n")
		script.WriteString("    [ -s \"$NVM_DIR/nvm.sh\" ] && \\. \"$NVM_DIR/nvm.sh\"\n")
		script.WriteString("    nvm install --lts\n")
		script.WriteString("fi\n\n")

	case ProjectPython:
		script.WriteString("# Check Python\n")
		script.WriteString("if ! command -v python3 &> /dev/null; then\n")
		script.WriteString("    echo '⚠️ Python not found. Installing...'\n")
		script.WriteString("    sudo apt-get update && sudo apt-get install -y python3 python3-pip python3-venv\n")
		script.WriteString("fi\n\n")

	case ProjectGo:
		script.WriteString("# Check Go\n")
		script.WriteString("if ! command -v go &> /dev/null; then\n")
		script.WriteString("    echo '⚠️ Go not found. Installing...'\n")
		script.WriteString("    GO_VERSION=$(curl -s https://go.dev/VERSION?m=text | head -1)\n")
		script.WriteString("    wget -q https://go.dev/dl/${GO_VERSION}.linux-amd64.tar.gz\n")
		script.WriteString("    sudo tar -C /usr/local -xzf ${GO_VERSION}.linux-amd64.tar.gz\n")
		script.WriteString("    export PATH=$PATH:/usr/local/go/bin\n")
		script.WriteString("fi\n\n")

	case ProjectPHP:
		script.WriteString("# Check PHP & Composer\n")
		script.WriteString("if ! command -v php &> /dev/null; then\n")
		script.WriteString("    echo '⚠️ PHP not found. Installing...'\n")
		script.WriteString("    sudo apt-get update && sudo apt-get install -y php php-cli php-mbstring php-xml php-curl\n")
		script.WriteString("fi\n")
		script.WriteString("if ! command -v composer &> /dev/null; then\n")
		script.WriteString("    echo '⚠️ Composer not found. Installing...'\n")
		script.WriteString("    curl -sS https://getcomposer.org/installer | php\n")
		script.WriteString("    sudo mv composer.phar /usr/local/bin/composer\n")
		script.WriteString("fi\n\n")

	case ProjectRuby:
		script.WriteString("# Check Ruby\n")
		script.WriteString("if ! command -v ruby &> /dev/null; then\n")
		script.WriteString("    echo '⚠️ Ruby not found. Installing...'\n")
		script.WriteString("    sudo apt-get update && sudo apt-get install -y ruby ruby-dev\n")
		script.WriteString("fi\n\n")

	case ProjectRust:
		script.WriteString("# Check Rust\n")
		script.WriteString("if ! command -v cargo &> /dev/null; then\n")
		script.WriteString("    echo '⚠️ Rust not found. Installing...'\n")
		script.WriteString("    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y\n")
		script.WriteString("    source $HOME/.cargo/env\n")
		script.WriteString("fi\n\n")
	}

	return script.String()
}

// generateDependencyInstall creates dependency installation commands
func (r *RepoSetup) generateDependencyInstall() string {
	var script strings.Builder

	script.WriteString("# Install dependencies\n")
	script.WriteString(fmt.Sprintf("cd %s\n", r.LocalPath))

	switch r.Project.Type {
	case ProjectNodeJS:
		switch r.Project.PackageManager {
		case "pnpm":
			script.WriteString("echo '📦 Installing with pnpm...'\n")
			script.WriteString("npm install -g pnpm 2>/dev/null || true\n")
			script.WriteString("pnpm install\n")
		case "yarn":
			script.WriteString("echo '📦 Installing with yarn...'\n")
			script.WriteString("npm install -g yarn 2>/dev/null || true\n")
			script.WriteString("yarn install\n")
		case "bun":
			script.WriteString("echo '📦 Installing with bun...'\n")
			script.WriteString("curl -fsSL https://bun.sh/install | bash\n")
			script.WriteString("bun install\n")
		default:
			script.WriteString("echo '📦 Installing with npm...'\n")
			script.WriteString("npm install\n")
		}

	case ProjectPython:
		script.WriteString("echo '📦 Setting up Python environment...'\n")
		script.WriteString("python3 -m venv venv\n")
		script.WriteString("source venv/bin/activate\n")
		switch r.Project.PackageManager {
		case "pipenv":
			script.WriteString("pip install pipenv && pipenv install\n")
		default:
			script.WriteString("pip install -r requirements.txt 2>/dev/null || pip install -e .\n")
		}

	case ProjectGo:
		script.WriteString("echo '📦 Installing Go dependencies...'\n")
		script.WriteString("go mod download\n")
		script.WriteString("go build ./...\n")

	case ProjectPHP:
		script.WriteString("echo '📦 Installing Composer dependencies...'\n")
		script.WriteString("composer install\n")
		if r.Project.Framework == "Laravel" {
			script.WriteString("cp .env.example .env 2>/dev/null || true\n")
			script.WriteString("php artisan key:generate 2>/dev/null || true\n")
		}

	case ProjectRuby:
		script.WriteString("echo '📦 Installing Ruby gems...'\n")
		script.WriteString("gem install bundler\n")
		script.WriteString("bundle install\n")

	case ProjectRust:
		script.WriteString("echo '📦 Building Rust project...'\n")
		script.WriteString("cargo build\n")
	}

	return script.String()
}

// GenerateCloneScript creates git clone command
func GenerateCloneScript(url, targetDir string) string {
	return fmt.Sprintf(`#!/bin/bash
set -e
echo "📥 Cloning repository..."
git clone --depth 1 %s %s
cd %s
echo "✅ Repository cloned!"
`, url, targetDir, targetDir)
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
