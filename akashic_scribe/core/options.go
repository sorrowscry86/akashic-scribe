package core

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
	SubtitlePosition   string // "Translation on Top" or "Translation on Bottom"

	// Dubbing options
	CreateDubbing  bool   // Whether to generate dubbed audio
	VoiceModel     string // e.g., "Kore"
	UseCustomVoice bool   // Whether to use a custom voice model

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
		result += "  Voice Model: " + s.VoiceModel + "\n"
		if s.UseCustomVoice {
			result += "  Custom Voice: Enabled\n"
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
