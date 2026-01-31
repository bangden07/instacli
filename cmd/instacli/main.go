package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/instacli/instacli/internal/tui"
	"github.com/instacli/instacli/internal/updater"
	"github.com/instacli/instacli/internal/version"
)

func main() {
	// Parse command line flags
	versionFlag := flag.Bool("version", false, "Print version information")
	updateFlag := flag.Bool("update", false, "Check for updates and self-update")
	checkUpdateFlag := flag.Bool("check-update", false, "Check for updates without installing")
	skipUpdateCheck := flag.Bool("no-update-check", false, "Skip automatic update check on startup")
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Printf("InstaCli v%s\n", version.Version)
		fmt.Printf("Build Date: %s\n", version.BuildDate)
		return
	}

	// Handle check-update flag
	if *checkUpdateFlag {
		fmt.Println("🔍 Checking for updates...")
		result := updater.CheckForUpdate()
		if result.Error != nil {
			fmt.Printf("❌ Error checking for updates: %v\n", result.Error)
			os.Exit(1)
		}
		if result.UpdateAvailable {
			fmt.Printf("✨ New version available: v%s (current: v%s)\n", result.LatestVersion, result.CurrentVersion)
			fmt.Println("   Run 'instacli --update' to update")
		} else {
			fmt.Printf("✅ You're running the latest version (v%s)\n", result.CurrentVersion)
		}
		return
	}

	// Handle update flag
	if *updateFlag {
		performUpdate()
		return
	}

	// Auto-check for updates on startup (unless skipped)
	if !*skipUpdateCheck {
		checkAndPromptUpdate()
	}

	// Run TUI
	p := tea.NewProgram(tui.NewApp(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// performUpdate handles the update process
func performUpdate() {
	fmt.Println("🔍 Checking for updates...")
	result := updater.CheckForUpdate()
	if result.Error != nil {
		fmt.Printf("❌ Error checking for updates: %v\n", result.Error)
		os.Exit(1)
	}
	if !result.UpdateAvailable {
		fmt.Printf("✅ You're already running the latest version (v%s)\n", result.CurrentVersion)
		return
	}
	fmt.Printf("✨ Updating from v%s to v%s...\n", result.CurrentVersion, result.LatestVersion)
	if err := updater.SelfUpdate(result.DownloadURL); err != nil {
		fmt.Printf("❌ Update failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Update successful! Restarting...")

	// Restart the application
	restartApp()
}

// checkAndPromptUpdate checks for updates and prompts user
func checkAndPromptUpdate() {
	fmt.Print("🔍 Checking for updates... ")
	result := updater.CheckForUpdate()

	if result.Error != nil {
		// Silently continue if check fails (network issues, etc.)
		fmt.Println("⚠️  (offline)")
		return
	}

	if !result.UpdateAvailable {
		fmt.Println("✅")
		return
	}

	// Update available - prompt user
	fmt.Println()
	fmt.Println()
	fmt.Printf("╔══════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║  🎉 NEW VERSION AVAILABLE!                                   ║\n")
	fmt.Printf("║                                                              ║\n")
	fmt.Printf("║  Current version: v%-10s                                ║\n", result.CurrentVersion)
	fmt.Printf("║  Latest version:  v%-10s                                ║\n", result.LatestVersion)
	fmt.Printf("║                                                              ║\n")
	fmt.Printf("╚══════════════════════════════════════════════════════════════╝\n")
	fmt.Println()
	fmt.Print("Do you want to update now? [Y/n]: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	// Default to yes if user just presses enter
	if input == "" || input == "y" || input == "yes" {
		fmt.Println()
		fmt.Printf("📥 Downloading v%s...\n", result.LatestVersion)

		if err := updater.SelfUpdate(result.DownloadURL); err != nil {
			fmt.Printf("❌ Update failed: %v\n", err)
			fmt.Println("Continuing with current version...")
			fmt.Println()
			return
		}

		fmt.Println("✅ Update successful! Restarting...")
		fmt.Println()

		// Restart the application
		restartApp()
	} else {
		fmt.Println()
		fmt.Printf("Continuing with v%s...\n", result.CurrentVersion)
		fmt.Println()
	}
}

// restartApp restarts the application after update
func restartApp() {
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Please restart instacli manually.\n")
		os.Exit(0)
	}

	// Start new instance with --no-update-check flag to prevent loop
	cmd := exec.Command(execPath, "--no-update-check")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("Failed to restart: %v\n", err)
		fmt.Printf("Please restart instacli manually.\n")
		os.Exit(0)
	}

	// Exit current instance
	os.Exit(0)
}
