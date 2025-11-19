package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParseYtDlpProgress tests yt-dlp progress parsing.
func TestParseYtDlpProgress(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected float64
		shouldOk bool
	}{
		{
			name:     "Basic progress",
			line:     "[download]  45.2% of 100.00MiB at 1.50MiB/s ETA 00:36",
			expected: 0.452,
			shouldOk: true,
		},
		{
			name:     "100% complete",
			line:     "[download] 100% of 100.00MiB in 01:23",
			expected: 1.0,
			shouldOk: true,
		},
		{
			name:     "Low percentage",
			line:     "[download]   1.5% of 100.00MiB at 1.50MiB/s ETA 00:36",
			expected: 0.015,
			shouldOk: true,
		},
		{
			name:     "Progress with extra spaces",
			line:     "[download]    75.8% of 200.00MiB",
			expected: 0.758,
			shouldOk: true,
		},
		{
			name:     "Zero progress",
			line:     "[download]  0.0% of 100.00MiB",
			expected: 0.0,
			shouldOk: true,
		},
		{
			name:     "Non-download line",
			line:     "[info] Downloading video...",
			expected: 0,
			shouldOk: false,
		},
		{
			name:     "Empty line",
			line:     "",
			expected: 0,
			shouldOk: false,
		},
		{
			name:     "Different bracket type",
			line:     "(download) 50% complete",
			expected: 0,
			shouldOk: false,
		},
		{
			name:     "Progress without percentage sign",
			line:     "[download] 50 of 100MiB",
			expected: 0,
			shouldOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			progress, ok := parseYtDlpProgress(tt.line)

			assert.Equal(tt.shouldOk, ok, "Expected ok=%v, got %v", tt.shouldOk, ok)
			if tt.shouldOk {
				assert.InDelta(tt.expected, progress, 0.001, "Progress value mismatch")
			}
		})
	}
}

// TestParseFfmpegProgress tests ffmpeg progress parsing.
func TestParseFfmpegProgress(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		duration float64
		expected float64
		shouldOk bool
	}{
		{
			name:     "Basic progress",
			line:     "frame=1234 fps=30 size=1024kB time=00:01:23.45 bitrate=1234.5kbits/s speed=2.5x",
			duration: 180.0, // 3 minutes
			expected: 0.463,  // 83.45 / 180
			shouldOk: true,
		},
		{
			name:     "Beginning of file",
			line:     "frame=10 fps=30 size=100kB time=00:00:01.00 bitrate=800kbits/s speed=1.0x",
			duration: 60.0,
			expected: 0.0166, // 1.0 / 60
			shouldOk: true,
		},
		{
			name:     "Mid-process",
			line:     "frame=5000 fps=25 size=50MB time=00:05:30.00 bitrate=1280kbits/s speed=2.0x",
			duration: 600.0, // 10 minutes
			expected: 0.55,   // 330 / 600
			shouldOk: true,
		},
		{
			name:     "Nearly complete",
			line:     "frame=7200 fps=24 size=100MB time=00:09:55.00 bitrate=1400kbits/s speed=1.5x",
			duration: 600.0,
			expected: 0.991, // 595 / 600
			shouldOk: true,
		},
		{
			name:     "Progress exceeds duration (capped at 1.0)",
			line:     "frame=8000 fps=24 size=110MB time=00:11:00.00 bitrate=1400kbits/s speed=1.5x",
			duration: 600.0,
			expected: 1.0, // Capped
			shouldOk: true,
		},
		{
			name:     "Zero duration (invalid)",
			line:     "frame=1000 fps=24 size=10MB time=00:01:00.00 bitrate=1400kbits/s speed=1.5x",
			duration: 0.0,
			expected: 0,
			shouldOk: false,
		},
		{
			name:     "Line without time field",
			line:     "frame=1000 fps=24 size=10MB bitrate=1400kbits/s speed=1.5x",
			duration: 600.0,
			expected: 0,
			shouldOk: false,
		},
		{
			name:     "Empty line",
			line:     "",
			duration: 600.0,
			expected: 0,
			shouldOk: false,
		},
		{
			name:     "Non-ffmpeg line",
			line:     "Starting transcoding process...",
			duration: 600.0,
			expected: 0,
			shouldOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			progress, ok := parseFfmpegProgress(tt.line, tt.duration)

			assert.Equal(tt.shouldOk, ok, "Expected ok=%v, got %v", tt.shouldOk, ok)
			if tt.shouldOk {
				assert.InDelta(tt.expected, progress, 0.01, "Progress value mismatch")
			}
		})
	}
}

// TestProgressUpdateStruct tests the ProgressUpdate struct.
func TestProgressUpdateStruct(t *testing.T) {
	assert := assert.New(t)

	// Test creating progress updates
	updates := []ProgressUpdate{
		{0.0, "Starting..."},
		{0.25, "Processing..."},
		{0.50, "Halfway done..."},
		{0.75, "Almost there..."},
		{1.0, "Complete!"},
	}

	assert.Len(updates, 5)

	// Test field access
	assert.Equal(0.0, updates[0].Percentage)
	assert.Equal("Starting...", updates[0].Message)

	assert.Equal(1.0, updates[4].Percentage)
	assert.Equal("Complete!", updates[4].Message)
}

// TestProgressSequence tests that progress updates are monotonically increasing.
func TestProgressSequence(t *testing.T) {
	assert := assert.New(t)

	engine := NewMockScribeEngine()
	ctx := context.Background()

	opts := ScribeOptions{
		InputFile:      "test.mp4",
		OriginLanguage: "English",
		TargetLanguage: "Spanish",
		CreateDubbing:  true,
		CreateSubtitles: true,
	}

	progressChan := make(chan ProgressUpdate, 20)

	// Start processing
	go func() {
		_ = engine.StartProcessing(ctx, opts, progressChan)
		close(progressChan)
	}()

	// Collect all progress updates
	var updates []ProgressUpdate
	for update := range progressChan {
		updates = append(updates, update)
	}

	// Verify we got updates
	assert.NotEmpty(updates, "Should receive progress updates")

	// Verify progress is monotonically increasing
	for i := 1; i < len(updates); i++ {
		assert.GreaterOrEqual(updates[i].Percentage, updates[i-1].Percentage,
			"Progress should be monotonically increasing")
	}

	// Verify final progress is 1.0
	assert.Equal(1.0, updates[len(updates)-1].Percentage, "Final progress should be 100%")
}

// TestProgressGranularity tests that we get sufficient progress updates.
func TestProgressGranularity(t *testing.T) {
	assert := assert.New(t)

	engine := NewMockScribeEngine()
	ctx := context.Background()

	opts := ScribeOptions{
		InputFile:      "test.mp4",
		OriginLanguage: "English",
		TargetLanguage: "Spanish",
		CreateDubbing:  true,
		CreateSubtitles: true,
	}

	progressChan := make(chan ProgressUpdate, 20)

	// Start processing
	go func() {
		_ = engine.StartProcessing(ctx, opts, progressChan)
		close(progressChan)
	}()

	// Collect all progress updates
	var updates []ProgressUpdate
	for update := range progressChan {
		updates = append(updates, update)
	}

	// Should receive at least 5 updates (start, transcribe, translate, dub, subtitles, complete)
	assert.GreaterOrEqual(len(updates), 5, "Should receive multiple progress updates")

	// Should start at 0
	assert.Equal(0.0, updates[0].Percentage, "Should start at 0%")

	// Should end at 1.0
	assert.Equal(1.0, updates[len(updates)-1].Percentage, "Should end at 100%")
}

// TestProgressMessages tests that progress messages are descriptive.
func TestProgressMessages(t *testing.T) {
	assert := assert.New(t)

	engine := NewMockScribeEngine()
	ctx := context.Background()

	opts := ScribeOptions{
		InputFile:      "test.mp4",
		OriginLanguage: "English",
		TargetLanguage: "Spanish",
		CreateDubbing:  true,
		CreateSubtitles: true,
	}

	progressChan := make(chan ProgressUpdate, 20)

	// Start processing
	go func() {
		_ = engine.StartProcessing(ctx, opts, progressChan)
		close(progressChan)
	}()

	// Collect all progress updates
	var updates []ProgressUpdate
	for update := range progressChan {
		updates = append(updates, update)
	}

	// Verify messages are not empty
	for _, update := range updates {
		assert.NotEmpty(update.Message, "Progress message should not be empty")
	}

	// Check for key stage messages
	messages := make([]string, len(updates))
	for i, update := range updates {
		messages[i] = update.Message
	}

	// Should contain references to key stages
	hasTranscribe := false
	hasTranslate := false
	hasDubbing := false
	hasSubtitles := false
	hasComplete := false

	for _, msg := range messages {
		if contains(msg, "Transcrib") {
			hasTranscribe = true
		}
		if contains(msg, "Translat") {
			hasTranslate = true
		}
		if contains(msg, "dub") || contains(msg, "Dub") {
			hasDubbing = true
		}
		if contains(msg, "subtitle") || contains(msg, "Subtitle") {
			hasSubtitles = true
		}
		if contains(msg, "complete") || contains(msg, "Complete") {
			hasComplete = true
		}
	}

	assert.True(hasTranscribe, "Should have transcription message")
	assert.True(hasTranslate, "Should have translation message")
	assert.True(hasDubbing, "Should have dubbing message")
	assert.True(hasSubtitles, "Should have subtitles message")
	assert.True(hasComplete, "Should have completion message")
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// BenchmarkProgressParsing benchmarks the progress parsing functions.
func BenchmarkParseYtDlpProgress(b *testing.B) {
	line := "[download]  45.2% of 100.00MiB at 1.50MiB/s ETA 00:36"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseYtDlpProgress(line)
	}
}

func BenchmarkParseFfmpegProgress(b *testing.B) {
	line := "frame=1234 fps=30 size=1024kB time=00:01:23.45 bitrate=1234.5kbits/s speed=2.5x"
	duration := 180.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseFfmpegProgress(line, duration)
	}
}
