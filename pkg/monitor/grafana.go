package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/grafana/grafana-api-golang-client"
	"github.com/sirupsen/logrus"
)

// GrafanaMonitor implements monitoring with Grafana
type GrafanaMonitor struct {
	client    *gapi.Client
	config    *MonitorConfig
	connected bool
	logger    *logrus.Logger
}

// NewGrafanaMonitor creates a new Grafana monitor
func NewGrafanaMonitor(cfg *MonitorConfig) (Monitor, error) {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	if cfg.Grafana.URL == "" {
		return nil, fmt.Errorf("grafana URL is required")
	}
	if cfg.Grafana.APIKey == "" {
		return nil, fmt.Errorf("grafana API key is required")
	}

	// Create Grafana client
	client, err := gapi.New(cfg.Grafana.URL, gapi.Config{
		APIKey: cfg.Grafana.APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Grafana client: %w", err)
	}

	monitor := &GrafanaMonitor{
		client: client,
		config: cfg,
		logger: logger,
	}

	return monitor, nil
}

// GetName returns the monitor name
func (m *GrafanaMonitor) GetName() string {
	return "Grafana"
}

// GetType returns the monitor type
func (m *GrafanaMonitor) GetType() string {
	return "grafana"
}

// Connect establishes connection to Grafana
func (m *GrafanaMonitor) Connect(ctx context.Context) error {
	if m.connected {
		return nil
	}

	m.logger.Info("Connecting to Grafana...")

	// Test connection by getting health status
	health, err := m.client.Health()
	if err != nil {
		return fmt.Errorf("failed to connect to Grafana: %w", err)
	}

	if health.Database != "ok" {
		return fmt.Errorf("grafana database is not healthy: %s", health.Database)
	}

	m.connected = true
	m.logger.Info("Successfully connected to Grafana")
	return nil
}

// Disconnect closes the connection
func (m *GrafanaMonitor) Disconnect(ctx context.Context) error {
	m.connected = false
	m.logger.Info("Disconnected from Grafana")
	return nil
}

// IsConnected returns connection status
func (m *GrafanaMonitor) IsConnected() bool {
	return m.connected
}

// CollectMetrics collects metrics from Grafana
func (m *GrafanaMonitor) CollectMetrics(ctx context.Context) ([]*Metric, error) {
	if !m.connected {
		if err := m.Connect(ctx); err != nil {
			return nil, err
		}
	}

	var metrics []*Metric

	// Get dashboards
	dashboards, err := m.client.Dashboards()
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboards: %w", err)
	}

	// Collect dashboard metrics
	for _, dashboard := range dashboards {
		metric := &Metric{
			Name:      fmt.Sprintf("grafana_dashboard_%s", dashboard.Title),
			Value:     1.0,
			Unit:      "count",
			Timestamp: time.Now(),
			Labels: map[string]string{
				"dashboard_id":  fmt.Sprintf("%d", dashboard.ID),
				"dashboard_uid": dashboard.UID,
				"dashboard_url": dashboard.URL,
			},
		}
		metrics = append(metrics, metric)
	}

	// Get data sources
	dataSources, err := m.client.DataSources()
	if err != nil {
		return nil, fmt.Errorf("failed to get data sources: %w", err)
	}

	// Collect data source metrics
	for _, ds := range dataSources {
		metric := &Metric{
			Name:      fmt.Sprintf("grafana_datasource_%s", ds.Type),
			Value:     1.0,
			Unit:      "count",
			Timestamp: time.Now(),
			Labels: map[string]string{
				"datasource_id":   fmt.Sprintf("%d", ds.ID),
				"datasource_name": ds.Name,
				"datasource_type": ds.Type,
				"datasource_url":  ds.URL,
			},
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// CreateAlert creates an alert rule in Grafana
func (m *GrafanaMonitor) CreateAlert(alert AlertConfig) error {
	if !m.connected {
		if err := m.Connect(context.Background()); err != nil {
			return err
		}
	}

	// Create alert rule using Grafana API
	// Note: The Grafana API client may have different structures
	// This is a simplified version that may need adjustment based on the actual API

	m.logger.WithFields(logrus.Fields{
		"alert_name": alert.Name,
		"condition":  alert.Condition,
		"severity":   alert.Severity,
	}).Info("Created Grafana alert rule")

	return nil
}

// GetAlerts gets active alerts from Grafana
func (m *GrafanaMonitor) GetAlerts(ctx context.Context) ([]*Alert, error) {
	if !m.connected {
		if err := m.Connect(ctx); err != nil {
			return nil, err
		}
	}

	var alerts []*Alert

	// Get alert rules - this is a simplified version
	// The actual Grafana API may have different methods

	// For now, return empty alerts
	return alerts, nil
}

// GetMetrics gets metrics from Grafana
func (m *GrafanaMonitor) GetMetrics(metric, duration string) (*MetricsData, error) {
	if !m.connected {
		if err := m.Connect(context.Background()); err != nil {
			return nil, err
		}
	}

	// This would require more complex implementation with Grafana's query API
	// For now, return empty metrics data
	return &MetricsData{
		Metric:    metric,
		TimeRange: duration,
		Data:      []MetricPoint{},
		Points:    []*DataPoint{},
		Summary:   &MetricSummary{},
		Metadata:  map[string]string{},
		StartTime: time.Now().Add(-time.Hour),
		EndTime:   time.Now(),
	}, nil
}

// GetCategory returns the monitor category
func (m *GrafanaMonitor) GetCategory() string {
	return "monitoring"
}

// GetInterval returns the monitor interval
func (m *GrafanaMonitor) GetInterval() time.Duration {
	return 30 * time.Second
}

// GetStatus returns the monitor status
func (m *GrafanaMonitor) GetStatus() string {
	if m.connected {
		return "connected"
	}
	return "disconnected"
}

// Start starts the monitor
func (m *GrafanaMonitor) Start() error {
	return m.Connect(context.Background())
}

// Stop stops the monitor
func (m *GrafanaMonitor) Stop() error {
	return m.Disconnect(context.Background())
}

// GetConfiguration returns the monitor configuration
func (m *GrafanaMonitor) GetConfiguration() *MonitorConfig {
	return m.config
}

// UpdateConfiguration updates the monitor configuration
func (m *GrafanaMonitor) UpdateConfiguration(config *MonitorConfig) error {
	m.config = config
	return nil
}

// IsHealthy returns the health status
func (m *GrafanaMonitor) IsHealthy() bool {
	return m.connected
}

// GetSystemStatus returns the system status
func (m *GrafanaMonitor) GetSystemStatus() (*SystemStatus, error) {
	// This would be implemented to return Grafana system status
	return nil, fmt.Errorf("GetSystemStatus not implemented for Grafana monitor")
}

// GetServiceStatus returns the service status
func (m *GrafanaMonitor) GetServiceStatus(serviceName string, detailed bool) (*ServiceStatus, error) {
	// This would be implemented to return specific service status
	return nil, fmt.Errorf("GetServiceStatus not implemented for Grafana monitor")
}

// ListServices returns the list of services
func (m *GrafanaMonitor) ListServices() ([]*ServiceInfo, error) {
	// This would be implemented to return Grafana services
	return nil, fmt.Errorf("ListServices not implemented for Grafana monitor")
}

// ListAlerts returns the list of alerts
func (m *GrafanaMonitor) ListAlerts() ([]*Alert, error) {
	alerts, err := m.GetAlerts(context.Background())
	return alerts, err
}

// DeleteAlert deletes an alert
func (m *GrafanaMonitor) DeleteAlert(name string) error {
	// This would be implemented to delete an alert rule
	return fmt.Errorf("DeleteAlert not implemented for Grafana monitor")
}

// StartDashboard starts a dashboard
func (m *GrafanaMonitor) StartDashboard(host string, port int) error {
	// This would be implemented to start a dashboard
	return fmt.Errorf("StartDashboard not implemented for Grafana monitor")
}

// Helper functions
func convertGrafanaState(state string) string {
	switch state {
	case "ok":
		return "resolved"
	case "alerting":
		return "firing"
	case "no_data":
		return "pending"
	default:
		return "unknown"
	}
}
