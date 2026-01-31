package installers

import (
	"fmt"
	"strings"
)

// ========================================
// CLAUDE CLI
// ========================================

type ClaudeCLIInstaller struct {
	BaseInstaller
}

func NewClaudeCLIInstaller() *ClaudeCLIInstaller {
	return &ClaudeCLIInstaller{
		BaseInstaller: NewBaseInstaller(
			"Claude CLI",
			"Anthropic's Claude Code - AI coding assistant",
			CategoryAICLI,
			"🤖",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (i *ClaudeCLIInstaller) Dependencies() []string {
	return []string{"node", "npm"}
}

func (i *ClaudeCLIInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	var script strings.Builder

	switch os {
	case OSLinux, OSMacOS:
		script.WriteString(`#!/bin/bash
set -e

echo "🤖 Installing Claude CLI..."

# Check if npm is available
if ! command -v npm &> /dev/null; then
    echo "❌ npm is required. Please install Node.js first."
    exit 1
fi

# Install Claude CLI globally
npm install -g @anthropic-ai/claude-code

echo ""
echo "✅ Claude CLI installed successfully!"
echo ""
echo "To get started:"
echo "  1. Run: claude"
echo "  2. Follow the authentication prompts"
echo "  3. Start coding with AI assistance!"
`)

	case OSWindows:
		script.WriteString(`# PowerShell script for Windows
Write-Host "🤖 Installing Claude CLI..." -ForegroundColor Cyan

# Check if npm is available
if (-not (Get-Command npm -ErrorAction SilentlyContinue)) {
    Write-Host "❌ npm is required. Please install Node.js first." -ForegroundColor Red
    exit 1
}

# Install Claude CLI globally
npm install -g @anthropic-ai/claude-code

Write-Host ""
Write-Host "✅ Claude CLI installed successfully!" -ForegroundColor Green
Write-Host "Run 'claude' to get started"
`)
	}

	return script.String()
}

// ========================================
// GEMINI CLI
// ========================================

type GeminiCLIInstaller struct {
	BaseInstaller
}

func NewGeminiCLIInstaller() *GeminiCLIInstaller {
	return &GeminiCLIInstaller{
		BaseInstaller: NewBaseInstaller(
			"Gemini CLI",
			"Google's Gemini CLI - AI coding assistant",
			CategoryAICLI,
			"💎",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (i *GeminiCLIInstaller) Dependencies() []string {
	return []string{"node", "npm"}
}

func (i *GeminiCLIInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	var script strings.Builder

	switch os {
	case OSLinux, OSMacOS:
		script.WriteString(`#!/bin/bash
set -e

echo "💎 Installing Gemini CLI..."

# Check if npm is available
if ! command -v npm &> /dev/null; then
    echo "❌ npm is required. Please install Node.js first."
    exit 1
fi

# Install Gemini CLI globally
npm install -g @anthropic-ai/gemini-cli

echo ""
echo "✅ Gemini CLI installed successfully!"
echo ""
echo "To get started:"
echo "  1. Run: gemini"
echo "  2. Set your API key: export GOOGLE_API_KEY=your_key"
echo "  3. Start coding with AI assistance!"
`)

	case OSWindows:
		script.WriteString(`# PowerShell script for Windows
Write-Host "💎 Installing Gemini CLI..." -ForegroundColor Cyan

npm install -g @anthropic-ai/gemini-cli

Write-Host "✅ Gemini CLI installed!" -ForegroundColor Green
`)
	}

	return script.String()
}

// ========================================
// CODEX CLI (OpenAI)
// ========================================

type CodexCLIInstaller struct {
	BaseInstaller
}

func NewCodexCLIInstaller() *CodexCLIInstaller {
	return &CodexCLIInstaller{
		BaseInstaller: NewBaseInstaller(
			"Codex CLI",
			"OpenAI's Codex CLI - AI coding assistant",
			CategoryAICLI,
			"🧠",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (i *CodexCLIInstaller) Dependencies() []string {
	return []string{"node", "npm"}
}

func (i *CodexCLIInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	var script strings.Builder

	switch os {
	case OSLinux, OSMacOS:
		script.WriteString(`#!/bin/bash
set -e

echo "🧠 Installing OpenAI Codex CLI..."

# Check if npm is available
if ! command -v npm &> /dev/null; then
    echo "❌ npm is required. Please install Node.js first."
    exit 1
fi

# Install Codex CLI globally
npm install -g @openai/codex

echo ""
echo "✅ Codex CLI installed successfully!"
echo ""
echo "To get started:"
echo "  1. Set your API key: export OPENAI_API_KEY=your_key"
echo "  2. Run: codex"
`)

	case OSWindows:
		script.WriteString(`# PowerShell script for Windows
Write-Host "🧠 Installing OpenAI Codex CLI..." -ForegroundColor Cyan

npm install -g @openai/codex

Write-Host "✅ Codex CLI installed!" -ForegroundColor Green
`)
	}

	return script.String()
}

// ========================================
// AIDER
// ========================================

type AiderInstaller struct {
	BaseInstaller
}

func NewAiderInstaller() *AiderInstaller {
	return &AiderInstaller{
		BaseInstaller: NewBaseInstaller(
			"Aider",
			"AI pair programming in your terminal",
			CategoryAICLI,
			"🔧",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (i *AiderInstaller) Dependencies() []string {
	return []string{"python3", "pip"}
}

func (i *AiderInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	var script strings.Builder

	switch os {
	case OSLinux, OSMacOS:
		script.WriteString(`#!/bin/bash
set -e

echo "🔧 Installing Aider..."

# Check if pip is available
if command -v pip3 &> /dev/null; then
    PIP=pip3
elif command -v pip &> /dev/null; then
    PIP=pip
else
    echo "❌ pip is required. Please install Python first."
    exit 1
fi

# Install aider
$PIP install aider-chat

echo ""
echo "✅ Aider installed successfully!"
echo ""
echo "To get started:"
echo "  1. Set your API key: export OPENAI_API_KEY=your_key"
echo "  2. Navigate to your project directory"
echo "  3. Run: aider"
echo ""
echo "Aider supports: OpenAI, Anthropic, Google, and local models"
`)

	case OSWindows:
		script.WriteString(`# PowerShell script for Windows
Write-Host "🔧 Installing Aider..." -ForegroundColor Cyan

pip install aider-chat

Write-Host "✅ Aider installed!" -ForegroundColor Green
Write-Host "Run 'aider' in your project directory to get started"
`)
	}

	return script.String()
}

// ========================================
// KILO CODE (VS Code Extension)
// ========================================

type KiloCodeInstaller struct {
	BaseInstaller
}

func NewKiloCodeInstaller() *KiloCodeInstaller {
	return &KiloCodeInstaller{
		BaseInstaller: NewBaseInstaller(
			"Kilo Code",
			"AI coding extension for VS Code/Cursor",
			CategoryAICLI,
			"⚡",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (i *KiloCodeInstaller) Dependencies() []string {
	return []string{}
}

func (i *KiloCodeInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	var script strings.Builder

	switch os {
	case OSLinux, OSMacOS:
		script.WriteString(`#!/bin/bash
set -e

echo "⚡ Installing Kilo Code..."

# Try VS Code first
if command -v code &> /dev/null; then
    echo "Installing for VS Code..."
    code --install-extension kilocode.kilocode
    echo "✅ Kilo Code installed in VS Code!"
fi

# Try Cursor
if command -v cursor &> /dev/null; then
    echo "Installing for Cursor..."
    cursor --install-extension kilocode.kilocode
    echo "✅ Kilo Code installed in Cursor!"
fi

echo ""
echo "🎉 Installation complete!"
echo "Restart your editor to activate Kilo Code"
`)

	case OSWindows:
		script.WriteString(`# PowerShell script for Windows
Write-Host "⚡ Installing Kilo Code..." -ForegroundColor Cyan

if (Get-Command code -ErrorAction SilentlyContinue) {
    code --install-extension kilocode.kilocode
    Write-Host "✅ Kilo Code installed in VS Code!" -ForegroundColor Green
}

Write-Host "Restart your editor to activate"
`)
	}

	return script.String()
}

// ========================================
// CONTINUE (VS Code Extension)
// ========================================

type ContinueInstaller struct {
	BaseInstaller
}

func NewContinueInstaller() *ContinueInstaller {
	return &ContinueInstaller{
		BaseInstaller: NewBaseInstaller(
			"Continue",
			"Open-source AI code assistant extension",
			CategoryAICLI,
			"➡️",
			[]OS{OSLinux, OSMacOS, OSWindows},
		),
	}
}

func (i *ContinueInstaller) Dependencies() []string {
	return []string{}
}

func (i *ContinueInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	var script strings.Builder

	script.WriteString(fmt.Sprintf(`#!/bin/bash
set -e

echo "➡️ Installing Continue..."

# Try VS Code
if command -v code &> /dev/null; then
    echo "Installing for VS Code..."
    code --install-extension continue.continue
    echo "✅ Continue installed in VS Code!"
fi

# Try Cursor  
if command -v cursor &> /dev/null; then
    echo "Installing for Cursor..."
    cursor --install-extension continue.continue
    echo "✅ Continue installed in Cursor!"
fi

echo ""
echo "🎉 Installation complete!"
echo "Configure your AI provider in Continue settings"
`))

	return script.String()
}

// ========================================
// INTERFACE IMPLEMENTATIONS
// ========================================

// ClaudeCLI interface implementations
func (i *ClaudeCLIInstaller) RequiredPackageManagers() []PackageManager { return nil }
func (i *ClaudeCLIInstaller) Install(executor Executor) error {
	return executor.RunWithProgress(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()), nil)
}
func (i *ClaudeCLIInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *ClaudeCLIInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("claude --version")
	return err == nil, nil
}
func (i *ClaudeCLIInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return "npm uninstall -g @anthropic-ai/claude-code"
}

// GeminiCLI interface implementations
func (i *GeminiCLIInstaller) RequiredPackageManagers() []PackageManager { return nil }
func (i *GeminiCLIInstaller) Install(executor Executor) error {
	return executor.RunWithProgress(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()), nil)
}
func (i *GeminiCLIInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *GeminiCLIInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("gemini --version")
	return err == nil, nil
}
func (i *GeminiCLIInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return "npm uninstall -g @google/generative-ai-cli"
}

// CodexCLI interface implementations
func (i *CodexCLIInstaller) RequiredPackageManagers() []PackageManager { return nil }
func (i *CodexCLIInstaller) Install(executor Executor) error {
	return executor.RunWithProgress(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()), nil)
}
func (i *CodexCLIInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *CodexCLIInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("codex --version")
	return err == nil, nil
}
func (i *CodexCLIInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return "npm uninstall -g @openai/codex"
}

// Aider interface implementations
func (i *AiderInstaller) RequiredPackageManagers() []PackageManager { return nil }
func (i *AiderInstaller) Install(executor Executor) error {
	return executor.RunWithProgress(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()), nil)
}
func (i *AiderInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *AiderInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("aider --version")
	return err == nil, nil
}
func (i *AiderInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return "pip uninstall -y aider-chat"
}

// KiloCode interface implementations
func (i *KiloCodeInstaller) RequiredPackageManagers() []PackageManager { return nil }
func (i *KiloCodeInstaller) Install(executor Executor) error {
	return executor.RunWithProgress(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()), nil)
}
func (i *KiloCodeInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *KiloCodeInstaller) IsInstalled(executor Executor) (bool, error) {
	return false, nil // Extension check not easily scriptable
}
func (i *KiloCodeInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return "code --uninstall-extension kilocode.kilocode"
}

// Continue interface implementations
func (i *ContinueInstaller) RequiredPackageManagers() []PackageManager { return nil }
func (i *ContinueInstaller) Install(executor Executor) error {
	return executor.RunWithProgress(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()), nil)
}
func (i *ContinueInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *ContinueInstaller) IsInstalled(executor Executor) (bool, error) {
	return false, nil // Extension check not easily scriptable
}
func (i *ContinueInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return "code --uninstall-extension continue.continue"
}

// Ensure all installers implement Installer interface
var _ Installer = (*ClaudeCLIInstaller)(nil)
var _ Installer = (*GeminiCLIInstaller)(nil)
var _ Installer = (*CodexCLIInstaller)(nil)
var _ Installer = (*AiderInstaller)(nil)
var _ Installer = (*KiloCodeInstaller)(nil)
var _ Installer = (*ContinueInstaller)(nil)
