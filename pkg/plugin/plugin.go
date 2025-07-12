package plugin

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Plugin interface defines the contract for AlloraCLI plugins
type Plugin interface {
	Name() string
	Version() string
	Description() string
	Initialize(config map[string]interface{}) error
	Execute(ctx context.Context, args []string) (interface{}, error)
	Cleanup() error
}

// PluginManager manages plugin lifecycle
type PluginManager struct {
	plugins     map[string]Plugin
	pluginPaths map[string]string
	config      *PluginConfig
	logger      *logrus.Logger
	mu          sync.RWMutex
}

// PluginConfig represents plugin configuration
type PluginConfig struct {
	Directory      string   `json:"directory" yaml:"directory"`
	AutoUpdate     bool     `json:"auto_update" yaml:"auto_update"`
	AllowedSources []string `json:"allowed_sources" yaml:"allowed_sources"`
	Timeout        int      `json:"timeout" yaml:"timeout"`
	MaxPlugins     int      `json:"max_plugins" yaml:"max_plugins"`
}

// PluginInfo represents plugin information
type PluginInfo struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"`
	Path        string                 `json:"path"`
	Config      map[string]interface{} `json:"config"`
	LoadedAt    time.Time              `json:"loaded_at"`
	LastUsed    time.Time              `json:"last_used"`
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(config *PluginConfig) *PluginManager {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	return &PluginManager{
		plugins:     make(map[string]Plugin),
		pluginPaths: make(map[string]string),
		config:      config,
		logger:      logger,
	}
}

// LoadPlugin loads a plugin from file
func (pm *PluginManager) LoadPlugin(pluginPath string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.logger.Infof("Loading plugin from: %s", pluginPath)

	// Check if plugin already loaded
	if _, exists := pm.pluginPaths[pluginPath]; exists {
		return fmt.Errorf("plugin already loaded: %s", pluginPath)
	}

	// Check max plugins limit
	if pm.config.MaxPlugins > 0 && len(pm.plugins) >= pm.config.MaxPlugins {
		return fmt.Errorf("maximum number of plugins reached: %d", pm.config.MaxPlugins)
	}

	// Load the plugin
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin: %w", err)
	}

	// Look for the plugin symbol
	symbol, err := p.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("plugin symbol not found: %w", err)
	}

	// Assert that the symbol implements our Plugin interface
	pluginInstance, ok := symbol.(Plugin)
	if !ok {
		return fmt.Errorf("plugin does not implement Plugin interface")
	}

	// Initialize the plugin
	if err := pluginInstance.Initialize(nil); err != nil {
		return fmt.Errorf("failed to initialize plugin: %w", err)
	}

	// Store the plugin
	name := pluginInstance.Name()
	pm.plugins[name] = pluginInstance
	pm.pluginPaths[pluginPath] = name

	pm.logger.Infof("Successfully loaded plugin: %s v%s", name, pluginInstance.Version())
	return nil
}

// LoadPluginsFromDirectory loads all plugins from a directory
func (pm *PluginManager) LoadPluginsFromDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		pm.logger.Warnf("Plugin directory does not exist: %s", dir)
		return nil
	}

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-shared objects
		if info.IsDir() || !strings.HasSuffix(path, ".so") {
			return nil
		}

		// Load the plugin
		if err := pm.LoadPlugin(path); err != nil {
			pm.logger.Warnf("Failed to load plugin %s: %v", path, err)
			// Continue loading other plugins
		}

		return nil
	})
}

// UnloadPlugin unloads a plugin
func (pm *PluginManager) UnloadPlugin(name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin not found: %s", name)
	}

	// Cleanup the plugin
	if err := plugin.Cleanup(); err != nil {
		pm.logger.Warnf("Plugin cleanup failed for %s: %v", name, err)
	}

	// Remove from maps
	delete(pm.plugins, name)
	for path, pluginName := range pm.pluginPaths {
		if pluginName == name {
			delete(pm.pluginPaths, path)
			break
		}
	}

	pm.logger.Infof("Unloaded plugin: %s", name)
	return nil
}

// GetPlugin retrieves a plugin by name
func (pm *PluginManager) GetPlugin(name string) (Plugin, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin not found: %s", name)
	}

	return plugin, nil
}

// ListPlugins returns a list of all loaded plugins
func (pm *PluginManager) ListPlugins() []*PluginInfo {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var plugins []*PluginInfo
	for name, plugin := range pm.plugins {
		info := &PluginInfo{
			Name:        name,
			Version:     plugin.Version(),
			Description: plugin.Description(),
			Status:      "loaded",
			LoadedAt:    time.Now(), // This should be tracked properly
		}
		plugins = append(plugins, info)
	}

	return plugins
}

// ExecutePlugin executes a plugin with given arguments
func (pm *PluginManager) ExecutePlugin(ctx context.Context, name string, args []string) (interface{}, error) {
	plugin, err := pm.GetPlugin(name)
	if err != nil {
		return nil, err
	}

	// Add timeout context
	if pm.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(pm.config.Timeout)*time.Second)
		defer cancel()
	}

	pm.logger.Infof("Executing plugin: %s with args: %v", name, args)

	result, err := plugin.Execute(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("plugin execution failed: %w", err)
	}

	return result, nil
}

// RefreshPlugins reloads all plugins from the configured directory
func (pm *PluginManager) RefreshPlugins() error {
	pm.logger.Info("Refreshing plugins...")

	// Unload all current plugins
	pm.mu.Lock()
	for name := range pm.plugins {
		if err := pm.UnloadPlugin(name); err != nil {
			pm.logger.Warnf("Failed to unload plugin %s: %v", name, err)
		}
	}
	pm.mu.Unlock()

	// Load plugins from directory
	if pm.config.Directory != "" {
		return pm.LoadPluginsFromDirectory(pm.config.Directory)
	}

	return nil
}

// ValidatePlugin validates a plugin before loading
func (pm *PluginManager) ValidatePlugin(pluginPath string) error {
	// Check file exists
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return fmt.Errorf("plugin file does not exist: %s", pluginPath)
	}

	// Check file extension
	if !strings.HasSuffix(pluginPath, ".so") {
		return fmt.Errorf("invalid plugin file extension: %s", pluginPath)
	}

	// Try to open plugin to validate it
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin for validation: %w", err)
	}

	// Look for the plugin symbol
	symbol, err := p.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("plugin symbol not found: %w", err)
	}

	// Assert that the symbol implements our Plugin interface
	if _, ok := symbol.(Plugin); !ok {
		return fmt.Errorf("plugin does not implement Plugin interface")
	}

	return nil
}

// InstallPlugin installs a plugin from a source
func (pm *PluginManager) InstallPlugin(source, name string) error {
	// Check if source is allowed
	if !pm.isSourceAllowed(source) {
		return fmt.Errorf("plugin source not allowed: %s", source)
	}

	// This would typically download and install the plugin
	// For now, return not implemented
	return fmt.Errorf("plugin installation not implemented")
}

// isSourceAllowed checks if a plugin source is allowed
func (pm *PluginManager) isSourceAllowed(source string) bool {
	if len(pm.config.AllowedSources) == 0 {
		return true // Allow all sources if none specified
	}

	for _, allowedSource := range pm.config.AllowedSources {
		if strings.Contains(source, allowedSource) {
			return true
		}
	}

	return false
}

// GetPluginConfig returns the plugin configuration
func (pm *PluginManager) GetPluginConfig() *PluginConfig {
	return pm.config
}

// UpdatePluginConfig updates the plugin configuration
func (pm *PluginManager) UpdatePluginConfig(config *PluginConfig) error {
	pm.config = config
	return nil
}

// Shutdown shuts down the plugin manager and cleans up all plugins
func (pm *PluginManager) Shutdown() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.logger.Info("Shutting down plugin manager...")

	for name, plugin := range pm.plugins {
		if err := plugin.Cleanup(); err != nil {
			pm.logger.Warnf("Failed to cleanup plugin %s: %v", name, err)
		}
	}

	pm.plugins = make(map[string]Plugin)
	pm.pluginPaths = make(map[string]string)

	pm.logger.Info("Plugin manager shut down successfully")
	return nil
}

// BasePlugin provides a basic implementation that plugins can embed
type BasePlugin struct {
	name        string
	version     string
	description string
	initialized bool
}

// NewBasePlugin creates a new base plugin
func NewBasePlugin(name, version, description string) *BasePlugin {
	return &BasePlugin{
		name:        name,
		version:     version,
		description: description,
	}
}

// Name returns the plugin name
func (bp *BasePlugin) Name() string {
	return bp.name
}

// Version returns the plugin version
func (bp *BasePlugin) Version() string {
	return bp.version
}

// Description returns the plugin description
func (bp *BasePlugin) Description() string {
	return bp.description
}

// Initialize initializes the plugin
func (bp *BasePlugin) Initialize(config map[string]interface{}) error {
	bp.initialized = true
	return nil
}

// Execute executes the plugin (should be overridden by actual plugins)
func (bp *BasePlugin) Execute(ctx context.Context, args []string) (interface{}, error) {
	return nil, fmt.Errorf("Execute method not implemented")
}

// Cleanup cleans up the plugin
func (bp *BasePlugin) Cleanup() error {
	bp.initialized = false
	return nil
}

// IsInitialized returns whether the plugin is initialized
func (bp *BasePlugin) IsInitialized() bool {
	return bp.initialized
}
