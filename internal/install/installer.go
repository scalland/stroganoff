package install

// ServiceInstaller defines the interface for service installation
type ServiceInstaller interface {
	Install() error
	Uninstall() error
	Start() error
	Stop() error
}
