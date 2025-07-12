package cloud

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestCloudManager(t *testing.T) {
	manager := NewCloudManager()

	// Test adding a provider
	provider := &MockCloudProvider{
		name:   "aws",
		region: "us-west-2",
		status: "connected",
	}

	err := manager.AddProvider(provider)
	if err != nil {
		t.Errorf("AddProvider() failed: %v", err)
	}

	// Test getting a provider
	retrievedProvider, err := manager.GetProvider("aws")
	if err != nil {
		t.Errorf("GetProvider() failed: %v", err)
	}

	if retrievedProvider.GetName() != "aws" {
		t.Errorf("Expected provider name 'aws', got '%s'", retrievedProvider.GetName())
	}

	// Test listing providers
	providers := manager.ListProviders()
	if len(providers) != 1 {
		t.Errorf("Expected 1 provider, got %d", len(providers))
	}
}

func TestResourceOperations(t *testing.T) {
	provider := &MockCloudProvider{
		name:   "aws",
		region: "us-west-2",
		status: "connected",
	}

	ctx := context.Background()

	// Test listing resources
	resources, err := provider.ListResources(ctx, "ec2")
	if err != nil {
		t.Errorf("ListResources() failed: %v", err)
	}

	if len(resources) == 0 {
		t.Error("Expected non-empty resources list")
	}

	// Test getting resource details
	resourceID := resources[0].ID
	details, err := provider.GetResourceDetails(ctx, resourceID)
	if err != nil {
		t.Errorf("GetResourceDetails() failed: %v", err)
	}

	if details.ID != resourceID {
		t.Errorf("Expected resource ID '%s', got '%s'", resourceID, details.ID)
	}

	// Test creating a resource
	createReq := &CreateResourceRequest{
		Type:   "ec2",
		Name:   "test-instance",
		Region: "us-west-2",
		Config: map[string]interface{}{
			"instance_type": "t3.micro",
			"image_id":      "ami-12345678",
		},
	}

	resource, err := provider.CreateResource(ctx, createReq)
	if err != nil {
		t.Errorf("CreateResource() failed: %v", err)
	}

	if resource.Name != "test-instance" {
		t.Errorf("Expected resource name 'test-instance', got '%s'", resource.Name)
	}

	// Test updating a resource
	updateReq := &UpdateResourceRequest{
		ID: resource.ID,
		Config: map[string]interface{}{
			"instance_type": "t3.small",
		},
	}

	updatedResource, err := provider.UpdateResource(ctx, updateReq)
	if err != nil {
		t.Errorf("UpdateResource() failed: %v", err)
	}

	if updatedResource.Config["instance_type"] != "t3.small" {
		t.Error("Expected resource to be updated")
	}

	// Test deleting a resource
	err = provider.DeleteResource(ctx, resource.ID)
	if err != nil {
		t.Errorf("DeleteResource() failed: %v", err)
	}
}

func TestCloudMetrics(t *testing.T) {
	provider := &MockCloudProvider{
		name:   "aws",
		region: "us-west-2",
		status: "connected",
	}

	ctx := context.Background()

	// Test getting metrics
	metricsReq := &MetricsRequest{
		ResourceID: "i-12345678",
		MetricName: "CPUUtilization",
		StartTime:  time.Now().Add(-1 * time.Hour),
		EndTime:    time.Now(),
		Period:     300, // 5 minutes
	}

	metrics, err := provider.GetMetrics(ctx, metricsReq)
	if err != nil {
		t.Errorf("GetMetrics() failed: %v", err)
	}

	if len(metrics.DataPoints) == 0 {
		t.Error("Expected non-empty metrics data points")
	}

	// Test getting cost information
	costReq := &CostRequest{
		StartTime: time.Now().Add(-24 * time.Hour),
		EndTime:   time.Now(),
		GroupBy:   "SERVICE",
	}

	cost, err := provider.GetCost(ctx, costReq)
	if err != nil {
		t.Errorf("GetCost() failed: %v", err)
	}

	if cost.Total <= 0 {
		t.Error("Expected positive cost total")
	}
}

func TestCloudProviderConfiguration(t *testing.T) {
	provider := &MockCloudProvider{
		name:   "aws",
		region: "us-west-2",
		status: "connected",
	}

	// Test getting initial configuration
	config := provider.GetConfiguration()
	if config.Region != "us-west-2" {
		t.Errorf("Expected region 'us-west-2', got '%s'", config.Region)
	}

	// Test updating configuration
	newConfig := &ProviderConfig{
		Region: "us-east-1",
		Credentials: map[string]string{
			"access_key": "test-key",
			"secret_key": "test-secret",
		},
	}

	err := provider.UpdateConfiguration(newConfig)
	if err != nil {
		t.Errorf("UpdateConfiguration() failed: %v", err)
	}

	updatedConfig := provider.GetConfiguration()
	if updatedConfig.Region != "us-east-1" {
		t.Errorf("Expected region 'us-east-1', got '%s'", updatedConfig.Region)
	}
}

func TestCloudProviderStatus(t *testing.T) {
	provider := &MockCloudProvider{
		name:   "aws",
		region: "us-west-2",
		status: "connected",
	}

	status := provider.GetStatus()
	if status.Status != "connected" {
		t.Errorf("Expected status 'connected', got '%s'", status.Status)
	}

	if status.LastCheck.IsZero() {
		t.Error("Expected non-zero last check time")
	}
}

func BenchmarkListResources(b *testing.B) {
	provider := &MockCloudProvider{
		name:   "aws",
		region: "us-west-2",
		status: "connected",
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := provider.ListResources(ctx, "ec2")
		if err != nil {
			b.Errorf("ListResources() failed: %v", err)
		}
	}
}

// MockCloudProvider is a test implementation of the CloudProvider interface
type MockCloudProvider struct {
	name           string
	region         string
	status         string
	config         *ProviderConfig
	providerStatus *ProviderStatus
	resources      map[string]*Resource
}

func (m *MockCloudProvider) GetName() string {
	return m.name
}

func (m *MockCloudProvider) GetType() string {
	return "aws"
}

func (m *MockCloudProvider) Connect(ctx context.Context) error {
	m.status = "connected"
	return nil
}

func (m *MockCloudProvider) Disconnect(ctx context.Context) error {
	m.status = "disconnected"
	return nil
}

func (m *MockCloudProvider) IsConnected() bool {
	return m.status == "connected"
}

func (m *MockCloudProvider) ListResources(ctx context.Context, resourceType string) ([]*Resource, error) {
	if m.resources == nil {
		m.resources = make(map[string]*Resource)
		// Add default resources
		m.resources["i-12345678"] = &Resource{
			ID:     "i-12345678",
			Name:   "test-instance",
			Type:   "ec2",
			Region: m.region,
			Status: "running",
			Config: map[string]interface{}{
				"instance_type": "t3.micro",
				"image_id":      "ami-12345678",
			},
			CreatedAt: time.Now().Add(-24 * time.Hour),
		}
		m.resources["vol-87654321"] = &Resource{
			ID:     "vol-87654321",
			Name:   "test-volume",
			Type:   "ebs",
			Region: m.region,
			Status: "available",
			Config: map[string]interface{}{
				"size": 20,
				"type": "gp3",
			},
			CreatedAt: time.Now().Add(-12 * time.Hour),
		}
	}

	resources := make([]*Resource, 0, len(m.resources))
	for _, resource := range m.resources {
		if resourceType == "" || resource.Type == resourceType {
			resources = append(resources, resource)
		}
	}
	return resources, nil
}

func (m *MockCloudProvider) GetResourceDetails(ctx context.Context, resourceID string) (*Resource, error) {
	if m.resources == nil {
		m.ListResources(ctx, "") // Initialize resources
	}

	if resource, exists := m.resources[resourceID]; exists {
		return resource, nil
	}

	return nil, fmt.Errorf("resource not found: %s", resourceID)
}

func (m *MockCloudProvider) CreateResource(ctx context.Context, req *CreateResourceRequest) (*Resource, error) {
	if m.resources == nil {
		m.resources = make(map[string]*Resource)
	}

	resourceID := fmt.Sprintf("i-%s", generateRandomID())
	resource := &Resource{
		ID:        resourceID,
		Name:      req.Name,
		Type:      req.Type,
		Region:    req.Region,
		Status:    "pending",
		Config:    req.Config,
		CreatedAt: time.Now(),
	}

	m.resources[resourceID] = resource
	return resource, nil
}

func (m *MockCloudProvider) UpdateResource(ctx context.Context, req *UpdateResourceRequest) (*Resource, error) {
	if m.resources == nil {
		m.ListResources(ctx, "") // Initialize resources
	}

	resource, exists := m.resources[req.ID]
	if !exists {
		return nil, fmt.Errorf("resource not found: %s", req.ID)
	}

	// Update the config
	if resource.Config == nil {
		resource.Config = make(map[string]interface{})
	}
	for key, value := range req.Config {
		resource.Config[key] = value
	}

	return resource, nil
}

func (m *MockCloudProvider) DeleteResource(ctx context.Context, resourceID string) error {
	if m.resources == nil {
		m.ListResources(ctx, "") // Initialize resources
	}

	if _, exists := m.resources[resourceID]; !exists {
		return fmt.Errorf("resource not found: %s", resourceID)
	}

	delete(m.resources, resourceID)
	return nil
}

// generateRandomID generates a random ID for testing
func generateRandomID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()%1000000)
}

func (m *MockCloudProvider) GetMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, error) {
	dataPoints := make([]*MetricDataPoint, 0)

	// Generate mock data points
	for i := 0; i < 12; i++ {
		dataPoints = append(dataPoints, &MetricDataPoint{
			Timestamp: req.StartTime.Add(time.Duration(i) * time.Minute * 5),
			Value:     float64(30 + i*2), // Mock CPU utilization
			Unit:      "Percent",
		})
	}

	return &MetricsResponse{
		MetricName: req.MetricName,
		DataPoints: dataPoints,
	}, nil
}

func (m *MockCloudProvider) GetCost(ctx context.Context, req *CostRequest) (*CostResponse, error) {
	return &CostResponse{
		Total:    45.67,
		Currency: "USD",
		Period: &CostPeriod{
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
		},
		BreakdownBy: map[string]float64{
			"EC2": 30.45,
			"EBS": 15.22,
		},
	}, nil
}

func (m *MockCloudProvider) GetConfiguration() *ProviderConfig {
	if m.config == nil {
		m.config = &ProviderConfig{
			Region: m.region,
			Credentials: map[string]string{
				"profile": "default",
			},
		}
	}
	return m.config
}

func (m *MockCloudProvider) UpdateConfiguration(config *ProviderConfig) error {
	m.config = config
	m.region = config.Region
	return nil
}

func (m *MockCloudProvider) GetStatus() *ProviderStatus {
	if m.providerStatus == nil {
		m.providerStatus = &ProviderStatus{
			Status:    m.status,
			LastCheck: time.Now(),
			Health:    "healthy",
		}
	}
	return m.providerStatus
}

func (m *MockCloudProvider) ValidateCredentials(ctx context.Context) error {
	return nil
}

func (m *MockCloudProvider) GetRegions(ctx context.Context) ([]string, error) {
	return []string{"us-west-2", "us-east-1", "eu-west-1"}, nil
}

func (m *MockCloudProvider) GetResourceTypes(ctx context.Context) ([]string, error) {
	return []string{"ec2", "ebs", "s3", "rds"}, nil
}
