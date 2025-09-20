package core

// ScribeEngine defines the interface for the core transcription and translation engine.
//
// Phase IV: The Animus (Backend Integration)
// Adds StartProcessing for full workflow and progress reporting.
type ScribeEngine interface {
	// Transcribe takes a video source (local path or URL) and returns the transcription.
	Transcribe(videoSource string) (string, error)

	// Translate takes a text and a target language, and returns the translation.
	Translate(text string, targetLanguage string) (string, error)

	// StartProcessing runs the full pipeline and reports progress.
	StartProcessing(options ScribeOptions, progress chan<- ProgressUpdate) error
}

// ProgressUpdate is sent over the progress channel to report backend status.
type ProgressUpdate struct {
	Percentage float64 // 0.0 to 1.0
	Message    string  // Status message
}
