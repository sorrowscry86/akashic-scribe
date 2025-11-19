package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration.
type Config struct {
	// Default audio settings
	DefaultVoiceModel    string  `json:"default_voice_model"`
	DefaultVoiceSpeed    float64 `json:"default_voice_speed"`
	DefaultVoiceStability float64 `json:"default_voice_stability"`
	DefaultAudioFormat   string  `json:"default_audio_format"`
	DefaultAudioQuality  string  `json:"default_audio_quality"`
	DefaultSampleRate    int     `json:"default_sample_rate"`
	DefaultBitRate       int     `json:"default_bit_rate"`
	DefaultChannels      int     `json:"default_channels"`

	// Default subtitle settings
	DefaultSubtitlePosition string `json:"default_subtitle_position"`
	DefaultBilingualSubtitles bool `json:"default_bilingual_subtitles"`

	// Performance settings
	MaxConcurrentJobs int  `json:"max_concurrent_jobs"`
	EnableCaching     bool `json:"enable_caching"`

	// Output settings
	DefaultOutputDir string `json:"default_output_dir"`
}

// DefaultConfig returns a configuration with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		// Audio defaults
		DefaultVoiceModel:     "alloy",
		DefaultVoiceSpeed:     1.0,
		DefaultVoiceStability: 0.5,
		DefaultAudioFormat:    "mp3",
		DefaultAudioQuality:   "high",
		DefaultSampleRate:     44100,
		DefaultBitRate:        192,
		DefaultChannels:       2,

		// Subtitle defaults
		DefaultSubtitlePosition:   "Translation on Top",
		DefaultBilingualSubtitles: true,

		// Performance defaults
		MaxConcurrentJobs: 2,
		EnableCaching:     true,

		// Output defaults
		DefaultOutputDir: "",
	}
}

// Validate checks if the configuration values are valid.
func (c *Config) Validate() error {
	// Validate voice model
	validModels := map[string]bool{
		"alloy": true, "echo": true, "fable": true,
		"onyx": true, "nova": true, "shimmer": true,
	}
	if !validModels[c.DefaultVoiceModel] {
		return fmt.Errorf("invalid default voice model: %s", c.DefaultVoiceModel)
	}

	// Validate voice speed
	if c.DefaultVoiceSpeed < 0.25 || c.DefaultVoiceSpeed > 4.0 {
		return fmt.Errorf("default voice speed must be between 0.25 and 4.0, got %.2f", c.DefaultVoiceSpeed)
	}

	// Validate voice stability
	if c.DefaultVoiceStability < 0 || c.DefaultVoiceStability > 1.0 {
		return fmt.Errorf("default voice stability must be between 0.0 and 1.0, got %.2f", c.DefaultVoiceStability)
	}

	// Validate audio format
	validFormats := map[string]bool{"mp3": true, "wav": true, "flac": true, "aac": true, "ogg": true}
	if !validFormats[c.DefaultAudioFormat] {
		return fmt.Errorf("invalid default audio format: %s", c.DefaultAudioFormat)
	}

	// Validate audio quality
	validQualities := map[string]bool{"low": true, "medium": true, "high": true, "lossless": true}
	if !validQualities[c.DefaultAudioQuality] {
		return fmt.Errorf("invalid default audio quality: %s", c.DefaultAudioQuality)
	}

	// Validate sample rate
	validSampleRates := map[int]bool{8000: true, 16000: true, 22050: true, 44100: true, 48000: true}
	if !validSampleRates[c.DefaultSampleRate] {
		return fmt.Errorf("invalid default sample rate: %d Hz", c.DefaultSampleRate)
	}

	// Validate bit rate
	if c.DefaultBitRate < 64 || c.DefaultBitRate > 320 {
		return fmt.Errorf("default bit rate must be between 64 and 320 kbps, got %d", c.DefaultBitRate)
	}

	// Validate channels
	if c.DefaultChannels != 1 && c.DefaultChannels != 2 {
		return fmt.Errorf("default channels must be 1 (mono) or 2 (stereo), got %d", c.DefaultChannels)
	}

	// Validate subtitle position
	validSubPositions := map[string]bool{
		"Translation on Top":    true,
		"Translation on Bottom": true,
	}
	if !validSubPositions[c.DefaultSubtitlePosition] {
		return fmt.Errorf("invalid default subtitle position: %s", c.DefaultSubtitlePosition)
	}

	// Validate max concurrent jobs
	if c.MaxConcurrentJobs < 1 || c.MaxConcurrentJobs > 10 {
		return fmt.Errorf("max concurrent jobs must be between 1 and 10, got %d", c.MaxConcurrentJobs)
	}

	return nil
}

// LoadConfig loads configuration from a JSON file.
// If the file doesn't exist, it returns the default configuration.
func LoadConfig(path string) (*Config, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return DefaultConfig(), nil
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to a JSON file.
func SaveConfig(config *Config, path string) error {
	// Validate before saving
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetDefaultConfigPath returns the default configuration file path.
// On Unix-like systems: ~/.config/akashic-scribe/config.json
// On Windows: %APPDATA%\akashic-scribe\config.json
func GetDefaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Determine config directory based on OS
	var configDir string
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		configDir = os.Getenv("XDG_CONFIG_HOME")
	} else if os.Getenv("APPDATA") != "" {
		configDir = os.Getenv("APPDATA")
	} else {
		configDir = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(configDir, "akashic-scribe", "config.json"), nil
}

// ApplyConfigToOptions applies configuration defaults to ScribeOptions.
func ApplyConfigToOptions(config *Config, opts *ScribeOptions) error {
	if config == nil {
		return errors.New("config cannot be nil")
	}

	if opts == nil {
		return errors.New("options cannot be nil")
	}

	// Apply audio defaults if not set
	if opts.VoiceModel == "" {
		opts.VoiceModel = config.DefaultVoiceModel
	}
	if opts.VoiceSpeed == 0 {
		opts.VoiceSpeed = config.DefaultVoiceSpeed
	}
	if opts.VoiceStability == 0 {
		opts.VoiceStability = config.DefaultVoiceStability
	}
	if opts.AudioFormat == "" {
		opts.AudioFormat = config.DefaultAudioFormat
	}
	if opts.AudioQuality == "" {
		opts.AudioQuality = config.DefaultAudioQuality
	}
	if opts.AudioSampleRate == 0 {
		opts.AudioSampleRate = config.DefaultSampleRate
	}
	if opts.AudioBitRate == 0 {
		opts.AudioBitRate = config.DefaultBitRate
	}
	if opts.AudioChannels == 0 {
		opts.AudioChannels = config.DefaultChannels
	}

	// Apply subtitle defaults
	if opts.SubtitlePosition == "" {
		opts.SubtitlePosition = config.DefaultSubtitlePosition
	}
	if !opts.BilingualSubtitles {
		opts.BilingualSubtitles = config.DefaultBilingualSubtitles
	}

	// Apply output directory default
	if opts.OutputDir == "" && config.DefaultOutputDir != "" {
		opts.OutputDir = config.DefaultOutputDir
	}

	return nil
}
