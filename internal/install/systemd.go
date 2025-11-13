package install

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

// SystemdInstaller handles systemd service installation on Linux
type SystemdInstaller struct {
	serviceName string
	binaryPath  string
	user        string
}

// NewSystemdInstaller creates a new systemd installer
func NewSystemdInstaller(serviceName, binaryPath, user string) *SystemdInstaller {
	if user == "" {
		user = "root"
	}
	return &SystemdInstaller{
		serviceName: serviceName,
		binaryPath:  binaryPath,
		user:        user,
	}
}

// Install creates and enables a systemd service
func (si *SystemdInstaller) Install() error {
	serviceFile := fmt.Sprintf("/etc/systemd/system/%s.service", si.serviceName)

	tmpl, err := template.New("systemd").Parse(`[Unit]
Description={{.ServiceName}} Service
After=network.target

[Service]
Type=simple
User={{.User}}
ExecStart={{.BinaryPath}}
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
`)
	if err != nil {
		return err
	}

	file, err := os.Create(serviceFile)
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.Execute(file, map[string]string{
		"ServiceName": si.serviceName,
		"User":        si.user,
		"BinaryPath":  si.binaryPath,
	})
	if err != nil {
		os.Remove(serviceFile)
		return err
	}

	// Reload systemd and enable service
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		os.Remove(serviceFile)
		return err
	}

	if err := exec.Command("systemctl", "enable", si.serviceName).Run(); err != nil {
		os.Remove(serviceFile)
		return err
	}

	return nil
}

// Uninstall removes the systemd service
func (si *SystemdInstaller) Uninstall() error {
	serviceFile := fmt.Sprintf("/etc/systemd/system/%s.service", si.serviceName)

	if err := exec.Command("systemctl", "disable", si.serviceName).Run(); err != nil {
		return err
	}

	if err := os.Remove(serviceFile); err != nil {
		return err
	}

	return exec.Command("systemctl", "daemon-reload").Run()
}

// Start starts the systemd service
func (si *SystemdInstaller) Start() error {
	return exec.Command("systemctl", "start", si.serviceName).Run()
}

// Stop stops the systemd service
func (si *SystemdInstaller) Stop() error {
	return exec.Command("systemctl", "stop", si.serviceName).Run()
}
