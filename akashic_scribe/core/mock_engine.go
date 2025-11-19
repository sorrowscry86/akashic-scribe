package core

import (
	"context"
	"fmt"
	"time"
)

// MockScribeEngine is a mock implementation of the ScribeEngine for testing and UI development.
type MockScribeEngine struct{}

// NewMockScribeEngine creates a new instance of the mock engine.
func NewMockScribeEngine() *MockScribeEngine {
	return &MockScribeEngine{}
}

// Transcribe simulates the transcription process.
func (m *MockScribeEngine) Transcribe(videoSource string) (string, error) {
	fmt.Printf("Mock Transcribe called with: %s\n", videoSource)
	time.Sleep(200 * time.Millisecond) // Simulate work
	transcription := "This is a mock transcription of the video."
	fmt.Println("Mock transcription finished.")
	return transcription, nil
}

// Translate simulates the translation process.
func (m *MockScribeEngine) Translate(text string, targetLanguage string) (string, error) {
	fmt.Printf("Mock Translate called for language '%s' with text: %s\n", targetLanguage, text)
	time.Sleep(100 * time.Millisecond) // Simulate work
	translation := fmt.Sprintf("This is a mock translation to '%s'.", targetLanguage)
	fmt.Println("Mock translation finished.")
	return translation, nil
}

// StartProcessing simulates the full processing pipeline with progress reporting.
// The operation can be cancelled via the provided context.
func (m *MockScribeEngine) StartProcessing(ctx context.Context, options ScribeOptions, progress chan<- ProgressUpdate) error {
	fmt.Printf("Mock StartProcessing called with options: %+v\n", options)

	// Check for cancellation before starting
	select {
	case <-ctx.Done():
		return fmt.Errorf("operation cancelled before start: %w", ctx.Err())
	default:
	}

	// Simulate the full processing pipeline with progress updates
	progress <- ProgressUpdate{0.0, "Starting processing..."}

	// Simulate work with cancellation check
	select {
	case <-time.After(50 * time.Millisecond):
	case <-ctx.Done():
		return fmt.Errorf("operation cancelled: %w", ctx.Err())
	}

	progress <- ProgressUpdate{0.2, "Transcribing audio..."}

	select {
	case <-time.After(100 * time.Millisecond):
	case <-ctx.Done():
		return fmt.Errorf("operation cancelled: %w", ctx.Err())
	}

	transcription, err := m.Transcribe(options.InputFile)
	if err != nil {
		return err
	}

	// Check for cancellation before translation
	select {
	case <-ctx.Done():
		return fmt.Errorf("operation cancelled: %w", ctx.Err())
	default:
	}

	progress <- ProgressUpdate{0.6, "Translating text..."}

	select {
	case <-time.After(100 * time.Millisecond):
	case <-ctx.Done():
		return fmt.Errorf("operation cancelled: %w", ctx.Err())
	}

	translation, err := m.Translate(transcription, options.TargetLanguage)
	if err != nil {
		return err
	}

	if options.CreateDubbing {
		// Check for cancellation before dubbing
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation cancelled: %w", ctx.Err())
		default:
		}

		progress <- ProgressUpdate{0.8, "Creating dubbed audio..."}

		select {
		case <-time.After(50 * time.Millisecond):
		case <-ctx.Done():
			return fmt.Errorf("operation cancelled: %w", ctx.Err())
		}
	}

	if options.CreateSubtitles {
		// Check for cancellation before subtitles
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation cancelled: %w", ctx.Err())
		default:
		}

		progress <- ProgressUpdate{0.9, "Generating subtitles..."}

		select {
		case <-time.After(50 * time.Millisecond):
		case <-ctx.Done():
			return fmt.Errorf("operation cancelled: %w", ctx.Err())
		}
	}

	// Final result with JSON structure
	result := fmt.Sprintf(`{
  "Transcription": "%s",
  "Translation": "%s"
}`, transcription, translation)

	progress <- ProgressUpdate{1.0, "Scribing complete.\n" + result}

	fmt.Println("Mock processing completed successfully.")
	return nil
}

// ProcessWithContext simulates the full processing pipeline and returns a structured result.
// This is used by the batch processor.
func (m *MockScribeEngine) ProcessWithContext(ctx context.Context, options ScribeOptions, progress chan<- ProgressUpdate) (*ScribeResult, error) {
	fmt.Printf("Mock ProcessWithContext called with options: %+v\n", options)

	// Check for cancellation before starting
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("operation cancelled before start: %w", ctx.Err())
	default:
	}

	// Simulate the full processing pipeline with progress updates
	progress <- ProgressUpdate{0.0, "Starting processing..."}

	// Simulate transcription
	select {
	case <-time.After(50 * time.Millisecond):
	case <-ctx.Done():
		return nil, fmt.Errorf("operation cancelled: %w", ctx.Err())
	}

	progress <- ProgressUpdate{0.3, "Transcribing audio..."}
	transcription, err := m.Transcribe(options.InputFile)
	if err != nil {
		return nil, err
	}

	// Check for cancellation
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("operation cancelled: %w", ctx.Err())
	default:
	}

	// Simulate translation
	progress <- ProgressUpdate{0.6, "Translating text..."}
	select {
	case <-time.After(50 * time.Millisecond):
	case <-ctx.Done():
		return nil, fmt.Errorf("operation cancelled: %w", ctx.Err())
	}

	translation, err := m.Translate(transcription, options.TargetLanguage)
	if err != nil {
		return nil, err
	}

	result := &ScribeResult{
		Transcription: transcription,
		Translation:   translation,
		OutputDir:     "/mock/output/dir",
	}

	if options.CreateDubbing {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("operation cancelled: %w", ctx.Err())
		default:
		}

		progress <- ProgressUpdate{0.8, "Creating dubbed audio..."}
		select {
		case <-time.After(30 * time.Millisecond):
		case <-ctx.Done():
			return nil, fmt.Errorf("operation cancelled: %w", ctx.Err())
		}

		result.DubbedAudio = "/mock/output/dir/dubbed_audio.mp3"
	}

	if options.CreateSubtitles {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("operation cancelled: %w", ctx.Err())
		default:
		}

		progress <- ProgressUpdate{0.9, "Generating subtitles..."}
		select {
		case <-time.After(30 * time.Millisecond):
		case <-ctx.Done():
			return nil, fmt.Errorf("operation cancelled: %w", ctx.Err())
		}

		format := options.SubtitleFormat
		if format == "" {
			format = "srt"
		}
		result.SubtitlesFile = fmt.Sprintf("/mock/output/dir/subtitles.%s", format)
	}

	progress <- ProgressUpdate{1.0, "Processing complete"}

	fmt.Println("Mock processing with context completed successfully.")
	return result, nil
}
