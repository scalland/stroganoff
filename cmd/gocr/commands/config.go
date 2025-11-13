package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/gocr/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "Show and manage application configuration",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.GetInstance().Get()

		fmt.Println("Server Configuration:")
		fmt.Printf("  Host: %s\n", cfg.Server.Host)
		fmt.Printf("  Port: %d\n", cfg.Server.Port)
		fmt.Printf("  Theme: %s\n", cfg.Server.Theme)

		fmt.Println("\nAPI Configuration:")
		fmt.Printf("  Rate Limit: %d requests\n", cfg.API.RateLimit)
		fmt.Printf("  Rate Limit Window: %d seconds\n", cfg.API.RateLimitWindow)
		fmt.Printf("  Auth Enabled: %v\n", cfg.API.AuthEnabled)

		fmt.Println("\nDatabase Configuration:")
		fmt.Printf("  Host: %s\n", cfg.Database.Host)
		fmt.Printf("  Port: %d\n", cfg.Database.Port)
		fmt.Printf("  Database: %s\n", cfg.Database.Database)

		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
}
