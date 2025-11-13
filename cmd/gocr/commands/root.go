package commands

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "stroganoff",
	Short: "A professional Go application with advanced features",
	Long: `stroganoff is a feature-rich Go CLI application with support for
configuration management, web interfaces, API endpoints, and more.`,
}

func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(upgradeCmd)
	RootCmd.AddCommand(installCmd)
	RootCmd.AddCommand(serveCmd)
	RootCmd.AddCommand(configCmd)
}
