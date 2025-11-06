package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// pluginManager is the default implementation of PluginManager
type pluginManager struct {
	plugins      map[string]Plugin
	pluginInfo   map[string]*PluginInfo
	capabilities map[PluginCapability][]string
	dataDir      string
	cacheDir     string
	mu           sync.RWMutex
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(dataDir, cacheDir string) PluginManager {
	if dataDir == "" {
		dataDir = filepath.Join(os.TempDir(), "akashic_scribe_plugins_data")
	}
	if cacheDir == "" {
		cacheDir = filepath.Join(os.TempDir(), "akashic_scribe_plugins_cache")
	}

	// Ensure directories exist
	os.MkdirAll(dataDir, 0755)
	os.MkdirAll(cacheDir, 0755)

	return &pluginManager{
		plugins:      make(map[string]Plugin),
		pluginInfo:   make(map[string]*PluginInfo),
		capabilities: make(map[PluginCapability][]string),
		dataDir:      dataDir,
		cacheDir:     cacheDir,
	}
}

// RegisterPlugin registers a plugin with the manager
func (pm *pluginManager) RegisterPlugin(plugin Plugin) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	id := plugin.ID()
	if id == "" {
		return fmt.Errorf("plugin ID cannot be empty")
	}

	if _, exists := pm.plugins[id]; exists {
		return fmt.Errorf("plugin %s is already registered", id)
	}

	// Create plugin info
	info := &PluginInfo{
		ID:           plugin.ID(),
		Name:         plugin.Name(),
		Version:      plugin.Version(),
		Description:  plugin.Description(),
		Author:       plugin.Author(),
		Capabilities: plugin.GetCapabilities(),
		Dependencies: plugin.GetDependencies(),
		Loaded:       false,
		Enabled:      false,
	}

	pm.pluginInfo[id] = info
	return nil
}

// UnregisterPlugin removes a plugin from the manager
func (pm *pluginManager) UnregisterPlugin(pluginID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.plugins[pluginID]; !exists {
		return fmt.Errorf("plugin %s not found", pluginID)
	}

	// Unload if loaded
	if pm.pluginInfo[pluginID].Loaded {
		plugin := pm.plugins[pluginID]
		if err := plugin.Shutdown(); err != nil {
			return fmt.Errorf("failed to shutdown plugin: %w", err)
		}
	}

	delete(pm.plugins, pluginID)
	delete(pm.pluginInfo, pluginID)

	// Remove from capability mappings
	for cap, ids := range pm.capabilities {
		newIDs := []string{}
		for _, id := range ids {
			if id != pluginID {
				newIDs = append(newIDs, id)
			}
		}
		pm.capabilities[cap] = newIDs
	}

	return nil
}

// LoadPlugin loads and initializes a plugin
func (pm *pluginManager) LoadPlugin(plugin Plugin) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	id := plugin.ID()

	// Check if already loaded
	if info, exists := pm.pluginInfo[id]; exists && info.Loaded {
		return fmt.Errorf("plugin %s is already loaded", id)
	}

	// Create plugin context
	context := newPluginContext(id, pm.dataDir, pm.cacheDir)

	// Initialize plugin
	if err := plugin.Initialize(context); err != nil {
		return fmt.Errorf("failed to initialize plugin %s: %w", id, err)
	}

	// Store plugin
	pm.plugins[id] = plugin

	// Update info
	if info, exists := pm.pluginInfo[id]; exists {
		info.Loaded = true
		info.LoadedAt = time.Now()
	} else {
		pm.pluginInfo[id] = &PluginInfo{
			ID:           plugin.ID(),
			Name:         plugin.Name(),
			Version:      plugin.Version(),
			Description:  plugin.Description(),
			Author:       plugin.Author(),
			Capabilities: plugin.GetCapabilities(),
			Dependencies: plugin.GetDependencies(),
			Loaded:       true,
			Enabled:      false,
			LoadedAt:     time.Now(),
		}
	}

	// Register capabilities
	for _, cap := range plugin.GetCapabilities() {
		pm.capabilities[cap] = append(pm.capabilities[cap], id)
	}

	return nil
}

// UnloadPlugin unloads a plugin
func (pm *pluginManager) UnloadPlugin(pluginID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin, exists := pm.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin %s not loaded", pluginID)
	}

	// Deactivate if enabled
	if info := pm.pluginInfo[pluginID]; info.Enabled {
		if err := plugin.Deactivate(); err != nil {
			return fmt.Errorf("failed to deactivate plugin: %w", err)
		}
		info.Enabled = false
	}

	// Shutdown plugin
	if err := plugin.Shutdown(); err != nil {
		return fmt.Errorf("failed to shutdown plugin: %w", err)
	}

	// Update info
	if info := pm.pluginInfo[pluginID]; info != nil {
		info.Loaded = false
	}

	// Remove from plugins map
	delete(pm.plugins, pluginID)

	// Remove from capability mappings
	for cap, ids := range pm.capabilities {
		newIDs := []string{}
		for _, id := range ids {
			if id != pluginID {
				newIDs = append(newIDs, id)
			}
		}
		pm.capabilities[cap] = newIDs
	}

	return nil
}

// EnablePlugin enables a loaded plugin
func (pm *pluginManager) EnablePlugin(pluginID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin, exists := pm.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin %s not loaded", pluginID)
	}

	info := pm.pluginInfo[pluginID]
	if info.Enabled {
		return nil // Already enabled
	}

	if err := plugin.Activate(); err != nil {
		return fmt.Errorf("failed to activate plugin: %w", err)
	}

	info.Enabled = true
	return nil
}

// DisablePlugin disables an enabled plugin
func (pm *pluginManager) DisablePlugin(pluginID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin, exists := pm.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin %s not loaded", pluginID)
	}

	info := pm.pluginInfo[pluginID]
	if !info.Enabled {
		return nil // Already disabled
	}

	if err := plugin.Deactivate(); err != nil {
		return fmt.Errorf("failed to deactivate plugin: %w", err)
	}

	info.Enabled = false
	return nil
}

// GetLoadedPlugins returns all loaded plugins
func (pm *pluginManager) GetLoadedPlugins() []Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugins := make([]Plugin, 0, len(pm.plugins))
	for _, plugin := range pm.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// GetAvailablePlugins returns info about all available plugins
func (pm *pluginManager) GetAvailablePlugins() []PluginInfo {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	infos := make([]PluginInfo, 0, len(pm.pluginInfo))
	for _, info := range pm.pluginInfo {
		infos = append(infos, *info)
	}
	return infos
}

// GetPluginByID retrieves a plugin by its ID
func (pm *pluginManager) GetPluginByID(pluginID string) (Plugin, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugin, exists := pm.plugins[pluginID]
	return plugin, exists
}

// GetPluginInfo retrieves plugin info by ID
func (pm *pluginManager) GetPluginInfo(pluginID string) (*PluginInfo, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	info, exists := pm.pluginInfo[pluginID]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", pluginID)
	}
	return info, nil
}

// SendMessage sends a message to a specific plugin
func (pm *pluginManager) SendMessage(senderID, targetID string, message interface{}) error {
	pm.mu.RLock()
	plugin, exists := pm.plugins[targetID]
	pm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("target plugin %s not found", targetID)
	}

	// For now, we don't have a message handling interface
	// This would be implemented when plugins need inter-plugin communication
	_ = plugin
	_ = senderID
	_ = message

	return nil
}

// BroadcastMessage sends a message to all loaded plugins
func (pm *pluginManager) BroadcastMessage(senderID string, message interface{}) error {
	pm.mu.RLock()
	plugins := make([]Plugin, 0, len(pm.plugins))
	for id, plugin := range pm.plugins {
		if id != senderID {
			plugins = append(plugins, plugin)
		}
	}
	pm.mu.RUnlock()

	// For now, we don't have a message handling interface
	_ = plugins
	_ = message

	return nil
}

// GetPluginsByCapability returns all plugins with a specific capability
func (pm *pluginManager) GetPluginsByCapability(capability PluginCapability) []Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	pluginIDs, exists := pm.capabilities[capability]
	if !exists {
		return []Plugin{}
	}

	plugins := make([]Plugin, 0, len(pluginIDs))
	for _, id := range pluginIDs {
		if plugin, ok := pm.plugins[id]; ok {
			plugins = append(plugins, plugin)
		}
	}
	return plugins
}

// GetCapabilities returns all available capabilities
func (pm *pluginManager) GetCapabilities() []PluginCapability {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	capabilities := make([]PluginCapability, 0, len(pm.capabilities))
	for cap := range pm.capabilities {
		capabilities = append(capabilities, cap)
	}
	return capabilities
}

// HealthCheckAll performs health checks on all loaded plugins
func (pm *pluginManager) HealthCheckAll() map[string]error {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	results := make(map[string]error)
	for id, plugin := range pm.plugins {
		results[id] = plugin.HealthCheck()
	}
	return results
}

// HealthCheck performs a health check on a specific plugin
func (pm *pluginManager) HealthCheck(pluginID string) error {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugin, exists := pm.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin %s not found", pluginID)
	}

	return plugin.HealthCheck()
}

// pluginContext implements PluginContext
type pluginContext struct {
	pluginID string
	config   map[string]interface{}
	dataDir  string
	cacheDir string
	mu       sync.RWMutex
}

func newPluginContext(pluginID, dataDir, cacheDir string) *pluginContext {
	pluginDataDir := filepath.Join(dataDir, pluginID)
	pluginCacheDir := filepath.Join(cacheDir, pluginID)

	os.MkdirAll(pluginDataDir, 0755)
	os.MkdirAll(pluginCacheDir, 0755)

	return &pluginContext{
		pluginID: pluginID,
		config:   make(map[string]interface{}),
		dataDir:  pluginDataDir,
		cacheDir: pluginCacheDir,
	}
}

func (pc *pluginContext) GetConfig() map[string]interface{} {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	// Return a copy
	config := make(map[string]interface{})
	for k, v := range pc.config {
		config[k] = v
	}
	return config
}

func (pc *pluginContext) SetConfig(key string, value interface{}) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.config[key] = value
	return nil
}

func (pc *pluginContext) LogInfo(message string) {
	fmt.Printf("[INFO] [%s] %s\n", pc.pluginID, message)
}

func (pc *pluginContext) LogWarning(message string) {
	fmt.Printf("[WARN] [%s] %s\n", pc.pluginID, message)
}

func (pc *pluginContext) LogError(message string) {
	fmt.Printf("[ERROR] [%s] %s\n", pc.pluginID, message)
}

func (pc *pluginContext) SendMessage(targetPluginID string, message interface{}) error {
	// Would be implemented with actual plugin manager reference
	return fmt.Errorf("not implemented")
}

func (pc *pluginContext) BroadcastMessage(message interface{}) error {
	// Would be implemented with actual plugin manager reference
	return fmt.Errorf("not implemented")
}

func (pc *pluginContext) GetPluginDataDir() string {
	return pc.dataDir
}

func (pc *pluginContext) GetPluginCacheDir() string {
	return pc.cacheDir
}
