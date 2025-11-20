package core

import "context"

// ScribeEngine defines the interface for the core transcription and translation engine.
//
// Phase IV: The Animus (Backend Integration)
// Adds StartProcessing for full workflow and progress reporting.
// Phase 4.2: Added context.Context support for graceful cancellation.
// Phase 5: Added ProcessWithContext for batch processing support.
type ScribeEngine interface {
	// Transcribe takes a video source (local path or URL) and returns the transcription.
	Transcribe(videoSource string) (string, error)

	// Translate takes a text and a target language, and returns the translation.
	Translate(text string, targetLanguage string) (string, error)

	// StartProcessing runs the full pipeline and reports progress.
	// The context can be used to cancel the operation at any time.
	StartProcessing(ctx context.Context, options ScribeOptions, progress chan<- ProgressUpdate) error

	// ProcessWithContext runs the full pipeline with context and returns a result.
	// This is used by the batch processor to get structured results.
	ProcessWithContext(ctx context.Context, options ScribeOptions, progress chan<- ProgressUpdate) (*ScribeResult, error)
}

// ProgressUpdate is sent over the progress channel to report backend status.
type ProgressUpdate struct {
	Percentage float64 // 0.0 to 1.0
	Message    string  // Status message
}

// ScribeResult contains the results of a completed transcription/translation job.
type ScribeResult struct {
	Transcription string `json:"transcription"`
	Translation   string `json:"translation"`
	DubbedAudio   string `json:"dubbed_audio,omitempty"`
	SubtitlesFile string `json:"subtitles_file,omitempty"`
	OutputDir     string `json:"output_dir"`
}
