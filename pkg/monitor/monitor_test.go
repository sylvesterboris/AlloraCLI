package monitor

import (
	"context"
	"testing"
	"time"
)

func TestMonitoringManager(t *testing.T) {
	manager := NewMonitoringManager()

	// Test adding a monitor
	monitor := &MockMonitor{
		name:     "cpu-monitor",
		category: "system",
		interval: 30 * time.Second,
		status:   "running",
	}

	err := manager.AddMonitor(monitor)
	if err != nil {
		t.Errorf("AddMonitor() failed: %v", err)
	}

	// Test getting a monitor
	retrievedMonitor, err := manager.GetMonitor("cpu-monitor")
	if err != nil {
		t.Errorf("GetMonitor() failed: %v", err)
	}

	if retrievedMonitor.GetName() != "cpu-monitor" {
		t.Errorf("Expected monitor name 'cpu-monitor', got '%s'", retrievedMonitor.GetName())
	}

	// Test listing monitors
	monitors := manager.ListMonitors()
	if len(monitors) != 1 {
		t.Errorf("Expected 1 monitor, got %d", len(monitors))
	}

	// Test starting monitoring
	err = manager.StartMonitoring("cpu-monitor")
	if err != nil {
		t.Errorf("StartMonitoring() failed: %v", err)
	}

	// Test stopping monitoring
	err = manager.StopMonitoring("cpu-monitor")
	if err != nil {
		t.Errorf("StopMonitoring() failed: %v", err)
	}
}

func TestMetricsCollection(t *testing.T) {
	monitor := &MockMonitor{
		name:     "memory-monitor",
		category: "system",
		interval: 10 * time.Second,
		status:   "running",
	}

	ctx := context.Background()

	// Test collecting metrics
	metrics, err := monitor.CollectMetrics(ctx)
	if err != nil {
		t.Errorf("CollectMetrics() failed: %v", err)
	}

	if len(metrics) == 0 {
		t.Error("Expected non-empty metrics collection")
	}

	// Verify metric structure
	for _, metric := range metrics {
		if metric.Name == "" {
			t.Error("Expected non-empty metric name")
		}
		if metric.Timestamp.IsZero() {
			t.Error("Expected non-zero metric timestamp")
		}
		if metric.Value == nil {
			t.Error("Expected non-nil metric value")
		}
	}
}

func TestAlerting(t *testing.T) {
	alertManager := NewAlertManager()

	// Test adding an alert rule
	rule := &AlertRule{
		Name:        "high-cpu-usage",
		Description: "CPU usage above 80%",
		Condition:   "cpu_usage > 80",
		Severity:    "warning",
		Actions:     []string{"email", "slack"},
		Enabled:     true,
	}

	err := alertManager.AddRule(rule)
	if err != nil {
		t.Errorf("AddRule() failed: %v", err)
	}

	// Test getting an alert rule
	retrievedRule, err := alertManager.GetRule("high-cpu-usage")
	if err != nil {
		t.Errorf("GetRule() failed: %v", err)
	}

	if retrievedRule.Name != "high-cpu-usage" {
		t.Errorf("Expected rule name 'high-cpu-usage', got '%s'", retrievedRule.Name)
	}

	// Test evaluating alerts
	ctx := context.Background()
	metrics := []*Metric{
		{
			Name:      "cpu_usage",
			Value:     85.5,
			Timestamp: time.Now(),
		},
	}

	alerts, err := alertManager.EvaluateRules(ctx, metrics)
	if err != nil {
		t.Errorf("EvaluateRules() failed: %v", err)
	}

	if len(alerts) == 0 {
		t.Error("Expected at least one alert to be triggered")
	}

	// Verify alert structure
	for _, alert := range alerts {
		if alert.RuleName == "" {
			t.Error("Expected non-empty alert rule name")
		}
		if alert.Severity == "" {
			t.Error("Expected non-empty alert severity")
		}
		if alert.Timestamp.IsZero() {
			t.Error("Expected non-zero alert timestamp")
		}
	}
}

func TestHealthCheck(t *testing.T) {
	healthChecker := NewHealthChecker()

	// Test adding a health check
	check := &HealthCheck{
		Name:        "database-check",
		Description: "Check database connectivity",
		Type:        "tcp",
		Target:      "localhost:5432",
		Interval:    30 * time.Second,
		Timeout:     5 * time.Second,
		Enabled:     true,
	}

	err := healthChecker.AddCheck(check)
	if err != nil {
		t.Errorf("AddCheck() failed: %v", err)
	}

	// Test running a health check
	ctx := context.Background()
	result, err := healthChecker.RunCheck(ctx, "database-check")
	if err != nil {
		t.Errorf("RunCheck() failed: %v", err)
	}

	if result.CheckName != "database-check" {
		t.Errorf("Expected check name 'database-check', got '%s'", result.CheckName)
	}

	if result.Status == "" {
		t.Error("Expected non-empty health check status")
	}

	// Test getting health check history
	history, err := healthChecker.GetHistory("database-check", 10)
	if err != nil {
		t.Errorf("GetHistory() failed: %v", err)
	}

	if len(history) == 0 {
		t.Error("Expected non-empty health check history")
	}
}

func TestDashboard(t *testing.T) {
	dashboard := NewDashboard()

	// Test adding a widget
	widget := &Widget{
		ID:          "cpu-widget",
		Title:       "CPU Usage",
		Type:        "line-chart",
		MetricQuery: "cpu_usage",
		Position:    Position{X: 0, Y: 0},
		Size:        Size{Width: 6, Height: 4},
	}

	err := dashboard.AddWidget(widget)
	if err != nil {
		t.Errorf("AddWidget() failed: %v", err)
	}

	// Test getting a widget
	retrievedWidget, err := dashboard.GetWidget("cpu-widget")
	if err != nil {
		t.Errorf("GetWidget() failed: %v", err)
	}

	if retrievedWidget.ID != "cpu-widget" {
		t.Errorf("Expected widget ID 'cpu-widget', got '%s'", retrievedWidget.ID)
	}

	// Test listing widgets
	widgets := dashboard.ListWidgets()
	if len(widgets) != 1 {
		t.Errorf("Expected 1 widget, got %d", len(widgets))
	}

	// Test generating dashboard data
	ctx := context.Background()
	data, err := dashboard.GenerateData(ctx)
	if err != nil {
		t.Errorf("GenerateData() failed: %v", err)
	}

	if len(data.Widgets) != 1 {
		t.Errorf("Expected 1 widget in dashboard data, got %d", len(data.Widgets))
	}
}

func BenchmarkMetricsCollection(b *testing.B) {
	monitor := &MockMonitor{
		name:     "benchmark-monitor",
		category: "system",
		interval: 1 * time.Second,
		status:   "running",
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := monitor.CollectMetrics(ctx)
		if err != nil {
			b.Errorf("CollectMetrics() failed: %v", err)
		}
	}
}

// MockMonitor is a test implementation of the Monitor interface
type MockMonitor struct {
	name     string
	category string
	interval time.Duration
	status   string
}

func (m *MockMonitor) GetName() string {
	return m.name
}

func (m *MockMonitor) GetCategory() string {
	return m.category
}

func (m *MockMonitor) GetInterval() time.Duration {
	return m.interval
}

func (m *MockMonitor) GetStatus() string {
	return m.status
}

func (m *MockMonitor) Start() error {
	m.status = "running"
	return nil
}

func (m *MockMonitor) Stop() error {
	m.status = "stopped"
	return nil
}

func (m *MockMonitor) CollectMetrics(ctx context.Context) ([]*Metric, error) {
	metrics := []*Metric{
		{
			Name:      "cpu_usage",
			Value:     float64(45.5),
			Unit:      "percent",
			Timestamp: time.Now(),
			Labels: map[string]string{
				"host":     "localhost",
				"category": m.category,
			},
		},
		{
			Name:      "memory_usage",
			Value:     float64(67.2),
			Unit:      "percent",
			Timestamp: time.Now(),
			Labels: map[string]string{
				"host":     "localhost",
				"category": m.category,
			},
		},
	}

	return metrics, nil
}

func (m *MockMonitor) GetConfiguration() *MonitorConfig {
	return &MonitorConfig{
		Name:     m.name,
		Category: m.category,
		Interval: m.interval,
		Enabled:  true,
	}
}

func (m *MockMonitor) UpdateConfiguration(config *MonitorConfig) error {
	m.name = config.Name
	m.category = config.Category
	m.interval = config.Interval
	return nil
}

func (m *MockMonitor) IsHealthy() bool {
	return m.status == "running"
}

func (m *MockMonitor) GetSystemStatus() (*SystemStatus, error) {
	return &SystemStatus{
		Overall:   "healthy",
		Timestamp: time.Now(),
	}, nil
}

func (m *MockMonitor) GetServiceStatus(serviceName string, detailed bool) (*ServiceStatus, error) {
	return &ServiceStatus{
		Name:   serviceName,
		Status: "healthy",
		Health: "healthy",
	}, nil
}

func (m *MockMonitor) ListServices() ([]*ServiceInfo, error) {
	return []*ServiceInfo{
		{
			Name:   "test-service",
			Status: "running",
		},
	}, nil
}

func (m *MockMonitor) GetMetrics(metric, duration string) (*MetricsData, error) {
	return &MetricsData{
		Metric:    metric,
		TimeRange: duration,
		Data: []MetricPoint{
			{
				Timestamp: time.Now(),
				Value:     42.0,
			},
		},
	}, nil
}

func (m *MockMonitor) CreateAlert(alert AlertConfig) error {
	// Mock implementation
	return nil
}

func (m *MockMonitor) ListAlerts() ([]*Alert, error) {
	return []*Alert{
		{
			RuleName:  "test-alert",
			Severity:  "warning",
			Message:   "Test alert message",
			Timestamp: time.Now(),
			Value:     42.0,
		},
	}, nil
}

func (m *MockMonitor) DeleteAlert(name string) error {
	// Mock implementation
	return nil
}

func (m *MockMonitor) StartDashboard(host string, port int) error {
	// Mock implementation
	return nil
}
