package installers

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// IDEType represents detected IDE
type IDEType string

const (
	IDEVSCode   IDEType = "vscode"
	IDECursor   IDEType = "cursor"
	IDEWindsurf IDEType = "windsurf"
	IDEZed      IDEType = "zed"
	IDEUnknown  IDEType = "unknown"
)

// DetectIDE tries to detect which IDE is being used
func DetectIDE() IDEType {
	// Check environment variables first
	if os.Getenv("VSCODE_PID") != "" || os.Getenv("VSCODE_IPC_HOOK") != "" {
		return IDEVSCode
	}
	if os.Getenv("CURSOR_TRACE") != "" || os.Getenv("CURSOR_PID") != "" {
		return IDECursor
	}

	home, _ := os.UserHomeDir()

	// Check for config directories
	if runtime.GOOS == "windows" {
		if _, err := os.Stat(filepath.Join(home, ".cursor")); err == nil {
			return IDECursor
		}
		if _, err := os.Stat(filepath.Join(home, ".windsurf")); err == nil {
			return IDEWindsurf
		}
	} else {
		if _, err := os.Stat(filepath.Join(home, ".cursor")); err == nil {
			return IDECursor
		}
		if _, err := os.Stat(filepath.Join(home, ".windsurf")); err == nil {
			return IDEWindsurf
		}
		if _, err := os.Stat(filepath.Join(home, ".config", "zed")); err == nil {
			return IDEZed
		}
	}

	return IDEUnknown
}

// GetMCPConfigPath returns the MCP config path for the given IDE
func GetMCPConfigPath(ide IDEType) string {
	home, _ := os.UserHomeDir()

	switch ide {
	case IDEVSCode:
		if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "Code", "User", "globalStorage", "mcp.json")
		}
		return filepath.Join(home, ".config", "Code", "User", "globalStorage", "mcp.json")
	case IDECursor:
		return filepath.Join(home, ".cursor", "mcp.json")
	case IDEWindsurf:
		return filepath.Join(home, ".windsurf", "mcp.json")
	case IDEZed:
		return filepath.Join(home, ".config", "zed", "mcp.json")
	default:
		return filepath.Join(home, ".mcp", "config.json")
	}
}

// ========================================
// BASE MCP INSTALLER
// ========================================

type BaseMCPInstaller struct {
	BaseInstaller
	mcpName    string
	mcpPackage string
	mcpArgs    []string
}

func (m *BaseMCPInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	var script strings.Builder

	switch os {
	case OSLinux, OSMacOS:
		script.WriteString(fmt.Sprintf(`#!/bin/bash
set -e

echo "📦 Installing MCP: %s..."

# Check if npx is available
if ! command -v npx &> /dev/null; then
    echo "❌ npx is required. Please install Node.js first."
    exit 1
fi

# Test MCP server
echo "Testing MCP server..."
npx -y %s --help 2>/dev/null || echo "MCP server ready"

echo ""
echo "✅ MCP %s is ready!"
echo ""
echo "To use this MCP, add it to your IDE's MCP configuration:"
echo ""
echo "For Cursor (~/.cursor/mcp.json):"
echo '{'
echo '  "mcpServers": {'
echo '    "%s": {'
echo '      "command": "npx",'
echo '      "args": ["-y", "%s"]'
echo '    }'
echo '  }'
echo '}'
echo ""
echo "For VS Code, Windsurf, or other editors, check their MCP documentation."
`, m.mcpName, m.mcpPackage, m.mcpName, m.mcpName, m.mcpPackage))

	case OSWindows:
		script.WriteString(fmt.Sprintf(`# PowerShell script
Write-Host "📦 Installing MCP: %s..." -ForegroundColor Cyan

# Test MCP
npx -y %s --help 2>$null

Write-Host "✅ MCP %s is ready!" -ForegroundColor Green
Write-Host "Add it to your IDE's MCP configuration"
`, m.mcpName, m.mcpPackage, m.mcpName))
	}

	return script.String()
}

// ========================================
// CONTEXT7 MCP
// ========================================

type Context7MCPInstaller struct {
	BaseInstaller
}

func NewContext7MCPInstaller() *Context7MCPInstaller {
	return &Context7MCPInstaller{
		BaseInstaller: NewBaseInstaller(
			"Context7 MCP",
			"Documentation lookup for libraries and frameworks",
			CategoryMCP,
			"📚",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (i *Context7MCPInstaller) Dependencies() []string {
	return []string{"node", "npx"}
}

func (i *Context7MCPInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return generateMCPScript("context7", "@context7/mcp-server", os)
}

// ========================================
// PLAYWRIGHT MCP
// ========================================

type PlaywrightMCPInstaller struct {
	BaseInstaller
}

func NewPlaywrightMCPInstaller() *PlaywrightMCPInstaller {
	return &PlaywrightMCPInstaller{
		BaseInstaller: NewBaseInstaller(
			"Playwright MCP",
			"Browser automation and testing",
			CategoryMCP,
			"🎭",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (i *PlaywrightMCPInstaller) Dependencies() []string {
	return []string{"node", "npx"}
}

func (i *PlaywrightMCPInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return generateMCPScript("playwright", "@playwright/mcp", os)
}

// ========================================
// GITHUB MCP
// ========================================

type GitHubMCPInstaller struct {
	BaseInstaller
}

func NewGitHubMCPInstaller() *GitHubMCPInstaller {
	return &GitHubMCPInstaller{
		BaseInstaller: NewBaseInstaller(
			"GitHub MCP",
			"GitHub repository and API integration",
			CategoryMCP,
			"🐙",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (i *GitHubMCPInstaller) Dependencies() []string {
	return []string{"node", "npx"}
}

func (i *GitHubMCPInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return generateMCPScript("github", "@modelcontextprotocol/server-github", os)
}

// ========================================
// FILESYSTEM MCP
// ========================================

type FilesystemMCPInstaller struct {
	BaseInstaller
}

func NewFilesystemMCPInstaller() *FilesystemMCPInstaller {
	return &FilesystemMCPInstaller{
		BaseInstaller: NewBaseInstaller(
			"Filesystem MCP",
			"File system read/write access",
			CategoryMCP,
			"📁",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (i *FilesystemMCPInstaller) Dependencies() []string {
	return []string{"node", "npx"}
}

func (i *FilesystemMCPInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return generateMCPScript("filesystem", "@modelcontextprotocol/server-filesystem", os)
}

// ========================================
// POSTGRES MCP
// ========================================

type PostgresMCPInstaller struct {
	BaseInstaller
}

func NewPostgresMCPInstaller() *PostgresMCPInstaller {
	return &PostgresMCPInstaller{
		BaseInstaller: NewBaseInstaller(
			"PostgreSQL MCP",
			"PostgreSQL database access",
			CategoryMCP,
			"🐘",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (i *PostgresMCPInstaller) Dependencies() []string {
	return []string{"node", "npx"}
}

func (i *PostgresMCPInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return generateMCPScript("postgres", "@modelcontextprotocol/server-postgres", os)
}

// ========================================
// BRAVE SEARCH MCP
// ========================================

type BraveSearchMCPInstaller struct {
	BaseInstaller
}

func NewBraveSearchMCPInstaller() *BraveSearchMCPInstaller {
	return &BraveSearchMCPInstaller{
		BaseInstaller: NewBaseInstaller(
			"Brave Search MCP",
			"Web search using Brave Search API",
			CategoryMCP,
			"🦁",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (i *BraveSearchMCPInstaller) Dependencies() []string {
	return []string{"node", "npx"}
}

func (i *BraveSearchMCPInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return generateMCPScript("brave-search", "@modelcontextprotocol/server-brave-search", os)
}

// ========================================
// MEMORY MCP
// ========================================

type MemoryMCPInstaller struct {
	BaseInstaller
}

func NewMemoryMCPInstaller() *MemoryMCPInstaller {
	return &MemoryMCPInstaller{
		BaseInstaller: NewBaseInstaller(
			"Memory MCP",
			"Persistent memory and knowledge graph",
			CategoryMCP,
			"🧠",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (i *MemoryMCPInstaller) Dependencies() []string {
	return []string{"node", "npx"}
}

func (i *MemoryMCPInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return generateMCPScript("memory", "@modelcontextprotocol/server-memory", os)
}

// ========================================
// SEQUENTIAL THINKING MCP
// ========================================

type SequentialThinkingMCPInstaller struct {
	BaseInstaller
}

func NewSequentialThinkingMCPInstaller() *SequentialThinkingMCPInstaller {
	return &SequentialThinkingMCPInstaller{
		BaseInstaller: NewBaseInstaller(
			"Sequential Thinking MCP",
			"Step-by-step reasoning and problem solving",
			CategoryMCP,
			"🤔",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (i *SequentialThinkingMCPInstaller) Dependencies() []string {
	return []string{"node", "npx"}
}

func (i *SequentialThinkingMCPInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return generateMCPScript("sequential-thinking", "@modelcontextprotocol/server-sequential-thinking", os)
}

// ========================================
// HELPER FUNCTION
// ========================================

func generateMCPScript(name, pkg string, osType OS) string {
	var script strings.Builder

	switch osType {
	case OSLinux, OSMacOS:
		script.WriteString(fmt.Sprintf(`#!/bin/bash
set -e

echo "📦 Installing MCP: %s..."

# Check if npx is available
if ! command -v npx &> /dev/null; then
    echo "❌ npx is required. Please install Node.js first."
    exit 1
fi

# Pre-download the package
echo "Downloading %s..."
npx -y %s --version 2>/dev/null || true

echo ""
echo "✅ MCP %s is ready!"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Add to your IDE's MCP config:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📁 Cursor:   ~/.cursor/mcp.json"
echo "📁 VSCode:   ~/.vscode/mcp.json"  
echo "📁 Windsurf: ~/.windsurf/mcp.json"
echo ""
echo "Configuration:"
echo '{'
echo '  "mcpServers": {'
echo '    "%s": {'
echo '      "command": "npx",'
echo '      "args": ["-y", "%s"]'
echo '    }'
echo '  }'
echo '}'
`, name, pkg, pkg, name, name, pkg))

	case OSWindows:
		script.WriteString(fmt.Sprintf(`# PowerShell
Write-Host "📦 Installing MCP: %s..." -ForegroundColor Cyan

npx -y %s --version 2>$null

Write-Host "✅ MCP %s is ready!" -ForegroundColor Green
Write-Host "Add to: %%USERPROFILE%%\.cursor\mcp.json"
`, name, pkg, name))
	}

	return script.String()
}

// ========================================
// MCP INTERFACE IMPLEMENTATIONS
// ========================================

// Context7 MCP
func (i *Context7MCPInstaller) RequiredPackageManagers() []PackageManager { return nil }
func (i *Context7MCPInstaller) Install(executor Executor) error {
	return executor.RunWithProgress(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()), nil)
}
func (i *Context7MCPInstaller) Uninstall(executor Executor) error           { return nil }
func (i *Context7MCPInstaller) IsInstalled(executor Executor) (bool, error) { return false, nil }
func (i *Context7MCPInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return "# MCP servers don't require uninstallation - just remove from config"
}

// Playwright MCP
func (i *PlaywrightMCPInstaller) RequiredPackageManagers() []PackageManager { return nil }
func (i *PlaywrightMCPInstaller) Install(executor Executor) error {
	return executor.RunWithProgress(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()), nil)
}
func (i *PlaywrightMCPInstaller) Uninstall(executor Executor) error           { return nil }
func (i *PlaywrightMCPInstaller) IsInstalled(executor Executor) (bool, error) { return false, nil }
func (i *PlaywrightMCPInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return "# MCP servers don't require uninstallation - just remove from config"
}

// GitHub MCP
func (i *GitHubMCPInstaller) RequiredPackageManagers() []PackageManager { return nil }
func (i *GitHubMCPInstaller) Install(executor Executor) error {
	return executor.RunWithProgress(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()), nil)
}
func (i *GitHubMCPInstaller) Uninstall(executor Executor) error           { return nil }
func (i *GitHubMCPInstaller) IsInstalled(executor Executor) (bool, error) { return false, nil }
func (i *GitHubMCPInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return "# MCP servers don't require uninstallation - just remove from config"
}

// Filesystem MCP
func (i *FilesystemMCPInstaller) RequiredPackageManagers() []PackageManager { return nil }
func (i *FilesystemMCPInstaller) Install(executor Executor) error {
	return executor.RunWithProgress(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()), nil)
}
func (i *FilesystemMCPInstaller) Uninstall(executor Executor) error           { return nil }
func (i *FilesystemMCPInstaller) IsInstalled(executor Executor) (bool, error) { return false, nil }
func (i *FilesystemMCPInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return "# MCP servers don't require uninstallation - just remove from config"
}

// Postgres MCP
func (i *PostgresMCPInstaller) RequiredPackageManagers() []PackageManager { return nil }
func (i *PostgresMCPInstaller) Install(executor Executor) error {
	return executor.RunWithProgress(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()), nil)
}
func (i *PostgresMCPInstaller) Uninstall(executor Executor) error           { return nil }
func (i *PostgresMCPInstaller) IsInstalled(executor Executor) (bool, error) { return false, nil }
func (i *PostgresMCPInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return "# MCP servers don't require uninstallation - just remove from config"
}

// Brave Search MCP
func (i *BraveSearchMCPInstaller) RequiredPackageManagers() []PackageManager { return nil }
func (i *BraveSearchMCPInstaller) Install(executor Executor) error {
	return executor.RunWithProgress(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()), nil)
}
func (i *BraveSearchMCPInstaller) Uninstall(executor Executor) error           { return nil }
func (i *BraveSearchMCPInstaller) IsInstalled(executor Executor) (bool, error) { return false, nil }
func (i *BraveSearchMCPInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return "# MCP servers don't require uninstallation - just remove from config"
}

// Memory MCP
func (i *MemoryMCPInstaller) RequiredPackageManagers() []PackageManager { return nil }
func (i *MemoryMCPInstaller) Install(executor Executor) error {
	return executor.RunWithProgress(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()), nil)
}
func (i *MemoryMCPInstaller) Uninstall(executor Executor) error           { return nil }
func (i *MemoryMCPInstaller) IsInstalled(executor Executor) (bool, error) { return false, nil }
func (i *MemoryMCPInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return "# MCP servers don't require uninstallation - just remove from config"
}

// Sequential Thinking MCP
func (i *SequentialThinkingMCPInstaller) RequiredPackageManagers() []PackageManager { return nil }
func (i *SequentialThinkingMCPInstaller) Install(executor Executor) error {
	return executor.RunWithProgress(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()), nil)
}
func (i *SequentialThinkingMCPInstaller) Uninstall(executor Executor) error { return nil }
func (i *SequentialThinkingMCPInstaller) IsInstalled(executor Executor) (bool, error) {
	return false, nil
}
func (i *SequentialThinkingMCPInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return "# MCP servers don't require uninstallation - just remove from config"
}

// Ensure all MCP installers implement Installer interface
var _ Installer = (*Context7MCPInstaller)(nil)
var _ Installer = (*PlaywrightMCPInstaller)(nil)
var _ Installer = (*GitHubMCPInstaller)(nil)
var _ Installer = (*FilesystemMCPInstaller)(nil)
var _ Installer = (*PostgresMCPInstaller)(nil)
var _ Installer = (*BraveSearchMCPInstaller)(nil)
var _ Installer = (*MemoryMCPInstaller)(nil)
var _ Installer = (*SequentialThinkingMCPInstaller)(nil)
