package cloud

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/option"
	"github.com/sirupsen/logrus"
)

// GCPProvider implements the CloudProvider interface for Google Cloud Platform
type GCPProvider struct {
	computeClient *compute.InstancesClient
	zonesClient   *compute.ZonesClient
	projectID     string
	config        *ProviderConfig
	connected     bool
	logger        *logrus.Logger
}

// NewGCPProvider creates a new GCP provider
func NewGCPProvider(cfg *ProviderConfig) (CloudProvider, error) {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	provider := &GCPProvider{
		config:    cfg,
		logger:    logger,
		projectID: cfg.ProjectID,
	}

	return provider, nil
}

// GetName returns the provider name
func (p *GCPProvider) GetName() string {
	return "Google Cloud Platform"
}

// GetType returns the provider type
func (p *GCPProvider) GetType() string {
	return "gcp"
}

// Connect establishes connection to GCP
func (p *GCPProvider) Connect(ctx context.Context) error {
	if p.connected {
		return nil
	}

	p.logger.Info("Connecting to Google Cloud Platform...")

	// Validate project ID
	if p.projectID == "" {
		return fmt.Errorf("project ID is required for GCP provider")
	}

	// Create compute client
	var opts []option.ClientOption
	if p.config.ServiceAccountPath != "" {
		opts = append(opts, option.WithCredentialsFile(p.config.ServiceAccountPath))
	}

	client, err := compute.NewInstancesRESTClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("failed to create GCP compute client: %w", err)
	}
	p.computeClient = client

	// Create zones client
	zonesClient, err := compute.NewZonesRESTClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("failed to create GCP zones client: %w", err)
	}
	p.zonesClient = zonesClient

	// Test connection
	if err := p.ValidateCredentials(ctx); err != nil {
		return fmt.Errorf("failed to validate GCP credentials: %w", err)
	}

	p.connected = true
	p.logger.Info("Successfully connected to Google Cloud Platform")
	return nil
}

// Disconnect closes the connection
func (p *GCPProvider) Disconnect(ctx context.Context) error {
	if p.computeClient != nil {
		p.computeClient.Close()
		p.computeClient = nil
	}
	if p.zonesClient != nil {
		p.zonesClient.Close()
		p.zonesClient = nil
	}
	p.connected = false
	p.logger.Info("Disconnected from Google Cloud Platform")
	return nil
}

// IsConnected returns connection status
func (p *GCPProvider) IsConnected() bool {
	return p.connected
}

// ValidateCredentials validates GCP credentials
func (p *GCPProvider) ValidateCredentials(ctx context.Context) error {
	if p.computeClient == nil {
		return fmt.Errorf("compute client not initialized")
	}

	// Try to list zones to test credentials
	req := &computepb.ListZonesRequest{
		Project: p.projectID,
	}

	it := p.zonesClient.List(ctx, req)
	_, err := it.Next()
	if err != nil {
		return fmt.Errorf("failed to validate credentials: %w", err)
	}

	return nil
}

// ListResources lists GCP resources
func (p *GCPProvider) ListResources(ctx context.Context, resourceType string) ([]*Resource, error) {
	if !p.connected {
		if err := p.Connect(ctx); err != nil {
			return nil, err
		}
	}

	switch strings.ToLower(resourceType) {
	case "instances", "vm", "vms":
		return p.listInstances(ctx)
	case "disks":
		return p.listDisks(ctx)
	case "networks":
		return p.listNetworks(ctx)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

// listInstances lists GCP compute instances
func (p *GCPProvider) listInstances(ctx context.Context) ([]*Resource, error) {
	var resources []*Resource

	// List all zones first
	zonesReq := &computepb.ListZonesRequest{
		Project: p.projectID,
	}

	zonesIt := p.zonesClient.List(ctx, zonesReq)
	for {
		zone, err := zonesIt.Next()
		if err != nil {
			break
		}

		// List instances in this zone
		req := &computepb.ListInstancesRequest{
			Project: p.projectID,
			Zone:    zone.GetName(),
		}

		it := p.computeClient.List(ctx, req)
		for {
			instance, err := it.Next()
			if err != nil {
				break
			}

			resource := &Resource{
				ID:       strconv.FormatUint(instance.GetId(), 10),
				Name:     instance.GetName(),
				Type:     "compute-instance",
				Provider: "gcp",
				Region:   p.getZoneRegion(zone.GetName()),
				State:    instance.GetStatus(),
				Status:   instance.GetStatus(),
				Created:  p.parseGCPTime(instance.GetCreationTimestamp()),
				Modified: time.Now(),
				Tags:     p.convertGCPLabels(instance.GetLabels()),
				Config: map[string]interface{}{
					"zone":          zone.GetName(),
					"machine_type":  p.getMachineType(instance.GetMachineType()),
					"cpu_platform":  instance.GetCpuPlatform(),
					"self_link":     instance.GetSelfLink(),
					"network_interfaces": len(instance.GetNetworkInterfaces()),
					"disks":         len(instance.GetDisks()),
					"can_ip_forward": instance.GetCanIpForward(),
					"scheduling":    p.getSchedulingInfo(instance.GetScheduling()),
				},
			}

			// Add metadata
			if instance.GetMetadata() != nil {
				resource.Metadata = p.convertGCPMetadata(instance.GetMetadata())
			}

			resources = append(resources, resource)
		}
	}

	return resources, nil
}

// listDisks lists GCP persistent disks
func (p *GCPProvider) listDisks(ctx context.Context) ([]*Resource, error) {
	// Note: This is a simplified implementation
	// In a real implementation, you would need to use the disk client
	return []*Resource{}, nil
}

// listNetworks lists GCP networks
func (p *GCPProvider) listNetworks(ctx context.Context) ([]*Resource, error) {
	// Note: This is a simplified implementation
	// In a real implementation, you would need to use the network client
	return []*Resource{}, nil
}

// GetResourceDetails gets detailed information about a resource
func (p *GCPProvider) GetResourceDetails(ctx context.Context, resourceID string) (*Resource, error) {
	if !p.connected {
		if err := p.Connect(ctx); err != nil {
			return nil, err
		}
	}

	// For GCP, we need to determine the zone and instance name
	// This is a simplified implementation
	return nil, fmt.Errorf("GetResourceDetails not fully implemented for GCP provider")
}

// Helper methods for GCP resource conversion
func (p *GCPProvider) getZoneRegion(zone string) string {
	// Extract region from zone (e.g., "us-central1-a" -> "us-central1")
	parts := strings.Split(zone, "-")
	if len(parts) >= 3 {
		return strings.Join(parts[:2], "-")
	}
	return zone
}

func (p *GCPProvider) getMachineType(machineType string) string {
	// Extract machine type from URL
	parts := strings.Split(machineType, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return machineType
}

func (p *GCPProvider) parseGCPTime(timestamp string) time.Time {
	// Parse GCP timestamp format (RFC3339)
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return time.Now()
	}
	return t
}

func (p *GCPProvider) convertGCPLabels(labels map[string]string) map[string]string {
	if labels == nil {
		return make(map[string]string)
	}
	return labels
}

func (p *GCPProvider) convertGCPMetadata(metadata *computepb.Metadata) map[string]string {
	result := make(map[string]string)
	if metadata == nil {
		return result
	}

	for _, item := range metadata.GetItems() {
		if item.GetKey() != "" && item.GetValue() != "" {
			result[item.GetKey()] = item.GetValue()
		}
	}
	return result
}

func (p *GCPProvider) getSchedulingInfo(scheduling *computepb.Scheduling) map[string]interface{} {
	if scheduling == nil {
		return make(map[string]interface{})
	}

	return map[string]interface{}{
		"automatic_restart":   scheduling.GetAutomaticRestart(),
		"on_host_maintenance": scheduling.GetOnHostMaintenance(),
		"preemptible":         scheduling.GetPreemptible(),
	}
}

// Additional methods to implement CloudProvider interface
func (p *GCPProvider) CreateResource(ctx context.Context, req *CreateResourceRequest) (*Resource, error) {
	return nil, fmt.Errorf("CreateResource not implemented for GCP provider")
}

func (p *GCPProvider) UpdateResource(ctx context.Context, req *UpdateResourceRequest) (*Resource, error) {
	return nil, fmt.Errorf("UpdateResource not implemented for GCP provider")
}

func (p *GCPProvider) DeleteResource(ctx context.Context, resourceID string) error {
	return fmt.Errorf("DeleteResource not implemented for GCP provider")
}

func (p *GCPProvider) GetMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, error) {
	return nil, fmt.Errorf("GetMetrics not implemented for GCP provider")
}

func (p *GCPProvider) GetCost(ctx context.Context, req *CostRequest) (*CostResponse, error) {
	return nil, fmt.Errorf("GetCost not implemented for GCP provider")
}

func (p *GCPProvider) GetConfiguration() *ProviderConfig {
	return p.config
}

func (p *GCPProvider) UpdateConfiguration(config *ProviderConfig) error {
	p.config = config
	if config.ProjectID != "" {
		p.projectID = config.ProjectID
	}
	return nil
}

func (p *GCPProvider) GetStatus() *ProviderStatus {
	status := &ProviderStatus{
		Name:      p.GetName(),
		Type:      p.GetType(),
		Connected: p.connected,
		LastCheck: time.Now(),
	}
	if p.connected {
		status.Status = "connected"
	} else {
		status.Status = "disconnected"
	}
	return status
}

func (p *GCPProvider) GetRegions(ctx context.Context) ([]string, error) {
	// GCP regions are well-known, return common ones
	return []string{
		"us-central1",
		"us-east1",
		"us-east4",
		"us-west1",
		"us-west2",
		"us-west3",
		"us-west4",
		"northamerica-northeast1",
		"northamerica-northeast2",
		"southamerica-east1",
		"southamerica-west1",
		"europe-central2",
		"europe-north1",
		"europe-southwest1",
		"europe-west1",
		"europe-west2",
		"europe-west3",
		"europe-west4",
		"europe-west6",
		"europe-west8",
		"europe-west9",
		"europe-west10",
		"europe-west12",
		"asia-east1",
		"asia-east2",
		"asia-northeast1",
		"asia-northeast2",
		"asia-northeast3",
		"asia-south1",
		"asia-south2",
		"asia-southeast1",
		"asia-southeast2",
		"australia-southeast1",
		"australia-southeast2",
		"me-central1",
		"me-west1",
		"africa-south1",
	}, nil
}

func (p *GCPProvider) GetResourceTypes(ctx context.Context) ([]string, error) {
	return []string{
		"instances",
		"vm",
		"vms",
		"disks",
		"networks",
	}, nil
}
