package main

import (
	"flag"
	"fmt"
	"os"

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
		fmt.Println("✅ Update successful! Please restart instacli.")
		return
	}

	// Run TUI
	p := tea.NewProgram(tui.NewApp(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
