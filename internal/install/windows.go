package install

import (
	"os/exec"
)

// WindowsServiceInstaller handles Windows Service installation
type WindowsServiceInstaller struct {
	serviceName string
	binaryPath  string
}

// NewWindowsServiceInstaller creates a new Windows service installer
func NewWindowsServiceInstaller(serviceName, binaryPath string) *WindowsServiceInstaller {
	return &WindowsServiceInstaller{
		serviceName: serviceName,
		binaryPath:  binaryPath,
	}
}

// Install creates a Windows service
func (wsi *WindowsServiceInstaller) Install() error {
	// Use nssm (Non-Sucking Service Manager) or sc command
	// This example uses sc.exe (built-in on Windows)
	return exec.Command("sc", "create", wsi.serviceName, "binPath=", wsi.binaryPath).Run()
}

// Uninstall removes the Windows service
func (wsi *WindowsServiceInstaller) Uninstall() error {
	return exec.Command("sc", "delete", wsi.serviceName).Run()
}

// Start starts the Windows service
func (wsi *WindowsServiceInstaller) Start() error {
	return exec.Command("sc", "start", wsi.serviceName).Run()
}

// Stop stops the Windows service
func (wsi *WindowsServiceInstaller) Stop() error {
	return exec.Command("sc", "stop", wsi.serviceName).Run()
}
