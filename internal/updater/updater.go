package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/instacli/instacli/internal/version"
)

const (
	repoOwner = "bangden07"
	repoName  = "instacli"
	apiURL    = "https://api.github.com/repos/" + repoOwner + "/" + repoName + "/releases/latest"
)

// ReleaseInfo contains GitHub release information
type ReleaseInfo struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

// Asset represents a release asset
type Asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
}

// UpdateResult contains the result of an update check
type UpdateResult struct {
	CurrentVersion  string
	LatestVersion   string
	UpdateAvailable bool
	DownloadURL     string
	Error           error
}

// CheckForUpdate checks if a new version is available
func CheckForUpdate() UpdateResult {
	result := UpdateResult{
		CurrentVersion: version.Version,
	}

	// Create HTTP client with timeout
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		result.Error = err
		return result
	}
	req.Header.Set("User-Agent", "instacli/"+version.Version)

	resp, err := client.Do(req)
	if err != nil {
		result.Error = err
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Error = fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
		return result
	}

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		result.Error = err
		return result
	}

	// Parse version (remove 'v' prefix)
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	result.LatestVersion = latestVersion

	// Compare versions
	if compareVersions(latestVersion, version.Version) > 0 {
		result.UpdateAvailable = true
		result.DownloadURL = getAssetURL(release.Assets)
	}

	return result
}

// SelfUpdate downloads and replaces the current binary
func SelfUpdate(downloadURL string) error {
	if downloadURL == "" {
		return fmt.Errorf("no download URL available for platform %s/%s. Check releases at https://github.com/%s/%s/releases",
			runtime.GOOS, runtime.GOARCH, repoOwner, repoName)
	}

	fmt.Printf("📥 Downloading from: %s\n", downloadURL)

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	fmt.Printf("📍 Current binary: %s\n", execPath)

	// Download new binary
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Read into memory first to ensure complete download
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read update data: %w", err)
	}

	fmt.Printf("📦 Downloaded %d bytes\n", len(data))

	if len(data) < 1000 {
		return fmt.Errorf("downloaded file too small (%d bytes), likely an error page", len(data))
	}

	// Create temp file in same directory
	tmpPath := execPath + ".new"
	if err := os.WriteFile(tmpPath, data, 0755); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Backup old binary
	backupPath := execPath + ".bak"
	os.Remove(backupPath) // Remove old backup if exists

	if err := os.Rename(execPath, backupPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to backup old binary: %w", err)
	}

	// Move new binary to executable path
	if err := os.Rename(tmpPath, execPath); err != nil {
		// Restore backup on failure
		os.Rename(backupPath, execPath)
		os.Remove(tmpPath)
		return fmt.Errorf("failed to install new binary: %w", err)
	}

	// Make executable (Unix)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(execPath, 0755); err != nil {
			fmt.Printf("⚠️ Warning: failed to set permissions: %v\n", err)
		}
	}

	// Remove backup
	os.Remove(backupPath)

	return nil
}

// getAssetURL finds the correct asset for the current platform
func getAssetURL(assets []Asset) string {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// Convert to our naming convention
	if osName == "darwin" {
		osName = "darwin"
	}

	expectedName := fmt.Sprintf("instacli-%s-%s", osName, arch)
	if osName == "windows" {
		expectedName += ".exe"
	}

	for _, asset := range assets {
		if asset.Name == expectedName {
			return asset.DownloadURL
		}
	}

	return ""
}

// compareVersions compares two semantic versions
// Returns: 1 if v1 > v2, -1 if v1 < v2, 0 if equal
func compareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	for i := 0; i < 3; i++ {
		var n1, n2 int
		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &n1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &n2)
		}

		if n1 > n2 {
			return 1
		}
		if n1 < n2 {
			return -1
		}
	}

	return 0
}
