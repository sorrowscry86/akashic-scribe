# Akashic Scribe Plugin System

**VoidCat RDC - Extensible Architecture**

## Overview

The Akashic Scribe plugin system provides a powerful, flexible framework for extending the application's capabilities. Built with Go's strong type system and interface-based design, it enables developers to create custom functionality while maintaining code quality and reliability.

## Table of Contents

- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [Available Plugins](#available-plugins)
- [Creating Plugins](#creating-plugins)
- [Plugin Manager](#plugin-manager)
- [Examples](#examples)
- [Documentation](#documentation)

---

## Quick Start

### Using the Plugin System

```go
import (
    "akashic_scribe/core"
    "akashic_scribe/plugins/audio_effects"
    "akashic_scribe/plugins/format_converter"
    "akashic_scribe/plugins/subtitle_styler"
    "akashic_scribe/plugins/elevenlabs"
)

// Create plugin manager
manager := core.NewPluginManager("", "")

// Create and register plugins
audioPlugin := audio_effects.NewAudioEffectsPlugin()
formatPlugin := format_converter.NewFormatConverterPlugin()
stylerPlugin := subtitle_styler.NewSubtitleStylerPlugin()
elevenLabsPlugin := elevenlabs.NewElevenLabsPlugin()

manager.LoadPlugin(audioPlugin)
manager.LoadPlugin(formatPlugin)
manager.LoadPlugin(stylerPlugin)
manager.LoadPlugin(elevenLabsPlugin)

manager.EnablePlugin(audioPlugin.ID())
manager.EnablePlugin(formatPlugin.ID())
manager.EnablePlugin(stylerPlugin.ID())
manager.EnablePlugin(elevenLabsPlugin.ID())

// Use plugins
audioPlugin.ApplyNoiseReduction("input.mp3", "output.mp3")
formatPlugin.ConvertFormat("input.srt", "output.vtt", "vtt", nil)
stylerPlugin.ApplyTheme("input.srt", "output.ass", "cinema")
elevenLabsPlugin.GenerateSpeech("Hello world", "speech.mp3", options)
```

### Running the Demo

```bash
cd akashic_scribe/plugins
go run plugin_demo.go
```

---

## Architecture

### Component Overview

```
┌──────────────────────────────────────────────┐
│         Core Plugin System                    │
│  ┌────────────────────────────────────────┐  │
│  │  Plugin Interface                      │  │
│  │  - Metadata (ID, Name, Version)        │  │
│  │  - Lifecycle (Init, Activate, Shutdown)│  │
│  │  - Capabilities                        │  │
│  │  - Health Monitoring                   │  │
│  └────────────────────────────────────────┘  │
│                                               │
│  ┌────────────────────────────────────────┐  │
│  │  Plugin Manager                        │  │
│  │  - Registration & Loading              │  │
│  │  - Capability-based Discovery          │  │
│  │  - Lifecycle Management                │  │
│  │  - Health Checks                       │  │
│  │  - Inter-plugin Communication          │  │
│  └────────────────────────────────────────┘  │
│                                               │
│  ┌────────────────────────────────────────┐  │
│  │  Plugin Context                        │  │
│  │  - Configuration Storage               │  │
│  │  - Logging (Info, Warning, Error)      │  │
│  │  - Resource Access (Data, Cache)       │  │
│  │  - Message Passing                     │  │
│  └────────────────────────────────────────┘  │
└──────────────────────────────────────────────┘
                      │
    ┌─────────────────┼─────────────────┬────────────────┐
    │                 │                 │                │
┌───v──────────┐  ┌──v───────────┐  ┌──v──────────┐  ┌─v───────────┐
│ Audio Effects│  │Format Convert│  │Subtitle Style│  │ ElevenLabs  │
│  Plugin      │  │   Plugin     │  │   Plugin     │  │   Plugin    │
└──────────────┘  └──────────────┘  └──────────────┘  └─────────────┘
```

### Plugin Lifecycle

1. **Registration**: Plugin is registered with the manager
2. **Loading**: Plugin is initialized with context
3. **Activation**: Plugin becomes active and operational
4. **Operation**: Plugin performs its functions
5. **Deactivation**: Plugin is temporarily disabled (optional)
6. **Shutdown**: Plugin releases resources and terminates

### Capability System

Plugins declare their capabilities for easy discovery:

- `audio_processing` - Audio manipulation and effects
- `video_processing` - Video file processing
- `transcription` - Speech-to-text conversion
- `translation` - Language translation
- `voice_synthesis` - Text-to-speech generation
- `subtitle_generation` - Subtitle file creation
- `subtitle_styling` - Visual subtitle theming
- `format_conversion` - File format conversion
- `ui_extension` - User interface extensions
- `theme_provider` - Custom theme provision
- `file_format` - Additional file format support
- `cloud_integration` - Cloud service integration
- `api_integration` - External API integration

---

## Available Plugins

### 1. Audio Effects Plugin

**ID**: `voidcat.audio_effects`
**Version**: 1.0.0
**Capabilities**: `audio_processing`

Professional audio processing suite with:
- Noise reduction
- Audio normalization
- Equalization (custom frequency bands)
- Reverb effects
- Bass and treble boost
- Dynamic range compression
- Fade in/out effects

**Quick Example:**
```go
plugin.ApplyPodcastPreset("input.mp3", "output.mp3")
```

[Full Documentation →](akashic_scribe/plugins/audio_effects/)

---

### 2. Format Converter Plugin

**ID**: `voidcat.format_converter`
**Version**: 1.0.0
**Capabilities**: `format_conversion`, `subtitle_generation`

Convert between subtitle formats:
- SRT (SubRip)
- VTT (WebVTT)
- ASS/SSA (Advanced SubStation Alpha)
- TXT (Plain text)

**Quick Example:**
```go
plugin.ConvertFormat("input.srt", "output.vtt", "vtt", nil)
```

[Full Documentation →](akashic_scribe/plugins/format_converter/)

---

### 3. Subtitle Styler Plugin

**ID**: `voidcat.subtitle_styler`
**Version**: 1.0.0
**Capabilities**: `subtitle_styling`, `subtitle_generation`

Apply professional themes to subtitles:
- 6 built-in themes (cinema, modern, elegant, anime, etc.)
- Custom fonts, colors, and effects
- Position control (top/center/bottom)
- Background and shadow options

**Quick Example:**
```go
plugin.ApplyTheme("input.srt", "output.ass", "cinema")
```

[Full Documentation →](akashic_scribe/plugins/subtitle_styler/)

---

### 4. ElevenLabs Plugin

**ID**: `voidcat.elevenlabs`
**Version**: 1.0.0
**Capabilities**: `voice_synthesis`, `api_integration`

AI-powered voice synthesis with ElevenLabs:
- Natural-sounding voice generation
- 20+ professional voice models
- Multilingual support
- Advanced voice parameter tuning
- Voice cloning support
- Multiple quality presets

**Quick Example:**
```go
plugin.GenerateWithPreset(text, voiceID, "natural", "output.mp3")
```

[Full Documentation →](akashic_scribe/plugins/elevenlabs/)

---

## Creating Plugins

### Minimal Plugin Template

```go
package myplugin

import "akashic_scribe/core"

type MyPlugin struct {
    *core.BasePlugin
}

func NewMyPlugin() core.Plugin {
    base := core.NewBasePlugin(
        "vendor.myplugin",        // Unique ID
        "My Plugin",              // Display name
        "1.0.0",                  // Version
        "Does something cool",    // Description
        "Your Name",              // Author
    )

    return &MyPlugin{
        BasePlugin: base,
    }
}

func (p *MyPlugin) GetCapabilities() []core.PluginCapability {
    return []core.PluginCapability{
        core.CapabilityAudioProcessing,
    }
}

// Implement specialized interface methods...
```

### Specialized Interfaces

#### Audio Processor

```go
type AudioProcessor interface {
    Plugin
    ProcessAudio(inputPath, outputPath string, options map[string]interface{}) error
    GetSupportedFormats() []string
}
```

#### Subtitle Processor

```go
type SubtitleProcessor interface {
    Plugin
    ProcessSubtitles(inputPath, outputPath string, options map[string]interface{}) error
    GetSupportedFormats() []string
}
```

#### Format Converter

```go
type FormatConverter interface {
    Plugin
    ConvertFormat(inputPath, outputPath, targetFormat string, options map[string]interface{}) error
    GetSupportedInputFormats() []string
    GetSupportedOutputFormats() []string
}
```

---

## Plugin Manager

### Creating a Manager

```go
// With custom directories
manager := core.NewPluginManager("/path/to/data", "/path/to/cache")

// With default temp directories
manager := core.NewPluginManager("", "")
```

### Loading Plugins

```go
plugin := NewMyPlugin()

// Register plugin
manager.RegisterPlugin(plugin)

// Load and initialize
manager.LoadPlugin(plugin)

// Enable/activate
manager.EnablePlugin(plugin.ID())
```

### Querying Plugins

```go
// Get all loaded plugins
plugins := manager.GetLoadedPlugins()

// Get plugin by ID
plugin, found := manager.GetPluginByID("voidcat.audio_effects")

// Query by capability
audioPlugins := manager.GetPluginsByCapability(core.CapabilityAudioProcessing)

// Get plugin information
info, err := manager.GetPluginInfo("voidcat.audio_effects")
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
err := manager.HealthCheck("voidcat.audio_effects")
```

---

## Examples

### Audio Processing Pipeline

```go
// Create pipeline: Download -> Extract -> Enhance -> Generate Speech
manager := core.NewPluginManager("", "")

// Load plugins
audioPlugin := audio_effects.NewAudioEffectsPlugin()
elevenLabsPlugin := elevenlabs.NewElevenLabsPlugin()

manager.LoadPlugin(audioPlugin)
manager.LoadPlugin(elevenLabsPlugin)
manager.EnablePlugin(audioPlugin.ID())
manager.EnablePlugin(elevenLabsPlugin.ID())

// Process
audioPlugin.ProcessAudio("raw.mp3", "enhanced.mp3", map[string]interface{}{
    "noise_reduction": true,
    "normalization":   true,
    "compression":     true,
})

elevenLabsPlugin.GenerateSpeech("Dubbed audio text", "dubbed.mp3", options)
```

### Subtitle Workflow

```go
// Convert format, apply styling, and sync with audio
formatPlugin := format_converter.NewFormatConverterPlugin()
stylerPlugin := subtitle_styler.NewSubtitleStylerPlugin()

// SRT -> VTT
formatPlugin.ConvertFormat("input.srt", "intermediate.vtt", "vtt", nil)

// VTT -> Styled ASS
stylerPlugin.ProcessSubtitles("intermediate.vtt", "final.ass", map[string]interface{}{
    "theme":     "cinema",
    "font_size": 24,
    "position":  "bottom",
})
```

### Multi-Plugin Processing

```go
// Use capability-based discovery
manager := core.NewPluginManager("", "")

// Load all plugins
// ...

// Get all audio processors
audioProcessors := manager.GetPluginsByCapability(core.CapabilityAudioProcessing)

for _, processor := range audioProcessors {
    if ap, ok := processor.(core.AudioProcessor); ok {
        ap.ProcessAudio(input, output, options)
    }
}
```

---

## Documentation

### Comprehensive Guides

1. **[Plugin Development Guide](PLUGIN_DEVELOPMENT_GUIDE.md)** - Complete guide to creating plugins
   - Architecture overview
   - Step-by-step tutorials
   - Best practices
   - Testing guidelines

2. **[API Documentation](API_DOCUMENTATION.md)** - Detailed API reference
   - Core interfaces
   - Plugin manager API
   - Context and configuration

3. **Plugin-Specific Docs**:
   - [Audio Effects Plugin](akashic_scribe/plugins/audio_effects/)
   - [Format Converter Plugin](akashic_scribe/plugins/format_converter/)
   - [Subtitle Styler Plugin](akashic_scribe/plugins/subtitle_styler/)
   - [ElevenLabs Plugin](akashic_scribe/plugins/elevenlabs/README.md)

### Quick References

**Plugin Interface Methods:**
```go
ID() string
Name() string
Version() string
Description() string
Author() string
Initialize(context PluginContext) error
Activate() error
Deactivate() error
Shutdown() error
GetCapabilities() []PluginCapability
GetDependencies() []PluginDependency
HealthCheck() error
```

**Plugin Manager Methods:**
```go
RegisterPlugin(plugin Plugin) error
LoadPlugin(plugin Plugin) error
UnloadPlugin(pluginID string) error
EnablePlugin(pluginID string) error
DisablePlugin(pluginID string) error
GetPluginByID(pluginID string) (Plugin, bool)
GetPluginsByCapability(capability PluginCapability) []Plugin
HealthCheckAll() map[string]error
```

---

## File Structure

```
akashic_scribe/
├── core/
│   ├── plugin.go                 # Plugin interfaces
│   ├── plugin_manager.go         # Plugin manager
│   └── ...
├── plugins/
│   ├── audio_effects/
│   │   └── audio_effects.go
│   ├── format_converter/
│   │   └── format_converter.go
│   ├── subtitle_styler/
│   │   └── subtitle_styler.go
│   ├── elevenlabs/
│   │   ├── elevenlabs.go
│   │   └── README.md
│   └── plugin_demo.go            # Demo application
├── PLUGIN_DEVELOPMENT_GUIDE.md   # Development guide
└── PLUGIN_SYSTEM.md              # This file
```

---

## Testing

### Running Plugin Tests

```bash
# Test all plugins
cd akashic_scribe/plugins
go test ./...

# Test specific plugin
go test ./audio_effects/...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

### Demo Application

```bash
# Run the plugin demo
cd akashic_scribe/plugins
go run plugin_demo.go
```

The demo shows:
- Plugin registration and loading
- Capability queries
- Health monitoring
- Example usage of each plugin

---

## Best Practices

### Plugin Development

1. **Single Responsibility**: Each plugin should focus on one primary capability
2. **Error Handling**: Always return descriptive errors
3. **Resource Cleanup**: Release resources in `Shutdown()`
4. **Logging**: Use appropriate log levels (Info, Warning, Error)
5. **Configuration**: Store settings in plugin context
6. **Testing**: Write comprehensive unit and integration tests

### Performance

1. **Lazy Loading**: Load resources only when needed
2. **Caching**: Use cache directory for temporary data
3. **Streaming**: Process large files in chunks
4. **Async Operations**: Use goroutines for long operations

### Security

1. **API Keys**: Never hardcode keys, use environment variables
2. **Input Validation**: Validate all user inputs
3. **Path Sanitization**: Sanitize file paths
4. **Rate Limiting**: Implement retry logic for APIs

---

## Troubleshooting

### Plugin Won't Load

**Check:**
- Plugin ID is unique
- `Initialize()` succeeds
- Required dependencies are available
- Check logs for specific errors

### Health Check Fails

**Check:**
- Plugin is properly initialized
- Plugin is enabled
- Required resources are available
- Plugin-specific requirements are met

### Plugin Not Found by Capability

**Check:**
- Plugin is loaded and enabled
- Capability is correctly declared
- Plugin manager has the plugin registered

---

## Roadmap

### Upcoming Features

- [ ] Hot-reload plugin support
- [ ] Plugin marketplace/registry
- [ ] Remote plugin loading
- [ ] Plugin dependencies resolution
- [ ] Enhanced inter-plugin messaging
- [ ] Plugin sandboxing
- [ ] Performance profiling tools
- [ ] Visual plugin builder

### Community Plugins

We welcome community-developed plugins! See [PLUGIN_DEVELOPMENT_GUIDE.md](PLUGIN_DEVELOPMENT_GUIDE.md) for details on creating and submitting plugins.

---

## Support

- **GitHub Issues**: [Report bugs or request features](https://github.com/sorrowscry86/akashic-scribe/issues)
- **Discussions**: [Community Q&A](https://github.com/sorrowscry86/akashic-scribe/discussions)
- **Email**: SorrowsCry86@voidcat.org
- **Documentation**: See `PLUGIN_DEVELOPMENT_GUIDE.md` for detailed guides

---

## License

The plugin system is part of Akashic Scribe and follows the same MIT license.

---

**© 2025 VoidCat RDC. All rights reserved.**

*Building Extensible Excellence*
