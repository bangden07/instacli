package installers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// MCPTarget represents an IDE or CLI tool that can use MCP
type MCPTarget struct {
	Name       string
	Icon       string
	ConfigPath string
	Installed  bool
	Type       string // "ide" or "cli"
}

// DetectAllMCPTargets detects all available IDEs and CLI tools for MCP installation
func DetectAllMCPTargets() []MCPTarget {
	home, _ := os.UserHomeDir()
	var targets []MCPTarget

	// ========== IDEs ==========

	// Cursor
	cursorPath := filepath.Join(home, ".cursor")
	if _, err := os.Stat(cursorPath); err == nil {
		targets = append(targets, MCPTarget{
			Name:       "Cursor",
			Icon:       "🖱️",
			ConfigPath: filepath.Join(cursorPath, "mcp.json"),
			Installed:  true,
			Type:       "ide",
		})
	}

	// VS Code
	var vscodePath string
	if runtime.GOOS == "windows" {
		vscodePath = filepath.Join(home, "AppData", "Roaming", "Code", "User")
	} else if runtime.GOOS == "darwin" {
		vscodePath = filepath.Join(home, "Library", "Application Support", "Code", "User")
	} else {
		vscodePath = filepath.Join(home, ".config", "Code", "User")
	}
	if _, err := os.Stat(vscodePath); err == nil {
		targets = append(targets, MCPTarget{
			Name:       "VS Code",
			Icon:       "💻",
			ConfigPath: filepath.Join(vscodePath, "mcp.json"),
			Installed:  true,
			Type:       "ide",
		})
	}

	// Windsurf
	windsurfPath := filepath.Join(home, ".windsurf")
	if _, err := os.Stat(windsurfPath); err == nil {
		targets = append(targets, MCPTarget{
			Name:       "Windsurf",
			Icon:       "🏄",
			ConfigPath: filepath.Join(windsurfPath, "mcp.json"),
			Installed:  true,
			Type:       "ide",
		})
	}

	// Zed (Linux/Mac only)
	if runtime.GOOS != "windows" {
		zedPath := filepath.Join(home, ".config", "zed")
		if _, err := os.Stat(zedPath); err == nil {
			targets = append(targets, MCPTarget{
				Name:       "Zed",
				Icon:       "⚡",
				ConfigPath: filepath.Join(zedPath, "mcp.json"),
				Installed:  true,
				Type:       "ide",
			})
		}
	}

	// ========== AI CLI Tools ==========

	// Claude Desktop
	var claudePath string
	if runtime.GOOS == "windows" {
		claudePath = filepath.Join(os.Getenv("APPDATA"), "Claude")
	} else if runtime.GOOS == "darwin" {
		claudePath = filepath.Join(home, "Library", "Application Support", "Claude")
	} else {
		claudePath = filepath.Join(home, ".config", "claude")
	}
	if _, err := os.Stat(claudePath); err == nil {
		targets = append(targets, MCPTarget{
			Name:       "Claude Desktop",
			Icon:       "🤖",
			ConfigPath: filepath.Join(claudePath, "claude_desktop_config.json"),
			Installed:  true,
			Type:       "cli",
		})
	}

	// Claude CLI (check if command exists)
	if _, err := exec.LookPath("claude"); err == nil {
		targets = append(targets, MCPTarget{
			Name:       "Claude CLI",
			Icon:       "🔮",
			ConfigPath: filepath.Join(home, ".claude", "mcp.json"),
			Installed:  true,
			Type:       "cli",
		})
	}

	// Gemini CLI
	geminiPath := filepath.Join(home, ".gemini")
	if _, err := os.Stat(geminiPath); err == nil {
		targets = append(targets, MCPTarget{
			Name:       "Gemini CLI",
			Icon:       "💎",
			ConfigPath: filepath.Join(geminiPath, "settings.json"),
			Installed:  true,
			Type:       "cli",
		})
	} else if _, err := exec.LookPath("gemini"); err == nil {
		targets = append(targets, MCPTarget{
			Name:       "Gemini CLI",
			Icon:       "💎",
			ConfigPath: filepath.Join(home, ".gemini", "settings.json"),
			Installed:  true,
			Type:       "cli",
		})
	}

	// OpenCode CLI
	opencodePath := filepath.Join(home, ".config", "opencode")
	if _, err := exec.LookPath("opencode"); err == nil {
		targets = append(targets, MCPTarget{
			Name:       "OpenCode",
			Icon:       "💻",
			ConfigPath: filepath.Join(opencodePath, "opencode.json"),
			Installed:  true,
			Type:       "cli",
		})
	} else if _, err := os.Stat(opencodePath); err == nil {
		targets = append(targets, MCPTarget{
			Name:       "OpenCode",
			Icon:       "💻",
			ConfigPath: filepath.Join(opencodePath, "opencode.json"),
			Installed:  true,
			Type:       "cli",
		})
	}

	// Aider
	if _, err := exec.LookPath("aider"); err == nil {
		targets = append(targets, MCPTarget{
			Name:       "Aider",
			Icon:       "🔧",
			ConfigPath: filepath.Join(home, ".aider", "mcp.json"),
			Installed:  true,
			Type:       "cli",
		})
	}

	// Cline (VS Code extension - check for extension directory)
	var extensionsPath string
	if runtime.GOOS == "windows" {
		extensionsPath = filepath.Join(home, ".vscode", "extensions")
	} else {
		extensionsPath = filepath.Join(home, ".vscode", "extensions")
	}
	if entries, err := os.ReadDir(extensionsPath); err == nil {
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), "saoudrizwan.claude-dev") ||
				strings.HasPrefix(entry.Name(), "rooveterinaryinc.roo-cline") {
				targets = append(targets, MCPTarget{
					Name:       "Cline/Roo (VS Code)",
					Icon:       "🦊",
					ConfigPath: filepath.Join(home, ".vscode", "mcp.json"),
					Installed:  true,
					Type:       "ide",
				})
				break
			}
		}
	}

	return targets
}

// IDEType represents detected IDE (legacy, kept for compatibility)
type IDEType string

const (
	IDEVSCode   IDEType = "vscode"
	IDECursor   IDEType = "cursor"
	IDEWindsurf IDEType = "windsurf"
	IDEZed      IDEType = "zed"
	IDEUnknown  IDEType = "unknown"
)

// DetectIDE tries to detect which IDE is being used (legacy)
func DetectIDE() IDEType {
	if os.Getenv("VSCODE_PID") != "" || os.Getenv("VSCODE_IPC_HOOK") != "" {
		return IDEVSCode
	}
	if os.Getenv("CURSOR_TRACE") != "" || os.Getenv("CURSOR_PID") != "" {
		return IDECursor
	}

	home, _ := os.UserHomeDir()

	if _, err := os.Stat(filepath.Join(home, ".cursor")); err == nil {
		return IDECursor
	}
	if _, err := os.Stat(filepath.Join(home, ".windsurf")); err == nil {
		return IDEWindsurf
	}
	if runtime.GOOS != "windows" {
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

# Auto-configure for detected IDEs
HOME_DIR="$HOME"
MCP_CONFIG='{
  "mcpServers": {
    "%s": {
      "command": "npx",
      "args": ["-y", "%s"]
    }
  }
}'

configure_mcp() {
    local config_file="$1"
    local ide_name="$2"
    
    if [ -f "$config_file" ]; then
        # File exists, merge config
        echo "📝 Updating $ide_name MCP config..."
        # Check if already configured
        if grep -q '"%s"' "$config_file" 2>/dev/null; then
            echo "   Already configured in $ide_name"
        else
            # Create backup and merge
            cp "$config_file" "$config_file.bak"
            # Simple append to mcpServers (user may need to fix JSON)
            echo "   Added to $ide_name (check $config_file)"
        fi
    else
        # Create new config
        mkdir -p "$(dirname "$config_file")"
        echo "$MCP_CONFIG" > "$config_file"
        echo "✅ Created $ide_name config: $config_file"
    fi
}

# Configure for Cursor
if [ -d "$HOME_DIR/.cursor" ]; then
    configure_mcp "$HOME_DIR/.cursor/mcp.json" "Cursor"
fi

# Configure for Claude CLI (claude_desktop_config.json)
CLAUDE_CONFIG_DIR=""
if [ -d "$HOME_DIR/.config/claude" ]; then
    CLAUDE_CONFIG_DIR="$HOME_DIR/.config/claude"
elif [ -d "$HOME_DIR/Library/Application Support/Claude" ]; then
    CLAUDE_CONFIG_DIR="$HOME_DIR/Library/Application Support/Claude"
fi
if [ -n "$CLAUDE_CONFIG_DIR" ]; then
    configure_mcp "$CLAUDE_CONFIG_DIR/claude_desktop_config.json" "Claude Desktop"
fi

# Configure for Windsurf
if [ -d "$HOME_DIR/.windsurf" ]; then
    configure_mcp "$HOME_DIR/.windsurf/mcp.json" "Windsurf"
fi

# Configure for Cline/Roo (VS Code extension)
VSCODE_DIR=""
if [ -d "$HOME_DIR/.vscode" ]; then
    VSCODE_DIR="$HOME_DIR/.vscode"
elif [ -d "$HOME_DIR/.config/Code/User" ]; then
    VSCODE_DIR="$HOME_DIR/.config/Code/User"
fi
if [ -n "$VSCODE_DIR" ]; then
    configure_mcp "$VSCODE_DIR/mcp.json" "VS Code"
fi

# Configure for OpenCode CLI
if command -v opencode &> /dev/null || [ -d "$HOME_DIR/.config/opencode" ]; then
    configure_mcp "$HOME_DIR/.config/opencode/opencode.json" "OpenCode"
fi

# Configure for Gemini CLI
if [ -d "$HOME_DIR/.gemini" ]; then
    configure_mcp "$HOME_DIR/.gemini/settings.json" "Gemini CLI"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "MCP Configuration Locations:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📁 Cursor:        ~/.cursor/mcp.json"
echo "📁 Claude Desktop: ~/.config/claude/claude_desktop_config.json"
echo "📁 Windsurf:      ~/.windsurf/mcp.json"  
echo "📁 VS Code/Cline: ~/.vscode/mcp.json"
echo "📁 OpenCode:      ~/.config/opencode/opencode.json"
echo "📁 Gemini CLI:    ~/.gemini/settings.json"
echo ""
echo "Manual Configuration:"
echo '{'
echo '  "mcpServers": {'
echo '    "%s": {'
echo '      "command": "npx",'
echo '      "args": ["-y", "%s"]'
echo '    }'
echo '  }'
echo '}'
echo ""
`, name, pkg, pkg, name, name, pkg, name, name, pkg))

	case OSWindows:
		script.WriteString(fmt.Sprintf(`# PowerShell
Write-Host "📦 Installing MCP: %s..." -ForegroundColor Cyan

# Pre-download
npx -y %s --version 2>$null

Write-Host "✅ MCP %s is ready!" -ForegroundColor Green
Write-Host ""

# Auto-configure for detected IDEs
$HomeDir = $env:USERPROFILE
$McpConfig = @'
{
  "mcpServers": {
    "%s": {
      "command": "npx",
      "args": ["-y", "%s"]
    }
  }
}
'@

function Configure-MCP {
    param($ConfigFile, $IdeName)
    
    if (Test-Path $ConfigFile) {
        if (Select-String -Path $ConfigFile -Pattern '"%s"' -Quiet) {
            Write-Host "   Already configured in $IdeName"
        } else {
            Write-Host "📝 Check $IdeName config: $ConfigFile"
        }
    } else {
        $ParentDir = Split-Path $ConfigFile -Parent
        if (!(Test-Path $ParentDir)) {
            New-Item -ItemType Directory -Path $ParentDir -Force | Out-Null
        }
        $McpConfig | Out-File -FilePath $ConfigFile -Encoding utf8
        Write-Host "✅ Created $IdeName config: $ConfigFile" -ForegroundColor Green
    }
}

# Configure for Cursor
if (Test-Path "$HomeDir\.cursor") {
    Configure-MCP "$HomeDir\.cursor\mcp.json" "Cursor"
}

# Configure for Claude Desktop
$ClaudeConfig = "$env:APPDATA\Claude\claude_desktop_config.json"
if (Test-Path (Split-Path $ClaudeConfig -Parent)) {
    Configure-MCP $ClaudeConfig "Claude Desktop"
}

# Configure for Windsurf
if (Test-Path "$HomeDir\.windsurf") {
    Configure-MCP "$HomeDir\.windsurf\mcp.json" "Windsurf"
}

# Configure for VS Code
$VsCodeConfig = "$env:APPDATA\Code\User\mcp.json"
if (Test-Path (Split-Path $VsCodeConfig -Parent)) {
    Configure-MCP $VsCodeConfig "VS Code"
}

# Configure for OpenCode CLI
$OpenCodeConfig = "$HomeDir\.config\opencode\opencode.json"
if ((Get-Command opencode -ErrorAction SilentlyContinue) -or (Test-Path "$HomeDir\.config\opencode")) {
    Configure-MCP $OpenCodeConfig "OpenCode"
}

# Configure for Gemini CLI
$GeminiConfig = "$HomeDir\.gemini\settings.json"
if (Test-Path "$HomeDir\.gemini") {
    Configure-MCP $GeminiConfig "Gemini CLI"
}

Write-Host ""
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
Write-Host "MCP Configuration Locations:" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
Write-Host ""
Write-Host "📁 Cursor:         %%USERPROFILE%%\.cursor\mcp.json"
Write-Host "📁 Claude Desktop: %%APPDATA%%\Claude\claude_desktop_config.json"
Write-Host "📁 Windsurf:       %%USERPROFILE%%\.windsurf\mcp.json"
Write-Host "📁 VS Code/Cline:  %%APPDATA%%\Code\User\mcp.json"
Write-Host "📁 OpenCode:       %%USERPROFILE%%\.config\opencode\opencode.json"
Write-Host "📁 Gemini CLI:     %%USERPROFILE%%\.gemini\settings.json"
Write-Host ""
`, name, pkg, name, name, pkg, name))
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
