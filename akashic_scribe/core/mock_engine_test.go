package core

import (
	"testing"
)

func TestMockScribeEngine(t *testing.T) {
	engine := NewMockScribeEngine()

	// Test Transcription
	t.Run("Transcribe", func(t *testing.T) {
		videoSource := "test_video.mp4"
		expectedTranscription := "This is a mock transcription of the video."

		transcription, err := engine.Transcribe(videoSource)
		if err != nil {
			t.Errorf("Transcribe() returned an unexpected error: %v", err)
		}

		if transcription != expectedTranscription {
			t.Errorf("Transcribe() returned %q, want %q", transcription, expectedTranscription)
		}
	})

	// Test Translation
	t.Run("Translate", func(t *testing.T) {
		textToTranslate := "Hello, world."
		targetLanguage := "Español"
		expectedTranslation := "This is a mock translation to 'Español'."

		translation, err := engine.Translate(textToTranslate, targetLanguage)
		if err != nil {
			t.Errorf("Translate() returned an unexpected error: %v", err)
		}

		if translation != expectedTranslation {
			t.Errorf("Translate() returned %q, want %q", translation, expectedTranslation)
		}
	})
}
