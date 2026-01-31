package executor

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// ToolStatus represents the installation status of a tool
type ToolStatus int

const (
	StatusNotInstalled ToolStatus = iota
	StatusInstalled
	StatusPartiallyInstalled
	StatusOutdated
	StatusUnknown
)

func (s ToolStatus) String() string {
	switch s {
	case StatusNotInstalled:
		return "Not Installed"
	case StatusInstalled:
		return "Installed"
	case StatusPartiallyInstalled:
		return "Partially Installed"
	case StatusOutdated:
		return "Outdated"
	default:
		return "Unknown"
	}
}

// ToolCheck represents the result of checking a tool
type ToolCheck struct {
	Name           string
	Status         ToolStatus
	CurrentVersion string
	LatestVersion  string
	Path           string
	Details        string
}

// SystemCheck holds all system checks
type SystemCheck struct {
	OS            string
	Arch          string
	PackageManager string
	IsRoot        bool
	Tools         map[string]*ToolCheck
}

// NewSystemCheck creates a new system checker
func NewSystemCheck() *SystemCheck {
	return &SystemCheck{
		OS:    runtime.GOOS,
		Arch:  runtime.GOARCH,
		Tools: make(map[string]*ToolCheck),
	}
}

// DetectAll runs all detection checks
func (sc *SystemCheck) DetectAll() {
	sc.detectPackageManager()
	sc.detectIsRoot()
	sc.detectCommonTools()
}

// detectPackageManager finds the system package manager
func (sc *SystemCheck) detectPackageManager() {
	managers := []struct {
		name string
		cmd  string
	}{
		{"apt", "apt"},
		{"apt-get", "apt-get"},
		{"yum", "yum"},
		{"dnf", "dnf"},
		{"pacman", "pacman"},
		{"brew", "brew"},
		{"choco", "choco"},
		{"winget", "winget"},
	}

	for _, m := range managers {
		if _, err := exec.LookPath(m.cmd); err == nil {
			sc.PackageManager = m.name
			return
		}
	}
	sc.PackageManager = "unknown"
}

// detectIsRoot checks if running as root/admin
func (sc *SystemCheck) detectIsRoot() {
	if sc.OS == "windows" {
		// Check admin on Windows
		cmd := exec.Command("net", "session")
		sc.IsRoot = cmd.Run() == nil
	} else {
		// Check root on Unix
		cmd := exec.Command("id", "-u")
		out, err := cmd.Output()
		sc.IsRoot = err == nil && strings.TrimSpace(string(out)) == "0"
	}
}

// detectCommonTools checks for commonly used development tools
func (sc *SystemCheck) detectCommonTools() {
	tools := []struct {
		name       string
		commands   []string
		versionArg string
	}{
		{"docker", []string{"docker"}, "--version"},
		{"docker-compose", []string{"docker-compose", "docker compose"}, "version"},
		{"node", []string{"node"}, "--version"},
		{"npm", []string{"npm"}, "--version"},
		{"nvm", []string{"nvm"}, "--version"},
		{"go", []string{"go"}, "version"},
		{"php", []string{"php"}, "--version"},
		{"composer", []string{"composer"}, "--version"},
		{"mysql", []string{"mysql"}, "--version"},
		{"nginx", []string{"nginx"}, "-v"},
		{"apache2", []string{"apache2", "httpd"}, "-v"},
		{"git", []string{"git"}, "--version"},
		{"curl", []string{"curl"}, "--version"},
		{"wget", []string{"wget"}, "--version"},
		{"python", []string{"python3", "python"}, "--version"},
		{"pip", []string{"pip3", "pip"}, "--version"},
		{"redis", []string{"redis-server"}, "--version"},
		{"postgresql", []string{"psql"}, "--version"},
		{"mongo", []string{"mongo", "mongod"}, "--version"},
		{"pm2", []string{"pm2"}, "--version"},
		{"certbot", []string{"certbot"}, "--version"},
		{"ufw", []string{"ufw"}, "version"},
	}

	for _, tool := range tools {
		check := sc.checkTool(tool.name, tool.commands, tool.versionArg)
		sc.Tools[tool.name] = check
	}
}

// checkTool checks if a specific tool is installed
func (sc *SystemCheck) checkTool(name string, commands []string, versionArg string) *ToolCheck {
	check := &ToolCheck{
		Name:   name,
		Status: StatusNotInstalled,
	}

	for _, cmdName := range commands {
		path, err := exec.LookPath(cmdName)
		if err != nil {
			continue
		}

		check.Path = path
		check.Status = StatusInstalled

		// Get version
		args := strings.Fields(versionArg)
		cmd := exec.Command(cmdName, args...)
		out, err := cmd.CombinedOutput()
		if err == nil {
			version := strings.TrimSpace(string(out))
			// Extract first line only
			if idx := strings.Index(version, "\n"); idx > 0 {
				version = version[:idx]
			}
			check.CurrentVersion = version
		}
		break
	}

	return check
}

// CheckTool checks a specific tool by name
func (sc *SystemCheck) CheckTool(name string) *ToolCheck {
	if check, ok := sc.Tools[name]; ok {
		return check
	}
	return &ToolCheck{Name: name, Status: StatusUnknown}
}

// IsInstalled returns true if the tool is installed
func (sc *SystemCheck) IsInstalled(name string) bool {
	if check, ok := sc.Tools[name]; ok {
		return check.Status == StatusInstalled
	}
	return false
}

// GetInstalledTools returns list of installed tools
func (sc *SystemCheck) GetInstalledTools() []*ToolCheck {
	var installed []*ToolCheck
	for _, check := range sc.Tools {
		if check.Status == StatusInstalled {
			installed = append(installed, check)
		}
	}
	return installed
}

// GetMissingTools returns list of tools that are not installed
func (sc *SystemCheck) GetMissingTools() []*ToolCheck {
	var missing []*ToolCheck
	for _, check := range sc.Tools {
		if check.Status == StatusNotInstalled {
			missing = append(missing, check)
		}
	}
	return missing
}

// Report generates a formatted report of system status
func (sc *SystemCheck) Report() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("System: %s/%s\n", sc.OS, sc.Arch))
	b.WriteString(fmt.Sprintf("Package Manager: %s\n", sc.PackageManager))
	b.WriteString(fmt.Sprintf("Root/Admin: %v\n", sc.IsRoot))
	b.WriteString("\n")

	installed := sc.GetInstalledTools()
	missing := sc.GetMissingTools()

	b.WriteString(fmt.Sprintf("✅ Installed (%d):\n", len(installed)))
	for _, t := range installed {
		version := t.CurrentVersion
		if len(version) > 50 {
			version = version[:50] + "..."
		}
		b.WriteString(fmt.Sprintf("   • %s: %s\n", t.Name, version))
	}

	b.WriteString(fmt.Sprintf("\n❌ Not Installed (%d):\n", len(missing)))
	for _, t := range missing {
		b.WriteString(fmt.Sprintf("   • %s\n", t.Name))
	}

	return b.String()
}

// CheckDependencies verifies if required dependencies for a tool are met
func (sc *SystemCheck) CheckDependencies(deps []string) (bool, []string) {
	var missing []string
	for _, dep := range deps {
		if !sc.IsInstalled(dep) {
			missing = append(missing, dep)
		}
	}
	return len(missing) == 0, missing
}

// PreflightCheck runs before installation
type PreflightResult struct {
	CanInstall      bool
	AlreadyInstalled bool
	MissingDeps     []string
	Warnings        []string
	Errors          []string
}

// Preflight runs pre-installation checks for a specific installer
func (sc *SystemCheck) Preflight(toolName string, dependencies []string) *PreflightResult {
	result := &PreflightResult{
		CanInstall: true,
	}

	// Check if already installed
	if sc.IsInstalled(toolName) {
		result.AlreadyInstalled = true
		result.Warnings = append(result.Warnings, 
			fmt.Sprintf("%s is already installed", toolName))
	}

	// Check dependencies
	for _, dep := range dependencies {
		if !sc.IsInstalled(dep) {
			result.MissingDeps = append(result.MissingDeps, dep)
		}
	}

	// Check package manager
	if sc.PackageManager == "unknown" {
		result.Warnings = append(result.Warnings, 
			"No known package manager detected")
	}

	// Check root for certain operations
	if sc.OS != "windows" && !sc.IsRoot {
		result.Warnings = append(result.Warnings, 
			"Not running as root - some installations may require sudo")
	}

	return result
}
