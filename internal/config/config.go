package config

import (
	"sync"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	API      APIConfig      `yaml:"api"`
	Database DatabaseConfig `yaml:"database"`
	Logging  LoggingConfig  `yaml:"logging"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Theme    string `yaml:"theme"`
	TLSCert  string `yaml:"tls_cert"`
	TLSKey   string `yaml:"tls_key"`
	ReadTimeout  int `yaml:"read_timeout"`
	WriteTimeout int `yaml:"write_timeout"`
}

// APIConfig holds API configuration
type APIConfig struct {
	RateLimit        int    `yaml:"rate_limit"`
	RateLimitWindow  int    `yaml:"rate_limit_window"`
	AuthEnabled      bool   `yaml:"auth_enabled"`
	AuthTokenHeader  string `yaml:"auth_token_header"`
	AllowedOrigins   []string `yaml:"allowed_origins"`
	CORSEnabled      bool   `yaml:"cors_enabled"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	OutputPath string `yaml:"output_path"`
}

// ConfigManager manages the configuration with singleton pattern
type ConfigManager struct {
	config   *Config
	mu       sync.RWMutex
	watchers []func(*Config)
}

var (
	instance *ConfigManager
	once     sync.Once
)

// GetInstance returns the singleton config manager instance
func GetInstance() *ConfigManager {
	once.Do(func() {
		instance = &ConfigManager{
			config:   &Config{},
			watchers: make([]func(*Config), 0),
		}
	})
	return instance
}

// Load loads configuration from YAML bytes
func (cm *ConfigManager) Load(data []byte) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return err
	}

	cm.config = cfg
	cm.notifyWatchers()
	return nil
}

// Get returns a copy of the current configuration
func (cm *ConfigManager) Get() *Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Return a copy to prevent external modifications
	configCopy := *cm.config
	return &configCopy
}

// GetServer returns the server configuration
func (cm *ConfigManager) GetServer() ServerConfig {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.Server
}

// GetAPI returns the API configuration
func (cm *ConfigManager) GetAPI() APIConfig {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.API
}

// Watch registers a watcher function that gets called on config changes
func (cm *ConfigManager) Watch(watcher func(*Config)) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.watchers = append(cm.watchers, watcher)
}

func (cm *ConfigManager) notifyWatchers() {
	for _, watcher := range cm.watchers {
		go watcher(cm.config)
	}
}
