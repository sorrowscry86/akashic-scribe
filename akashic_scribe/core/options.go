package core

import "fmt"

// ScribeOptions represents all user configuration for the transcription and translation process.
//
// === BACKEND INTEGRATION PHASE ===
// This struct is the single source of truth for all user-selected options.
// It is passed directly to the backend engine for processing.
//
// TODOs for integration:
//   - Ensure all fields are populated before backend handoff.
//   - Add any new options required by backend features.
//   - Keep this struct in sync with backend expectations.
type ScribeOptions struct {
	// Input sources
	InputFile string // Full path to the local video file.
	InputURL  string // URL of the video to be downloaded.

	// Language configuration
	OriginLanguage string // e.g., "en-US"
	TargetLanguage string // e.g., "ja-JP"

	// Subtitle options
	CreateSubtitles    bool   // Whether to generate subtitles
	BilingualSubtitles bool   // Whether to include both languages in subtitles
	SubtitlePosition   string // "top" or "bottom" - position of translation in bilingual mode
	SubtitleFormat     string // "srt" or "vtt" (default "srt")

	// Dubbing options
	CreateDubbing   bool    // Whether to generate dubbed audio
	VoiceModel      string  // e.g., "alloy", "echo", "fable", "onyx", "nova", "shimmer"
	UseCustomVoice  bool    // Whether to use a custom voice model
	CustomVoicePath string  // Path to custom voice file for voice cloning
	VoiceSpeed      float64 // Speech speed (0.25 to 4.0, default 1.0)
	VoicePitch      float64 // Voice pitch adjustment (-20 to 20 semitones, default 0)
	VoiceStability  float64 // Voice stability (0.0 to 1.0, higher = more stable, default 0.5)
	AudioFormat     string  // Output format: "mp3", "wav", "flac", "aac", "ogg" (default "mp3")
	AudioQuality    string  // Quality: "low", "medium", "high", "lossless" (default "high")
	AudioSampleRate int     // Sample rate in Hz (8000, 16000, 22050, 44100, 48000, default 44100)
	AudioBitRate    int     // Bit rate in kbps (64, 128, 192, 256, 320, default 192)
	NormalizeAudio  bool    // Whether to normalize audio levels (default true)
	RemoveSilence   bool    // Whether to remove long silences (default false)
	AudioChannels   int     // Number of audio channels: 1 (mono) or 2 (stereo), default 2

	// Output configuration
	OutputDir string // Optional. If empty, defaults to the input file directory or a sensible default.
}

// String returns a formatted representation of the ScribeOptions for debugging
func (s ScribeOptions) String() string {
	result := "=== SCRIBE OPTIONS ===\n"
	result += "Input Configuration:\n"
	if s.InputFile != "" {
		result += "  File: " + s.InputFile + "\n"
	}
	if s.InputURL != "" {
		result += "  URL: " + s.InputURL + "\n"
	}

	result += "Language Configuration:\n"
	result += "  Origin: " + s.OriginLanguage + "\n"
	result += "  Target: " + s.TargetLanguage + "\n"

	result += "Subtitle Options:\n"
	if s.CreateSubtitles {
		result += "  Create Subtitles: Enabled\n"
		if s.SubtitleFormat != "" {
			result += "  Format: " + s.SubtitleFormat + "\n"
		}
		if s.BilingualSubtitles {
			result += "  Bilingual: Enabled (" + s.SubtitlePosition + ")\n"
		} else {
			result += "  Bilingual: Disabled\n"
		}
	} else {
		result += "  Create Subtitles: Disabled\n"
	}

	result += "Dubbing Options:\n"
	if s.CreateDubbing {
		result += "  Create Dubbing: Enabled\n"
		if s.UseCustomVoice && s.CustomVoicePath != "" {
			result += "  Custom Voice: " + s.CustomVoicePath + "\n"
		} else if s.VoiceModel != "" {
			result += "  Voice Model: " + s.VoiceModel + "\n"
		}
		if s.VoiceSpeed > 0 {
			result += fmt.Sprintf("  Voice Speed: %.2fx\n", s.VoiceSpeed)
		}
		if s.VoicePitch != 0 {
			result += fmt.Sprintf("  Voice Pitch: %+.1f semitones\n", s.VoicePitch)
		}
		if s.VoiceStability > 0 {
			result += fmt.Sprintf("  Voice Stability: %.2f\n", s.VoiceStability)
		}
		if s.AudioFormat != "" {
			result += "  Audio Format: " + s.AudioFormat + "\n"
		}
		if s.AudioQuality != "" {
			result += "  Audio Quality: " + s.AudioQuality + "\n"
		}
		if s.AudioSampleRate > 0 {
			result += fmt.Sprintf("  Sample Rate: %d Hz\n", s.AudioSampleRate)
		}
		if s.AudioBitRate > 0 {
			result += fmt.Sprintf("  Bit Rate: %d kbps\n", s.AudioBitRate)
		}
		if s.AudioChannels > 0 {
			result += fmt.Sprintf("  Channels: %d\n", s.AudioChannels)
		}
		if s.NormalizeAudio {
			result += "  Normalize Audio: Enabled\n"
		}
		if s.RemoveSilence {
			result += "  Remove Silence: Enabled\n"
		}
	} else {
		result += "  Create Dubbing: Disabled\n"
	}

	if s.OutputDir != "" {
		result += "Output Directory:\n"
		result += "  Path: " + s.OutputDir + "\n"
	}

	return result
}
