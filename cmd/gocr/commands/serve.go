package commands

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/gocr/internal/config"
	"github.com/yourusername/gocr/internal/monitor"
	"github.com/yourusername/gocr/internal/web"
)

var (
	configFile string
	webHost    string
	webPort    int
	webTheme   string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Long:  "Start the GOCR server with HTTP API and web interface",
	RunE: func(cmd *cobra.Command, args []string) error {
		return startServer()
	},
}

func init() {
	serveCmd.Flags().StringVar(&configFile, "config", "config.yaml", "Configuration file path")
	serveCmd.Flags().StringVar(&webHost, "host", "localhost", "Server host")
	serveCmd.Flags().IntVar(&webPort, "port", 8080, "Server port")
	serveCmd.Flags().StringVar(&webTheme, "theme", "default", "Theme name")
}

func startServer() error {
	// Load configuration
	loader, err := config.NewLoader(configFile)
	if err != nil && configFile != "config.yaml" {
		return fmt.Errorf("failed to create config loader: %w", err)
	}

	// Try to load config file, but continue if it doesn't exist
	if err := loader.Load(); err != nil {
		fmt.Printf("Warning: Could not load config file: %v\n", err)
	}

	// Start watching for config changes
	if err := loader.StartWatching(); err != nil {
		fmt.Printf("Warning: Could not watch config file: %v\n", err)
	}
	defer loader.Stop()

	// Override config from command line flags
	cfg := config.GetInstance().Get()
	if webHost != "localhost" {
		cfg.Server.Host = webHost
	}
	if webPort != 8080 {
		cfg.Server.Port = webPort
	}
	if webTheme != "default" {
		cfg.Server.Theme = webTheme
	}

	// Update config
	if err := config.GetInstance().Load(marshalConfig(cfg)); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	// Initialize monitor
	appMonitor := monitor.NewMonitor(10 * time.Second)
	defer appMonitor.Stop()

	// Create and start server
	fmt.Printf("Starting GOCR server on %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("Theme: %s\n", cfg.Server.Theme)

	server := web.NewServer()

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		fmt.Println("\nShutting down server...")
		server.Stop()
		os.Exit(0)
	}()

	return server.Run()
}

func marshalConfig(cfg *config.Config) []byte {
	// Simple YAML marshaling (would be better with proper YAML library)
	data := fmt.Sprintf(`server:
  host: %s
  port: %d
  theme: %s
api:
  rate_limit: %d
  rate_limit_window: %d
  auth_enabled: %v
`, cfg.Server.Host, cfg.Server.Port, cfg.Server.Theme,
		cfg.API.RateLimit, cfg.API.RateLimitWindow, cfg.API.AuthEnabled)
	return []byte(data)
}
