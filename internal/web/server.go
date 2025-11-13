package web

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/stroganoff/internal/config"
	"github.com/yourusername/stroganoff/pkg/auth"
	"github.com/yourusername/stroganoff/pkg/ratelimit"
)

// Server represents the HTTP server
type Server struct {
	engine      *gin.Engine
	limiter     *ratelimit.Limiter
	authenticator *auth.Authenticator
	config      *config.Config
}

// NewServer creates a new HTTP server
func NewServer() *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	server := &Server{
		engine:        engine,
		limiter:       ratelimit.NewLimiter(),
		authenticator: auth.NewAuthenticator(),
		config:        config.GetInstance().Get(),
	}

	server.setupMiddleware()
	server.setupRoutes()

	return server
}

func (s *Server) setupMiddleware() {
	// CORS middleware
	s.engine.Use(s.corsMiddleware())

	// Security headers middleware
	s.engine.Use(s.securityHeadersMiddleware())

	// Rate limiting middleware
	s.engine.Use(s.rateLimitMiddleware())

	// Authentication middleware
	s.engine.Use(s.authMiddleware())

	// Logging and recovery
	s.engine.Use(gin.Logger())
	s.engine.Use(gin.Recovery())
}

func (s *Server) setupRoutes() {
	// Health check endpoint
	s.engine.GET("/health", s.healthHandler)

	// API routes
	api := s.engine.Group("/api")
	{
		api.GET("/heartbeat", s.heartbeatHandler)
		api.GET("/metrics", s.metricsHandler)
		api.POST("/auth/token", s.createTokenHandler)
	}

	// Web interface routes
	web := s.engine.Group("")
	{
		web.GET("/", s.indexHandler)
		web.GET("/static/*filepath", s.staticFilesHandler)
	}
}

// corsMiddleware handles CORS headers
func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.GetInstance().GetAPI()

		if cfg.CORSEnabled {
			origin := c.Request.Header.Get("Origin")
			isAllowed := false

			for _, allowedOrigin := range cfg.AllowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					isAllowed = true
					break
				}
			}

			if isAllowed || len(cfg.AllowedOrigins) == 0 {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
				c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// securityHeadersMiddleware adds security headers to prevent hotlinking and scraping
func (s *Server) securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Writer.Header().Set("X-Frame-Options", "DENY")

		// Prevent MIME sniffing
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		// Content Security Policy - prevent inline scripts and restrict resource loading
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'")

		// Referrer Policy - prevent referrer leaking
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions Policy - restrict API access
		c.Writer.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}

// rateLimitMiddleware applies rate limiting
func (s *Server) rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip rate limiting for health checks
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		identifier := c.ClientIP()
		if !s.limiter.Allow(identifier) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// authMiddleware applies authentication
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.GetInstance().GetAPI()

		// Skip auth for public endpoints
		publicPaths := map[string]bool{
			"/health":            true,
			"/":                  true,
		}

		if publicPaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		if cfg.AuthEnabled {
			authHeader := c.Request.Header.Get("Authorization")
			token := auth.ExtractToken(authHeader)

			if !s.authenticator.ValidateToken(token) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized",
				})
				c.Abort()
				return
			}

			c.Set("token", token)
		}

		c.Next()
	}
}

// Handler functions
func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

func (s *Server) heartbeatHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"timestamp": time.Now().Unix(),
		"status":    "alive",
	})
}

func (s *Server) metricsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"uptime":     getUptime(),
		"goroutines": runtime.NumGoroutine(),
	})
}

func (s *Server) createTokenHandler(c *gin.Context) {
	var req struct {
		Scopes   []string `json:"scopes"`
		Duration int      `json:"duration"` // in seconds
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	duration := time.Duration(req.Duration) * time.Second
	if duration == 0 {
		duration = 24 * time.Hour
	}

	token := s.authenticator.CreateToken(req.Scopes, duration)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (s *Server) indexHandler(c *gin.Context) {
	cfg := config.GetInstance().GetServer()

	// Render HTML from embedded files
	theme := cfg.Theme
	if theme == "" {
		theme = "default"
	}

	html, err := getThemeFile(theme, "pages/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading page")
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", html)
}

func (s *Server) staticFilesHandler(c *gin.Context) {
	cfg := config.GetInstance().GetServer()
	theme := cfg.Theme
	if theme == "" {
		theme = "default"
	}

	filepath := c.Param("filepath")
	if filepath == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	data, err := getThemeFile(theme, "static"+filepath)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Data(http.StatusOK, getContentType(filepath), data)
}

// Run starts the HTTP server
func (s *Server) Run() error {
	cfg := config.GetInstance().GetServer()
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	return s.engine.Run(addr)
}

// Stop gracefully stops the server
func (s *Server) Stop() error {
	s.limiter.Stop()
	return nil
}

var startTime = time.Now()

func getUptime() int64 {
	return int64(time.Since(startTime).Seconds())
}

func getContentType(filename string) string {
	// Simplified content type detection
	if strings.HasSuffix(filename, ".css") {
		return "text/css"
	}
	if strings.HasSuffix(filename, ".js") {
		return "application/javascript"
	}
	if strings.HasSuffix(filename, ".png") {
		return "image/png"
	}
	if strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") {
		return "image/jpeg"
	}
	return "application/octet-stream"
}
