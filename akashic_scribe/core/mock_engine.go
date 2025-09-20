package core

import (
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
func (m *MockScribeEngine) StartProcessing(options ScribeOptions, progress chan<- ProgressUpdate) error {
	fmt.Printf("Mock StartProcessing called with options: %+v\n", options)

	// Simulate the full processing pipeline with progress updates
	progress <- ProgressUpdate{0.0, "Starting processing..."}
	time.Sleep(50 * time.Millisecond)

	progress <- ProgressUpdate{0.2, "Transcribing audio..."}
	time.Sleep(100 * time.Millisecond)

	transcription, err := m.Transcribe(options.InputFile)
	if err != nil {
		return err
	}

	progress <- ProgressUpdate{0.6, "Translating text..."}
	time.Sleep(100 * time.Millisecond)

	translation, err := m.Translate(transcription, options.TargetLanguage)
	if err != nil {
		return err
	}

	if options.CreateDubbing {
		progress <- ProgressUpdate{0.8, "Creating dubbed audio..."}
		time.Sleep(50 * time.Millisecond)
	}

	if options.CreateSubtitles {
		progress <- ProgressUpdate{0.9, "Generating subtitles..."}
		time.Sleep(50 * time.Millisecond)
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
