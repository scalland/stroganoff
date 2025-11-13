package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

// LaunchdInstaller handles launchd service installation on macOS
type LaunchdInstaller struct {
	serviceName string
	binaryPath  string
}

// NewLaunchdInstaller creates a new launchd installer
func NewLaunchdInstaller(serviceName, binaryPath string) *LaunchdInstaller {
	return &LaunchdInstaller{
		serviceName: serviceName,
		binaryPath:  binaryPath,
	}
}

// Install creates and enables a launchd service
func (li *LaunchdInstaller) Install() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	launchdDir := filepath.Join(homeDir, "Library/LaunchAgents")
	os.MkdirAll(launchdDir, 0755)

	plistFile := filepath.Join(launchdDir, fmt.Sprintf("com.%s.plist", li.serviceName))

	tmpl, err := template.New("launchd").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.{{.ServiceName}}</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{.BinaryPath}}</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/var/log/{{.ServiceName}}.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/{{.ServiceName}}.error.log</string>
</dict>
</plist>
`)
	if err != nil {
		return err
	}

	file, err := os.Create(plistFile)
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.Execute(file, map[string]string{
		"ServiceName": li.serviceName,
		"BinaryPath":  li.binaryPath,
	})
	if err != nil {
		os.Remove(plistFile)
		return err
	}

	// Load the service
	return exec.Command("launchctl", "load", plistFile).Run()
}

// Uninstall removes the launchd service
func (li *LaunchdInstaller) Uninstall() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	plistFile := filepath.Join(homeDir, "Library/LaunchAgents", fmt.Sprintf("com.%s.plist", li.serviceName))

	if err := exec.Command("launchctl", "unload", plistFile).Run(); err != nil {
		return err
	}

	return os.Remove(plistFile)
}

// Start starts the launchd service
func (li *LaunchdInstaller) Start() error {
	return exec.Command("launchctl", "start", fmt.Sprintf("com.%s", li.serviceName)).Run()
}

// Stop stops the launchd service
func (li *LaunchdInstaller) Stop() error {
	return exec.Command("launchctl", "stop", fmt.Sprintf("com.%s", li.serviceName)).Run()
}
