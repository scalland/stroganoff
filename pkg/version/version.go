package version

import (
	"fmt"
)

// These variables are set during build time
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

// Info holds version information
type Info struct {
	Version   string
	Commit    string
	BuildDate string
}

// Get returns the version info
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
	}
}

// String returns a formatted version string
func (i Info) String() string {
	return fmt.Sprintf("Version: %s\nCommit: %s\nBuild Date: %s", i.Version, i.Commit, i.BuildDate)
}

// GetVersion returns just the version number
func GetVersion() string {
	return Version
}
