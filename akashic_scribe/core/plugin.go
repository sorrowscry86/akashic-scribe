package core

import (
	"errors"
	"time"
)

// PluginCapability defines what a plugin can do
type PluginCapability string

const (
	CapabilityAudioProcessing    PluginCapability = "audio_processing"
	CapabilityVideoProcessing    PluginCapability = "video_processing"
	CapabilityTranscription      PluginCapability = "transcription"
	CapabilityTranslation        PluginCapability = "translation"
	CapabilityVoiceSynthesis     PluginCapability = "voice_synthesis"
	CapabilitySubtitleGeneration PluginCapability = "subtitle_generation"
	CapabilitySubtitleStyling    PluginCapability = "subtitle_styling"
	CapabilityFormatConversion   PluginCapability = "format_conversion"
	CapabilityUIExtension        PluginCapability = "ui_extension"
	CapabilityThemeProvider      PluginCapability = "theme_provider"
	CapabilityFileFormat         PluginCapability = "file_format"
	CapabilityCloudIntegration   PluginCapability = "cloud_integration"
	CapabilityAPIIntegration     PluginCapability = "api_integration"
)

// Plugin defines the interface for Akashic Scribe plugins
type Plugin interface {
	// Metadata
	ID() string
	Name() string
	Version() string
	Description() string
	Author() string

	// Lifecycle
	Initialize(context PluginContext) error
	Activate() error
	Deactivate() error
	Shutdown() error

	// Capabilities
	GetCapabilities() []PluginCapability
	GetDependencies() []PluginDependency

	// Health
	HealthCheck() error
}

// PluginDependency describes plugin dependencies
type PluginDependency struct {
	PluginID    string `json:"plugin_id"`
	Version     string `json:"version"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

// PluginInfo provides metadata about available plugins
type PluginInfo struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Version      string              `json:"version"`
	Description  string              `json:"description"`
	Author       string              `json:"author"`
	License      string              `json:"license,omitempty"`
	Website      string              `json:"website,omitempty"`
	Capabilities []PluginCapability  `json:"capabilities"`
	Dependencies []PluginDependency  `json:"dependencies"`
	FilePath     string              `json:"file_path,omitempty"`
	Loaded       bool                `json:"loaded"`
	Enabled      bool                `json:"enabled"`
	LoadedAt     time.Time           `json:"loaded_at,omitempty"`
}

// PluginContext provides access to application services
type PluginContext interface {
	// Configuration
	GetConfig() map[string]interface{}
	SetConfig(key string, value interface{}) error

	// Logging
	LogInfo(message string)
	LogWarning(message string)
	LogError(message string)

	// Plugin communication
	SendMessage(targetPluginID string, message interface{}) error
	BroadcastMessage(message interface{}) error

	// Resource access
	GetPluginDataDir() string
	GetPluginCacheDir() string
}

// PluginManager handles plugin lifecycle and communication
type PluginManager interface {
	// Plugin registration
	RegisterPlugin(plugin Plugin) error
	UnregisterPlugin(pluginID string) error

	// Plugin lifecycle
	LoadPlugin(plugin Plugin) error
	UnloadPlugin(pluginID string) error
	EnablePlugin(pluginID string) error
	DisablePlugin(pluginID string) error

	// Plugin discovery
	GetLoadedPlugins() []Plugin
	GetAvailablePlugins() []PluginInfo
	GetPluginByID(pluginID string) (Plugin, bool)
	GetPluginInfo(pluginID string) (*PluginInfo, error)

	// Plugin communication
	SendMessage(senderID, targetID string, message interface{}) error
	BroadcastMessage(senderID string, message interface{}) error

	// Plugin capabilities
	GetPluginsByCapability(capability PluginCapability) []Plugin
	GetCapabilities() []PluginCapability

	// Health monitoring
	HealthCheckAll() map[string]error
	HealthCheck(pluginID string) error
}

// AudioProcessor is an interface for audio processing plugins
type AudioProcessor interface {
	Plugin
	ProcessAudio(inputPath, outputPath string, options map[string]interface{}) error
	GetSupportedFormats() []string
}

// SubtitleProcessor is an interface for subtitle processing plugins
type SubtitleProcessor interface {
	Plugin
	ProcessSubtitles(inputPath, outputPath string, options map[string]interface{}) error
	GetSupportedFormats() []string
}

// FormatConverter is an interface for format conversion plugins
type FormatConverter interface {
	Plugin
	ConvertFormat(inputPath, outputPath, targetFormat string, options map[string]interface{}) error
	GetSupportedInputFormats() []string
	GetSupportedOutputFormats() []string
}

// BasePlugin provides common plugin functionality
type BasePlugin struct {
	id          string
	name        string
	version     string
	description string
	author      string
	context     PluginContext
	enabled     bool
}

// NewBasePlugin creates a new base plugin with common fields
func NewBasePlugin(id, name, version, description, author string) *BasePlugin {
	return &BasePlugin{
		id:          id,
		name:        name,
		version:     version,
		description: description,
		author:      author,
		enabled:     false,
	}
}

func (p *BasePlugin) ID() string          { return p.id }
func (p *BasePlugin) Name() string        { return p.name }
func (p *BasePlugin) Version() string     { return p.version }
func (p *BasePlugin) Description() string { return p.description }
func (p *BasePlugin) Author() string      { return p.author }

func (p *BasePlugin) Initialize(context PluginContext) error {
	p.context = context
	p.context.LogInfo("Plugin initialized: " + p.name)
	return nil
}

func (p *BasePlugin) Activate() error {
	p.enabled = true
	p.context.LogInfo("Plugin activated: " + p.name)
	return nil
}

func (p *BasePlugin) Deactivate() error {
	p.enabled = false
	p.context.LogInfo("Plugin deactivated: " + p.name)
	return nil
}

func (p *BasePlugin) Shutdown() error {
	p.enabled = false
	p.context.LogInfo("Plugin shutdown: " + p.name)
	return nil
}

func (p *BasePlugin) GetDependencies() []PluginDependency {
	return []PluginDependency{}
}

func (p *BasePlugin) HealthCheck() error {
	if p.context == nil {
		return errors.New("plugin context not initialized")
	}
	if !p.enabled {
		return errors.New("plugin not enabled")
	}
	return nil
}

func (p *BasePlugin) GetContext() PluginContext {
	return p.context
}

func (p *BasePlugin) IsEnabled() bool {
	return p.enabled
}
