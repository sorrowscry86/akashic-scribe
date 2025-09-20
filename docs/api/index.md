# Akashic Scribe API Documentation

This document provides comprehensive information about the Akashic Scribe API, including interfaces, data structures, and usage examples.

## Table of Contents

1. [Core Engine API](#core-engine-api)
2. [GUI Integration API](#gui-integration-api)
3. [Plugin System API](#plugin-system-api)
4. [Configuration API](#configuration-api)
5. [Examples](#examples)

## Core Engine API

The Core Engine API provides the main functionality for transcription, translation, and dubbing operations.

### ScribeEngine Interface

```go
type ScribeEngine interface {
    // Transcribe a video file or URL to text
    Transcribe(options *ScribeOptions) (*TranscriptionResult, error)
    
    // Translate transcribed text to target languages
    Translate(transcription *TranscriptionResult, options *ScribeOptions) (*TranslationResult, error)
    
    // Generate subtitles from transcription and translation
    GenerateSubtitles(translation *TranslationResult, options *ScribeOptions) (*SubtitleResult, error)
    
    // Generate dubbed audio from translation
    GenerateDubbing(translation *TranslationResult, options *ScribeOptions) (*DubbingResult, error)
    
    // Process a video end-to-end (transcribe, translate, subtitle, dub)
    ProcessVideo(options *ScribeOptions) (*ProcessingResult, error)
    
    // Get progress information for ongoing operations
    GetProgress() ProgressInfo
    
    // Cancel ongoing operations
    Cancel() error
}
```

### Data Structures

#### ScribeOptions

```go
type ScribeOptions struct {
    // Input options
    InputType      InputType // File or URL
    InputPath      string    // File path or URL
    
    // Language options
    SourceLanguage string
    TargetLanguages []string
    
    // Subtitle options
    SubtitleFormat SubtitleFormat // SRT, VTT, ASS
    SubtitleStyles SubtitleStyles
    
    // Dubbing options
    DubbingEnabled bool
    VoiceModel     string
    AudioQuality   AudioQuality
    
    // Output options
    OutputDirectory string
    OutputFileName  string
    
    // Processing options
    ProcessingQuality ProcessingQuality // Speed, Balanced, Quality, Ultra
    BatchProcessing   bool
}
```

## GUI Integration API

The GUI Integration API allows for seamless integration between the user interface and the core engine.

### ProgressCallback

```go
type ProgressCallback func(progress ProgressInfo)

type ProgressInfo struct {
    Stage           ProcessingStage
    PercentComplete float64
    CurrentItem     string
    TotalItems      int
    CurrentItemIndex int
    EstimatedTimeRemaining time.Duration
    Error           error
}
```

## Plugin System API

The Plugin System API enables third-party extensions to enhance Akashic Scribe's functionality.

### Plugin Interface

```go
type Plugin interface {
    // Get plugin information
    GetInfo() PluginInfo
    
    // Initialize the plugin
    Initialize() error
    
    // Shutdown the plugin
    Shutdown() error
}

type PluginInfo struct {
    ID          string
    Name        string
    Version     string
    Description string
    Author      string
    Website     string
    Capabilities []PluginCapability
}
```

## Configuration API

The Configuration API provides access to user preferences and application settings.

```go
type ConfigManager interface {
    // Get a configuration value
    Get(key string) interface{}
    
    // Set a configuration value
    Set(key string, value interface{}) error
    
    // Save configuration to disk
    Save() error
    
    // Load configuration from disk
    Load() error
    
    // Reset configuration to defaults
    Reset() error
}
```

## Examples

### Basic Usage Example

```go
package main

import (
    "akashic_scribe/core"
    "fmt"
)

func main() {
    // Create a new engine
    engine := core.NewRealScribeEngine()
    
    // Configure options
    options := &core.ScribeOptions{
        InputType:      core.InputTypeFile,
        InputPath:      "/path/to/video.mp4",
        SourceLanguage: "en",
        TargetLanguages: []string{"es", "fr", "de"},
        SubtitleFormat: core.SubtitleFormatSRT,
        DubbingEnabled: true,
        VoiceModel:     "neural-standard",
        OutputDirectory: "/path/to/output",
    }
    
    // Process the video
    result, err := engine.ProcessVideo(options)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    // Print results
    fmt.Printf("Transcription: %d characters\n", len(result.Transcription.Text))
    fmt.Printf("Translations: %d languages\n", len(result.Translations))
    fmt.Printf("Subtitles: %d files\n", len(result.SubtitleFiles))
    fmt.Printf("Dubbed Audio: %d files\n", len(result.DubbingFiles))
}
```

For more examples and detailed API documentation, please refer to the specific API sections in this directory.