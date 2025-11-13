package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/gocr/pkg/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display version, commit hash, and build date information",
	Run: func(cmd *cobra.Command, args []string) {
		info := version.Get()
		fmt.Println(info.String())
	},
}
