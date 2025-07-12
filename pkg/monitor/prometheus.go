package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
)

// PrometheusMonitor implements Monitor interface for Prometheus
type PrometheusMonitor struct {
	client    api.Client
	api       v1.API
	endpoint  string
	config    *MonitorConfig
	logger    *logrus.Logger
	connected bool
}

// NewPrometheusMonitor creates a new Prometheus monitor
func NewPrometheusMonitor(endpoint string, config *MonitorConfig) (*PrometheusMonitor, error) {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	client, err := api.NewClient(api.Config{
		Address: endpoint,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus client: %w", err)
	}

	return &PrometheusMonitor{
		client:   client,
		api:      v1.NewAPI(client),
		endpoint: endpoint,
		config:   config,
		logger:   logger,
	}, nil
}

// GetName returns the monitor name
func (m *PrometheusMonitor) GetName() string {
	return "Prometheus"
}

// GetCategory returns the monitor category
func (m *PrometheusMonitor) GetCategory() string {
	return "metrics"
}

// GetInterval returns the monitoring interval
func (m *PrometheusMonitor) GetInterval() time.Duration {
	if m.config != nil && m.config.Interval > 0 {
		return time.Duration(m.config.Interval) * time.Second
	}
	return 30 * time.Second
}

// GetStatus returns the monitor status
func (m *PrometheusMonitor) GetStatus() string {
	if m.connected {
		return "connected"
	}
	return "disconnected"
}

// Start starts the monitor
func (m *PrometheusMonitor) Start() error {
	m.logger.Info("Starting Prometheus monitor...")
	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, _, err := m.api.Query(ctx, "up", time.Now())
	if err != nil {
		return fmt.Errorf("failed to connect to Prometheus: %w", err)
	}

	m.connected = true
	m.logger.Info("Prometheus monitor started successfully")
	return nil
}

// Stop stops the monitor
func (m *PrometheusMonitor) Stop() error {
	m.connected = false
	m.logger.Info("Prometheus monitor stopped")
	return nil
}

// CollectMetrics collects metrics from Prometheus
func (m *PrometheusMonitor) CollectMetrics(ctx context.Context) ([]*Metric, error) {
	if !m.connected {
		if err := m.Start(); err != nil {
			return nil, err
		}
	}

	var metrics []*Metric

	// Common system metrics to collect
	queries := map[string]string{
		"cpu_usage":    `100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)`,
		"memory_usage": `(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100`,
		"disk_usage":   `100 - ((node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"}) * 100)`,
		"network_rx":   `irate(node_network_receive_bytes_total{device!~"lo|veth.*|docker.*"}[5m])`,
		"network_tx":   `irate(node_network_transmit_bytes_total{device!~"lo|veth.*|docker.*"}[5m])`,
	}

	for name, query := range queries {
		result, warnings, err := m.api.Query(ctx, query, time.Now())
		if err != nil {
			m.logger.Warnf("Failed to query %s: %v", name, err)
			continue
		}

		if len(warnings) > 0 {
			m.logger.Warnf("Query warnings for %s: %v", name, warnings)
		}

		// Convert Prometheus result to our metrics
		prometheusMetrics := m.convertPrometheusResult(name, result)
		metrics = append(metrics, prometheusMetrics...)
	}

	return metrics, nil
}

// convertPrometheusResult converts Prometheus query result to our metric format
func (m *PrometheusMonitor) convertPrometheusResult(name string, result model.Value) []*Metric {
	var metrics []*Metric

	switch result.Type() {
	case model.ValVector:
		vector := result.(model.Vector)
		for _, sample := range vector {
			metric := &Metric{
				Name:      name,
				Value:     float64(sample.Value),
				Unit:      m.getMetricUnit(name),
				Timestamp: sample.Timestamp.Time(),
				Labels:    m.convertLabels(sample.Metric),
			}
			metrics = append(metrics, metric)
		}
	case model.ValScalar:
		scalar := result.(*model.Scalar)
		metric := &Metric{
			Name:      name,
			Value:     float64(scalar.Value),
			Unit:      m.getMetricUnit(name),
			Timestamp: scalar.Timestamp.Time(),
			Labels:    make(map[string]string),
		}
		metrics = append(metrics, metric)
	}

	return metrics
}

// convertLabels converts Prometheus labels to our format
func (m *PrometheusMonitor) convertLabels(metric model.Metric) map[string]string {
	labels := make(map[string]string)
	for k, v := range metric {
		labels[string(k)] = string(v)
	}
	return labels
}

// getMetricUnit returns the unit for a metric
func (m *PrometheusMonitor) getMetricUnit(name string) string {
	switch name {
	case "cpu_usage", "memory_usage", "disk_usage":
		return "percent"
	case "network_rx", "network_tx":
		return "bytes/sec"
	default:
		return ""
	}
}

// GetConfiguration returns the monitor configuration
func (m *PrometheusMonitor) GetConfiguration() *MonitorConfig {
	return m.config
}

// UpdateConfiguration updates the monitor configuration
func (m *PrometheusMonitor) UpdateConfiguration(config *MonitorConfig) error {
	m.config = config
	return nil
}

// IsHealthy returns the health status
func (m *PrometheusMonitor) IsHealthy() bool {
	return m.connected
}

// GetSystemStatus returns overall system status
func (m *PrometheusMonitor) GetSystemStatus() (*SystemStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	metrics, err := m.CollectMetrics(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate overall system health
	status := &SystemStatus{
		Overall:   "healthy",
		Timestamp: time.Now(),
		Services:  []*ServiceStatus{},
		Resources: &ResourceUsage{
			CPU:     &CPUUsage{},
			Memory:  &MemoryUsage{},
			Disk:    &DiskUsage{},
			Network: &NetworkUsage{},
		},
		Alerts:   []*ActiveAlert{},
		Uptime:   time.Since(time.Now().Add(-24 * time.Hour)), // Mock uptime
		Metadata: make(map[string]string),
	}

	// Process metrics to populate status
	for _, metric := range metrics {
		switch metric.Name {
		case "cpu_usage":
			if val, ok := metric.Value.(float64); ok {
				status.Resources.CPU.Usage = val
				if val > 80 {
					status.Overall = "warning"
				}
			}
		case "memory_usage":
			if val, ok := metric.Value.(float64); ok {
				status.Resources.Memory.Usage = val
				if val > 85 {
					status.Overall = "warning"
				}
			}
		case "disk_usage":
			if val, ok := metric.Value.(float64); ok {
				status.Resources.Disk.Usage = val
				if val > 90 {
					status.Overall = "critical"
				}
			}
		}
	}

	return status, nil
}

// GetServiceStatus returns service status
func (m *PrometheusMonitor) GetServiceStatus(serviceName string, detailed bool) (*ServiceStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Query service-specific metrics
	query := fmt.Sprintf(`up{job="%s"}`, serviceName)
	result, _, err := m.api.Query(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to query service status: %w", err)
	}

	status := &ServiceStatus{
		Name:      serviceName,
		Status:    "unknown",
		Health:    "unknown",
		Uptime:    0,
		CPU:       0,
		Memory:    0,
		Disk:      0,
		Network:   &NetworkMetrics{},
		Endpoints: []*EndpointStatus{},
		LastCheck: time.Now(),
		Metadata:  make(map[string]string),
	}

	// Process result
	if vector, ok := result.(model.Vector); ok && len(vector) > 0 {
		if float64(vector[0].Value) == 1 {
			status.Status = "running"
			status.Health = "healthy"
		} else {
			status.Status = "down"
			status.Health = "unhealthy"
		}
	}

	return status, nil
}

// ListServices lists monitored services
func (m *PrometheusMonitor) ListServices() ([]*ServiceInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Query all available jobs
	query := `group by (job) (up)`
	result, _, err := m.api.Query(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var services []*ServiceInfo

	if vector, ok := result.(model.Vector); ok {
		for _, sample := range vector {
			if job, exists := sample.Metric["job"]; exists {
				service := &ServiceInfo{
					Name:        string(job),
					Type:        "service",
					Status:      "monitored",
					Description: fmt.Sprintf("Service monitored by Prometheus: %s", job),
				}
				services = append(services, service)
			}
		}
	}

	return services, nil
}

// GetMetrics returns historical metrics
func (m *PrometheusMonitor) GetMetrics(metric, duration string) (*MetricsData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Parse duration
	dur, err := time.ParseDuration(duration)
	if err != nil {
		return nil, fmt.Errorf("invalid duration: %w", err)
	}

	// Query range
	endTime := time.Now()
	startTime := endTime.Add(-dur)
	step := dur / 100 // 100 data points

	result, _, err := m.api.QueryRange(ctx, metric, v1.Range{
		Start: startTime,
		End:   endTime,
		Step:  step,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics: %w", err)
	}

	data := &MetricsData{
		Metric:    metric,
		StartTime: startTime,
		EndTime:   endTime,
		Points:    []*DataPoint{},
	}

	// Convert result
	if matrix, ok := result.(model.Matrix); ok {
		for _, sampleStream := range matrix {
			for _, pair := range sampleStream.Values {
				point := &DataPoint{
					Timestamp: pair.Timestamp.Time(),
					Value:     float64(pair.Value),
					Labels:    m.convertLabels(sampleStream.Metric),
				}
				data.Points = append(data.Points, point)
			}
		}
	}

	return data, nil
}

// CreateAlert creates a new alert
func (m *PrometheusMonitor) CreateAlert(alert AlertConfig) error {
	// This would typically integrate with Prometheus Alertmanager
	return fmt.Errorf("CreateAlert not implemented for Prometheus monitor")
}

// ListAlerts lists active alerts
func (m *PrometheusMonitor) ListAlerts() ([]*Alert, error) {
	// This would typically query Prometheus Alertmanager
	return []*Alert{}, nil
}

// DeleteAlert deletes an alert
func (m *PrometheusMonitor) DeleteAlert(name string) error {
	// This would typically integrate with Prometheus Alertmanager
	return fmt.Errorf("DeleteAlert not implemented for Prometheus monitor")
}

// StartDashboard starts a monitoring dashboard
func (m *PrometheusMonitor) StartDashboard(host string, port int) error {
	// This would typically start a web dashboard
	return fmt.Errorf("StartDashboard not implemented for Prometheus monitor")
}
