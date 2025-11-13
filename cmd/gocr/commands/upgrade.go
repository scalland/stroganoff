package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/stroganoff/internal/upgrade"
	"github.com/yourusername/stroganoff/pkg/version"
)

var (
	upgradeVersion string
	upgradeToken   string
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade to a new version",
	Long: `Download and install a new version of gocr from Github releases.
Can upgrade to a specific version or the latest available version.
Supports both public and private repositories with authentication token.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performUpgrade()
	},
}

func init() {
	upgradeCmd.Flags().StringVar(&upgradeVersion, "version", "latest", "Version to upgrade to (default: latest)")
	upgradeCmd.Flags().StringVar(&upgradeToken, "token", "", "Github token for private repositories (optional)")
}

func performUpgrade() error {
	fmt.Printf("Current version: %s\n", version.GetVersion())

	// Get the latest or specified release
	client := upgrade.NewGithubClient(upgradeToken)

	var releaseInfo *upgrade.Release
	var err error

	if upgradeVersion == "latest" {
		releaseInfo, err = client.GetLatestRelease("yourusername", "gocr")
	} else {
		releaseInfo, err = client.GetRelease("yourusername", "gocr", upgradeVersion)
	}

	if err != nil {
		return fmt.Errorf("failed to fetch release: %w", err)
	}

	fmt.Printf("Available version: %s\n", releaseInfo.TagName)

	if releaseInfo.TagName == "v"+version.GetVersion() {
		fmt.Println("Already at the latest version")
		return nil
	}

	// Find the appropriate asset for current OS/ARCH
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	assetName := fmt.Sprintf("gocr-*-%s-%s", goos, goarch)
	if goos == "windows" {
		assetName = fmt.Sprintf("gocr-*-%s-%s.exe", goos, goarch)
	}

	var downloadURL string
	for _, asset := range releaseInfo.Assets {
		if matchesPattern(asset.Name, assetName) {
			downloadURL = asset.DownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no compatible binary found for %s/%s", goos, goarch)
	}

	fmt.Printf("Downloading %s...\n", filepath.Base(downloadURL))

	// Download the binary
	binaryPath, err := downloadBinary(downloadURL, upgradeToken)
	if err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}

	// Get the path of the current executable
	currentExePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Backup the current binary
	backupPath := currentExePath + ".bak"
	if err := os.Rename(currentExePath, backupPath); err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	// Move new binary to the current location
	if err := os.Rename(binaryPath, currentExePath); err != nil {
		// Restore backup on failure
		os.Rename(backupPath, currentExePath)
		return fmt.Errorf("failed to install new binary: %w", err)
	}

	// Make it executable
	if err := os.Chmod(currentExePath, 0755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	// Remove backup
	os.Remove(backupPath)

	fmt.Printf("Successfully upgraded to %s\n", releaseInfo.TagName)
	return nil
}

func downloadBinary(url, token string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "gocr-*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

func matchesPattern(filename, pattern string) bool {
	// Simple glob pattern matching for downloads
	parts := strings.Split(pattern, "*")
	for i, part := range parts {
		if part == "" {
			continue
		}
		if i == 0 && !strings.HasPrefix(filename, part) {
			return false
		}
		if i == len(parts)-1 && !strings.HasSuffix(filename, part) {
			return false
		}
		if i > 0 && i < len(parts)-1 && !strings.Contains(filename, part) {
			return false
		}
	}
	return true
}
