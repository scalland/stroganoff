package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/yourusername/stroganoff/internal/install"
)

var (
	installService string
	installUser    string
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install gocr as a service",
	Long: `Install the gocr application as a system service.
Supports systemd (Linux), launchd (macOS), and Windows Service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performInstall()
	},
}

func init() {
	installCmd.Flags().StringVar(&installService, "service", "gocr", "Service name")
	installCmd.Flags().StringVar(&installUser, "user", "", "User to run service as (Linux only)")
}

func performInstall() error {
	// Get the path of the current executable
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	absPath, err := filepath.Abs(exePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if running as root/admin
	if !isAdmin() {
		return fmt.Errorf("this command requires administrator/root privileges")
	}

	goos := runtime.GOOS
	var installer install.ServiceInstaller

	switch goos {
	case "linux":
		installer = install.NewSystemdInstaller(installService, absPath, installUser)
	case "darwin":
		installer = install.NewLaunchdInstaller(installService, absPath)
	case "windows":
		installer = install.NewWindowsServiceInstaller(installService, absPath)
	default:
		return fmt.Errorf("unsupported operating system: %s", goos)
	}

	if err := installer.Install(); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	fmt.Printf("Service '%s' installed successfully on %s\n", installService, goos)
	fmt.Printf("Start service: ")
	switch goos {
	case "linux":
		fmt.Println("sudo systemctl start " + installService)
	case "darwin":
		fmt.Println("launchctl start " + installService)
	case "windows":
		fmt.Println("net start " + installService)
	}

	return nil
}

func isAdmin() bool {
	goos := runtime.GOOS

	switch goos {
	case "windows":
		_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
		return err == nil
	case "linux", "darwin":
		return os.Getuid() == 0
	}

	return false
}
