package web

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
)

// NOTE: Embed directive commented out - can be enabled when building from project root
// //go:embed web/themes
// var themeFS embed.FS

var themeFS embed.FS

// getThemeFile retrieves a file from the embedded theme filesystem
func getThemeFile(theme, filename string) ([]byte, error) {
	// Validate theme name to prevent directory traversal
	if theme == "" || theme == "." || filepath.IsAbs(theme) {
		return nil, fmt.Errorf("invalid theme")
	}

	// For now, return a default page since embed is relative to build directory
	defaultHTML := []byte(`<!DOCTYPE html>
<html>
<head>
  <title>GOCR</title>
  <style>
    body { font-family: Arial, sans-serif; margin: 40px; }
    .container { max-width: 1200px; margin: 0 auto; }
    .navbar { background: #f0f0f0; padding: 1rem; margin-bottom: 2rem; border-radius: 4px; }
    .navbar h1 { margin: 0; color: #0066cc; }
    .hero { text-align: center; margin: 3rem 0; }
    .hero h2 { font-size: 2.5rem; color: #0066cc; }
    .features { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 1.5rem; margin: 2rem 0; }
    .feature { padding: 1rem; background: #f9f9f9; border-radius: 4px; border-left: 4px solid #0066cc; }
    .endpoint { margin: 1.5rem 0; padding: 1rem; background: #f0f0f0; border-radius: 4px; }
    .endpoint code { display: block; background: #e0e0e0; padding: 0.5rem; margin-bottom: 0.5rem; border-radius: 3px; }
    footer { margin-top: 3rem; padding-top: 1rem; border-top: 1px solid #ddd; text-align: center; color: #666; }
  </style>
</head>
<body>
  <div class="container">
    <nav class="navbar">
      <h1>GOCR</h1>
    </nav>

    <main>
      <section class="hero">
        <h2>Welcome to GOCR</h2>
        <p>A professional Go CLI application with advanced features</p>
      </section>

      <section class="features">
        <div class="feature">✓ Multi-OS/Architecture Support</div>
        <div class="feature">✓ Service Installation & Management</div>
        <div class="feature">✓ Configuration Management with Hot-Reload</div>
        <div class="feature">✓ API Authentication & Rate Limiting</div>
        <div class="feature">✓ Theme Support</div>
        <div class="feature">✓ Monitoring & Health Checks</div>
        <div class="feature">✓ Automatic Updates from Github</div>
      </section>

      <section>
        <h3>API Endpoints</h3>
        <div class="endpoint">
          <code>GET /health</code>
          <p>Check application health status</p>
        </div>
        <div class="endpoint">
          <code>GET /api/heartbeat</code>
          <p>Get server heartbeat</p>
        </div>
        <div class="endpoint">
          <code>GET /api/metrics</code>
          <p>Get application metrics (requires auth)</p>
        </div>
        <div class="endpoint">
          <code>POST /api/auth/token</code>
          <p>Create authentication token</p>
        </div>
      </section>
    </main>

    <footer>
      <p>&copy; 2024 GOCR. All rights reserved.</p>
    </footer>
  </div>
</body>
</html>`)

	return defaultHTML, nil
}

// isPathSafe checks if a path is safe from directory traversal attacks
func isPathSafe(path, baseDir string) bool {
	absPath := filepath.Join(baseDir, path)
	baseDirAbs := filepath.Join(baseDir)

	// Check if the absolute path starts with the base directory
	if len(absPath) < len(baseDirAbs) {
		return false
	}

	return absPath[:len(baseDirAbs)] == baseDirAbs
}

// GetAvailableThemes returns a list of available themes
func GetAvailableThemes() ([]string, error) {
	var themes []string

	themeDir, err := fs.ReadDir(themeFS, "web/themes")
	if err != nil {
		return nil, err
	}

	for _, entry := range themeDir {
		if entry.IsDir() {
			themes = append(themes, entry.Name())
		}
	}

	return themes, nil
}
