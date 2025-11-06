package audio_effects

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"akashic_scribe/core"
)

// AudioEffectsPlugin provides advanced audio processing capabilities
type AudioEffectsPlugin struct {
	*core.BasePlugin
	supportedFormats []string
}

// NewAudioEffectsPlugin creates a new audio effects plugin
func NewAudioEffectsPlugin() core.Plugin {
	base := core.NewBasePlugin(
		"voidcat.audio_effects",
		"Audio Effects Suite",
		"1.0.0",
		"Advanced audio processing with noise reduction, normalization, equalization, and more",
		"VoidCat RDC",
	)

	plugin := &AudioEffectsPlugin{
		BasePlugin: base,
		supportedFormats: []string{
			"mp3", "wav", "flac", "aac", "ogg", "m4a", "wma",
		},
	}

	return plugin
}

// GetCapabilities returns the plugin's capabilities
func (p *AudioEffectsPlugin) GetCapabilities() []core.PluginCapability {
	return []core.PluginCapability{
		core.CapabilityAudioProcessing,
	}
}

// ProcessAudio processes audio with the specified effects
// Supported options:
//   - "noise_reduction": bool - Enable noise reduction (default: false)
//   - "normalization": bool - Enable audio normalization (default: false)
//   - "equalization": map[string]float64 - EQ settings (frequencies to gain in dB)
//   - "reverb": float64 - Reverb amount 0.0-1.0 (default: 0.0)
//   - "bass_boost": float64 - Bass boost in dB (default: 0.0)
//   - "treble_boost": float64 - Treble boost in dB (default: 0.0)
//   - "compression": bool - Enable dynamic range compression (default: false)
//   - "fade_in": float64 - Fade in duration in seconds (default: 0.0)
//   - "fade_out": float64 - Fade out duration in seconds (default: 0.0)
//   - "volume": float64 - Volume adjustment factor, 1.0 = no change (default: 1.0)
func (p *AudioEffectsPlugin) ProcessAudio(inputPath, outputPath string, options map[string]interface{}) error {
	ctx := p.GetContext()

	// Validate input file
	if _, err := os.Stat(inputPath); err != nil {
		return fmt.Errorf("input file not found: %w", err)
	}

	// Check if ffmpeg is available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found - required for audio processing: %w", err)
	}

	ctx.LogInfo(fmt.Sprintf("Processing audio: %s -> %s", inputPath, outputPath))

	// Build ffmpeg filter chain
	filters := p.buildFilterChain(options)

	// Build ffmpeg command
	args := []string{"-i", inputPath}

	if len(filters) > 0 {
		args = append(args, "-af", strings.Join(filters, ","))
	}

	// Output settings
	args = append(args, "-y", outputPath)

	// Execute ffmpeg
	cmd := exec.Command("ffmpeg", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	ctx.LogInfo(fmt.Sprintf("Executing: ffmpeg %s", strings.Join(args, " ")))

	if err := cmd.Run(); err != nil {
		ctx.LogError(fmt.Sprintf("FFmpeg error: %s", stderr.String()))
		return fmt.Errorf("audio processing failed: %w", err)
	}

	// Verify output
	if _, err := os.Stat(outputPath); err != nil {
		return fmt.Errorf("output file not created: %w", err)
	}

	ctx.LogInfo("Audio processing completed successfully")
	return nil
}

// buildFilterChain constructs the ffmpeg audio filter chain based on options
func (p *AudioEffectsPlugin) buildFilterChain(options map[string]interface{}) []string {
	var filters []string

	// Fade in
	if fadeIn, ok := options["fade_in"].(float64); ok && fadeIn > 0 {
		filters = append(filters, fmt.Sprintf("afade=t=in:st=0:d=%.2f", fadeIn))
	}

	// Noise reduction using highpass and lowpass filters
	if noiseReduction, ok := options["noise_reduction"].(bool); ok && noiseReduction {
		// Remove low frequency rumble (below 80 Hz)
		filters = append(filters, "highpass=f=80")
		// Remove high frequency hiss (above 12000 Hz)
		filters = append(filters, "lowpass=f=12000")
		// Apply gentle noise gate
		filters = append(filters, "afftdn=nf=-25")
	}

	// Bass boost
	if bassBoost, ok := options["bass_boost"].(float64); ok && bassBoost != 0 {
		filters = append(filters, fmt.Sprintf("bass=g=%.1f:f=100:w=100", bassBoost))
	}

	// Treble boost
	if trebleBoost, ok := options["treble_boost"].(float64); ok && trebleBoost != 0 {
		filters = append(filters, fmt.Sprintf("treble=g=%.1f:f=8000:w=2000", trebleBoost))
	}

	// Equalization
	if eq, ok := options["equalization"].(map[string]interface{}); ok {
		for freqStr, gainVal := range eq {
			freq, err := strconv.ParseFloat(freqStr, 64)
			if err != nil {
				continue
			}
			gain, ok := gainVal.(float64)
			if !ok {
				continue
			}
			filters = append(filters, fmt.Sprintf("equalizer=f=%.0f:t=h:w=200:g=%.1f", freq, gain))
		}
	}

	// Dynamic range compression
	if compression, ok := options["compression"].(bool); ok && compression {
		// Gentle compression with 3:1 ratio
		filters = append(filters, "acompressor=threshold=-20dB:ratio=3:attack=5:release=50")
	}

	// Reverb
	if reverb, ok := options["reverb"].(float64); ok && reverb > 0 {
		if reverb > 1.0 {
			reverb = 1.0
		}
		// Convert reverb amount to delay and decay
		decay := reverb * 0.5
		filters = append(filters, fmt.Sprintf("aecho=0.8:0.88:%.0f:%.2f", 60.0*reverb, decay))
	}

	// Volume adjustment
	if volume, ok := options["volume"].(float64); ok && volume != 1.0 {
		filters = append(filters, fmt.Sprintf("volume=%.2f", volume))
	}

	// Normalization (should be near the end)
	if normalization, ok := options["normalization"].(bool); ok && normalization {
		filters = append(filters, "loudnorm=I=-16:TP=-1.5:LRA=11")
	}

	// Fade out
	if fadeOut, ok := options["fade_out"].(float64); ok && fadeOut > 0 {
		// Note: For fade out, we need to know the duration
		// This is a simplified version that assumes duration will be calculated
		filters = append(filters, fmt.Sprintf("afade=t=out:d=%.2f", fadeOut))
	}

	return filters
}

// GetSupportedFormats returns the audio formats supported by this plugin
func (p *AudioEffectsPlugin) GetSupportedFormats() []string {
	return p.supportedFormats
}

// ApplyNoiseReduction applies noise reduction to an audio file
func (p *AudioEffectsPlugin) ApplyNoiseReduction(inputPath, outputPath string) error {
	return p.ProcessAudio(inputPath, outputPath, map[string]interface{}{
		"noise_reduction": true,
		"normalization":   true,
	})
}

// ApplyNormalization applies audio normalization
func (p *AudioEffectsPlugin) ApplyNormalization(inputPath, outputPath string) error {
	return p.ProcessAudio(inputPath, outputPath, map[string]interface{}{
		"normalization": true,
	})
}

// ApplyVocalEnhancement applies settings optimized for vocal clarity
func (p *AudioEffectsPlugin) ApplyVocalEnhancement(inputPath, outputPath string) error {
	return p.ProcessAudio(inputPath, outputPath, map[string]interface{}{
		"noise_reduction": true,
		"normalization":   true,
		"compression":     true,
		"equalization": map[string]interface{}{
			"200":  -2.0, // Reduce mud
			"800":  1.0,  // Slight boost for warmth
			"3000": 3.0,  // Boost presence
			"8000": 2.0,  // Boost clarity
		},
	})
}

// ApplyPodcastPreset applies settings optimized for podcasts
func (p *AudioEffectsPlugin) ApplyPodcastPreset(inputPath, outputPath string) error {
	return p.ProcessAudio(inputPath, outputPath, map[string]interface{}{
		"noise_reduction": true,
		"normalization":   true,
		"compression":     true,
		"bass_boost":      3.0,
		"equalization": map[string]interface{}{
			"100":  -3.0, // Reduce rumble
			"3000": 2.0,  // Boost speech clarity
		},
	})
}

// ApplyMusicEnhancement applies settings optimized for music
func (p *AudioEffectsPlugin) ApplyMusicEnhancement(inputPath, outputPath string) error {
	return p.ProcessAudio(inputPath, outputPath, map[string]interface{}{
		"normalization": true,
		"bass_boost":    2.0,
		"treble_boost":  1.5,
	})
}

// ConvertFormat converts audio to a different format
func (p *AudioEffectsPlugin) ConvertFormat(inputPath, outputFormat string) (string, error) {
	ctx := p.GetContext()

	// Generate output path
	outputPath := filepath.Join(
		ctx.GetPluginDataDir(),
		fmt.Sprintf("%s.%s", filepath.Base(inputPath[:len(inputPath)-len(filepath.Ext(inputPath))]), outputFormat),
	)

	// Simple format conversion without effects
	return outputPath, p.ProcessAudio(inputPath, outputPath, map[string]interface{}{})
}

// GetEffectPresets returns available effect presets
func (p *AudioEffectsPlugin) GetEffectPresets() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"clean_vocal": {
			"noise_reduction": true,
			"normalization":   true,
			"compression":     true,
			"equalization": map[string]interface{}{
				"200":  -2.0,
				"3000": 3.0,
				"8000": 2.0,
			},
		},
		"podcast": {
			"noise_reduction": true,
			"normalization":   true,
			"compression":     true,
			"bass_boost":      3.0,
		},
		"music": {
			"normalization": true,
			"bass_boost":    2.0,
			"treble_boost":  1.5,
		},
		"radio_voice": {
			"noise_reduction": true,
			"normalization":   true,
			"compression":     true,
			"bass_boost":      5.0,
			"equalization": map[string]interface{}{
				"100":  -5.0,
				"3000": 4.0,
			},
		},
	}
}

// Ensure AudioEffectsPlugin implements the AudioProcessor interface
var _ core.AudioProcessor = (*AudioEffectsPlugin)(nil)
