# Plugin Development Guide

**Akashic Scribe Plugin System - VoidCat RDC**

## Table of Contents

1. [Introduction](#introduction)
2. [Plugin Architecture](#plugin-architecture)
3. [Getting Started](#getting-started)
4. [Plugin Interface](#plugin-interface)
5. [Plugin Manager](#plugin-manager)
6. [Specialized Plugin Types](#specialized-plugin-types)
7. [Creating Your First Plugin](#creating-your-first-plugin)
8. [Built-in Example Plugins](#built-in-example-plugins)
9. [Best Practices](#best-practices)
10. [Testing Plugins](#testing-plugins)
11. [Deployment](#deployment)
12. [Troubleshooting](#troubleshooting)

---

## Introduction

The Akashic Scribe plugin system provides a powerful and flexible way to extend the application's functionality. Plugins can add new capabilities for audio processing, subtitle formatting, format conversion, and more.

### Key Features

- **Capability-Based Discovery**: Plugins declare their capabilities for easy discovery
- **Lifecycle Management**: Complete control over plugin initialization, activation, and shutdown
- **Inter-Plugin Communication**: Plugins can send messages to each other
- **Health Monitoring**: Built-in health check system
- **Isolated Context**: Each plugin has its own data and cache directories
- **Type-Safe Interfaces**: Strongly-typed Go interfaces for plugin development

### Plugin Capabilities

Available plugin capabilities:

- `audio_processing` - Process and enhance audio
- `video_processing` - Process video files
- `transcription` - Speech-to-text conversion
- `translation` - Text translation between languages
- `voice_synthesis` - Text-to-speech generation
- `subtitle_generation` - Generate subtitle files
- `subtitle_styling` - Apply visual themes to subtitles
- `format_conversion` - Convert between file formats
- `ui_extension` - Extend the user interface
- `theme_provider` - Provide custom themes
- `file_format` - Support additional file formats
- `cloud_integration` - Integrate with cloud services
- `api_integration` - Integrate with external APIs

---

## Plugin Architecture

### Component Overview

```
┌─────────────────────────────────────────────┐
│           Application Core                   │
│  (ScribeEngine, State Management, GUI)       │
└──────────────┬──────────────────────────────┘
               │
┌──────────────v──────────────────────────────┐
│         Plugin Manager                       │
│  - Plugin Registration                       │
│  - Lifecycle Management                      │
│  - Capability Queries                        │
│  - Health Monitoring                         │
└──────────────┬──────────────────────────────┘
               │
    ┌──────────┴──────────┬──────────────┐
    │                     │              │
┌───v───────────┐  ┌─────v────────┐  ┌──v─────────────┐
│ Plugin A      │  │ Plugin B     │  │ Plugin C       │
│ (Audio)       │  │ (Subtitle)   │  │ (Format)       │
└───────────────┘  └──────────────┘  └────────────────┘
```

### Plugin Lifecycle

1. **Registration**: Plugin is registered with the manager
2. **Loading**: Plugin is loaded and initialized with context
3. **Activation**: Plugin is activated and ready to use
4. **Operation**: Plugin performs its functions
5. **Deactivation**: Plugin is temporarily disabled
6. **Shutdown**: Plugin releases resources and terminates

---

## Getting Started

### Prerequisites

- Go 1.20 or higher
- Akashic Scribe source code
- Basic understanding of Go interfaces

### Project Structure

```
akashic_scribe/
├── core/
│   ├── plugin.go              # Plugin interfaces
│   ├── plugin_manager.go      # Plugin manager implementation
│   └── ...
└── plugins/
    ├── your_plugin/
    │   ├── your_plugin.go     # Plugin implementation
    │   ├── your_plugin_test.go
    │   └── README.md
    └── plugin_demo.go          # Demo/example usage
```

---

## Plugin Interface

### Core Plugin Interface

Every plugin must implement the `Plugin` interface:

```go
type Plugin interface {
    // Metadata
    ID() string          // Unique identifier (e.g., "voidcat.my_plugin")
    Name() string        // Human-readable name
    Version() string     // Semantic version (e.g., "1.0.0")
    Description() string // Brief description
    Author() string      // Plugin author

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
```

### Plugin Context

The `PluginContext` provides access to application services:

```go
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
    GetPluginDataDir() string   // Persistent data storage
    GetPluginCacheDir() string  // Temporary cache storage
}
```

### Base Plugin

Use `BasePlugin` for common functionality:

```go
type MyPlugin struct {
    *core.BasePlugin
    // Your custom fields
}

func NewMyPlugin() core.Plugin {
    base := core.NewBasePlugin(
        "voidcat.my_plugin",      // ID
        "My Awesome Plugin",       // Name
        "1.0.0",                   // Version
        "Does amazing things",     // Description
        "VoidCat RDC",            // Author
    )

    return &MyPlugin{
        BasePlugin: base,
    }
}
```

---

## Plugin Manager

### Creating a Plugin Manager

```go
import "akashic_scribe/core"

// Create manager with custom directories
manager := core.NewPluginManager("/path/to/data", "/path/to/cache")

// Or use defaults (temp directories)
manager := core.NewPluginManager("", "")
```

### Registering Plugins

```go
plugin := NewMyPlugin()

// Register plugin
if err := manager.RegisterPlugin(plugin); err != nil {
    log.Fatal(err)
}
```

### Loading and Activating

```go
// Load plugin (initializes with context)
if err := manager.LoadPlugin(plugin); err != nil {
    log.Fatal(err)
}

// Enable plugin (activates it)
if err := manager.EnablePlugin(plugin.ID()); err != nil {
    log.Fatal(err)
}
```

### Querying Plugins

```go
// Get all loaded plugins
plugins := manager.GetLoadedPlugins()

// Get plugin by ID
plugin, found := manager.GetPluginByID("voidcat.my_plugin")

// Get plugins by capability
audioPlugins := manager.GetPluginsByCapability(core.CapabilityAudioProcessing)

// Get plugin info
info, err := manager.GetPluginInfo("voidcat.my_plugin")
```

### Health Monitoring

```go
// Check all plugins
healthMap := manager.HealthCheckAll()
for id, err := range healthMap {
    if err != nil {
        log.Printf("Plugin %s unhealthy: %v", id, err)
    }
}

// Check specific plugin
if err := manager.HealthCheck("voidcat.my_plugin"); err != nil {
    log.Printf("Health check failed: %v", err)
}
```

---

## Specialized Plugin Types

### Audio Processor

For audio processing plugins:

```go
type AudioProcessor interface {
    Plugin
    ProcessAudio(inputPath, outputPath string, options map[string]interface{}) error
    GetSupportedFormats() []string
}
```

### Subtitle Processor

For subtitle processing plugins:

```go
type SubtitleProcessor interface {
    Plugin
    ProcessSubtitles(inputPath, outputPath string, options map[string]interface{}) error
    GetSupportedFormats() []string
}
```

### Format Converter

For format conversion plugins:

```go
type FormatConverter interface {
    Plugin
    ConvertFormat(inputPath, outputPath, targetFormat string, options map[string]interface{}) error
    GetSupportedInputFormats() []string
    GetSupportedOutputFormats() []string
}
```

---

## Creating Your First Plugin

### Step 1: Define Your Plugin

```go
package myplugin

import "akashic_scribe/core"

type MyPlugin struct {
    *core.BasePlugin
    // Add custom fields
    customSetting string
}

func NewMyPlugin() core.Plugin {
    base := core.NewBasePlugin(
        "voidcat.my_plugin",
        "My First Plugin",
        "1.0.0",
        "A simple example plugin",
        "Your Name",
    )

    return &MyPlugin{
        BasePlugin: base,
        customSetting: "default",
    }
}
```

### Step 2: Implement Capabilities

```go
func (p *MyPlugin) GetCapabilities() []core.PluginCapability {
    return []core.PluginCapability{
        core.CapabilityAudioProcessing,
    }
}
```

### Step 3: Override Lifecycle Methods (Optional)

```go
func (p *MyPlugin) Initialize(context core.PluginContext) error {
    // Call base initialization
    if err := p.BasePlugin.Initialize(context); err != nil {
        return err
    }

    // Custom initialization
    p.customSetting = "initialized"
    context.LogInfo("Custom initialization complete")

    return nil
}

func (p *MyPlugin) Activate() error {
    if err := p.BasePlugin.Activate(); err != nil {
        return err
    }

    // Custom activation logic
    p.GetContext().LogInfo("Plugin is now active!")

    return nil
}
```

### Step 4: Implement Specialized Interface

```go
func (p *MyPlugin) ProcessAudio(inputPath, outputPath string, options map[string]interface{}) error {
    ctx := p.GetContext()
    ctx.LogInfo("Processing audio: " + inputPath)

    // Your processing logic here
    // ...

    ctx.LogInfo("Audio processing complete")
    return nil
}

func (p *MyPlugin) GetSupportedFormats() []string {
    return []string{"mp3", "wav", "flac"}
}

// Verify interface implementation
var _ core.AudioProcessor = (*MyPlugin)(nil)
```

### Step 5: Add Tests

```go
package myplugin

import (
    "testing"
    "akashic_scribe/core"
)

func TestMyPlugin(t *testing.T) {
    plugin := NewMyPlugin()

    // Test metadata
    if plugin.ID() != "voidcat.my_plugin" {
        t.Errorf("Expected ID voidcat.my_plugin, got %s", plugin.ID())
    }

    // Test capabilities
    caps := plugin.GetCapabilities()
    if len(caps) == 0 {
        t.Error("Plugin should have at least one capability")
    }

    // Test lifecycle
    manager := core.NewPluginManager("", "")
    if err := manager.LoadPlugin(plugin); err != nil {
        t.Fatalf("Failed to load plugin: %v", err)
    }

    if err := plugin.HealthCheck(); err != nil {
        t.Errorf("Health check failed: %v", err)
    }
}
```

---

## Built-in Example Plugins

### 1. Audio Effects Plugin

**ID**: `voidcat.audio_effects`
**Capability**: `audio_processing`

Advanced audio processing with:
- Noise reduction
- Normalization
- Equalization
- Reverb
- Bass/treble boost
- Dynamic range compression
- Fade in/out effects

**Example Usage:**

```go
plugin := audio_effects.NewAudioEffectsPlugin()

// Apply noise reduction and normalization
err := plugin.ProcessAudio("input.mp3", "output.mp3", map[string]interface{}{
    "noise_reduction": true,
    "normalization":   true,
    "compression":     true,
})

// Use preset
err = plugin.ApplyPodcastPreset("input.mp3", "output.mp3")

// Get available presets
presets := plugin.GetEffectPresets()
// Returns: clean_vocal, podcast, music, radio_voice
```

**Available Presets:**

- **clean_vocal**: Optimized for voice clarity
- **podcast**: Professional podcast settings
- **music**: Music enhancement
- **radio_voice**: Radio broadcast quality

### 2. Format Converter Plugin

**ID**: `voidcat.format_converter`
**Capability**: `format_conversion`, `subtitle_generation`

Convert between subtitle formats:
- SRT (SubRip)
- VTT (WebVTT)
- ASS/SSA (Advanced SubStation Alpha)
- TXT (Plain text)

**Example Usage:**

```go
plugin := format_converter.NewFormatConverterPlugin()

// Convert SRT to VTT
err := plugin.ConvertFormat("input.srt", "output.vtt", "vtt", nil)

// Convert VTT to ASS
err = plugin.ConvertFormat("input.vtt", "output.ass", "ass", nil)

// Get supported formats
inputFormats := plugin.GetSupportedInputFormats()
outputFormats := plugin.GetSupportedOutputFormats()
```

**Format Support:**

| Format | Input | Output | Notes |
|--------|-------|--------|-------|
| SRT | ✅ | ✅ | Most common format |
| VTT | ✅ | ✅ | HTML5 standard |
| ASS | ✅ | ✅ | Advanced styling |
| SSA | ✅ | ❌ | Converted to ASS |
| TXT | ✅ | ✅ | Plain text |

### 3. Subtitle Styler Plugin

**ID**: `voidcat.subtitle_styler`
**Capability**: `subtitle_styling`, `subtitle_generation`

Apply beautiful themes to subtitles:
- Professional styling presets
- Custom fonts, colors, and effects
- Position control (top/center/bottom)
- Background options

**Example Usage:**

```go
plugin := subtitle_styler.NewSubtitleStylerPlugin()

// Apply cinema theme
err := plugin.ApplyTheme("input.srt", "output.ass", "cinema")

// Custom styling
err = plugin.ProcessSubtitles("input.srt", "output.ass", map[string]interface{}{
    "theme":          "modern",
    "position":       "bottom",
    "font_size":      24,
    "add_background": true,
})

// Get available themes
themes := plugin.GetAvailableThemes()
// Returns: default, cinema, modern, elegant, bold_yellow, anime
```

**Available Themes:**

| Theme | Font | Size | Style | Best For |
|-------|------|------|-------|----------|
| **default** | Arial | 20 | White, Black outline | General use |
| **cinema** | Trebuchet MS | 24 | Bold, White | Movies |
| **modern** | Segoe UI | 22 | Light gray, Semi-transparent | Modern videos |
| **elegant** | Georgia | 20 | Italic, Light yellow | Artistic content |
| **bold_yellow** | Arial | 24 | Bold, Yellow | High visibility |
| **anime** | Arial | 22 | Bold, White | Anime/animation |

---

## Best Practices

### Plugin Design

1. **Single Responsibility**: Each plugin should focus on one primary capability
2. **Fail Gracefully**: Handle errors without crashing the host application
3. **Resource Management**: Clean up resources in `Shutdown()`
4. **Stateless Operations**: Avoid maintaining state between operations when possible
5. **Configuration**: Use the plugin context for configuration storage

### Error Handling

```go
func (p *MyPlugin) ProcessData(input string) error {
    ctx := p.GetContext()

    // Validate input
    if input == "" {
        ctx.LogError("Input cannot be empty")
        return fmt.Errorf("invalid input: empty string")
    }

    // Process with error handling
    if err := p.doProcessing(input); err != nil {
        ctx.LogError(fmt.Sprintf("Processing failed: %v", err))
        return fmt.Errorf("processing failed: %w", err)
    }

    ctx.LogInfo("Processing completed successfully")
    return nil
}
```

### Logging

Use appropriate log levels:

```go
ctx := p.GetContext()

// Informational messages
ctx.LogInfo("Starting audio processing")

// Warnings (non-fatal issues)
ctx.LogWarning("Input format is deprecated, consider upgrading")

// Errors (operation failures)
ctx.LogError("Failed to open file: permission denied")
```

### Performance

1. **Lazy Loading**: Load resources only when needed
2. **Caching**: Use the cache directory for temporary data
3. **Streaming**: Process large files in chunks
4. **Async Operations**: Consider using goroutines for long operations

```go
func (p *MyPlugin) ProcessLargeFile(inputPath string) error {
    // Use cache for temporary files
    cacheDir := p.GetContext().GetPluginCacheDir()
    tempFile := filepath.Join(cacheDir, "temp_processing.dat")
    defer os.Remove(tempFile)

    // Process in chunks
    // ...

    return nil
}
```

### Versioning

Follow semantic versioning:
- **Major**: Breaking changes
- **Minor**: New features (backward compatible)
- **Patch**: Bug fixes

```go
Version() string {
    return "1.2.3" // Major.Minor.Patch
}
```

---

## Testing Plugins

### Unit Tests

```go
func TestPluginInitialization(t *testing.T) {
    plugin := NewMyPlugin()

    // Test metadata
    assert.Equal(t, "voidcat.my_plugin", plugin.ID())
    assert.Equal(t, "1.0.0", plugin.Version())

    // Test capabilities
    caps := plugin.GetCapabilities()
    assert.NotEmpty(t, caps)
}

func TestPluginLifecycle(t *testing.T) {
    manager := core.NewPluginManager("", "")
    plugin := NewMyPlugin()

    // Test load
    err := manager.LoadPlugin(plugin)
    assert.NoError(t, err)

    // Test enable
    err = manager.EnablePlugin(plugin.ID())
    assert.NoError(t, err)

    // Test health
    err = manager.HealthCheck(plugin.ID())
    assert.NoError(t, err)

    // Test disable
    err = manager.DisablePlugin(plugin.ID())
    assert.NoError(t, err)
}
```

### Integration Tests

```go
func TestPluginProcessing(t *testing.T) {
    manager := core.NewPluginManager("", "")
    plugin := NewMyPlugin()

    // Setup
    manager.LoadPlugin(plugin)
    manager.EnablePlugin(plugin.ID())

    // Test actual processing
    err := plugin.ProcessData("test input")
    assert.NoError(t, err)

    // Cleanup
    manager.UnloadPlugin(plugin.ID())
}
```

### Running Tests

```bash
# Run all plugin tests
cd akashic_scribe/plugins
go test ./...

# Run with coverage
go test -cover ./...

# Verbose output
go test -v ./...

# Specific plugin
go test ./audio_effects/...
```

---

## Deployment

### Installing Plugins

Plugins are compiled into the main application. To add a new plugin:

1. Create plugin directory in `akashic_scribe/plugins/`
2. Implement plugin interface
3. Import in main application
4. Register with plugin manager

### Distribution

For distributing standalone plugins:

```bash
# Build plugin as shared library (if using plugin packages)
go build -buildmode=plugin -o myplugin.so ./plugins/myplugin

# Or include in main build
go build -o akashic_scribe .
```

### Configuration

Plugin configuration can be stored using the context:

```go
func (p *MyPlugin) Initialize(context core.PluginContext) error {
    // Load config
    config := context.GetConfig()

    if val, ok := config["my_setting"]; ok {
        p.mySetting = val.(string)
    }

    // Save config
    context.SetConfig("last_run", time.Now())

    return p.BasePlugin.Initialize(context)
}
```

---

## Troubleshooting

### Common Issues

#### Plugin Won't Load

**Symptom**: `LoadPlugin()` returns an error

**Solutions**:
1. Check plugin ID is unique
2. Verify `Initialize()` doesn't return error
3. Check dependencies are available
4. Review logs for specific error messages

#### Plugin Health Check Fails

**Symptom**: `HealthCheck()` returns an error

**Solutions**:
1. Ensure plugin is properly initialized
2. Check plugin is enabled
3. Verify required resources are available
4. Review plugin-specific health check logic

#### Plugin Not Found by Capability

**Symptom**: `GetPluginsByCapability()` returns empty list

**Solutions**:
1. Verify plugin is loaded and enabled
2. Check capability is correctly declared in `GetCapabilities()`
3. Ensure plugin manager has the plugin registered

### Debugging

Enable verbose logging:

```go
// In plugin implementation
func (p *MyPlugin) ProcessData(input string) error {
    ctx := p.GetContext()
    ctx.LogInfo(fmt.Sprintf("Processing started: %s", input))

    // ... processing ...

    ctx.LogInfo("Processing step 1 complete")
    // ... more steps ...

    return nil
}
```

Check plugin health:

```go
manager := core.NewPluginManager("", "")
// ... load plugins ...

// Debug health check
healthMap := manager.HealthCheckAll()
for id, err := range healthMap {
    if err != nil {
        fmt.Printf("❌ %s: %v\n", id, err)
    } else {
        fmt.Printf("✅ %s: healthy\n", id)
    }
}
```

### Getting Help

- **GitHub Issues**: https://github.com/sorrowscry86/akashic-scribe/issues
- **Documentation**: See `API_DOCUMENTATION.md` for detailed API reference
- **Examples**: Study built-in plugins in `akashic_scribe/plugins/`
- **Community**: GitHub Discussions

---

## Advanced Topics

### Plugin Dependencies

Declare dependencies on other plugins:

```go
func (p *MyPlugin) GetDependencies() []core.PluginDependency {
    return []core.PluginDependency{
        {
            PluginID:    "voidcat.audio_effects",
            Version:     "1.0.0",
            Required:    true,
            Description: "Requires audio effects for processing",
        },
    }
}
```

### Custom Plugin Types

Define your own plugin interfaces:

```go
type TranslationPlugin interface {
    core.Plugin
    Translate(text, sourceLang, targetLang string) (string, error)
    GetSupportedLanguages() []string
}
```

### Plugin Communication

Send messages between plugins (future feature):

```go
// Send to specific plugin
ctx.SendMessage("voidcat.other_plugin", myMessage)

// Broadcast to all
ctx.BroadcastMessage(myMessage)
```

---

## Appendix

### Plugin Checklist

Before releasing your plugin, ensure:

- [ ] Implements `Plugin` interface completely
- [ ] Has unique ID (format: `vendor.plugin_name`)
- [ ] Follows semantic versioning
- [ ] Has comprehensive error handling
- [ ] Includes unit tests
- [ ] Has proper documentation
- [ ] Handles cleanup in `Shutdown()`
- [ ] Uses logging appropriately
- [ ] Declares capabilities correctly
- [ ] Has example usage

### Resources

- **Core Interfaces**: `akashic_scribe/core/plugin.go`
- **Plugin Manager**: `akashic_scribe/core/plugin_manager.go`
- **Example Plugins**: `akashic_scribe/plugins/`
- **Demo Application**: `akashic_scribe/plugins/plugin_demo.go`
- **API Documentation**: `API_DOCUMENTATION.md`

---

**© 2025 VoidCat RDC. All rights reserved.**

*Building Excellence Through Extensibility*
