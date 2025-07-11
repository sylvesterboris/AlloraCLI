package cloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/sirupsen/logrus"
)

// AzureProvider implements the CloudProvider interface for Azure
type AzureProvider struct {
	credential      azcore.TokenCredential
	computeClient   *armcompute.VirtualMachinesClient
	networkClient   *armnetwork.VirtualNetworksClient
	resourceClient  *armresources.Client
	subscriptionID  string
	config          *ProviderConfig
	connected       bool
	logger          *logrus.Logger
}

// NewAzureProvider creates a new Azure provider
func NewAzureProvider(cfg *ProviderConfig) (CloudProvider, error) {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	provider := &AzureProvider{
		config:         cfg,
		logger:         logger,
		subscriptionID: cfg.SubscriptionID,
	}

	return provider, nil
}

// GetName returns the provider name
func (p *AzureProvider) GetName() string {
	return "Azure"
}

// GetType returns the provider type
func (p *AzureProvider) GetType() string {
	return "azure"
}

// Connect establishes connection to Azure
func (p *AzureProvider) Connect(ctx context.Context) error {
	if p.connected {
		return nil
	}

	p.logger.Info("Connecting to Azure...")

	// Create credential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("failed to create Azure credential: %w", err)
	}
	p.credential = cred

	// Validate subscription ID
	if p.subscriptionID == "" {
		return fmt.Errorf("subscription ID is required for Azure provider")
	}

	// Create clients
	clientFactory, err := armcompute.NewClientFactory(p.subscriptionID, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create Azure compute client factory: %w", err)
	}
	p.computeClient = clientFactory.NewVirtualMachinesClient()

	networkClientFactory, err := armnetwork.NewClientFactory(p.subscriptionID, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create Azure network client factory: %w", err)
	}
	p.networkClient = networkClientFactory.NewVirtualNetworksClient()

	resourceClientFactory, err := armresources.NewClientFactory(p.subscriptionID, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create Azure resource client factory: %w", err)
	}
	p.resourceClient = resourceClientFactory.NewClient()

	// Test connection
	if err := p.ValidateCredentials(ctx); err != nil {
		return fmt.Errorf("failed to validate Azure credentials: %w", err)
	}

	p.connected = true
	p.logger.Info("Successfully connected to Azure")
	return nil
}

// Disconnect closes the connection
func (p *AzureProvider) Disconnect(ctx context.Context) error {
	p.connected = false
	p.logger.Info("Disconnected from Azure")
	return nil
}

// IsConnected returns connection status
func (p *AzureProvider) IsConnected() bool {
	return p.connected
}

// ValidateCredentials validates Azure credentials
func (p *AzureProvider) ValidateCredentials(ctx context.Context) error {
	if p.resourceClient == nil {
		return fmt.Errorf("resource client not initialized")
	}

	// Try to list resource groups to test credentials
	pager := p.resourceClient.NewListPager(nil)
	if !pager.More() {
		return fmt.Errorf("failed to access subscription resources")
	}

	_, err := pager.NextPage(ctx)
	if err != nil {
		return fmt.Errorf("failed to validate credentials: %w", err)
	}

	return nil
}

// ListResources lists Azure resources
func (p *AzureProvider) ListResources(ctx context.Context, resourceType string) ([]*Resource, error) {
	if !p.connected {
		if err := p.Connect(ctx); err != nil {
			return nil, err
		}
	}

	switch strings.ToLower(resourceType) {
	case "vm", "virtualmachines", "vms":
		return p.listVirtualMachines(ctx)
	case "vnets", "virtualnetworks", "networks":
		return p.listVirtualNetworks(ctx)
	case "resourcegroups", "rg":
		return p.listResourceGroups(ctx)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

// listVirtualMachines lists Azure virtual machines
func (p *AzureProvider) listVirtualMachines(ctx context.Context) ([]*Resource, error) {
	var resources []*Resource

	// List all resource groups first
	rgPager := p.resourceClient.NewListPager(nil)
	for rgPager.More() {
		page, err := rgPager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list resource groups: %w", err)
		}

		for _, rg := range page.Value {
			if rg.Name == nil {
				continue
			}

			// List VMs in this resource group
			vmPager := p.computeClient.NewListPager(*rg.Name, nil)
			for vmPager.More() {
				vmPage, err := vmPager.NextPage(ctx)
				if err != nil {
					p.logger.Warnf("Failed to list VMs in resource group %s: %v", *rg.Name, err)
					continue
				}

				for _, vm := range vmPage.Value {
					if vm.Name == nil || vm.ID == nil {
						continue
					}

					resource := &Resource{
						ID:       *vm.ID,
						Name:     *vm.Name,
						Type:     "virtual-machine",
						Provider: "azure",
						Region:   p.getStringValue(vm.Location),
						State:    p.getVMState(vm),
						Status:   p.getVMState(vm),
						Created:  time.Now(), // Azure doesn't provide creation time in list operation
						Modified: time.Now(),
						Tags:     p.convertAzureTags(vm.Tags),
						Config: map[string]interface{}{
							"resource_group": *rg.Name,
							"vm_size":        p.getVMSize(vm),
							"os_type":        p.getOSType(vm),
							"location":       p.getStringValue(vm.Location),
						},
					}
					resources = append(resources, resource)
				}
			}
		}
	}

	return resources, nil
}

// listVirtualNetworks lists Azure virtual networks
func (p *AzureProvider) listVirtualNetworks(ctx context.Context) ([]*Resource, error) {
	var resources []*Resource

	// List all resource groups first
	rgPager := p.resourceClient.NewListPager(nil)
	for rgPager.More() {
		page, err := rgPager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list resource groups: %w", err)
		}

		for _, rg := range page.Value {
			if rg.Name == nil {
				continue
			}

			// List VNets in this resource group
			vnetPager := p.networkClient.NewListPager(*rg.Name, nil)
			for vnetPager.More() {
				vnetPage, err := vnetPager.NextPage(ctx)
				if err != nil {
					p.logger.Warnf("Failed to list VNets in resource group %s: %v", *rg.Name, err)
					continue
				}

				for _, vnet := range vnetPage.Value {
					if vnet.Name == nil || vnet.ID == nil {
						continue
					}

					resource := &Resource{
						ID:       *vnet.ID,
						Name:     *vnet.Name,
						Type:     "virtual-network",
						Provider: "azure",
						Region:   p.getStringValue(vnet.Location),
						State:    p.getVNetState(vnet),
						Status:   p.getVNetState(vnet),
						Created:  time.Now(),
						Modified: time.Now(),
						Tags:     p.convertAzureTags(vnet.Tags),
						Config: map[string]interface{}{
							"resource_group":  *rg.Name,
							"location":        p.getStringValue(vnet.Location),
							"address_spaces":  p.getAddressSpaces(vnet),
							"subnets_count":   p.getSubnetsCount(vnet),
						},
					}
					resources = append(resources, resource)
				}
			}
		}
	}

	return resources, nil
}

// listResourceGroups lists Azure resource groups
func (p *AzureProvider) listResourceGroups(ctx context.Context) ([]*Resource, error) {
	var resources []*Resource

	pager := p.resourceClient.NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list resource groups: %w", err)
		}

		for _, rg := range page.Value {
			if rg.Name == nil || rg.ID == nil {
				continue
			}

			resource := &Resource{
				ID:       *rg.ID,
				Name:     *rg.Name,
				Type:     "resource-group",
				Provider: "azure",
				Region:   p.getStringValue(rg.Location),
				State:    p.getGenericResourceState(rg),
				Status:   p.getGenericResourceState(rg),
				Created:  time.Now(),
				Modified: time.Now(),
				Tags:     p.convertAzureTags(rg.Tags),
				Config: map[string]interface{}{
					"location":         p.getStringValue(rg.Location),
					"provisioning_state": p.getGenericResourceProvisioningState(rg),
				},
			}
			resources = append(resources, resource)
		}
	}

	return resources, nil
}

// GetResourceDetails gets detailed information about a resource
func (p *AzureProvider) GetResourceDetails(ctx context.Context, resourceID string) (*Resource, error) {
	if !p.connected {
		if err := p.Connect(ctx); err != nil {
			return nil, err
		}
	}

	// Parse resource ID to determine type
	if strings.Contains(resourceID, "/virtualMachines/") {
		return p.getVirtualMachineDetails(ctx, resourceID)
	} else if strings.Contains(resourceID, "/virtualNetworks/") {
		return p.getVirtualNetworkDetails(ctx, resourceID)
	} else if strings.Contains(resourceID, "/resourceGroups/") && !strings.Contains(resourceID, "/providers/") {
		return p.getResourceGroupDetails(ctx, resourceID)
	}

	return nil, fmt.Errorf("unsupported resource type for ID: %s", resourceID)
}

// getVirtualMachineDetails gets detailed VM information
func (p *AzureProvider) getVirtualMachineDetails(ctx context.Context, resourceID string) (*Resource, error) {
	// Parse resource ID to get resource group and VM name
	parts := strings.Split(resourceID, "/")
	if len(parts) < 9 {
		return nil, fmt.Errorf("invalid resource ID format: %s", resourceID)
	}

	resourceGroup := parts[4]
	vmName := parts[8]

	resp, err := p.computeClient.Get(ctx, resourceGroup, vmName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get VM details: %w", err)
	}

	vm := resp.VirtualMachine
	resource := &Resource{
		ID:       *vm.ID,
		Name:     *vm.Name,
		Type:     "virtual-machine",
		Provider: "azure",
		Region:   p.getStringValue(vm.Location),
		State:    p.getVMState(&vm),
		Status:   p.getVMState(&vm),
		Created:  time.Now(),
		Modified: time.Now(),
		Tags:     p.convertAzureTags(vm.Tags),
		Config: map[string]interface{}{
			"resource_group":     resourceGroup,
			"vm_size":           p.getVMSize(&vm),
			"os_type":           p.getOSType(&vm),
			"location":          p.getStringValue(vm.Location),
			"provisioning_state": p.getVMProvisioningState(&vm),
		},
	}

	return resource, nil
}

// Helper methods for Azure resource conversion
func (p *AzureProvider) getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func (p *AzureProvider) getVMState(vm *armcompute.VirtualMachine) string {
	if vm.Properties != nil && vm.Properties.ProvisioningState != nil {
		return *vm.Properties.ProvisioningState
	}
	return "unknown"
}

func (p *AzureProvider) getVMSize(vm *armcompute.VirtualMachine) string {
	if vm.Properties != nil && vm.Properties.HardwareProfile != nil && vm.Properties.HardwareProfile.VMSize != nil {
		return string(*vm.Properties.HardwareProfile.VMSize)
	}
	return "unknown"
}

func (p *AzureProvider) getOSType(vm *armcompute.VirtualMachine) string {
	if vm.Properties != nil && vm.Properties.StorageProfile != nil && vm.Properties.StorageProfile.OSDisk != nil && vm.Properties.StorageProfile.OSDisk.OSType != nil {
		return string(*vm.Properties.StorageProfile.OSDisk.OSType)
	}
	return "unknown"
}

func (p *AzureProvider) getVMProvisioningState(vm *armcompute.VirtualMachine) string {
	if vm.Properties != nil && vm.Properties.ProvisioningState != nil {
		return *vm.Properties.ProvisioningState
	}
	return "unknown"
}

func (p *AzureProvider) getVNetState(vnet *armnetwork.VirtualNetwork) string {
	if vnet.Properties != nil && vnet.Properties.ProvisioningState != nil {
		return string(*vnet.Properties.ProvisioningState)
	}
	return "unknown"
}

func (p *AzureProvider) getAddressSpaces(vnet *armnetwork.VirtualNetwork) []string {
	if vnet.Properties != nil && vnet.Properties.AddressSpace != nil && vnet.Properties.AddressSpace.AddressPrefixes != nil {
		var spaces []string
		for _, prefix := range vnet.Properties.AddressSpace.AddressPrefixes {
			if prefix != nil {
				spaces = append(spaces, *prefix)
			}
		}
		return spaces
	}
	return []string{}
}

func (p *AzureProvider) getSubnetsCount(vnet *armnetwork.VirtualNetwork) int {
	if vnet.Properties != nil && vnet.Properties.Subnets != nil {
		return len(vnet.Properties.Subnets)
	}
	return 0
}

func (p *AzureProvider) getResourceGroupState(rg *armresources.ResourceGroup) string {
	if rg.Properties != nil && rg.Properties.ProvisioningState != nil {
		return *rg.Properties.ProvisioningState
	}
	return "unknown"
}

func (p *AzureProvider) getResourceGroupProvisioningState(rg *armresources.ResourceGroup) string {
	if rg.Properties != nil && rg.Properties.ProvisioningState != nil {
		return *rg.Properties.ProvisioningState
	}
	return "unknown"
}

func (p *AzureProvider) getGenericResourceState(rg *armresources.GenericResourceExpanded) string {
	if rg.Properties != nil {
		if state, ok := rg.Properties.(map[string]interface{}); ok {
			if provisioningState, ok := state["provisioningState"].(string); ok {
				return provisioningState
			}
		}
	}
	return "unknown"
}

func (p *AzureProvider) getGenericResourceProvisioningState(rg *armresources.GenericResourceExpanded) string {
	if rg.Properties != nil {
		if state, ok := rg.Properties.(map[string]interface{}); ok {
			if provisioningState, ok := state["provisioningState"].(string); ok {
				return provisioningState
			}
		}
	}
	return "unknown"
}

func (p *AzureProvider) convertAzureTags(tags map[string]*string) map[string]string {
	result := make(map[string]string)
	for k, v := range tags {
		if v != nil {
			result[k] = *v
		}
	}
	return result
}

// Additional methods to implement CloudProvider interface
func (p *AzureProvider) CreateResource(ctx context.Context, req *CreateResourceRequest) (*Resource, error) {
	return nil, fmt.Errorf("CreateResource not implemented for Azure provider")
}

func (p *AzureProvider) UpdateResource(ctx context.Context, req *UpdateResourceRequest) (*Resource, error) {
	return nil, fmt.Errorf("UpdateResource not implemented for Azure provider")
}

func (p *AzureProvider) DeleteResource(ctx context.Context, resourceID string) error {
	return fmt.Errorf("DeleteResource not implemented for Azure provider")
}

func (p *AzureProvider) GetMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, error) {
	return nil, fmt.Errorf("GetMetrics not implemented for Azure provider")
}

func (p *AzureProvider) GetCost(ctx context.Context, req *CostRequest) (*CostResponse, error) {
	return nil, fmt.Errorf("GetCost not implemented for Azure provider")
}

func (p *AzureProvider) GetConfiguration() *ProviderConfig {
	return p.config
}

func (p *AzureProvider) UpdateConfiguration(config *ProviderConfig) error {
	p.config = config
	if config.SubscriptionID != "" {
		p.subscriptionID = config.SubscriptionID
	}
	return nil
}

func (p *AzureProvider) GetStatus() *ProviderStatus {
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

func (p *AzureProvider) GetRegions(ctx context.Context) ([]string, error) {
	// Azure regions are well-known, return common ones
	return []string{
		"eastus",
		"eastus2",
		"westus",
		"westus2",
		"westus3",
		"centralus",
		"northcentralus",
		"southcentralus",
		"westcentralus",
		"canadacentral",
		"canadaeast",
		"brazilsouth",
		"northeurope",
		"westeurope",
		"uksouth",
		"ukwest",
		"francecentral",
		"francesouth",
		"germanynorth",
		"germanywestcentral",
		"norwayeast",
		"norwaywest",
		"switzerlandnorth",
		"switzerlandwest",
		"swedencentral",
		"swedensouth",
		"uaenorth",
		"uaecentral",
		"southafricanorth",
		"southafricawest",
		"australiaeast",
		"australiasoutheast",
		"australiacentral",
		"australiacentral2",
		"eastasia",
		"southeastasia",
		"japaneast",
		"japanwest",
		"koreacentral",
		"koreasouth",
		"southindia",
		"centralindia",
		"westindia",
		"jioindiawest",
		"jioindiacentral",
	}, nil
}

func (p *AzureProvider) GetResourceTypes(ctx context.Context) ([]string, error) {
	return []string{
		"vm",
		"virtualmachines",
		"vms",
		"vnets",
		"virtualnetworks",
		"networks",
		"resourcegroups",
		"rg",
	}, nil
}

// Helper methods for additional resource details
func (p *AzureProvider) getVirtualNetworkDetails(ctx context.Context, resourceID string) (*Resource, error) {
	// Parse resource ID to get resource group and VNet name
	parts := strings.Split(resourceID, "/")
	if len(parts) < 9 {
		return nil, fmt.Errorf("invalid resource ID format: %s", resourceID)
	}

	resourceGroup := parts[4]
	vnetName := parts[8]

	resp, err := p.networkClient.Get(ctx, resourceGroup, vnetName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get VNet details: %w", err)
	}

	vnet := resp.VirtualNetwork
	resource := &Resource{
		ID:       *vnet.ID,
		Name:     *vnet.Name,
		Type:     "virtual-network",
		Provider: "azure",
		Region:   p.getStringValue(vnet.Location),
		State:    p.getVNetState(&vnet),
		Status:   p.getVNetState(&vnet),
		Created:  time.Now(),
		Modified: time.Now(),
		Tags:     p.convertAzureTags(vnet.Tags),
		Config: map[string]interface{}{
			"resource_group":     resourceGroup,
			"location":          p.getStringValue(vnet.Location),
			"address_spaces":    p.getAddressSpaces(&vnet),
			"subnets_count":     p.getSubnetsCount(&vnet),
			"provisioning_state": p.getVNetState(&vnet),
		},
	}

	return resource, nil
}

func (p *AzureProvider) getResourceGroupDetails(ctx context.Context, resourceID string) (*Resource, error) {
	// For simplicity, we'll use the list API to find the resource
	resources, err := p.listResourceGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list resource groups: %w", err)
	}

	for _, resource := range resources {
		if resource.ID == resourceID {
			return resource, nil
		}
	}

	return nil, fmt.Errorf("resource group not found: %s", resourceID)
}
