package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	assert := assert.New(t)

	config := DefaultConfig()

	// Test defaults are set
	assert.Equal("alloy", config.DefaultVoiceModel)
	assert.Equal(1.0, config.DefaultVoiceSpeed)
	assert.Equal(0.5, config.DefaultVoiceStability)
	assert.Equal("mp3", config.DefaultAudioFormat)
	assert.Equal("high", config.DefaultAudioQuality)
	assert.Equal(44100, config.DefaultSampleRate)
	assert.Equal(192, config.DefaultBitRate)
	assert.Equal(2, config.DefaultChannels)
	assert.Equal("Translation on Top", config.DefaultSubtitlePosition)
	assert.True(config.DefaultBilingualSubtitles)
	assert.Equal(2, config.MaxConcurrentJobs)
	assert.True(config.EnableCaching)

	// Validate default config
	err := config.Validate()
	assert.NoError(err, "Default config should be valid")
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		modify    func(*Config)
		shouldErr bool
		errMsg    string
	}{
		{
			name:      "Valid config",
			modify:    func(c *Config) {},
			shouldErr: false,
		},
		{
			name: "Invalid voice model",
			modify: func(c *Config) {
				c.DefaultVoiceModel = "invalid"
			},
			shouldErr: true,
			errMsg:    "invalid default voice model",
		},
		{
			name: "Voice speed too low",
			modify: func(c *Config) {
				c.DefaultVoiceSpeed = 0.1
			},
			shouldErr: true,
			errMsg:    "default voice speed must be between",
		},
		{
			name: "Voice speed too high",
			modify: func(c *Config) {
				c.DefaultVoiceSpeed = 5.0
			},
			shouldErr: true,
			errMsg:    "default voice speed must be between",
		},
		{
			name: "Invalid audio format",
			modify: func(c *Config) {
				c.DefaultAudioFormat = "xyz"
			},
			shouldErr: true,
			errMsg:    "invalid default audio format",
		},
		{
			name: "Invalid audio quality",
			modify: func(c *Config) {
				c.DefaultAudioQuality = "ultra"
			},
			shouldErr: true,
			errMsg:    "invalid default audio quality",
		},
		{
			name: "Invalid sample rate",
			modify: func(c *Config) {
				c.DefaultSampleRate = 12345
			},
			shouldErr: true,
			errMsg:    "invalid default sample rate",
		},
		{
			name: "Bit rate too low",
			modify: func(c *Config) {
				c.DefaultBitRate = 32
			},
			shouldErr: true,
			errMsg:    "default bit rate must be between",
		},
		{
			name: "Bit rate too high",
			modify: func(c *Config) {
				c.DefaultBitRate = 500
			},
			shouldErr: true,
			errMsg:    "default bit rate must be between",
		},
		{
			name: "Invalid channels",
			modify: func(c *Config) {
				c.DefaultChannels = 5
			},
			shouldErr: true,
			errMsg:    "default channels must be 1 (mono) or 2 (stereo)",
		},
		{
			name: "Invalid subtitle position",
			modify: func(c *Config) {
				c.DefaultSubtitlePosition = "Invalid Position"
			},
			shouldErr: true,
			errMsg:    "invalid default subtitle position",
		},
		{
			name: "Max concurrent jobs too low",
			modify: func(c *Config) {
				c.MaxConcurrentJobs = 0
			},
			shouldErr: true,
			errMsg:    "max concurrent jobs must be between",
		},
		{
			name: "Max concurrent jobs too high",
			modify: func(c *Config) {
				c.MaxConcurrentJobs = 20
			},
			shouldErr: true,
			errMsg:    "max concurrent jobs must be between",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			config := DefaultConfig()
			tt.modify(config)

			err := config.Validate()
			if tt.shouldErr {
				assert.Error(err)
				assert.Contains(err.Error(), tt.errMsg)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Create temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.json")

	// Create a config with custom values
	config := DefaultConfig()
	config.DefaultVoiceModel = "nova"
	config.DefaultVoiceSpeed = 1.5
	config.MaxConcurrentJobs = 4

	// Save config
	err := SaveConfig(config, configPath)
	require.NoError(err, "Should save config successfully")

	// Verify file exists
	_, err = os.Stat(configPath)
	assert.NoError(err, "Config file should exist")

	// Load config
	loadedConfig, err := LoadConfig(configPath)
	require.NoError(err, "Should load config successfully")

	// Verify loaded values match saved values
	assert.Equal(config.DefaultVoiceModel, loadedConfig.DefaultVoiceModel)
	assert.Equal(config.DefaultVoiceSpeed, loadedConfig.DefaultVoiceSpeed)
	assert.Equal(config.MaxConcurrentJobs, loadedConfig.MaxConcurrentJobs)
}

func TestLoadConfigNonExistent(t *testing.T) {
	assert := assert.New(t)

	// Load non-existent config should return default
	config, err := LoadConfig("/nonexistent/path/config.json")
	assert.NoError(err, "Should not error on non-existent file")
	assert.NotNil(config)

	// Should be default config
	defaultConfig := DefaultConfig()
	assert.Equal(defaultConfig.DefaultVoiceModel, config.DefaultVoiceModel)
	assert.Equal(defaultConfig.DefaultVoiceSpeed, config.DefaultVoiceSpeed)
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	assert := assert.New(t)

	// Create temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid-config.json")

	// Write invalid JSON
	err := os.WriteFile(configPath, []byte("{invalid json}"), 0644)
	assert.NoError(err)

	// Load should fail
	_, err = LoadConfig(configPath)
	assert.Error(err)
	assert.Contains(err.Error(), "failed to parse config file")
}

func TestLoadConfigInvalidValues(t *testing.T) {
	assert := assert.New(t)

	// Create temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid-values-config.json")

	// Write config with invalid values
	invalidJSON := `{
		"default_voice_model": "invalid_model",
		"default_voice_speed": 1.0
	}`
	err := os.WriteFile(configPath, []byte(invalidJSON), 0644)
	assert.NoError(err)

	// Load should fail validation
	_, err = LoadConfig(configPath)
	assert.Error(err)
	assert.Contains(err.Error(), "invalid configuration")
}

func TestSaveConfigInvalid(t *testing.T) {
	assert := assert.New(t)

	// Create invalid config
	config := DefaultConfig()
	config.DefaultVoiceModel = "invalid"

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.json")

	// Save should fail validation
	err := SaveConfig(config, configPath)
	assert.Error(err)
	assert.Contains(err.Error(), "invalid configuration")
}

func TestGetDefaultConfigPath(t *testing.T) {
	assert := assert.New(t)

	path, err := GetDefaultConfigPath()
	assert.NoError(err)
	assert.NotEmpty(path)
	assert.Contains(path, "akashic-scribe")
	assert.Contains(path, "config.json")
}

func TestApplyConfigToOptions(t *testing.T) {
	assert := assert.New(t)

	config := DefaultConfig()
	config.DefaultVoiceModel = "echo"
	config.DefaultVoiceSpeed = 1.2
	config.DefaultAudioFormat = "wav"

	// Test with empty options
	opts := &ScribeOptions{}
	err := ApplyConfigToOptions(config, opts)
	assert.NoError(err)

	assert.Equal("echo", opts.VoiceModel)
	assert.Equal(1.2, opts.VoiceSpeed)
	assert.Equal("wav", opts.AudioFormat)

	// Test with partially filled options (should not override)
	opts2 := &ScribeOptions{
		VoiceModel: "nova", // Already set
	}
	err = ApplyConfigToOptions(config, opts2)
	assert.NoError(err)

	assert.Equal("nova", opts2.VoiceModel) // Should keep original
	assert.Equal(1.2, opts2.VoiceSpeed)    // Should apply from config
}

func TestApplyConfigToOptionsNil(t *testing.T) {
	assert := assert.New(t)

	config := DefaultConfig()

	// Nil options
	err := ApplyConfigToOptions(config, nil)
	assert.Error(err)
	assert.Contains(err.Error(), "options cannot be nil")

	// Nil config
	opts := &ScribeOptions{}
	err = ApplyConfigToOptions(nil, opts)
	assert.Error(err)
	assert.Contains(err.Error(), "config cannot be nil")
}

func TestConfigEdgeCases(t *testing.T) {
	assert := assert.New(t)

	// Test boundary values
	config := DefaultConfig()

	// Minimum valid voice speed
	config.DefaultVoiceSpeed = 0.25
	assert.NoError(config.Validate())

	// Maximum valid voice speed
	config.DefaultVoiceSpeed = 4.0
	assert.NoError(config.Validate())

	// Minimum valid bit rate
	config.DefaultBitRate = 64
	assert.NoError(config.Validate())

	// Maximum valid bit rate
	config.DefaultBitRate = 320
	assert.NoError(config.Validate())

	// Mono audio
	config.DefaultChannels = 1
	assert.NoError(config.Validate())

	// Stereo audio
	config.DefaultChannels = 2
	assert.NoError(config.Validate())
}
