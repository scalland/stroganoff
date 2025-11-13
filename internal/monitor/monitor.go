package monitor

import (
	"runtime"
	"sync"
	"time"
)

// Metrics holds application metrics
type Metrics struct {
	Uptime       int64
	Goroutines   int
	MemStats     runtime.MemStats
	RequestCount int64
	ErrorCount   int64
	LastUpdated  time.Time
}

// Monitor tracks application metrics
type Monitor struct {
	mu              sync.RWMutex
	metrics         *Metrics
	startTime       time.Time
	requestCount    int64
	errorCount      int64
	updateInterval  time.Duration
	stopCh          chan struct{}
}

// NewMonitor creates a new application monitor
func NewMonitor(updateInterval time.Duration) *Monitor {
	monitor := &Monitor{
		metrics: &Metrics{
			LastUpdated: time.Now(),
		},
		startTime:      time.Now(),
		updateInterval: updateInterval,
		stopCh:         make(chan struct{}),
	}

	go monitor.updateLoop()
	return monitor
}

// GetMetrics returns a copy of the current metrics
func (m *Monitor) GetMetrics() Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics := *m.metrics
	return metrics
}

// RecordRequest increments the request counter
func (m *Monitor) RecordRequest() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestCount++
	m.metrics.RequestCount = m.requestCount
}

// RecordError increments the error counter
func (m *Monitor) RecordError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorCount++
	m.metrics.ErrorCount = m.errorCount
}

func (m *Monitor) updateLoop() {
	ticker := time.NewTicker(m.updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopCh:
			return
		case <-ticker.C:
			m.updateMetrics()
		}
	}
}

func (m *Monitor) updateMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics.Uptime = int64(time.Since(m.startTime).Seconds())
	m.metrics.Goroutines = runtime.NumGoroutine()
	runtime.ReadMemStats(&m.metrics.MemStats)
	m.metrics.LastUpdated = time.Now()
}

// Stop stops the monitor
func (m *Monitor) Stop() {
	close(m.stopCh)
}

// HealthStatus represents the health status of the application
type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp int64                  `json:"timestamp"`
	Metrics   map[string]interface{} `json:"metrics"`
	Checks    []HealthCheck          `json:"checks"`
}

// HealthCheck represents an individual health check
type HealthCheck struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// HealthChecker performs health checks
type HealthChecker struct {
	checks map[string]func() error
	mu     sync.RWMutex
}

// NewHealthChecker creates a new health checker
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: make(map[string]func() error),
	}
}

// RegisterCheck registers a new health check
func (hc *HealthChecker) RegisterCheck(name string, check func() error) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.checks[name] = check
}

// Check runs all health checks
func (hc *HealthChecker) Check() HealthStatus {
	hc.mu.RLock()
	checks := hc.checks
	hc.mu.RUnlock()

	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().Unix(),
		Checks:    make([]HealthCheck, 0),
	}

	for name, checkFn := range checks {
		check := HealthCheck{Name: name}

		if err := checkFn(); err != nil {
			check.Status = "unhealthy"
			check.Error = err.Error()
			status.Status = "degraded"
		} else {
			check.Status = "healthy"
		}

		status.Checks = append(status.Checks, check)
	}

	return status
}
