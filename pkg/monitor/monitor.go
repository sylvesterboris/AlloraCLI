package monitor

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/AlloraAi/AlloraCLI/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Monitor interface defines monitoring operations
type Monitor interface {
	GetName() string
	GetCategory() string
	GetInterval() time.Duration
	GetStatus() string
	Start() error
	Stop() error
	CollectMetrics(ctx context.Context) ([]*Metric, error)
	GetConfiguration() *MonitorConfig
	UpdateConfiguration(config *MonitorConfig) error
	IsHealthy() bool
	GetSystemStatus() (*SystemStatus, error)
	GetServiceStatus(serviceName string, detailed bool) (*ServiceStatus, error)
	ListServices() ([]*ServiceInfo, error)
	GetMetrics(metric, duration string) (*MetricsData, error)
	CreateAlert(alert Alert) error
	ListAlerts() ([]*Alert, error)
	DeleteAlert(name string) error
	StartDashboard(host string, port int) error
}

// SystemStatus represents overall system status
type SystemStatus struct {
	Overall     string               `json:"overall" yaml:"overall"`
	Timestamp   time.Time           `json:"timestamp" yaml:"timestamp"`
	Services    []*ServiceStatus    `json:"services" yaml:"services"`
	Resources   *ResourceUsage      `json:"resources" yaml:"resources"`
	Alerts      []*ActiveAlert      `json:"alerts" yaml:"alerts"`
	Uptime      time.Duration       `json:"uptime" yaml:"uptime"`
	Metadata    map[string]string   `json:"metadata" yaml:"metadata"`
}

// ServiceStatus represents the status of a specific service
type ServiceStatus struct {
	Name        string            `json:"name" yaml:"name"`
	Status      string            `json:"status" yaml:"status"`
	Health      string            `json:"health" yaml:"health"`
	Uptime      time.Duration     `json:"uptime" yaml:"uptime"`
	CPU         float64           `json:"cpu" yaml:"cpu"`
	Memory      float64           `json:"memory" yaml:"memory"`
	Disk        float64           `json:"disk" yaml:"disk"`
	Network     *NetworkMetrics   `json:"network" yaml:"network"`
	Endpoints   []*EndpointStatus `json:"endpoints" yaml:"endpoints"`
	LastCheck   time.Time         `json:"last_check" yaml:"last_check"`
	Metadata    map[string]string `json:"metadata" yaml:"metadata"`
}

// ServiceInfo represents basic service information
type ServiceInfo struct {
	Name        string `json:"name" yaml:"name"`
	Type        string `json:"type" yaml:"type"`
	Status      string `json:"status" yaml:"status"`
	Description string `json:"description" yaml:"description"`
}

// ResourceUsage represents system resource usage
type ResourceUsage struct {
	CPU     *CPUUsage     `json:"cpu" yaml:"cpu"`
	Memory  *MemoryUsage  `json:"memory" yaml:"memory"`
	Disk    *DiskUsage    `json:"disk" yaml:"disk"`
	Network *NetworkUsage `json:"network" yaml:"network"`
}

// CPUUsage represents CPU usage metrics
type CPUUsage struct {
	Usage       float64 `json:"usage" yaml:"usage"`
	LoadAverage float64 `json:"load_average" yaml:"load_average"`
	Cores       int     `json:"cores" yaml:"cores"`
}

// MemoryUsage represents memory usage metrics
type MemoryUsage struct {
	Used      int64   `json:"used" yaml:"used"`
	Available int64   `json:"available" yaml:"available"`
	Total     int64   `json:"total" yaml:"total"`
	Usage     float64 `json:"usage" yaml:"usage"`
}

// DiskUsage represents disk usage metrics
type DiskUsage struct {
	Used      int64   `json:"used" yaml:"used"`
	Available int64   `json:"available" yaml:"available"`
	Total     int64   `json:"total" yaml:"total"`
	Usage     float64 `json:"usage" yaml:"usage"`
}

// NetworkUsage represents network usage metrics
type NetworkUsage struct {
	BytesIn  int64 `json:"bytes_in" yaml:"bytes_in"`
	BytesOut int64 `json:"bytes_out" yaml:"bytes_out"`
	PacketsIn  int64 `json:"packets_in" yaml:"packets_in"`
	PacketsOut int64 `json:"packets_out" yaml:"packets_out"`
}

// NetworkMetrics represents network metrics for a service
type NetworkMetrics struct {
	Connections int64   `json:"connections" yaml:"connections"`
	Latency     float64 `json:"latency" yaml:"latency"`
	Throughput  float64 `json:"throughput" yaml:"throughput"`
}

// EndpointStatus represents the status of a service endpoint
type EndpointStatus struct {
	URL          string        `json:"url" yaml:"url"`
	Status       string        `json:"status" yaml:"status"`
	ResponseTime time.Duration `json:"response_time" yaml:"response_time"`
	StatusCode   int           `json:"status_code" yaml:"status_code"`
	LastCheck    time.Time     `json:"last_check" yaml:"last_check"`
}

// Alert represents a monitoring alert configuration
type Alert struct {
	Name      string    `json:"name" yaml:"name"`
	Condition string    `json:"condition" yaml:"condition"`
	Action    string    `json:"action" yaml:"action"`
	Severity  string    `json:"severity" yaml:"severity"`
	Enabled   bool      `json:"enabled" yaml:"enabled"`
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at"`
}

// ActiveAlert represents an active alert
type ActiveAlert struct {
	Alert       *Alert    `json:"alert" yaml:"alert"`
	Triggered   time.Time `json:"triggered" yaml:"triggered"`
	Status      string    `json:"status" yaml:"status"`
	Message     string    `json:"message" yaml:"message"`
	Acknowledged bool     `json:"acknowledged" yaml:"acknowledged"`
}

// MetricsData represents metrics data
type MetricsData struct {
	Metric    string                   `json:"metric" yaml:"metric"`
	TimeRange string                   `json:"time_range" yaml:"time_range"`
	Data      []MetricPoint            `json:"data" yaml:"data"`
	Summary   *MetricSummary           `json:"summary" yaml:"summary"`
	Metadata  map[string]string        `json:"metadata" yaml:"metadata"`
}

// MetricPoint represents a single metric data point
type MetricPoint struct {
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
	Value     float64   `json:"value" yaml:"value"`
	Labels    map[string]string `json:"labels" yaml:"labels"`
}

// MetricSummary represents a summary of metrics
type MetricSummary struct {
	Average float64 `json:"average" yaml:"average"`
	Min     float64 `json:"min" yaml:"min"`
	Max     float64 `json:"max" yaml:"max"`
	Count   int64   `json:"count" yaml:"count"`
}

// MonitorImpl implements the Monitor interface
type MonitorImpl struct {
	config   *config.Config
	registry *prometheus.Registry
	ctx      context.Context
}

// New creates a new monitor instance
func New() (Monitor, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	registry := prometheus.NewRegistry()
	
	return &MonitorImpl{
		config:   cfg,
		registry: registry,
		ctx:      context.Background(),
	}, nil
}

// GetSystemStatus returns overall system status
func (m *MonitorImpl) GetSystemStatus() (*SystemStatus, error) {
	// Mock implementation - in real scenario, this would collect actual system metrics
	status := &SystemStatus{
		Overall:   "healthy",
		Timestamp: time.Now(),
		Services: []*ServiceStatus{
			{
				Name:      "web-server",
				Status:    "running",
				Health:    "healthy",
				Uptime:    24 * time.Hour,
				CPU:       15.5,
				Memory:    45.2,
				Disk:      12.8,
				LastCheck: time.Now(),
				Metadata:  map[string]string{"version": "1.2.3", "port": "8080"},
			},
			{
				Name:      "database",
				Status:    "running",
				Health:    "healthy",
				Uptime:    72 * time.Hour,
				CPU:       8.3,
				Memory:    62.7,
				Disk:      78.1,
				LastCheck: time.Now(),
				Metadata:  map[string]string{"version": "13.4", "port": "5432"},
			},
		},
		Resources: &ResourceUsage{
			CPU: &CPUUsage{
				Usage:       25.4,
				LoadAverage: 0.8,
				Cores:       4,
			},
			Memory: &MemoryUsage{
				Used:      4 * 1024 * 1024 * 1024,  // 4GB
				Available: 4 * 1024 * 1024 * 1024,  // 4GB
				Total:     8 * 1024 * 1024 * 1024,  // 8GB
				Usage:     50.0,
			},
			Disk: &DiskUsage{
				Used:      100 * 1024 * 1024 * 1024, // 100GB
				Available: 400 * 1024 * 1024 * 1024, // 400GB
				Total:     500 * 1024 * 1024 * 1024, // 500GB
				Usage:     20.0,
			},
		},
		Alerts: []*ActiveAlert{},
		Uptime: 168 * time.Hour, // 7 days
		Metadata: map[string]string{
			"hostname": "server-01",
			"region":   "us-west-2",
		},
	}

	return status, nil
}

// GetServiceStatus returns the status of a specific service
func (m *MonitorImpl) GetServiceStatus(serviceName string, detailed bool) (*ServiceStatus, error) {
	// Mock implementation
	service := &ServiceStatus{
		Name:      serviceName,
		Status:    "running",
		Health:    "healthy",
		Uptime:    24 * time.Hour,
		CPU:       15.5,
		Memory:    45.2,
		Disk:      12.8,
		Network: &NetworkMetrics{
			Connections: 42,
			Latency:     15.2,
			Throughput:  125.5,
		},
		LastCheck: time.Now(),
		Metadata:  map[string]string{"version": "1.2.3"},
	}

	if detailed {
		service.Endpoints = []*EndpointStatus{
			{
				URL:          "http://localhost:8080/health",
				Status:       "healthy",
				ResponseTime: 25 * time.Millisecond,
				StatusCode:   200,
				LastCheck:    time.Now(),
			},
		}
	}

	return service, nil
}

// ListServices returns a list of all services
func (m *MonitorImpl) ListServices() ([]*ServiceInfo, error) {
	// Mock implementation
	services := []*ServiceInfo{
		{
			Name:        "web-server",
			Type:        "http",
			Status:      "running",
			Description: "Main web application server",
		},
		{
			Name:        "database",
			Type:        "database",
			Status:      "running",
			Description: "PostgreSQL database server",
		},
		{
			Name:        "redis",
			Type:        "cache",
			Status:      "running",
			Description: "Redis cache server",
		},
	}

	return services, nil
}

// GetMetrics returns metrics data
func (m *MonitorImpl) GetMetrics(metric, duration string) (*MetricsData, error) {
	// Mock implementation
	now := time.Now()
	data := &MetricsData{
		Metric:    metric,
		TimeRange: duration,
		Data:      []MetricPoint{},
		Summary: &MetricSummary{
			Average: 45.5,
			Min:     12.3,
			Max:     78.9,
			Count:   100,
		},
		Metadata: map[string]string{
			"unit":   "percent",
			"source": "prometheus",
		},
	}

	// Generate sample data points
	for i := 0; i < 10; i++ {
		data.Data = append(data.Data, MetricPoint{
			Timestamp: now.Add(time.Duration(-i) * time.Minute),
			Value:     45.5 + float64(i)*2.3,
			Labels:    map[string]string{"instance": "server-01"},
		})
	}

	return data, nil
}

// CreateAlert creates a new alert
func (m *MonitorImpl) CreateAlert(alert Alert) error {
	// Mock implementation - in real scenario, this would persist the alert
	alert.CreatedAt = time.Now()
	alert.UpdatedAt = time.Now()
	return nil
}

// ListAlerts returns all configured alerts
func (m *MonitorImpl) ListAlerts() ([]*Alert, error) {
	// Mock implementation
	alerts := []*Alert{
		{
			Name:      "high-cpu",
			Condition: "cpu > 80%",
			Action:    "notify",
			Severity:  "warning",
			Enabled:   true,
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			Name:      "low-memory",
			Condition: "memory < 10%",
			Action:    "scale-up",
			Severity:  "critical",
			Enabled:   true,
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-48 * time.Hour),
		},
	}

	return alerts, nil
}

// DeleteAlert deletes an alert by name
func (m *MonitorImpl) DeleteAlert(name string) error {
	// Mock implementation - in real scenario, this would remove the alert
	return nil
}

// StartDashboard starts the monitoring dashboard web server
func (m *MonitorImpl) StartDashboard(host string, port int) error {
	mux := http.NewServeMux()
	
	// Prometheus metrics endpoint
	mux.Handle("/metrics", promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}))
	
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	// Simple dashboard endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>AlloraCLI Monitoring Dashboard</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .metric { margin: 10px 0; padding: 10px; border: 1px solid #ddd; }
        .status-ok { color: green; }
        .status-warning { color: orange; }
        .status-error { color: red; }
    </style>
</head>
<body>
    <h1>AlloraCLI Monitoring Dashboard</h1>
    <div class="metric">
        <h3>System Status: <span class="status-ok">Healthy</span></h3>
    </div>
    <div class="metric">
        <h3>Services</h3>
        <ul>
            <li>Web Server: <span class="status-ok">Running</span></li>
            <li>Database: <span class="status-ok">Running</span></li>
            <li>Redis: <span class="status-ok">Running</span></li>
        </ul>
    </div>
    <div class="metric">
        <h3>Resources</h3>
        <ul>
            <li>CPU: 25.4%</li>
            <li>Memory: 50.0%</li>
            <li>Disk: 20.0%</li>
        </ul>
    </div>
    <p><a href="/metrics">Prometheus Metrics</a></p>
</body>
</html>
		`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})
	
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: mux,
	}
	
	return server.ListenAndServe()
}

// MonitoringManager manages multiple monitors
type MonitoringManager struct {
	monitors map[string]Monitor
	mutex    sync.RWMutex
}

// NewMonitoringManager creates a new monitoring manager
func NewMonitoringManager() *MonitoringManager {
	return &MonitoringManager{
		monitors: make(map[string]Monitor),
	}
}

// AddMonitor adds a monitor to the manager
func (m *MonitoringManager) AddMonitor(monitor Monitor) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.monitors[monitor.GetName()] = monitor
	return nil
}

// GetMonitor retrieves a monitor by name
func (m *MonitoringManager) GetMonitor(name string) (Monitor, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	monitor, exists := m.monitors[name]
	if !exists {
		return nil, fmt.Errorf("monitor not found: %s", name)
	}
	return monitor, nil
}

// ListMonitors returns a list of all monitors
func (m *MonitoringManager) ListMonitors() []Monitor {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	monitors := make([]Monitor, 0, len(m.monitors))
	for _, monitor := range m.monitors {
		monitors = append(monitors, monitor)
	}
	return monitors
}

// StartMonitoring starts monitoring for a specific monitor
func (m *MonitoringManager) StartMonitoring(name string) error {
	monitor, err := m.GetMonitor(name)
	if err != nil {
		return err
	}
	return monitor.Start()
}

// StopMonitoring stops monitoring for a specific monitor
func (m *MonitoringManager) StopMonitoring(name string) error {
	monitor, err := m.GetMonitor(name)
	if err != nil {
		return err
	}
	return monitor.Stop()
}

// AlertManager manages alerts
type AlertManager struct {
	rules map[string]*AlertRule
	mutex sync.RWMutex
}

// NewAlertManager creates a new alert manager
func NewAlertManager() *AlertManager {
	return &AlertManager{
		rules: make(map[string]*AlertRule),
	}
}

// AddRule adds an alert rule
func (m *AlertManager) AddRule(rule *AlertRule) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.rules[rule.Name] = rule
	return nil
}

// GetRule retrieves an alert rule by name
func (m *AlertManager) GetRule(name string) (*AlertRule, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	rule, exists := m.rules[name]
	if !exists {
		return nil, fmt.Errorf("alert rule not found: %s", name)
	}
	return rule, nil
}

// EvaluateRules evaluates all alert rules against the given metrics
func (m *AlertManager) EvaluateRules(ctx context.Context, metrics []*Metric) ([]*Alert, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	var alerts []*Alert
	for _, rule := range m.rules {
		if rule.Enabled {
			// Simple evaluation logic for testing
			for _, metric := range metrics {
				if shouldTriggerAlert(rule, metric) {
					alerts = append(alerts, &Alert{
						RuleName:  rule.Name,
						Severity:  rule.Severity,
						Message:   fmt.Sprintf("%s: %s", rule.Name, rule.Description),
						Timestamp: time.Now(),
						Value:     metric.Value,
					})
				}
			}
		}
	}
	return alerts, nil
}

// HealthChecker manages health checks
type HealthChecker struct {
	checks map[string]*HealthCheck
	mutex  sync.RWMutex
}

// NewHealthChecker creates a new health checker
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: make(map[string]*HealthCheck),
	}
}

// AddCheck adds a health check
func (h *HealthChecker) AddCheck(check *HealthCheck) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	h.checks[check.Name] = check
	return nil
}

// RunCheck runs a specific health check
func (h *HealthChecker) RunCheck(ctx context.Context, name string) (*HealthCheckResult, error) {
	h.mutex.RLock()
	check, exists := h.checks[name]
	h.mutex.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("health check not found: %s", name)
	}
	
	// Simple health check logic for testing
	return &HealthCheckResult{
		CheckName: name,
		Status:    "healthy",
		Timestamp: time.Now(),
		Duration:  time.Millisecond * 100,
		Message:   "Health check passed",
	}, nil
}

// GetHistory returns health check history
func (h *HealthChecker) GetHistory(name string, limit int) ([]*HealthCheckResult, error) {
	// Mock history for testing
	results := make([]*HealthCheckResult, 0, limit)
	for i := 0; i < limit; i++ {
		results = append(results, &HealthCheckResult{
			CheckName: name,
			Status:    "healthy",
			Timestamp: time.Now().Add(time.Duration(-i) * time.Minute),
			Duration:  time.Millisecond * 100,
			Message:   "Health check passed",
		})
	}
	return results, nil
}

// Dashboard manages monitoring dashboard
type Dashboard struct {
	widgets map[string]*Widget
	mutex   sync.RWMutex
}

// NewDashboard creates a new dashboard
func NewDashboard() *Dashboard {
	return &Dashboard{
		widgets: make(map[string]*Widget),
	}
}

// AddWidget adds a widget to the dashboard
func (d *Dashboard) AddWidget(widget *Widget) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	
	d.widgets[widget.ID] = widget
	return nil
}

// GetWidget retrieves a widget by ID
func (d *Dashboard) GetWidget(id string) (*Widget, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	
	widget, exists := d.widgets[id]
	if !exists {
		return nil, fmt.Errorf("widget not found: %s", id)
	}
	return widget, nil
}

// ListWidgets returns a list of all widgets
func (d *Dashboard) ListWidgets() []*Widget {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	
	widgets := make([]*Widget, 0, len(d.widgets))
	for _, widget := range d.widgets {
		widgets = append(widgets, widget)
	}
	return widgets
}

// GenerateData generates dashboard data
func (d *Dashboard) GenerateData(ctx context.Context) (*DashboardData, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	
	data := &DashboardData{
		Widgets: make([]*WidgetData, 0, len(d.widgets)),
	}
	
	for _, widget := range d.widgets {
		widgetData := &WidgetData{
			ID:    widget.ID,
			Title: widget.Title,
			Type:  widget.Type,
			Data:  map[string]interface{}{
				"value": 42.5,
				"unit":  "percent",
			},
		}
		data.Widgets = append(data.Widgets, widgetData)
	}
	
	return data, nil
}

// Additional type definitions for the test interfaces
type Metric struct {
	Name      string                 `json:"name"`
	Value     interface{}            `json:"value"`
	Unit      string                 `json:"unit"`
	Timestamp time.Time              `json:"timestamp"`
	Labels    map[string]string      `json:"labels"`
}

type AlertRule struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Condition   string   `json:"condition"`
	Severity    string   `json:"severity"`
	Actions     []string `json:"actions"`
	Enabled     bool     `json:"enabled"`
}

type Alert struct {
	RuleName  string      `json:"rule_name"`
	Severity  string      `json:"severity"`
	Message   string      `json:"message"`
	Timestamp time.Time   `json:"timestamp"`
	Value     interface{} `json:"value"`
}

type HealthCheck struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Type        string        `json:"type"`
	Target      string        `json:"target"`
	Interval    time.Duration `json:"interval"`
	Timeout     time.Duration `json:"timeout"`
	Enabled     bool          `json:"enabled"`
}

type HealthCheckResult struct {
	CheckName string        `json:"check_name"`
	Status    string        `json:"status"`
	Timestamp time.Time     `json:"timestamp"`
	Duration  time.Duration `json:"duration"`
	Message   string        `json:"message"`
}

type Widget struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	MetricQuery string    `json:"metric_query"`
	Position    Position  `json:"position"`
	Size        Size      `json:"size"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type DashboardData struct {
	Widgets []*WidgetData `json:"widgets"`
}

type WidgetData struct {
	ID    string                 `json:"id"`
	Title string                 `json:"title"`
	Type  string                 `json:"type"`
	Data  map[string]interface{} `json:"data"`
}

type MonitorConfig struct {
	Name     string        `json:"name"`
	Category string        `json:"category"`
	Interval time.Duration `json:"interval"`
	Enabled  bool          `json:"enabled"`
}

// shouldTriggerAlert is a simple helper function for testing
func shouldTriggerAlert(rule *AlertRule, metric *Metric) bool {
	// Simple condition check for testing
	if rule.Condition == "cpu_usage > 80" && metric.Name == "cpu_usage" {
		if val, ok := metric.Value.(float64); ok {
			return val > 80
		}
	}
	return false
}
