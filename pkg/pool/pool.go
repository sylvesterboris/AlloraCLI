package pool

import (
	"context"
	"fmt"
	"sync"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/go-redis/redis/v8"
)

// ConnectionPool manages connection pools for different services
type ConnectionPool struct {
	awsEC2Pool        *AWS_EC2Pool
	azureComputePool  *AzureComputePool
	azureResourcePool *AzureResourcePool
	gcpComputePool    *GCPComputePool
	redisPool         *RedisPool
	mutex             sync.RWMutex
}

// AWS EC2 Pool
type AWS_EC2Pool struct {
	clients []*ec2.Client
	current int
	mutex   sync.RWMutex
	config  aws.Config
}

// Azure Compute Pool
type AzureComputePool struct {
	clients []*armcompute.VirtualMachinesClient
	current int
	mutex   sync.RWMutex
	cred    azcore.TokenCredential
	subID   string
}

// Azure Resource Pool
type AzureResourcePool struct {
	clients []*armresources.Client
	current int
	mutex   sync.RWMutex
	cred    azcore.TokenCredential
	subID   string
}

// GCP Compute Pool
type GCPComputePool struct {
	clients []*compute.InstancesClient
	current int
	mutex   sync.RWMutex
	ctx     context.Context
}

// Redis Pool
type RedisPool struct {
	clients []*redis.Client
	current int
	mutex   sync.RWMutex
	options *redis.Options
}

// NewConnectionPool creates a new connection pool
func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{}
}

// InitializeAWSPool initializes AWS EC2 connection pool
func (p *ConnectionPool) InitializeAWSPool(config aws.Config, poolSize int) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	pool := &AWS_EC2Pool{
		clients: make([]*ec2.Client, poolSize),
		config:  config,
	}

	for i := 0; i < poolSize; i++ {
		client := ec2.NewFromConfig(config)
		pool.clients[i] = client
	}

	p.awsEC2Pool = pool
	return nil
}

// InitializeAzureComputePool initializes Azure Compute connection pool
func (p *ConnectionPool) InitializeAzureComputePool(cred azcore.TokenCredential, subscriptionID string, poolSize int) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	pool := &AzureComputePool{
		clients: make([]*armcompute.VirtualMachinesClient, poolSize),
		cred:    cred,
		subID:   subscriptionID,
	}

	for i := 0; i < poolSize; i++ {
		client, err := armcompute.NewVirtualMachinesClient(subscriptionID, cred, nil)
		if err != nil {
			return fmt.Errorf("failed to create Azure compute client: %w", err)
		}
		pool.clients[i] = client
	}

	p.azureComputePool = pool
	return nil
}

// InitializeAzureResourcePool initializes Azure Resource connection pool
func (p *ConnectionPool) InitializeAzureResourcePool(cred azcore.TokenCredential, subscriptionID string, poolSize int) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	pool := &AzureResourcePool{
		clients: make([]*armresources.Client, poolSize),
		cred:    cred,
		subID:   subscriptionID,
	}

	for i := 0; i < poolSize; i++ {
		client, err := armresources.NewClient(subscriptionID, cred, nil)
		if err != nil {
			return fmt.Errorf("failed to create Azure resource client: %w", err)
		}
		pool.clients[i] = client
	}

	p.azureResourcePool = pool
	return nil
}

// InitializeGCPPool initializes GCP Compute connection pool
func (p *ConnectionPool) InitializeGCPPool(ctx context.Context, poolSize int) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	pool := &GCPComputePool{
		clients: make([]*compute.InstancesClient, poolSize),
		ctx:     ctx,
	}

	for i := 0; i < poolSize; i++ {
		client, err := compute.NewInstancesRESTClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to create GCP compute client: %w", err)
		}
		pool.clients[i] = client
	}

	p.gcpComputePool = pool
	return nil
}

// InitializeRedisPool initializes Redis connection pool
func (p *ConnectionPool) InitializeRedisPool(options *redis.Options, poolSize int) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	pool := &RedisPool{
		clients: make([]*redis.Client, poolSize),
		options: options,
	}

	for i := 0; i < poolSize; i++ {
		client := redis.NewClient(options)
		pool.clients[i] = client
	}

	p.redisPool = pool
	return nil
}

// GetAWSClient returns an AWS EC2 client from the pool
func (p *ConnectionPool) GetAWSClient() (*ec2.Client, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.awsEC2Pool == nil {
		return nil, fmt.Errorf("AWS pool not initialized")
	}

	p.awsEC2Pool.mutex.Lock()
	defer p.awsEC2Pool.mutex.Unlock()

	client := p.awsEC2Pool.clients[p.awsEC2Pool.current]
	p.awsEC2Pool.current = (p.awsEC2Pool.current + 1) % len(p.awsEC2Pool.clients)

	return client, nil
}

// GetAzureComputeClient returns an Azure Compute client from the pool
func (p *ConnectionPool) GetAzureComputeClient() (*armcompute.VirtualMachinesClient, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.azureComputePool == nil {
		return nil, fmt.Errorf("Azure Compute pool not initialized")
	}

	p.azureComputePool.mutex.Lock()
	defer p.azureComputePool.mutex.Unlock()

	client := p.azureComputePool.clients[p.azureComputePool.current]
	p.azureComputePool.current = (p.azureComputePool.current + 1) % len(p.azureComputePool.clients)

	return client, nil
}

// GetAzureResourceClient returns an Azure Resource client from the pool
func (p *ConnectionPool) GetAzureResourceClient() (*armresources.Client, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.azureResourcePool == nil {
		return nil, fmt.Errorf("Azure Resource pool not initialized")
	}

	p.azureResourcePool.mutex.Lock()
	defer p.azureResourcePool.mutex.Unlock()

	client := p.azureResourcePool.clients[p.azureResourcePool.current]
	p.azureResourcePool.current = (p.azureResourcePool.current + 1) % len(p.azureResourcePool.clients)

	return client, nil
}

// GetGCPClient returns a GCP Compute client from the pool
func (p *ConnectionPool) GetGCPClient() (*compute.InstancesClient, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.gcpComputePool == nil {
		return nil, fmt.Errorf("GCP pool not initialized")
	}

	p.gcpComputePool.mutex.Lock()
	defer p.gcpComputePool.mutex.Unlock()

	client := p.gcpComputePool.clients[p.gcpComputePool.current]
	p.gcpComputePool.current = (p.gcpComputePool.current + 1) % len(p.gcpComputePool.clients)

	return client, nil
}

// GetRedisClient returns a Redis client from the pool
func (p *ConnectionPool) GetRedisClient() (*redis.Client, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.redisPool == nil {
		return nil, fmt.Errorf("Redis pool not initialized")
	}

	p.redisPool.mutex.Lock()
	defer p.redisPool.mutex.Unlock()

	client := p.redisPool.clients[p.redisPool.current]
	p.redisPool.current = (p.redisPool.current + 1) % len(p.redisPool.clients)

	return client, nil
}

// HealthCheck performs health check on all pools
func (p *ConnectionPool) HealthCheck(ctx context.Context) error {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// Check AWS pool
	if p.awsEC2Pool != nil {
		client, err := p.GetAWSClient()
		if err != nil {
			return fmt.Errorf("AWS pool health check failed: %w", err)
		}

		// Simple health check - describe regions
		_, err = client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
		if err != nil {
			return fmt.Errorf("AWS pool health check failed: %w", err)
		}
	}

	// Check Redis pool
	if p.redisPool != nil {
		client, err := p.GetRedisClient()
		if err != nil {
			return fmt.Errorf("Redis pool health check failed: %w", err)
		}

		// Simple health check - ping
		err = client.Ping(ctx).Err()
		if err != nil {
			return fmt.Errorf("Redis pool health check failed: %w", err)
		}
	}

	return nil
}

// GetPoolStats returns statistics about the connection pools
func (p *ConnectionPool) GetPoolStats() map[string]interface{} {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	stats := make(map[string]interface{})

	if p.awsEC2Pool != nil {
		stats["aws_ec2_pool_size"] = len(p.awsEC2Pool.clients)
		stats["aws_ec2_current_index"] = p.awsEC2Pool.current
	}

	if p.azureComputePool != nil {
		stats["azure_compute_pool_size"] = len(p.azureComputePool.clients)
		stats["azure_compute_current_index"] = p.azureComputePool.current
	}

	if p.azureResourcePool != nil {
		stats["azure_resource_pool_size"] = len(p.azureResourcePool.clients)
		stats["azure_resource_current_index"] = p.azureResourcePool.current
	}

	if p.gcpComputePool != nil {
		stats["gcp_compute_pool_size"] = len(p.gcpComputePool.clients)
		stats["gcp_compute_current_index"] = p.gcpComputePool.current
	}

	if p.redisPool != nil {
		stats["redis_pool_size"] = len(p.redisPool.clients)
		stats["redis_current_index"] = p.redisPool.current
	}

	return stats
}

// Close closes all connections in the pools
func (p *ConnectionPool) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Close GCP connections
	if p.gcpComputePool != nil {
		for _, client := range p.gcpComputePool.clients {
			client.Close()
		}
	}

	// Close Redis connections
	if p.redisPool != nil {
		for _, client := range p.redisPool.clients {
			client.Close()
		}
	}

	return nil
}

// Global connection pool instance
var GlobalConnectionPool = NewConnectionPool()

// InitializeGlobalPools initializes all global connection pools
func InitializeGlobalPools(ctx context.Context) error {
	// Initialize with default pool sizes
	poolSize := 5

	// Initialize Redis pool (if Redis is available)
	redisOptions := &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}

	if err := GlobalConnectionPool.InitializeRedisPool(redisOptions, poolSize); err != nil {
		// Redis is optional, log but don't fail
		fmt.Printf("Warning: Failed to initialize Redis pool: %v\n", err)
	}

	return nil
}

// GetGlobalConnectionPool returns the global connection pool instance
func GetGlobalConnectionPool() *ConnectionPool {
	return GlobalConnectionPool
}
