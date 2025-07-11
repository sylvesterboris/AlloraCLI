package cloud

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AlloraAi/AlloraCLI/pkg/config"
)

// CloudService interface defines cloud provider operations
type CloudService interface {
	ListResources(ctx context.Context, provider string, resourceType string) ([]Resource, error)
	CreateResource(ctx context.Context, provider string, spec ResourceSpec) (*Resource, error)
	UpdateResource(ctx context.Context, provider string, resourceID string, spec ResourceSpec) (*Resource, error)
	DeleteResource(ctx context.Context, provider string, resourceID string) error
	GetResourceDetails(ctx context.Context, provider string, resourceID string) (*ResourceDetails, error)
	GetCostAnalysis(ctx context.Context, provider string, options CostOptions) (*CostAnalysis, error)
	OptimizeResources(ctx context.Context, provider string, options OptimizeOptions) (*OptimizationResult, error)
	MonitorHealth(ctx context.Context, provider string) (<-chan HealthEvent, error)
}

// CloudProvider interface defines cloud provider operations
type CloudProvider interface {
	GetName() string
	GetType() string
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	IsConnected() bool
	ListResources(ctx context.Context, resourceType string) ([]*Resource, error)
	GetResourceDetails(ctx context.Context, resourceID string) (*Resource, error)
	CreateResource(ctx context.Context, req *CreateResourceRequest) (*Resource, error)
	UpdateResource(ctx context.Context, req *UpdateResourceRequest) (*Resource, error)
	DeleteResource(ctx context.Context, resourceID string) error
	GetMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, error)
	GetCost(ctx context.Context, req *CostRequest) (*CostResponse, error)
	GetConfiguration() *ProviderConfig
	UpdateConfiguration(config *ProviderConfig) error
	GetStatus() *ProviderStatus
	ValidateCredentials(ctx context.Context) error
	GetRegions(ctx context.Context) ([]string, error)
	GetResourceTypes(ctx context.Context) ([]string, error)
}

// Resource represents a cloud resource
type Resource struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Provider  string                 `json:"provider"`
	Region    string                 `json:"region"`
	State     string                 `json:"state"`
	Status    string                 `json:"status"`
	Config    map[string]interface{} `json:"config"`
	CreatedAt time.Time              `json:"created_at"`
	Created   time.Time              `json:"created"`
	Modified  time.Time              `json:"modified"`
	Tags      map[string]string      `json:"tags"`
	Metadata  map[string]string      `json:"metadata"`
	Cost      *CostInfo              `json:"cost,omitempty"`
}

// ResourceSpec defines the specification for creating/updating resources
type ResourceSpec struct {
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Region       string            `json:"region"`
	Configuration map[string]interface{} `json:"configuration"`
	Tags         map[string]string `json:"tags"`
}

// ResourceDetails provides detailed information about a resource
type ResourceDetails struct {
	Resource       Resource               `json:"resource"`
	Configuration  map[string]interface{} `json:"configuration"`
	Dependencies   []string               `json:"dependencies"`
	SecurityGroups []SecurityGroup        `json:"security_groups"`
	NetworkInfo    NetworkInfo            `json:"network_info"`
	Monitoring     MonitoringInfo         `json:"monitoring"`
	Backup         BackupInfo             `json:"backup"`
}

// CostInfo provides cost information for a resource
type CostInfo struct {
	Monthly     float64 `json:"monthly"`
	Daily       float64 `json:"daily"`
	Currency    string  `json:"currency"`
	LastUpdated time.Time `json:"last_updated"`
}

// SecurityGroup represents a security group
type SecurityGroup struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Rules       []SecurityGroupRule   `json:"rules"`
}

// SecurityGroupRule represents a security group rule
type SecurityGroupRule struct {
	Protocol   string `json:"protocol"`
	Port       string `json:"port"`
	Source     string `json:"source"`
	Direction  string `json:"direction"`
	Action     string `json:"action"`
}

// NetworkInfo provides network information
type NetworkInfo struct {
	VPC        string   `json:"vpc"`
	Subnet     string   `json:"subnet"`
	PrivateIPs []string `json:"private_ips"`
	PublicIPs  []string `json:"public_ips"`
	DNS        []string `json:"dns"`
}

// MonitoringInfo provides monitoring information
type MonitoringInfo struct {
	Enabled       bool                   `json:"enabled"`
	Metrics       []MetricInfo           `json:"metrics"`
	Alerts        []AlertInfo            `json:"alerts"`
	Dashboards    []string               `json:"dashboards"`
}

// MetricInfo represents a metric
type MetricInfo struct {
	Name        string    `json:"name"`
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Threshold   float64   `json:"threshold"`
}

// AlertInfo represents an alert
type AlertInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Severity    string    `json:"severity"`
	Description string    `json:"description"`
	Triggered   time.Time `json:"triggered"`
}

// BackupInfo provides backup information
type BackupInfo struct {
	Enabled      bool      `json:"enabled"`
	Schedule     string    `json:"schedule"`
	LastBackup   time.Time `json:"last_backup"`
	NextBackup   time.Time `json:"next_backup"`
	RetentionDays int      `json:"retention_days"`
}

// CostOptions defines options for cost analysis
type CostOptions struct {
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Granularity  string    `json:"granularity"`
	GroupBy      []string  `json:"group_by"`
	ResourceType string    `json:"resource_type"`
}

// CostAnalysis provides cost analysis results
type CostAnalysis struct {
	TotalCost      float64         `json:"total_cost"`
	Currency       string          `json:"currency"`
	Period         string          `json:"period"`
	Breakdown      []CostBreakdown `json:"breakdown"`
	Trends         []CostTrend     `json:"trends"`
	Recommendations []CostRecommendation `json:"recommendations"`
}

// CostBreakdown provides cost breakdown by category
type CostBreakdown struct {
	Category    string  `json:"category"`
	Cost        float64 `json:"cost"`
	Percentage  float64 `json:"percentage"`
	ResourceCount int   `json:"resource_count"`
}

// CostTrend represents cost trend data
type CostTrend struct {
	Date   time.Time `json:"date"`
	Cost   float64   `json:"cost"`
	Change float64   `json:"change"`
}

// CostRecommendation represents a cost optimization recommendation
type CostRecommendation struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Savings     float64 `json:"potential_savings"`
	Effort      string  `json:"effort"`
	Risk        string  `json:"risk"`
	Actions     []string `json:"actions"`
}

// OptimizeOptions defines options for resource optimization
type OptimizeOptions struct {
	ResourceTypes []string `json:"resource_types"`
	Criteria      []string `json:"criteria"`
	DryRun        bool     `json:"dry_run"`
}

// OptimizationResult provides optimization results
type OptimizationResult struct {
	ID              string                    `json:"id"`
	Timestamp       time.Time                 `json:"timestamp"`
	Status          string                    `json:"status"`
	Recommendations []OptimizationRecommendation `json:"recommendations"`
	PotentialSavings float64                  `json:"potential_savings"`
	RiskAssessment  string                    `json:"risk_assessment"`
}

// OptimizationRecommendation represents an optimization recommendation
type OptimizationRecommendation struct {
	ResourceID   string            `json:"resource_id"`
	Type         string            `json:"type"`
	Current      map[string]interface{} `json:"current"`
	Recommended  map[string]interface{} `json:"recommended"`
	Savings      float64           `json:"savings"`
	Confidence   float64           `json:"confidence"`
	Actions      []string          `json:"actions"`
}

// HealthEvent represents a health monitoring event
type HealthEvent struct {
	ID          string            `json:"id"`
	Timestamp   time.Time         `json:"timestamp"`
	Provider    string            `json:"provider"`
	ResourceID  string            `json:"resource_id"`
	Type        string            `json:"type"`
	Status      string            `json:"status"`
	Message     string            `json:"message"`
	Severity    string            `json:"severity"`
	Details     map[string]string `json:"details"`
}

// CreateResourceRequest represents a request to create a resource
type CreateResourceRequest struct {
	Type   string                 `json:"type"`
	Name   string                 `json:"name"`
	Region string                 `json:"region"`
	Config map[string]interface{} `json:"config"`
}

// UpdateResourceRequest represents a request to update a resource
type UpdateResourceRequest struct {
	ID     string                 `json:"id"`
	Config map[string]interface{} `json:"config"`
}

// MetricsRequest represents a request for metrics
type MetricsRequest struct {
	ResourceID string    `json:"resource_id"`
	MetricName string    `json:"metric_name"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Period     int       `json:"period"`
}

// MetricsResponse represents a metrics response
type MetricsResponse struct {
	MetricName string              `json:"metric_name"`
	DataPoints []*MetricDataPoint  `json:"data_points"`
}

// MetricDataPoint represents a single metric data point
type MetricDataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
}

// CostRequest represents a request for cost information
type CostRequest struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	GroupBy   string    `json:"group_by"`
}

// CostResponse represents a cost response
type CostResponse struct {
	Total       float64            `json:"total"`
	Currency    string             `json:"currency"`
	Period      *CostPeriod        `json:"period"`
	BreakdownBy map[string]float64 `json:"breakdown_by"`
}

// CostPeriod represents a cost period
type CostPeriod struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// ProviderConfig represents cloud provider configuration
type ProviderConfig struct {
	Region             string            `json:"region"`
	Credentials        map[string]string `json:"credentials"`
	// AWS specific
	Profile            string            `json:"profile,omitempty"`
	// Azure specific
	SubscriptionID     string            `json:"subscription_id,omitempty"`
	TenantID           string            `json:"tenant_id,omitempty"`
	// GCP specific
	ProjectID          string            `json:"project_id,omitempty"`
	ServiceAccountPath string            `json:"service_account_path,omitempty"`
}

// ProviderStatus represents cloud provider status
type ProviderStatus struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Connected bool      `json:"connected"`
	LastCheck time.Time `json:"last_check"`
	Health    string    `json:"health"`
}

// DefaultCloudService provides a default implementation
type DefaultCloudService struct {
	config    *config.Config
	providers map[string]CloudProvider
	mu        sync.RWMutex
}

// NewCloudService creates a new cloud service
func NewCloudService(cfg *config.Config) CloudService {
	service := &DefaultCloudService{
		config:    cfg,
		providers: make(map[string]CloudProvider),
	}
	
	// Initialize providers based on configuration
	service.initializeProviders()
	
	return service
}

// initializeProviders initializes cloud providers based on configuration
func (c *DefaultCloudService) initializeProviders() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Initialize AWS provider if configured
	if awsConfig := c.getProviderConfig("aws"); awsConfig != nil {
		if provider, err := NewAWSProvider(awsConfig); err == nil {
			c.providers["aws"] = provider
		}
	}
	
	// Initialize Azure provider if configured
	if azureConfig := c.getProviderConfig("azure"); azureConfig != nil {
		if provider, err := NewAzureProvider(azureConfig); err == nil {
			c.providers["azure"] = provider
		}
	}
	
	// Initialize GCP provider if configured
	if gcpConfig := c.getProviderConfig("gcp"); gcpConfig != nil {
		if provider, err := NewGCPProvider(gcpConfig); err == nil {
			c.providers["gcp"] = provider
		}
	}
}

// getProviderConfig extracts provider configuration from main config
func (c *DefaultCloudService) getProviderConfig(provider string) *ProviderConfig {
	// This is a simplified implementation
	// In a real implementation, you would extract from c.config
	cfg := &ProviderConfig{
		Region:      "us-west-2", // Default region
		Credentials: make(map[string]string),
	}
	
	switch provider {
	case "aws":
		cfg.Profile = "default"
		return cfg
	case "azure":
		cfg.SubscriptionID = "" // Should be loaded from config
		cfg.TenantID = ""       // Should be loaded from config
		return cfg
	case "gcp":
		cfg.ProjectID = ""                // Should be loaded from config
		cfg.ServiceAccountPath = ""       // Should be loaded from config
		return cfg
	}
	
	return nil
}

// getProvider returns the cloud provider for the given name
func (c *DefaultCloudService) getProvider(name string) (CloudProvider, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	provider, exists := c.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found or not configured", name)
	}
	
	return provider, nil
}

// ListResources lists resources from the specified provider
func (c *DefaultCloudService) ListResources(ctx context.Context, provider string, resourceType string) ([]Resource, error) {
	// Try to use real provider first
	if cloudProvider, err := c.getProvider(provider); err == nil {
		resources, err := cloudProvider.ListResources(ctx, resourceType)
		if err == nil {
			// Convert []*Resource to []Resource
			var result []Resource
			for _, res := range resources {
				if res != nil {
					result = append(result, *res)
				}
			}
			return result, nil
		}
		// If real provider fails, log warning and fall back to mock
		fmt.Printf("Warning: Real provider %s failed: %v. Using mock data.\n", provider, err)
	}
	
	// Fallback to mock implementation
	switch provider {
	case "aws":
		return c.listAWSResources(ctx, resourceType)
	case "azure":
		return c.listAzureResources(ctx, resourceType)
	case "gcp":
		return c.listGCPResources(ctx, resourceType)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// CreateResource creates a new resource
func (c *DefaultCloudService) CreateResource(ctx context.Context, provider string, spec ResourceSpec) (*Resource, error) {
	// Mock implementation
	resource := &Resource{
		ID:       fmt.Sprintf("%s-%d", spec.Name, time.Now().Unix()),
		Name:     spec.Name,
		Type:     spec.Type,
		Provider: provider,
		Region:   spec.Region,
		Status:   "creating",
		Created:  time.Now(),
		Modified: time.Now(),
		Tags:     spec.Tags,
		Metadata: map[string]string{
			"created_by": "allora-cli",
		},
	}
	
	return resource, nil
}

// UpdateResource updates an existing resource
func (c *DefaultCloudService) UpdateResource(ctx context.Context, provider string, resourceID string, spec ResourceSpec) (*Resource, error) {
	// Mock implementation
	resource := &Resource{
		ID:       resourceID,
		Name:     spec.Name,
		Type:     spec.Type,
		Provider: provider,
		Region:   spec.Region,
		Status:   "updating",
		Created:  time.Now().Add(-24 * time.Hour),
		Modified: time.Now(),
		Tags:     spec.Tags,
		Metadata: map[string]string{
			"updated_by": "allora-cli",
		},
	}
	
	return resource, nil
}

// DeleteResource deletes a resource
func (c *DefaultCloudService) DeleteResource(ctx context.Context, provider string, resourceID string) error {
	// Mock implementation - would call cloud provider API
	return nil
}

// GetResourceDetails gets detailed information about a resource
func (c *DefaultCloudService) GetResourceDetails(ctx context.Context, provider string, resourceID string) (*ResourceDetails, error) {
	// Mock implementation
	details := &ResourceDetails{
		Resource: Resource{
			ID:       resourceID,
			Name:     "example-resource",
			Type:     "compute",
			Provider: provider,
			Region:   "us-west-2",
			Status:   "running",
			Created:  time.Now().Add(-24 * time.Hour),
			Modified: time.Now(),
			Tags: map[string]string{
				"environment": "production",
				"team":        "devops",
			},
		},
		Configuration: map[string]interface{}{
			"instance_type": "t3.medium",
			"storage":       "20GB",
			"network":       "vpc-123456",
		},
		Dependencies: []string{"db-001", "lb-001"},
		SecurityGroups: []SecurityGroup{
			{
				ID:          "sg-123456",
				Name:        "web-server",
				Description: "Web server security group",
				Rules: []SecurityGroupRule{
					{
						Protocol:  "tcp",
						Port:      "80",
						Source:    "0.0.0.0/0",
						Direction: "inbound",
						Action:    "allow",
					},
				},
			},
		},
		NetworkInfo: NetworkInfo{
			VPC:        "vpc-123456",
			Subnet:     "subnet-123456",
			PrivateIPs: []string{"10.0.1.100"},
			PublicIPs:  []string{"54.123.45.67"},
			DNS:        []string{"example.com"},
		},
		Monitoring: MonitoringInfo{
			Enabled: true,
			Metrics: []MetricInfo{
				{
					Name:      "cpu_utilization",
					Value:     45.5,
					Unit:      "percent",
					Timestamp: time.Now(),
					Threshold: 80.0,
				},
			},
			Alerts: []AlertInfo{},
			Dashboards: []string{"resource-dashboard"},
		},
		Backup: BackupInfo{
			Enabled:       true,
			Schedule:      "daily",
			LastBackup:    time.Now().Add(-6 * time.Hour),
			NextBackup:    time.Now().Add(18 * time.Hour),
			RetentionDays: 30,
		},
	}
	
	return details, nil
}

// GetCostAnalysis provides cost analysis
func (c *DefaultCloudService) GetCostAnalysis(ctx context.Context, provider string, options CostOptions) (*CostAnalysis, error) {
	// Mock implementation
	analysis := &CostAnalysis{
		TotalCost: 1250.75,
		Currency:  "USD",
		Period:    "monthly",
		Breakdown: []CostBreakdown{
			{
				Category:      "compute",
				Cost:          800.50,
				Percentage:    64.0,
				ResourceCount: 10,
			},
			{
				Category:      "storage",
				Cost:          250.25,
				Percentage:    20.0,
				ResourceCount: 5,
			},
			{
				Category:      "network",
				Cost:          200.00,
				Percentage:    16.0,
				ResourceCount: 3,
			},
		},
		Trends: []CostTrend{
			{
				Date:   time.Now().Add(-30 * 24 * time.Hour),
				Cost:   1100.00,
				Change: -13.6,
			},
			{
				Date:   time.Now(),
				Cost:   1250.75,
				Change: 13.6,
			},
		},
		Recommendations: []CostRecommendation{
			{
				ID:          "cost-rec-001",
				Type:        "rightsizing",
				Title:       "Resize over-provisioned instances",
				Description: "Several instances are consistently under-utilized",
				Savings:     200.00,
				Effort:      "low",
				Risk:        "low",
				Actions:     []string{"Resize t3.large to t3.medium instances"},
			},
		},
	}
	
	return analysis, nil
}

// OptimizeResources optimizes cloud resources
func (c *DefaultCloudService) OptimizeResources(ctx context.Context, provider string, options OptimizeOptions) (*OptimizationResult, error) {
	// Mock implementation
	result := &OptimizationResult{
		ID:        "opt-001",
		Timestamp: time.Now(),
		Status:    "completed",
		Recommendations: []OptimizationRecommendation{
			{
				ResourceID: "i-123456789",
				Type:       "rightsizing",
				Current: map[string]interface{}{
					"instance_type": "t3.large",
					"vcpus":         2,
					"memory":        8,
				},
				Recommended: map[string]interface{}{
					"instance_type": "t3.medium",
					"vcpus":         1,
					"memory":        4,
				},
				Savings:    150.00,
				Confidence: 0.85,
				Actions:    []string{"Stop instance", "Change instance type", "Start instance"},
			},
		},
		PotentialSavings: 300.00,
		RiskAssessment:   "low",
	}
	
	return result, nil
}

// MonitorHealth monitors cloud resource health
func (c *DefaultCloudService) MonitorHealth(ctx context.Context, provider string) (<-chan HealthEvent, error) {
	events := make(chan HealthEvent, 100)
	
	// Mock implementation - would integrate with cloud provider health APIs
	go func() {
		defer close(events)
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				event := HealthEvent{
					ID:         fmt.Sprintf("health-%d", time.Now().Unix()),
					Timestamp:  time.Now(),
					Provider:   provider,
					ResourceID: "resource-001",
					Type:       "status_check",
					Status:     "healthy",
					Message:    "Resource is operating normally",
					Severity:   "info",
					Details: map[string]string{
						"cpu_usage":    "45%",
						"memory_usage": "60%",
						"disk_usage":   "30%",
					},
				}
				select {
				case events <- event:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	
	return events, nil
}

// Helper methods for different cloud providers

func (c *DefaultCloudService) listAWSResources(ctx context.Context, resourceType string) ([]Resource, error) {
	// Mock AWS resources
	resources := []Resource{
		{
			ID:       "i-1234567890abcdef0",
			Name:     "web-server-1",
			Type:     "ec2-instance",
			Provider: "aws",
			Region:   "us-west-2",
			Status:   "running",
			Created:  time.Now().Add(-72 * time.Hour),
			Modified: time.Now().Add(-1 * time.Hour),
			Tags: map[string]string{
				"Environment": "production",
				"Team":        "web",
			},
			Cost: &CostInfo{
				Monthly:     85.50,
				Daily:       2.85,
				Currency:    "USD",
				LastUpdated: time.Now(),
			},
		},
		{
			ID:       "vol-1234567890abcdef0",
			Name:     "web-server-1-root",
			Type:     "ebs-volume",
			Provider: "aws",
			Region:   "us-west-2",
			Status:   "in-use",
			Created:  time.Now().Add(-72 * time.Hour),
			Modified: time.Now().Add(-72 * time.Hour),
			Tags: map[string]string{
				"Environment": "production",
				"Team":        "web",
			},
			Cost: &CostInfo{
				Monthly:     10.00,
				Daily:       0.33,
				Currency:    "USD",
				LastUpdated: time.Now(),
			},
		},
	}
	
	return resources, nil
}

func (c *DefaultCloudService) listAzureResources(ctx context.Context, resourceType string) ([]Resource, error) {
	// Mock Azure resources
	resources := []Resource{
		{
			ID:       "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myRG/providers/Microsoft.Compute/virtualMachines/myVM",
			Name:     "myVM",
			Type:     "virtual-machine",
			Provider: "azure",
			Region:   "eastus",
			Status:   "running",
			Created:  time.Now().Add(-48 * time.Hour),
			Modified: time.Now().Add(-2 * time.Hour),
			Tags: map[string]string{
				"environment": "production",
				"team":        "backend",
			},
			Cost: &CostInfo{
				Monthly:     95.00,
				Daily:       3.17,
				Currency:    "USD",
				LastUpdated: time.Now(),
			},
		},
	}
	
	return resources, nil
}

func (c *DefaultCloudService) listGCPResources(ctx context.Context, resourceType string) ([]Resource, error) {
	// Mock GCP resources
	resources := []Resource{
		{
			ID:       "projects/my-project/zones/us-central1-a/instances/my-instance",
			Name:     "my-instance",
			Type:     "compute-instance",
			Provider: "gcp",
			Region:   "us-central1",
			Status:   "running",
			Created:  time.Now().Add(-24 * time.Hour),
			Modified: time.Now().Add(-30 * time.Minute),
			Tags: map[string]string{
				"env":  "prod",
				"team": "data",
			},
			Cost: &CostInfo{
				Monthly:     75.00,
				Daily:       2.50,
				Currency:    "USD",
				LastUpdated: time.Now(),
			},
		},
	}
	
	return resources, nil
}

// CloudManager manages multiple cloud providers
type CloudManager struct {
	providers map[string]CloudProvider
	mutex     sync.RWMutex
}

// NewCloudManager creates a new cloud manager
func NewCloudManager() *CloudManager {
	return &CloudManager{
		providers: make(map[string]CloudProvider),
	}
}

// AddProvider adds a cloud provider to the manager
func (m *CloudManager) AddProvider(provider CloudProvider) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.providers[provider.GetName()] = provider
	return nil
}

// GetProvider retrieves a cloud provider by name
func (m *CloudManager) GetProvider(name string) (CloudProvider, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	provider, exists := m.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider not found: %s", name)
	}
	return provider, nil
}

// ListProviders returns a list of all cloud providers
func (m *CloudManager) ListProviders() []CloudProvider {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	providers := make([]CloudProvider, 0, len(m.providers))
	for _, provider := range m.providers {
		providers = append(providers, provider)
	}
	return providers
}
