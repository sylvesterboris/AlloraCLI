package cloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/sirupsen/logrus"
)

// AWSProvider implements the CloudProvider interface for AWS
type AWSProvider struct {
	ec2Client *ec2.Client
	stsClient *sts.Client
	config    *ProviderConfig
	connected bool
	logger    *logrus.Logger
}

// NewAWSProvider creates a new AWS provider
func NewAWSProvider(cfg *ProviderConfig) (CloudProvider, error) {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	provider := &AWSProvider{
		config: cfg,
		logger: logger,
	}

	return provider, nil
}

// GetName returns the provider name
func (p *AWSProvider) GetName() string {
	return "AWS"
}

// GetType returns the provider type
func (p *AWSProvider) GetType() string {
	return "aws"
}

// Connect establishes connection to AWS
func (p *AWSProvider) Connect(ctx context.Context) error {
	if p.connected {
		return nil
	}

	p.logger.Info("Connecting to AWS...")

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Override region if specified
	if p.config.Region != "" {
		cfg.Region = p.config.Region
	}

	// Create EC2 client
	p.ec2Client = ec2.NewFromConfig(cfg)
	p.stsClient = sts.NewFromConfig(cfg)

	// Test connection
	if err := p.ValidateCredentials(ctx); err != nil {
		return fmt.Errorf("failed to validate AWS credentials: %w", err)
	}

	p.connected = true
	p.logger.Info("Successfully connected to AWS")
	return nil
}

// Disconnect closes the connection
func (p *AWSProvider) Disconnect(ctx context.Context) error {
	p.connected = false
	p.logger.Info("Disconnected from AWS")
	return nil
}

// IsConnected returns connection status
func (p *AWSProvider) IsConnected() bool {
	return p.connected
}

// ValidateCredentials validates AWS credentials
func (p *AWSProvider) ValidateCredentials(ctx context.Context) error {
	if p.stsClient == nil {
		return fmt.Errorf("STS client not initialized")
	}

	_, err := p.stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return fmt.Errorf("failed to get caller identity: %w", err)
	}

	return nil
}

// ListResources lists AWS resources
func (p *AWSProvider) ListResources(ctx context.Context, resourceType string) ([]*Resource, error) {
	if !p.connected {
		if err := p.Connect(ctx); err != nil {
			return nil, err
		}
	}

	switch strings.ToLower(resourceType) {
	case "ec2", "instances":
		return p.listEC2Instances(ctx)
	case "volumes", "ebs":
		return p.listEBSVolumes(ctx)
	case "security-groups", "sg":
		return p.listSecurityGroups(ctx)
	case "vpcs", "vpc":
		return p.listVPCs(ctx)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

// listEC2Instances lists EC2 instances
func (p *AWSProvider) listEC2Instances(ctx context.Context) ([]*Resource, error) {
	input := &ec2.DescribeInstancesInput{}
	result, err := p.ec2Client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %w", err)
	}

	var resources []*Resource
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			resource := &Resource{
				ID:       aws.ToString(instance.InstanceId),
				Name:     p.getInstanceName(instance),
				Type:     "ec2-instance",
				Provider: "aws",
				Region:   aws.ToString(instance.Placement.AvailabilityZone),
				State:    string(instance.State.Name),
				Status:   string(instance.State.Name),
				Created:  aws.ToTime(instance.LaunchTime),
				Modified: time.Now(),
				Tags:     p.convertEC2Tags(instance.Tags),
				Config: map[string]interface{}{
					"instance_type":   string(instance.InstanceType),
					"architecture":    string(instance.Architecture),
					"platform":        aws.ToString(instance.PlatformDetails),
					"vpc_id":          aws.ToString(instance.VpcId),
					"subnet_id":       aws.ToString(instance.SubnetId),
					"public_ip":       aws.ToString(instance.PublicIpAddress),
					"private_ip":      aws.ToString(instance.PrivateIpAddress),
					"security_groups": p.getSecurityGroupNames(instance.SecurityGroups),
				},
			}
			resources = append(resources, resource)
		}
	}

	return resources, nil
}

// listEBSVolumes lists EBS volumes
func (p *AWSProvider) listEBSVolumes(ctx context.Context) ([]*Resource, error) {
	input := &ec2.DescribeVolumesInput{}
	result, err := p.ec2Client.DescribeVolumes(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe volumes: %w", err)
	}

	var resources []*Resource
	for _, volume := range result.Volumes {
		resource := &Resource{
			ID:       aws.ToString(volume.VolumeId),
			Name:     p.getVolumeName(volume),
			Type:     "ebs-volume",
			Provider: "aws",
			Region:   aws.ToString(volume.AvailabilityZone),
			State:    string(volume.State),
			Status:   string(volume.State),
			Created:  aws.ToTime(volume.CreateTime),
			Modified: time.Now(),
			Tags:     p.convertEBSVolumeTags(volume.Tags),
			Config: map[string]interface{}{
				"volume_type": string(volume.VolumeType),
				"size":        aws.ToInt32(volume.Size),
				"iops":        aws.ToInt32(volume.Iops),
				"throughput":  aws.ToInt32(volume.Throughput),
				"encrypted":   aws.ToBool(volume.Encrypted),
				"snapshot_id": aws.ToString(volume.SnapshotId),
			},
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

// listSecurityGroups lists security groups
func (p *AWSProvider) listSecurityGroups(ctx context.Context) ([]*Resource, error) {
	input := &ec2.DescribeSecurityGroupsInput{}
	result, err := p.ec2Client.DescribeSecurityGroups(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe security groups: %w", err)
	}

	var resources []*Resource
	for _, sg := range result.SecurityGroups {
		resource := &Resource{
			ID:       aws.ToString(sg.GroupId),
			Name:     aws.ToString(sg.GroupName),
			Type:     "security-group",
			Provider: "aws",
			Region:   "", // Security groups don't have a specific region in the response
			State:    "available",
			Status:   "available",
			Created:  time.Now(), // AWS doesn't provide creation time for security groups
			Modified: time.Now(),
			Tags:     p.convertSecurityGroupTags(sg.Tags),
			Config: map[string]interface{}{
				"description": aws.ToString(sg.Description),
				"vpc_id":      aws.ToString(sg.VpcId),
				"owner_id":    aws.ToString(sg.OwnerId),
				"rules_count": len(sg.IpPermissions) + len(sg.IpPermissionsEgress),
			},
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

// listVPCs lists VPCs
func (p *AWSProvider) listVPCs(ctx context.Context) ([]*Resource, error) {
	input := &ec2.DescribeVpcsInput{}
	result, err := p.ec2Client.DescribeVpcs(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe VPCs: %w", err)
	}

	var resources []*Resource
	for _, vpc := range result.Vpcs {
		resource := &Resource{
			ID:       aws.ToString(vpc.VpcId),
			Name:     p.getVPCName(vpc),
			Type:     "vpc",
			Provider: "aws",
			Region:   "", // VPCs don't have a specific region in the response
			State:    string(vpc.State),
			Status:   string(vpc.State),
			Created:  time.Now(), // AWS doesn't provide creation time for VPCs
			Modified: time.Now(),
			Tags:     p.convertVPCTags(vpc.Tags),
			Config: map[string]interface{}{
				"cidr_block":           aws.ToString(vpc.CidrBlock),
				"dhcp_options_id":      aws.ToString(vpc.DhcpOptionsId),
				"instance_tenancy":     string(vpc.InstanceTenancy),
				"is_default":           aws.ToBool(vpc.IsDefault),
				"ipv6_cidr_block_sets": len(vpc.Ipv6CidrBlockAssociationSet),
				"owner_id":             aws.ToString(vpc.OwnerId),
			},
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

// GetResourceDetails gets detailed information about a resource
func (p *AWSProvider) GetResourceDetails(ctx context.Context, resourceID string) (*Resource, error) {
	if !p.connected {
		if err := p.Connect(ctx); err != nil {
			return nil, err
		}
	}

	// Try to identify resource type based on ID format
	if strings.HasPrefix(resourceID, "i-") {
		return p.getEC2InstanceDetails(ctx, resourceID)
	} else if strings.HasPrefix(resourceID, "vol-") {
		return p.getEBSVolumeDetails(ctx, resourceID)
	} else if strings.HasPrefix(resourceID, "sg-") {
		return p.getSecurityGroupDetails(ctx, resourceID)
	} else if strings.HasPrefix(resourceID, "vpc-") {
		return p.getVPCDetails(ctx, resourceID)
	}

	return nil, fmt.Errorf("unsupported resource ID format: %s", resourceID)
}

// getEC2InstanceDetails gets detailed EC2 instance information
func (p *AWSProvider) getEC2InstanceDetails(ctx context.Context, instanceID string) (*Resource, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}
	result, err := p.ec2Client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instance %s: %w", instanceID, err)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance %s not found", instanceID)
	}

	instance := result.Reservations[0].Instances[0]
	resource := &Resource{
		ID:       aws.ToString(instance.InstanceId),
		Name:     p.getInstanceName(instance),
		Type:     "ec2-instance",
		Provider: "aws",
		Region:   aws.ToString(instance.Placement.AvailabilityZone),
		State:    string(instance.State.Name),
		Status:   string(instance.State.Name),
		Created:  aws.ToTime(instance.LaunchTime),
		Modified: time.Now(),
		Tags:     p.convertEC2Tags(instance.Tags),
		Config: map[string]interface{}{
			"instance_type":     string(instance.InstanceType),
			"architecture":      string(instance.Architecture),
			"platform":          aws.ToString(instance.PlatformDetails),
			"vpc_id":            aws.ToString(instance.VpcId),
			"subnet_id":         aws.ToString(instance.SubnetId),
			"public_ip":         aws.ToString(instance.PublicIpAddress),
			"private_ip":        aws.ToString(instance.PrivateIpAddress),
			"public_dns":        aws.ToString(instance.PublicDnsName),
			"private_dns":       aws.ToString(instance.PrivateDnsName),
			"security_groups":   p.getSecurityGroupNames(instance.SecurityGroups),
			"key_name":          aws.ToString(instance.KeyName),
			"image_id":          aws.ToString(instance.ImageId),
			"monitoring":        string(instance.Monitoring.State),
			"source_dest_check": aws.ToBool(instance.SourceDestCheck),
			"virtualization":    string(instance.VirtualizationType),
			"root_device_name":  aws.ToString(instance.RootDeviceName),
			"root_device_type":  string(instance.RootDeviceType),
		},
	}

	return resource, nil
}

// Helper methods for converting AWS types to our types
func (p *AWSProvider) getInstanceName(instance types.Instance) string {
	for _, tag := range instance.Tags {
		if aws.ToString(tag.Key) == "Name" {
			return aws.ToString(tag.Value)
		}
	}
	return aws.ToString(instance.InstanceId)
}

func (p *AWSProvider) getVolumeName(volume types.Volume) string {
	for _, tag := range volume.Tags {
		if aws.ToString(tag.Key) == "Name" {
			return aws.ToString(tag.Value)
		}
	}
	return aws.ToString(volume.VolumeId)
}

func (p *AWSProvider) getVPCName(vpc types.Vpc) string {
	for _, tag := range vpc.Tags {
		if aws.ToString(tag.Key) == "Name" {
			return aws.ToString(tag.Value)
		}
	}
	return aws.ToString(vpc.VpcId)
}

func (p *AWSProvider) convertEC2Tags(tags []types.Tag) map[string]string {
	result := make(map[string]string)
	for _, tag := range tags {
		result[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}
	return result
}

func (p *AWSProvider) convertEBSVolumeTags(tags []types.Tag) map[string]string {
	result := make(map[string]string)
	for _, tag := range tags {
		result[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}
	return result
}

func (p *AWSProvider) convertSecurityGroupTags(tags []types.Tag) map[string]string {
	result := make(map[string]string)
	for _, tag := range tags {
		result[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}
	return result
}

func (p *AWSProvider) convertVPCTags(tags []types.Tag) map[string]string {
	result := make(map[string]string)
	for _, tag := range tags {
		result[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}
	return result
}

func (p *AWSProvider) getSecurityGroupNames(groups []types.GroupIdentifier) []string {
	var names []string
	for _, group := range groups {
		names = append(names, aws.ToString(group.GroupName))
	}
	return names
}

// Additional methods to implement CloudProvider interface
func (p *AWSProvider) CreateResource(ctx context.Context, req *CreateResourceRequest) (*Resource, error) {
	return nil, fmt.Errorf("CreateResource not implemented for AWS provider")
}

func (p *AWSProvider) UpdateResource(ctx context.Context, req *UpdateResourceRequest) (*Resource, error) {
	return nil, fmt.Errorf("UpdateResource not implemented for AWS provider")
}

func (p *AWSProvider) DeleteResource(ctx context.Context, resourceID string) error {
	return fmt.Errorf("DeleteResource not implemented for AWS provider")
}

func (p *AWSProvider) GetMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, error) {
	return nil, fmt.Errorf("GetMetrics not implemented for AWS provider")
}

func (p *AWSProvider) GetCost(ctx context.Context, req *CostRequest) (*CostResponse, error) {
	return nil, fmt.Errorf("GetCost not implemented for AWS provider")
}

func (p *AWSProvider) GetConfiguration() *ProviderConfig {
	return p.config
}

func (p *AWSProvider) UpdateConfiguration(config *ProviderConfig) error {
	p.config = config
	return nil
}

func (p *AWSProvider) GetStatus() *ProviderStatus {
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

func (p *AWSProvider) GetRegions(ctx context.Context) ([]string, error) {
	input := &ec2.DescribeRegionsInput{}
	result, err := p.ec2Client.DescribeRegions(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe regions: %w", err)
	}

	var regions []string
	for _, region := range result.Regions {
		regions = append(regions, aws.ToString(region.RegionName))
	}

	return regions, nil
}

func (p *AWSProvider) GetResourceTypes(ctx context.Context) ([]string, error) {
	return []string{
		"ec2",
		"instances",
		"volumes",
		"ebs",
		"security-groups",
		"sg",
		"vpcs",
		"vpc",
	}, nil
}

// Helper methods for additional resource details
func (p *AWSProvider) getEBSVolumeDetails(ctx context.Context, volumeID string) (*Resource, error) {
	input := &ec2.DescribeVolumesInput{
		VolumeIds: []string{volumeID},
	}
	result, err := p.ec2Client.DescribeVolumes(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe volume %s: %w", volumeID, err)
	}

	if len(result.Volumes) == 0 {
		return nil, fmt.Errorf("volume %s not found", volumeID)
	}

	volume := result.Volumes[0]
	resource := &Resource{
		ID:       aws.ToString(volume.VolumeId),
		Name:     p.getVolumeName(volume),
		Type:     "ebs-volume",
		Provider: "aws",
		Region:   aws.ToString(volume.AvailabilityZone),
		State:    string(volume.State),
		Status:   string(volume.State),
		Created:  aws.ToTime(volume.CreateTime),
		Modified: time.Now(),
		Tags:     p.convertEBSVolumeTags(volume.Tags),
		Config: map[string]interface{}{
			"volume_type": string(volume.VolumeType),
			"size":        aws.ToInt32(volume.Size),
			"iops":        aws.ToInt32(volume.Iops),
			"throughput":  aws.ToInt32(volume.Throughput),
			"encrypted":   aws.ToBool(volume.Encrypted),
			"snapshot_id": aws.ToString(volume.SnapshotId),
			"kms_key_id":  aws.ToString(volume.KmsKeyId),
		},
	}

	return resource, nil
}

func (p *AWSProvider) getSecurityGroupDetails(ctx context.Context, groupID string) (*Resource, error) {
	input := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{groupID},
	}
	result, err := p.ec2Client.DescribeSecurityGroups(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe security group %s: %w", groupID, err)
	}

	if len(result.SecurityGroups) == 0 {
		return nil, fmt.Errorf("security group %s not found", groupID)
	}

	sg := result.SecurityGroups[0]
	resource := &Resource{
		ID:       aws.ToString(sg.GroupId),
		Name:     aws.ToString(sg.GroupName),
		Type:     "security-group",
		Provider: "aws",
		Region:   "",
		State:    "available",
		Status:   "available",
		Created:  time.Now(),
		Modified: time.Now(),
		Tags:     p.convertSecurityGroupTags(sg.Tags),
		Config: map[string]interface{}{
			"description":    aws.ToString(sg.Description),
			"vpc_id":         aws.ToString(sg.VpcId),
			"owner_id":       aws.ToString(sg.OwnerId),
			"inbound_rules":  len(sg.IpPermissions),
			"outbound_rules": len(sg.IpPermissionsEgress),
		},
	}

	return resource, nil
}

func (p *AWSProvider) getVPCDetails(ctx context.Context, vpcID string) (*Resource, error) {
	input := &ec2.DescribeVpcsInput{
		VpcIds: []string{vpcID},
	}
	result, err := p.ec2Client.DescribeVpcs(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe VPC %s: %w", vpcID, err)
	}

	if len(result.Vpcs) == 0 {
		return nil, fmt.Errorf("VPC %s not found", vpcID)
	}

	vpc := result.Vpcs[0]
	resource := &Resource{
		ID:       aws.ToString(vpc.VpcId),
		Name:     p.getVPCName(vpc),
		Type:     "vpc",
		Provider: "aws",
		Region:   "",
		State:    string(vpc.State),
		Status:   string(vpc.State),
		Created:  time.Now(),
		Modified: time.Now(),
		Tags:     p.convertVPCTags(vpc.Tags),
		Config: map[string]interface{}{
			"cidr_block":           aws.ToString(vpc.CidrBlock),
			"dhcp_options_id":      aws.ToString(vpc.DhcpOptionsId),
			"instance_tenancy":     string(vpc.InstanceTenancy),
			"is_default":           aws.ToBool(vpc.IsDefault),
			"ipv6_cidr_block_sets": len(vpc.Ipv6CidrBlockAssociationSet),
			"owner_id":             aws.ToString(vpc.OwnerId),
		},
	}

	return resource, nil
}
