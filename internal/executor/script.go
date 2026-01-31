package executor

import (
	"fmt"
	goos "os"
	"path/filepath"
	"time"

	"github.com/instacli/instacli/internal/installers"
)

// ScriptGenerator generates installation scripts for manual use
type ScriptGenerator struct {
	outputDir string
}

// NewScriptGenerator creates a new script generator
func NewScriptGenerator(outputDir string) *ScriptGenerator {
	if outputDir == "" {
		outputDir = "./scripts/generated"
	}
	return &ScriptGenerator{
		outputDir: outputDir,
	}
}

// GenerateScript generates a script file for an installer
func (g *ScriptGenerator) GenerateScript(installer installers.Installer, targetOS installers.OS, pm installers.PackageManager) (string, error) {
	script := installer.GenerateInstallScript(targetOS, pm)
	if script == "" {
		return "", fmt.Errorf("no script available for %s on %s", installer.Name(), targetOS)
	}

	// Ensure output directory exists
	if err := ensureDir(g.outputDir); err != nil {
		return "", err
	}

	// Generate filename
	timestamp := time.Now().Format("20060102-150405")
	ext := ".sh"
	if targetOS == installers.OSWindows {
		ext = ".ps1"
	}
	filename := fmt.Sprintf("%s_%s_%s%s", sanitizeName(installer.Name()), targetOS, timestamp, ext)
	filePath := filepath.Join(g.outputDir, filename)

	// Write script
	err := goos.WriteFile(filePath, []byte(script), 0755)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// GenerateMultiScript generates a combined script for multiple installers
func (g *ScriptGenerator) GenerateMultiScript(installersList []installers.Installer, targetOS installers.OS, pm installers.PackageManager) (string, error) {
	var combinedScript string

	// Add header
	if targetOS == installers.OSWindows {
		combinedScript = "# InstaCli Generated Script\n# Generated at " + time.Now().Format(time.RFC3339) + "\n\n"
	} else {
		combinedScript = "#!/bin/bash\n# InstaCli Generated Script\n# Generated at " + time.Now().Format(time.RFC3339) + "\nset -e\n\n"
	}

	// Add each installer's script
	for _, inst := range installersList {
		script := inst.GenerateInstallScript(targetOS, pm)
		if script != "" {
			combinedScript += fmt.Sprintf("\n# ========== %s ==========\n", inst.Name())
			combinedScript += script + "\n"
		}
	}

	// Ensure output directory exists
	if err := ensureDir(g.outputDir); err != nil {
		return "", err
	}

	// Generate filename
	timestamp := time.Now().Format("20060102-150405")
	ext := ".sh"
	if targetOS == installers.OSWindows {
		ext = ".ps1"
	}
	filename := fmt.Sprintf("combined_%s_%s%s", targetOS, timestamp, ext)
	filePath := filepath.Join(g.outputDir, filename)

	// Write script
	err := goos.WriteFile(filePath, []byte(combinedScript), 0755)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// Helper functions
func ensureDir(dir string) error {
	return goos.MkdirAll(dir, 0755)
}

func sanitizeName(name string) string {
	result := ""
	for _, c := range name {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			result += string(c)
		} else if c == ' ' || c == '-' {
			result += "_"
		}
	}
	return result
}
