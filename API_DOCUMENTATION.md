# Akashic Scribe API Documentation

**VoidCat RDC - Technical API Reference**

This document provides comprehensive technical information about the APIs, interfaces, and extension points used in Akashic Scribe.

## Table of Contents

- [Overview](#overview)
- [Core Engine Interface](#core-engine-interface)
- [ScribeOptions Configuration](#scribeoptions-configuration)
- [GUI Components API](#gui-components-api)
- [State Management API](#state-management-api)
- [Plugin Development API](#plugin-development-api)
- [Theme Customization API](#theme-customization-api)
- [Event System](#event-system)
- [Extension Points](#extension-points)
- [Error Handling](#error-handling)
- [Performance Monitoring](#performance-monitoring)
- [Integration Examples](#integration-examples)

## Overview

Akashic Scribe exposes a comprehensive API surface designed for extensibility, integration, and customization. The API follows VoidCat RDC's design principles of clean interfaces, strong typing, and predictable behavior.

### API Design Principles

1. **Interface Segregation**: Focused, single-purpose interfaces
2. **Dependency Injection**: Constructor-based dependency injection
3. **Error Transparency**: Comprehensive error information
4. **Async-First**: Non-blocking operations with progress reporting
5. **Type Safety**: Strong typing throughout the API surface

### Compatibility Promise

- **Major Version**: Breaking changes only with major version increments
- **Minor Version**: Backward-compatible feature additions
- **Patch Version**: Bug fixes and non-breaking improvements
- **API Versioning**: Explicit versioning for all public interfaces

## Core Engine Interface

### ScribeEngine Interface

The primary interface for all video processing operations:

```go
// ScribeEngine defines the contract for video processing operations.
// All implementations must be thread-safe and handle concurrent requests.
type ScribeEngine interface {
    // Core Processing Operations
    
    // Transcribe converts video audio to text in the original language.
    // videoSource can be a local file path or a URL.
    // Returns transcribed text or error if processing fails.
    Transcribe(videoSource string) (string, error)
    
    // Translate converts text from origin to target language.
    // Uses intelligent language detection if not specified.
    Translate(text string, targetLanguage string) (string, error)
    
    // GenerateSubtitles creates formatted subtitle files.
    // Supports SRT, VTT, and ASS formats based on options.
    GenerateSubtitles(transcription, translation string, options *ScribeOptions) (string, error)
    
    // GenerateDubbing creates dubbed audio using specified voice model.
    // Returns path to generated audio file or error.
    GenerateDubbing(translation string, voiceModel string, options *ScribeOptions) (string, error)
    
    // Workflow Management
    
    // StartProcessing initiates the complete workflow.
    // Progress is reported through the progressChan channel.
    // Returns final result path or error.
    StartProcessing(options *ScribeOptions, progressChan chan<- float64) (string, error)
    
    // CancelProcessing stops ongoing processing operations.
    // Returns immediately; actual cancellation may take time.
    CancelProcessing() error
    
    // Configuration and Status
    
    // GetSupportedLanguages returns list of supported language codes.
    GetSupportedLanguages() []LanguageInfo
    
    // GetAvailableVoiceModels returns available voice synthesis models.
    GetAvailableVoiceModels() []VoiceModelInfo
    
    // GetProcessingStatus returns current processing status.
    GetProcessingStatus() ProcessingStatus
    
    // Health and Diagnostics
    
    // HealthCheck verifies engine is ready for operations.
    HealthCheck() error
    
    // GetMetrics returns processing performance metrics.
    GetMetrics() ProcessingMetrics
}
```

### Supporting Data Structures

```go
// LanguageInfo provides details about supported languages
type LanguageInfo struct {
    Code        string   `json:"code"`         // RFC 5646 language tag (e.g., "en-US")
    Name        string   `json:"name"`         // Display name (e.g., "English (United States)")
    NativeName  string   `json:"native_name"`  // Native language name
    Variants    []string `json:"variants"`     // Regional variants
    Confidence  float64  `json:"confidence"`   // Detection confidence (0.0-1.0)
    Features    []string `json:"features"`     // Supported features (transcription, translation, etc.)
}

// VoiceModelInfo describes available voice synthesis models
type VoiceModelInfo struct {
    ID           string            `json:"id"`           // Unique model identifier
    Name         string            `json:"name"`         // Display name
    Language     string            `json:"language"`     // Supported language code
    Gender       string            `json:"gender"`       // Voice gender (male, female, neutral)
    Style        string            `json:"style"`        // Voice style (casual, formal, energetic, etc.)
    Quality      VoiceQuality      `json:"quality"`      // Quality rating
    Parameters   []VoiceParameter  `json:"parameters"`   // Adjustable parameters
    SampleRate   int               `json:"sample_rate"`  // Audio sample rate
    IsCustom     bool              `json:"is_custom"`    // Whether this is a custom model
    Metadata     map[string]string `json:"metadata"`     // Additional metadata
}

// VoiceQuality represents voice model quality levels
type VoiceQuality string

const (
    VoiceQualityBasic       VoiceQuality = "basic"       // Basic quality, fast processing
    VoiceQualityStandard    VoiceQuality = "standard"    // Standard quality, balanced
    VoiceQualityPremium     VoiceQuality = "premium"     // High quality, slower processing
    VoiceQualityUltra       VoiceQuality = "ultra"       // Highest quality, longest processing
)

// VoiceParameter defines adjustable voice model parameters
type VoiceParameter struct {
    Name        string      `json:"name"`         // Parameter name
    Type        string      `json:"type"`         // Parameter type (float, int, string, bool)
    Default     interface{} `json:"default"`      // Default value
    Min         interface{} `json:"min"`          // Minimum value (for numeric types)
    Max         interface{} `json:"max"`          // Maximum value (for numeric types)
    Options     []string    `json:"options"`      // Valid options (for enum types)
    Description string      `json:"description"`  // Human-readable description
}

// ProcessingStatus provides real-time status information
type ProcessingStatus struct {
    IsActive     bool      `json:"is_active"`      // Whether processing is currently active
    CurrentStage string    `json:"current_stage"`  // Current processing stage
    Progress     float64   `json:"progress"`       // Overall progress (0.0-1.0)
    StageProgress float64  `json:"stage_progress"` // Current stage progress (0.0-1.0)
    StartTime    time.Time `json:"start_time"`     // Processing start time
    ETA          time.Time `json:"eta"`            // Estimated completion time
    Message      string    `json:"message"`        // Current status message
    ErrorCount   int       `json:"error_count"`    // Number of recoverable errors
}

// ProcessingMetrics provides performance and quality metrics
type ProcessingMetrics struct {
    ProcessingTime    time.Duration `json:"processing_time"`     // Total processing time
    TranscriptionTime time.Duration `json:"transcription_time"`  // Time spent on transcription
    TranslationTime   time.Duration `json:"translation_time"`    // Time spent on translation
    DubbingTime       time.Duration `json:"dubbing_time"`        // Time spent on dubbing
    
    // Quality metrics
    TranscriptionConfidence float64 `json:"transcription_confidence"` // Average confidence score
    TranslationQuality      float64 `json:"translation_quality"`      // Translation quality score
    AudioQuality           float64 `json:"audio_quality"`             // Generated audio quality
    
    // Resource usage
    PeakMemoryUsage   int64 `json:"peak_memory_usage"`   // Peak memory usage in bytes
    AverageCPUUsage   float64 `json:"average_cpu_usage"`   // Average CPU usage percentage
    TotalDiskIO       int64 `json:"total_disk_io"`       // Total disk I/O in bytes
}
```

### Engine Implementation

```go
// Factory function for creating engine instances
func NewScribeEngine(config EngineConfig) ScribeEngine {
    return &realScribeEngine{
        config:          config,
        transcriber:     NewTranscriber(config.TranscriptionConfig),
        translator:      NewTranslator(config.TranslationConfig),
        dubber:          NewDubber(config.DubbingConfig),
        progressTracker: NewProgressTracker(),
        metrics:         NewMetricsCollector(),
    }
}

// EngineConfig configures the ScribeEngine instance
type EngineConfig struct {
    // Service configurations
    TranscriptionConfig TranscriptionConfig `json:"transcription_config"`
    TranslationConfig   TranslationConfig   `json:"translation_config"`
    DubbingConfig       DubbingConfig       `json:"dubbing_config"`
    
    // Performance settings
    MaxConcurrentJobs int           `json:"max_concurrent_jobs"` // Maximum parallel processing jobs
    TimeoutDuration   time.Duration `json:"timeout_duration"`   // Processing timeout
    RetryAttempts     int           `json:"retry_attempts"`     // Number of retry attempts
    
    // Quality settings
    QualityPreset     QualityPreset `json:"quality_preset"`      // Overall quality preset
    EnableMetrics     bool          `json:"enable_metrics"`      // Whether to collect metrics
    EnableCaching     bool          `json:"enable_caching"`      // Whether to enable result caching
    
    // Security settings
    AllowURLProcessing bool     `json:"allow_url_processing"` // Whether to allow URL inputs
    AllowedDomains     []string `json:"allowed_domains"`      // Whitelist of allowed domains
    MaxFileSize        int64    `json:"max_file_size"`        // Maximum input file size
}

// QualityPreset defines processing quality levels
type QualityPreset string

const (
    QualityPresetSpeed    QualityPreset = "speed"    // Optimized for speed
    QualityPresetBalanced QualityPreset = "balanced" // Balanced speed and quality
    QualityPresetQuality  QualityPreset = "quality"  // Optimized for quality
    QualityPresetUltra    QualityPreset = "ultra"    // Maximum quality, no speed constraints
)
```

## ScribeOptions Configuration

### Primary Configuration Structure

```go
// ScribeOptions represents all user configuration for processing operations.
// This structure is the single source of truth for processing parameters.
type ScribeOptions struct {
    // Input Configuration
    InputFile string `json:"input_file" validate:"required_without=InputURL,file"`
    InputURL  string `json:"input_url" validate:"required_without=InputFile,url"`
    
    // Language Configuration
    OriginLanguage string `json:"origin_language" validate:"required,language_code"`
    TargetLanguage string `json:"target_language" validate:"required,language_code"`
    
    // Processing Options
    ProcessingQuality QualityPreset `json:"processing_quality"`
    Priority          int           `json:"priority" validate:"min=1,max=10"`
    
    // Subtitle Configuration
    SubtitleOptions SubtitleOptions `json:"subtitle_options"`
    
    // Dubbing Configuration
    DubbingOptions DubbingOptions `json:"dubbing_options"`
    
    // Output Configuration
    OutputOptions OutputOptions `json:"output_options"`
    
    // Advanced Configuration
    AdvancedOptions AdvancedOptions `json:"advanced_options"`
    
    // Metadata
    Metadata ProcessingMetadata `json:"metadata"`
}

// SubtitleOptions configures subtitle generation
type SubtitleOptions struct {
    Enabled           bool            `json:"enabled"`
    Format            SubtitleFormat  `json:"format"`
    BilingualMode     bool            `json:"bilingual_mode"`
    Position          SubtitlePosition `json:"position"`
    Styling           SubtitleStyling `json:"styling"`
    TimingPrecision   time.Duration   `json:"timing_precision"`
    MaxLineLength     int             `json:"max_line_length"`
    MaxLinesPerEntry  int             `json:"max_lines_per_entry"`
    IncludeSpeakerIDs bool            `json:"include_speaker_ids"`
}

// SubtitleFormat defines supported subtitle formats
type SubtitleFormat string

const (
    SubtitleFormatSRT SubtitleFormat = "srt" // SubRip Text format
    SubtitleFormatVTT SubtitleFormat = "vtt" // WebVTT format
    SubtitleFormatASS SubtitleFormat = "ass" // Advanced SubStation Alpha format
)

// SubtitlePosition defines subtitle placement
type SubtitlePosition string

const (
    SubtitlePositionBottom SubtitlePosition = "bottom"
    SubtitlePositionTop    SubtitlePosition = "top"
    SubtitlePositionCustom SubtitlePosition = "custom"
)

// SubtitleStyling configures subtitle appearance
type SubtitleStyling struct {
    FontFamily     string  `json:"font_family"`
    FontSize       int     `json:"font_size"`
    FontColor      string  `json:"font_color"`
    BackgroundColor string `json:"background_color"`
    OutlineColor   string  `json:"outline_color"`
    OutlineWidth   int     `json:"outline_width"`
    Bold           bool    `json:"bold"`
    Italic         bool    `json:"italic"`
    Alignment      string  `json:"alignment"`
}

// DubbingOptions configures audio dubbing generation
type DubbingOptions struct {
    Enabled         bool                   `json:"enabled"`
    VoiceModel      string                 `json:"voice_model"`
    UseCustomVoice  bool                   `json:"use_custom_voice"`
    CustomVoicePath string                 `json:"custom_voice_path"`
    VoiceParameters map[string]interface{} `json:"voice_parameters"`
    AudioFormat     AudioFormat            `json:"audio_format"`
    Quality         AudioQuality           `json:"quality"`
    SampleRate      int                    `json:"sample_rate"`
    BitRate         int                    `json:"bit_rate"`
    Channels        int                    `json:"channels"`
    NormalizeAudio  bool                   `json:"normalize_audio"`
    RemoveOriginal  bool                   `json:"remove_original"`
}

// AudioFormat defines supported audio output formats
type AudioFormat string

const (
    AudioFormatMP3  AudioFormat = "mp3"
    AudioFormatWAV  AudioFormat = "wav"
    AudioFormatFLAC AudioFormat = "flac"
    AudioFormatAAC  AudioFormat = "aac"
    AudioFormatOGG  AudioFormat = "ogg"
)

// AudioQuality defines audio quality levels
type AudioQuality string

const (
    AudioQualityLow      AudioQuality = "low"      // 64-128 kbps
    AudioQualityMedium   AudioQuality = "medium"   // 128-192 kbps
    AudioQualityHigh     AudioQuality = "high"     // 192-320 kbps
    AudioQualityLossless AudioQuality = "lossless" // FLAC or WAV
)

// OutputOptions configures output generation
type OutputOptions struct {
    OutputDirectory   string            `json:"output_directory"`
    FileNamingPattern string            `json:"file_naming_pattern"`
    CreateProject     bool              `json:"create_project"`
    ProjectName       string            `json:"project_name"`
    IncludeOriginal   bool              `json:"include_original"`
    CompressOutput    bool              `json:"compress_output"`
    Metadata          map[string]string `json:"metadata"`
}

// AdvancedOptions provides fine-grained control
type AdvancedOptions struct {
    // Processing control
    ChunkSize          time.Duration          `json:"chunk_size"`
    OverlapDuration    time.Duration          `json:"overlap_duration"`
    ParallelProcessing bool                   `json:"parallel_processing"`
    MaxWorkers         int                    `json:"max_workers"`
    
    // Quality control
    ConfidenceThreshold float64               `json:"confidence_threshold"`
    NoiseReduction      bool                  `json:"noise_reduction"`
    AudioFilters        []AudioFilter         `json:"audio_filters"`
    
    // Custom processing
    PreProcessingHooks  []ProcessingHook      `json:"pre_processing_hooks"`
    PostProcessingHooks []ProcessingHook      `json:"post_processing_hooks"`
    CustomParameters    map[string]interface{} `json:"custom_parameters"`
}

// ProcessingMetadata provides context about the processing request
type ProcessingMetadata struct {
    JobID          string            `json:"job_id"`
    UserID         string            `json:"user_id"`
    CreatedAt      time.Time         `json:"created_at"`
    RequestSource  string            `json:"request_source"`
    Tags           []string          `json:"tags"`
    CustomMetadata map[string]string `json:"custom_metadata"`
}
```

### Configuration Validation

```go
// Validate performs comprehensive validation of ScribeOptions
func (s *ScribeOptions) Validate() error {
    validator := validator.New()
    
    // Register custom validation functions
    validator.RegisterValidation("language_code", validateLanguageCode)
    validator.RegisterValidation("file", validateFileExists)
    
    if err := validator.Struct(s); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    // Custom validation logic
    if err := s.validateInputSources(); err != nil {
        return err
    }
    
    if err := s.validateLanguageCompatibility(); err != nil {
        return err
    }
    
    if err := s.validateOutputConfiguration(); err != nil {
        return err
    }
    
    return nil
}

// validateInputSources ensures at least one input source is specified
func (s *ScribeOptions) validateInputSources() error {
    if s.InputFile == "" && s.InputURL == "" {
        return errors.New("either input file or input URL must be specified")
    }
    
    if s.InputFile != "" && s.InputURL != "" {
        return errors.New("cannot specify both input file and input URL")
    }
    
    return nil
}

// ToJSON serializes ScribeOptions to JSON
func (s *ScribeOptions) ToJSON() ([]byte, error) {
    return json.MarshalIndent(s, "", "  ")
}

// FromJSON deserializes ScribeOptions from JSON
func (s *ScribeOptions) FromJSON(data []byte) error {
    if err := json.Unmarshal(data, s); err != nil {
        return fmt.Errorf("failed to unmarshal ScribeOptions: %w", err)
    }
    
    return s.Validate()
}

// Clone creates a deep copy of ScribeOptions
func (s *ScribeOptions) Clone() *ScribeOptions {
    data, _ := s.ToJSON()
    clone := &ScribeOptions{}
    clone.FromJSON(data)
    return clone
}
```

## GUI Components API

### Layout Management

```go
// LayoutManager handles UI component composition and management
type LayoutManager interface {
    // CreateMainLayout constructs the primary UI structure
    CreateMainLayout(window fyne.Window, engine ScribeEngine) fyne.CanvasObject
    
    // CreateStepCard creates individual workflow step cards
    CreateStepCard(step WorkflowStep, state *GUIState) fyne.CanvasObject
    
    // UpdateLayout refreshes layout based on state changes
    UpdateLayout(state *GUIState) error
    
    // GetComponent retrieves specific UI components by ID
    GetComponent(componentID string) fyne.CanvasObject
}

// WorkflowStep represents individual steps in the processing workflow
type WorkflowStep struct {
    ID          string                 `json:"id"`
    Title       string                 `json:"title"`
    Description string                 `json:"description"`
    Icon        fyne.Resource          `json:"-"`
    Enabled     bool                   `json:"enabled"`
    Completed   bool                   `json:"completed"`
    Order       int                    `json:"order"`
    Components  []ComponentDefinition  `json:"components"`
}

// ComponentDefinition defines UI component configuration
type ComponentDefinition struct {
    Type       ComponentType          `json:"type"`
    ID         string                 `json:"id"`
    Label      string                 `json:"label"`
    Properties map[string]interface{} `json:"properties"`
    Validators []ValidationRule       `json:"validators"`
    Layout     LayoutProperties       `json:"layout"`
}

// ComponentType defines supported UI component types
type ComponentType string

const (
    ComponentTypeButton     ComponentType = "button"
    ComponentTypeEntry      ComponentType = "entry"
    ComponentTypeSelect     ComponentType = "select"
    ComponentTypeCheck      ComponentType = "check"
    ComponentTypeSlider     ComponentType = "slider"
    ComponentTypeProgress   ComponentType = "progress"
    ComponentTypeLabel      ComponentType = "label"
    ComponentTypeContainer  ComponentType = "container"
    ComponentTypeCustom     ComponentType = "custom"
)

// ValidationRule defines UI validation rules
type ValidationRule struct {
    Type        ValidationType `json:"type"`
    Parameters  []interface{}  `json:"parameters"`
    Message     string         `json:"message"`
    OnValidate  func(interface{}) error `json:"-"`
}

// ValidationType defines validation rule types
type ValidationType string

const (
    ValidationTypeRequired ValidationType = "required"
    ValidationTypeMinLength ValidationType = "min_length"
    ValidationTypeMaxLength ValidationType = "max_length"
    ValidationTypePattern   ValidationType = "pattern"
    ValidationTypeRange     ValidationType = "range"
    ValidationTypeCustom    ValidationType = "custom"
)
```

### Component Factory

```go
// ComponentFactory creates UI components based on definitions
type ComponentFactory interface {
    // CreateComponent creates a component from definition
    CreateComponent(def ComponentDefinition, state *GUIState) (fyne.CanvasObject, error)
    
    // RegisterComponentType registers custom component types
    RegisterComponentType(componentType ComponentType, creator ComponentCreator)
    
    // GetSupportedTypes returns list of supported component types
    GetSupportedTypes() []ComponentType
}

// ComponentCreator function type for creating custom components
type ComponentCreator func(def ComponentDefinition, state *GUIState) (fyne.CanvasObject, error)

// Built-in component creators
func CreateButtonComponent(def ComponentDefinition, state *GUIState) (fyne.CanvasObject, error) {
    button := widget.NewButton(def.Label, nil)
    
    // Configure button properties
    if text, ok := def.Properties["text"].(string); ok {
        button.SetText(text)
    }
    
    if disabled, ok := def.Properties["disabled"].(bool); ok {
        if disabled {
            button.Disable()
        }
    }
    
    // Set up event handlers
    if onTap, ok := def.Properties["onTap"].(func()); ok {
        button.OnTapped = onTap
    }
    
    return button, nil
}

func CreateSelectComponent(def ComponentDefinition, state *GUIState) (fyne.CanvasObject, error) {
    options, ok := def.Properties["options"].([]string)
    if !ok {
        return nil, errors.New("select component requires 'options' property")
    }
    
    selectWidget := widget.NewSelect(options, nil)
    
    // Configure select properties
    if placeholder, ok := def.Properties["placeholder"].(string); ok {
        selectWidget.PlaceHolder = placeholder
    }
    
    if selected, ok := def.Properties["selected"].(string); ok {
        selectWidget.SetSelected(selected)
    }
    
    // Set up event handlers
    if onChanged, ok := def.Properties["onChanged"].(func(string)); ok {
        selectWidget.OnChanged = onChanged
    }
    
    return selectWidget, nil
}
```

### Theme System

```go
// ThemeManager handles application theming
type ThemeManager interface {
    // GetCurrentTheme returns the currently active theme
    GetCurrentTheme() fyne.Theme
    
    // SetTheme applies a new theme to the application
    SetTheme(theme fyne.Theme) error
    
    // GetAvailableThemes returns list of available themes
    GetAvailableThemes() []ThemeInfo
    
    // CreateCustomTheme creates a theme from configuration
    CreateCustomTheme(config ThemeConfig) (fyne.Theme, error)
    
    // RegisterTheme registers a custom theme
    RegisterTheme(theme fyne.Theme, info ThemeInfo) error
}

// ThemeInfo provides metadata about a theme
type ThemeInfo struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Author      string            `json:"author"`
    Version     string            `json:"version"`
    Variant     fyne.ThemeVariant `json:"variant"`
    Preview     fyne.Resource     `json:"-"`
    Colors      map[string]string `json:"colors"`
}

// ThemeConfig configures custom theme creation
type ThemeConfig struct {
    Name        string                           `json:"name"`
    BaseTheme   string                           `json:"base_theme"`
    Colors      map[fyne.ThemeColorName]color.Color `json:"-"`
    Fonts       map[fyne.TextStyle]fyne.Resource  `json:"-"`
    Sizes       map[fyne.ThemeSizeName]float32    `json:"-"`
    Icons       map[fyne.ThemeIconName]fyne.Resource `json:"-"`
}

// AkashicTheme implements VoidCat RDC's custom theme
type AkashicTheme struct {
    variant fyne.ThemeVariant
}

func NewAkashicTheme() fyne.Theme {
    return &AkashicTheme{variant: fyne.ThemeVariantDark}
}

func (t *AkashicTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
    // VoidCat RDC brand colors
    switch name {
    case theme.ColorNamePrimary:
        if variant == fyne.ThemeVariantLight {
            return color.RGBA{102, 51, 153, 255} // VoidCat Purple
        }
        return color.RGBA{153, 102, 204, 255} // Lighter Purple for dark theme
    case theme.ColorNameBackground:
        if variant == fyne.ThemeVariantLight {
            return color.RGBA{248, 248, 248, 255} // Light background
        }
        return color.RGBA{24, 24, 24, 255} // Dark background
    case theme.ColorNameForeground:
        if variant == fyne.ThemeVariantLight {
            return color.RGBA{33, 33, 33, 255} // Dark text
        }
        return color.RGBA{255, 255, 255, 255} // Light text
    default:
        return theme.DefaultTheme().Color(name, variant)
    }
}

func (t *AkashicTheme) Font(style fyne.TextStyle) fyne.Resource {
    // Use custom fonts if available
    if style.Bold {
        return resourceFontBold // Custom bold font
    }
    if style.Italic {
        return resourceFontItalic // Custom italic font
    }
    return resourceFontRegular // Custom regular font
}

func (t *AkashicTheme) Size(name fyne.ThemeSizeName) float32 {
    switch name {
    case theme.SizeNameText:
        return 14
    case theme.SizeNameCaptionText:
        return 12
    case theme.SizeNameHeadingText:
        return 24
    case theme.SizeNameSubHeadingText:
        return 18
    default:
        return theme.DefaultTheme().Size(name)
    }
}
```

## State Management API

### Centralized State

```go
// StateManager handles application state coordination
type StateManager interface {
    // State access
    GetState() *GUIState
    UpdateState(updater StateUpdater) error
    
    // State persistence
    SaveState() error
    LoadState() error
    
    // State observation
    Subscribe(observer StateObserver) Subscription
    Unsubscribe(subscription Subscription)
    
    // State validation
    ValidateState() error
    ResetState() error
}

// GUIState represents the complete application state
type GUIState struct {
    // Workflow state
    CurrentStep     int                    `json:"current_step"`
    WorkflowHistory []WorkflowHistoryEntry `json:"workflow_history"`
    
    // Processing state
    IsProcessing   bool                   `json:"is_processing"`
    Progress       float64                `json:"progress"`
    CurrentStage   string                 `json:"current_stage"`
    StatusMessage  string                 `json:"status_message"`
    ProcessingJob  *ProcessingJob         `json:"processing_job"`
    
    // Configuration state
    Options         *ScribeOptions         `json:"options"`
    PresetConfigs   map[string]*ScribeOptions `json:"preset_configs"`
    RecentFiles     []string               `json:"recent_files"`
    
    // UI state
    SelectedInputTab    int                    `json:"selected_input_tab"`
    WindowGeometry      WindowGeometry         `json:"window_geometry"`
    PanelVisibility     map[string]bool        `json:"panel_visibility"`
    ComponentState      map[string]interface{} `json:"component_state"`
    
    // User preferences
    Theme              string                 `json:"theme"`
    Language           string                 `json:"language"`
    AutoSave           bool                   `json:"auto_save"`
    Notifications      NotificationSettings   `json:"notifications"`
    
    // Internal state (not serialized)
    observers []StateObserver `json:"-"`
    mu        sync.RWMutex    `json:"-"`
}

// StateUpdater function type for atomic state updates
type StateUpdater func(*GUIState) error

// StateObserver interface for state change notifications
type StateObserver interface {
    OnStateChanged(state *GUIState, changes StateChanges)
}

// StateChanges describes what changed in the state
type StateChanges struct {
    Fields    []string               `json:"fields"`
    Previous  map[string]interface{} `json:"previous"`
    Current   map[string]interface{} `json:"current"`
    Timestamp time.Time              `json:"timestamp"`
}

// Subscription represents a state subscription
type Subscription interface {
    ID() string
    IsActive() bool
    Cancel()
}

// ProcessingJob represents an active processing job
type ProcessingJob struct {
    ID            string                 `json:"id"`
    StartTime     time.Time              `json:"start_time"`
    Options       *ScribeOptions         `json:"options"`
    Progress      float64                `json:"progress"`
    CurrentStage  string                 `json:"current_stage"`
    Stages        []ProcessingStage      `json:"stages"`
    Results       ProcessingResults      `json:"results"`
    Errors        []ProcessingError      `json:"errors"`
    Metadata      map[string]interface{} `json:"metadata"`
}

// ProcessingStage represents a stage in the processing pipeline
type ProcessingStage struct {
    Name        string    `json:"name"`
    Progress    float64   `json:"progress"`
    StartTime   time.Time `json:"start_time"`
    EndTime     time.Time `json:"end_time"`
    Status      StageStatus `json:"status"`
    Message     string    `json:"message"`
    ElapsedTime time.Duration `json:"elapsed_time"`
}

// StageStatus represents the status of a processing stage
type StageStatus string

const (
    StageStatusPending    StageStatus = "pending"
    StageStatusInProgress StageStatus = "in_progress"
    StageStatusCompleted  StageStatus = "completed"
    StageStatusFailed     StageStatus = "failed"
    StageStatusSkipped    StageStatus = "skipped"
)
```

### State Management Implementation

```go
// realStateManager implements StateManager
type realStateManager struct {
    state       *GUIState
    observers   []StateObserver
    persistence StatePersistence
    validator   StateValidator
    mu          sync.RWMutex
}

func NewStateManager(persistence StatePersistence) StateManager {
    return &realStateManager{
        state:       NewDefaultGUIState(),
        observers:   make([]StateObserver, 0),
        persistence: persistence,
        validator:   NewStateValidator(),
    }
}

func (sm *realStateManager) UpdateState(updater StateUpdater) error {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    // Create a copy for comparison
    previousState := sm.state.Clone()
    
    // Apply the update
    if err := updater(sm.state); err != nil {
        return fmt.Errorf("state update failed: %w", err)
    }
    
    // Validate the new state
    if err := sm.validator.ValidateState(sm.state); err != nil {
        // Revert to previous state
        sm.state = previousState
        return fmt.Errorf("state validation failed: %w", err)
    }
    
    // Determine what changed
    changes := sm.calculateChanges(previousState, sm.state)
    
    // Notify observers
    sm.notifyObservers(changes)
    
    // Auto-save if enabled
    if sm.state.AutoSave {
        go sm.persistence.SaveState(sm.state)
    }
    
    return nil
}

func (sm *realStateManager) Subscribe(observer StateObserver) Subscription {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    subscription := &stateSubscription{
        id:       generateSubscriptionID(),
        observer: observer,
        manager:  sm,
        active:   true,
    }
    
    sm.observers = append(sm.observers, observer)
    
    return subscription
}

// Helper methods for state management
func (sm *realStateManager) notifyObservers(changes StateChanges) {
    for _, observer := range sm.observers {
        go func(obs StateObserver) {
            defer func() {
                if r := recover(); r != nil {
                    log.Printf("State observer panicked: %v", r)
                }
            }()
            obs.OnStateChanged(sm.state, changes)
        }(observer)
    }
}
```

## Plugin Development API

### Plugin Interface

```go
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

// PluginContext provides access to application services
type PluginContext interface {
    // Core services
    GetEngine() ScribeEngine
    GetStateManager() StateManager
    GetEventBus() EventBus
    
    // UI services
    GetLayoutManager() LayoutManager
    GetThemeManager() ThemeManager
    
    // Utility services
    GetLogger() Logger
    GetConfig() Config
    GetResourceManager() ResourceManager
    
    // Plugin communication
    GetPluginManager() PluginManager
    SendMessage(targetPluginID string, message interface{}) error
    BroadcastMessage(message interface{}) error
}

// PluginCapability defines what a plugin can do
type PluginCapability string

const (
    CapabilityAudioProcessing   PluginCapability = "audio_processing"
    CapabilityVideoProcessing   PluginCapability = "video_processing"
    CapabilityTranscription     PluginCapability = "transcription"
    CapabilityTranslation       PluginCapability = "translation"
    CapabilityVoiceSynthesis    PluginCapability = "voice_synthesis"
    CapabilitySubtitleGeneration PluginCapability = "subtitle_generation"
    CapabilityUIExtension       PluginCapability = "ui_extension"
    CapabilityThemeProvider     PluginCapability = "theme_provider"
    CapabilityFileFormat        PluginCapability = "file_format"
    CapabilityCloudIntegration  PluginCapability = "cloud_integration"
    CapabilityAPIIntegration    PluginCapability = "api_integration"
)

// PluginDependency describes plugin dependencies
type PluginDependency struct {
    PluginID       string `json:"plugin_id"`
    Version        string `json:"version"`
    Required       bool   `json:"required"`
    Description    string `json:"description"`
}
```

### Plugin Manager

```go
// PluginManager handles plugin lifecycle and communication
type PluginManager interface {
    // Plugin registration
    RegisterPlugin(plugin Plugin) error
    UnregisterPlugin(pluginID string) error
    
    // Plugin lifecycle
    LoadPlugin(pluginPath string) (Plugin, error)
    UnloadPlugin(pluginID string) error
    EnablePlugin(pluginID string) error
    DisablePlugin(pluginID string) error
    
    // Plugin discovery
    GetLoadedPlugins() []Plugin
    GetAvailablePlugins() []PluginInfo
    GetPluginByID(pluginID string) (Plugin, bool)
    
    // Plugin communication
    SendMessage(senderID, targetID string, message interface{}) error
    BroadcastMessage(senderID string, message interface{}) error
    
    // Plugin capabilities
    GetPluginsByCapability(capability PluginCapability) []Plugin
    GetCapabilities() []PluginCapability
}

// PluginInfo provides metadata about available plugins
type PluginInfo struct {
    ID           string              `json:"id"`
    Name         string              `json:"name"`
    Version      string              `json:"version"`
    Description  string              `json:"description"`
    Author       string              `json:"author"`
    License      string              `json:"license"`
    Website      string              `json:"website"`
    Capabilities []PluginCapability  `json:"capabilities"`
    Dependencies []PluginDependency  `json:"dependencies"`
    FilePath     string              `json:"file_path"`
    Loaded       bool                `json:"loaded"`
    Enabled      bool                `json:"enabled"`
    LastUpdated  time.Time           `json:"last_updated"`
}
```

### Example Plugin Implementation

```go
// NoiseReductionPlugin demonstrates audio processing plugin
type NoiseReductionPlugin struct {
    id      string
    context PluginContext
    config  NoiseReductionConfig
}

func NewNoiseReductionPlugin() Plugin {
    return &NoiseReductionPlugin{
        id: "voidcat.noise_reduction",
    }
}

func (p *NoiseReductionPlugin) ID() string {
    return p.id
}

func (p *NoiseReductionPlugin) Name() string {
    return "Advanced Noise Reduction"
}

func (p *NoiseReductionPlugin) Version() string {
    return "1.0.0"
}

func (p *NoiseReductionPlugin) Description() string {
    return "Provides advanced noise reduction capabilities using AI-powered algorithms"
}

func (p *NoiseReductionPlugin) Author() string {
    return "VoidCat RDC"
}

func (p *NoiseReductionPlugin) GetCapabilities() []PluginCapability {
    return []PluginCapability{
        CapabilityAudioProcessing,
    }
}

func (p *NoiseReductionPlugin) Initialize(context PluginContext) error {
    p.context = context
    
    // Load configuration
    config, err := p.loadConfig()
    if err != nil {
        return fmt.Errorf("failed to load plugin configuration: %w", err)
    }
    p.config = config
    
    // Register audio processor
    engine := context.GetEngine()
    if processor, ok := engine.(AudioProcessor); ok {
        processor.RegisterAudioFilter(NewNoiseReductionFilter(p.config))
    }
    
    return nil
}

func (p *NoiseReductionPlugin) Activate() error {
    // Plugin-specific activation logic
    p.context.GetLogger().Info("Noise reduction plugin activated")
    return nil
}

func (p *NoiseReductionPlugin) Deactivate() error {
    // Plugin-specific deactivation logic
    p.context.GetLogger().Info("Noise reduction plugin deactivated")
    return nil
}

func (p *NoiseReductionPlugin) Shutdown() error {
    // Cleanup resources
    return nil
}

func (p *NoiseReductionPlugin) HealthCheck() error {
    // Verify plugin is functioning correctly
    if p.context == nil {
        return errors.New("plugin context not initialized")
    }
    return nil
}

// NoiseReductionFilter implements the actual audio processing
type NoiseReductionFilter struct {
    config NoiseReductionConfig
}

func (f *NoiseReductionFilter) Name() string {
    return "Noise Reduction"
}

func (f *NoiseReductionFilter) Description() string {
    return "Removes background noise from audio using advanced algorithms"
}

func (f *NoiseReductionFilter) Parameters() []VoiceParameter {
    return []VoiceParameter{
        {
            Name:        "intensity",
            Type:        "float",
            Default:     0.5,
            Min:         0.0,
            Max:         1.0,
            Description: "Noise reduction intensity",
        },
        {
            Name:        "preserve_speech",
            Type:        "bool",
            Default:     true,
            Description: "Preserve speech quality during noise reduction",
        },
    }
}

func (f *NoiseReductionFilter) Apply(audioData []byte, params map[string]interface{}) ([]byte, error) {
    intensity := f.getFloatParam(params, "intensity", 0.5)
    preserveSpeech := f.getBoolParam(params, "preserve_speech", true)
    
    // Implement noise reduction algorithm
    return f.processAudio(audioData, intensity, preserveSpeech)
}
```

---

## 📞 Support & Contact

- **GitHub Issues**: [Report bugs or request features](https://github.com/sorrowscry86/akashic-scribe/issues)
- **Discussions**: [Community discussions and Q&A](https://github.com/sorrowscry86/akashic-scribe/discussions)
- **Developer**: [@sorrowscry86](https://github.com/sorrowscry86)
- **Project**: VoidCat RDC - Akashic Scribe
- **Contact**: Wykeve Freeman (Sorrow Eternal) - SorrowsCry86@voidcat.org
- **Organization**: VoidCat RDC
- **Support Development**: CashApp $WykeveTF

---

**© 2024 VoidCat RDC, LLC. All rights reserved.**

*Technical Excellence Through Clean Architecture*