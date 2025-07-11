package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// Cache interface defines caching operations
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Clear(ctx context.Context) error
	GetWithTTL(ctx context.Context, key string) ([]byte, time.Duration, error)
	SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetJSON(ctx context.Context, key string, dest interface{}) error
}

// MemoryCache implements an in-memory cache
type MemoryCache struct {
	data   map[string]*CacheItem
	mutex  sync.RWMutex
	maxTTL time.Duration
}

// CacheItem represents a cached item
type CacheItem struct {
	Value      []byte
	Expiration time.Time
	CreatedAt  time.Time
}

// RedisCache implements Redis-based cache
type RedisCache struct {
	client *redis.Client
	prefix string
}

// NewMemoryCache creates a new in-memory cache
func NewMemoryCache(maxTTL time.Duration) *MemoryCache {
	cache := &MemoryCache{
		data:   make(map[string]*CacheItem),
		maxTTL: maxTTL,
	}
	
	// Start cleanup goroutine
	go cache.cleanup()
	
	return cache
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(addr, password string, db int, prefix string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	
	return &RedisCache{
		client: client,
		prefix: prefix,
	}
}

// Memory Cache Implementation

// Get retrieves a value from memory cache
func (c *MemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	item, exists := c.data[key]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	
	// Check expiration
	if time.Now().After(item.Expiration) {
		c.mutex.RUnlock()
		c.mutex.Lock()
		delete(c.data, key)
		c.mutex.Unlock()
		c.mutex.RLock()
		return nil, fmt.Errorf("key expired: %s", key)
	}
	
	return item.Value, nil
}

// Set stores a value in memory cache
func (c *MemoryCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	// Use maxTTL if expiration is longer or zero
	if expiration == 0 || expiration > c.maxTTL {
		expiration = c.maxTTL
	}
	
	c.data[key] = &CacheItem{
		Value:      value,
		Expiration: time.Now().Add(expiration),
		CreatedAt:  time.Now(),
	}
	
	return nil
}

// Delete removes a value from memory cache
func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	delete(c.data, key)
	return nil
}

// Exists checks if a key exists in memory cache
func (c *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	item, exists := c.data[key]
	if !exists {
		return false, nil
	}
	
	// Check expiration
	if time.Now().After(item.Expiration) {
		return false, nil
	}
	
	return true, nil
}

// Clear removes all values from memory cache
func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.data = make(map[string]*CacheItem)
	return nil
}

// GetWithTTL retrieves a value with remaining TTL
func (c *MemoryCache) GetWithTTL(ctx context.Context, key string) ([]byte, time.Duration, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	item, exists := c.data[key]
	if !exists {
		return nil, 0, fmt.Errorf("key not found: %s", key)
	}
	
	ttl := time.Until(item.Expiration)
	if ttl <= 0 {
		return nil, 0, fmt.Errorf("key expired: %s", key)
	}
	
	return item.Value, ttl, nil
}

// SetJSON stores a JSON-encoded value
func (c *MemoryCache) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	
	return c.Set(ctx, key, data, expiration)
}

// GetJSON retrieves and decodes a JSON value
func (c *MemoryCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := c.Get(ctx, key)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, dest)
}

// cleanup removes expired items periodically
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, item := range c.data {
			if now.After(item.Expiration) {
				delete(c.data, key)
			}
		}
		c.mutex.Unlock()
	}
}

// Redis Cache Implementation

// Get retrieves a value from Redis cache
func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	key = c.prefix + key
	result, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("key not found: %s", key)
		}
		return nil, fmt.Errorf("failed to get key: %w", err)
	}
	
	return result, nil
}

// Set stores a value in Redis cache
func (c *RedisCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	key = c.prefix + key
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Delete removes a value from Redis cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	key = c.prefix + key
	return c.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in Redis cache
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	key = c.prefix + key
	count, err := c.client.Exists(ctx, key).Result()
	return count > 0, err
}

// Clear removes all values from Redis cache with the prefix
func (c *RedisCache) Clear(ctx context.Context) error {
	pattern := c.prefix + "*"
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	
	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}
	
	return nil
}

// GetWithTTL retrieves a value with remaining TTL from Redis
func (c *RedisCache) GetWithTTL(ctx context.Context, key string) ([]byte, time.Duration, error) {
	key = c.prefix + key
	
	pipe := c.client.Pipeline()
	getCmd := pipe.Get(ctx, key)
	ttlCmd := pipe.TTL(ctx, key)
	
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, 0, fmt.Errorf("failed to execute pipeline: %w", err)
	}
	
	if err == redis.Nil {
		return nil, 0, fmt.Errorf("key not found: %s", key)
	}
	
	value, err := getCmd.Bytes()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get value: %w", err)
	}
	
	ttl, err := ttlCmd.Result()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get TTL: %w", err)
	}
	
	return value, ttl, nil
}

// SetJSON stores a JSON-encoded value in Redis
func (c *RedisCache) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	
	return c.Set(ctx, key, data, expiration)
}

// GetJSON retrieves and decodes a JSON value from Redis
func (c *RedisCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := c.Get(ctx, key)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, dest)
}

// CacheManager manages multiple cache instances
type CacheManager struct {
	caches map[string]Cache
	mutex  sync.RWMutex
}

// NewCacheManager creates a new cache manager
func NewCacheManager() *CacheManager {
	return &CacheManager{
		caches: make(map[string]Cache),
	}
}

// AddCache adds a cache instance
func (m *CacheManager) AddCache(name string, cache Cache) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.caches[name] = cache
}

// GetCache retrieves a cache instance
func (m *CacheManager) GetCache(name string) (Cache, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	cache, exists := m.caches[name]
	if !exists {
		return nil, fmt.Errorf("cache not found: %s", name)
	}
	
	return cache, nil
}

// GetDefaultCache returns the default cache instance
func (m *CacheManager) GetDefaultCache() Cache {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	// Return the first cache or create a memory cache
	for _, cache := range m.caches {
		return cache
	}
	
	// Create default memory cache
	return NewMemoryCache(24 * time.Hour)
}
