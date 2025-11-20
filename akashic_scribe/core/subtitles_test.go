package core

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestSubtitleGenerator_AddSegment(t *testing.T) {
	sg := NewSubtitleGenerator()

	sg.AddSegment(0, 3*time.Second, "Hello world", "Hola mundo")
	sg.AddSegment(3*time.Second, 6*time.Second, "How are you?", "¿Cómo estás?")

	if len(sg.segments) != 2 {
		t.Errorf("Expected 2 segments, got %d", len(sg.segments))
	}

	if sg.segments[0].Index != 1 {
		t.Errorf("First segment index should be 1, got %d", sg.segments[0].Index)
	}

	if sg.segments[1].Index != 2 {
		t.Errorf("Second segment index should be 2, got %d", sg.segments[1].Index)
	}
}

func TestSubtitleGenerator_GenerateSRT(t *testing.T) {
	sg := NewSubtitleGenerator()
	sg.AddSegment(0, 3*time.Second, "Hello world", "")

	srt := sg.GenerateSRT(false, "bottom")

	// Check for SRT format components
	if !strings.Contains(srt, "1\n") {
		t.Error("SRT should contain sequence number")
	}

	if !strings.Contains(srt, "00:00:00,000 --> 00:00:03,000") {
		t.Error("SRT should contain proper timestamp format")
	}

	if !strings.Contains(srt, "Hello world") {
		t.Error("SRT should contain subtitle text")
	}
}

func TestSubtitleGenerator_GenerateSRT_Bilingual(t *testing.T) {
	sg := NewSubtitleGenerator()
	sg.AddSegment(0, 3*time.Second, "Hello world", "Hola mundo")

	// Test with translation on top
	srtTop := sg.GenerateSRT(true, "top")
	lines := strings.Split(srtTop, "\n")

	// Find the text lines (after timestamp)
	var textLines []string
	foundTimestamp := false
	for _, line := range lines {
		if strings.Contains(line, "-->") {
			foundTimestamp = true
			continue
		}
		if foundTimestamp && line != "" && !strings.HasPrefix(line, "1") {
			textLines = append(textLines, line)
			if len(textLines) == 2 {
				break
			}
		}
	}

	if len(textLines) < 2 {
		t.Fatalf("Expected at least 2 text lines for bilingual subtitles, got %d", len(textLines))
	}

	if textLines[0] != "Hello world" {
		t.Errorf("First line should be translation 'Hello world', got '%s'", textLines[0])
	}

	if textLines[1] != "Hola mundo" {
		t.Errorf("Second line should be original 'Hola mundo', got '%s'", textLines[1])
	}

	// Test with translation on bottom
	srtBottom := sg.GenerateSRT(true, "bottom")
	lines = strings.Split(srtBottom, "\n")

	textLines = []string{}
	foundTimestamp = false
	for _, line := range lines {
		if strings.Contains(line, "-->") {
			foundTimestamp = true
			continue
		}
		if foundTimestamp && line != "" && !strings.HasPrefix(line, "1") {
			textLines = append(textLines, line)
			if len(textLines) == 2 {
				break
			}
		}
	}

	if textLines[0] != "Hola mundo" {
		t.Errorf("First line should be original 'Hola mundo', got '%s'", textLines[0])
	}

	if textLines[1] != "Hello world" {
		t.Errorf("Second line should be translation 'Hello world', got '%s'", textLines[1])
	}
}

func TestSubtitleGenerator_GenerateVTT(t *testing.T) {
	sg := NewSubtitleGenerator()
	sg.AddSegment(0, 3*time.Second, "Hello world", "")

	vtt := sg.GenerateVTT(false, "bottom")

	// Check for VTT format components
	if !strings.HasPrefix(vtt, "WEBVTT\n") {
		t.Error("VTT should start with WEBVTT header")
	}

	if !strings.Contains(vtt, "00:00:00.000 --> 00:00:03.000") {
		t.Error("VTT should contain proper timestamp format with dots")
	}

	if !strings.Contains(vtt, "Hello world") {
		t.Error("VTT should contain subtitle text")
	}
}

func TestFormatSRTTimestamp(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{0, "00:00:00,000"},
		{3 * time.Second, "00:00:03,000"},
		{1*time.Minute + 30*time.Second, "00:01:30,000"},
		{1*time.Hour + 2*time.Minute + 3*time.Second + 456*time.Millisecond, "01:02:03,456"},
	}

	for _, tt := range tests {
		result := formatSRTTimestamp(tt.duration)
		if result != tt.expected {
			t.Errorf("formatSRTTimestamp(%v) = %s, expected %s", tt.duration, result, tt.expected)
		}
	}
}

func TestFormatVTTTimestamp(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{0, "00:00:00.000"},
		{3 * time.Second, "00:00:03.000"},
		{1*time.Minute + 30*time.Second, "00:01:30.000"},
		{1*time.Hour + 2*time.Minute + 3*time.Second + 456*time.Millisecond, "01:02:03.456"},
	}

	for _, tt := range tests {
		result := formatVTTTimestamp(tt.duration)
		if result != tt.expected {
			t.Errorf("formatVTTTimestamp(%v) = %s, expected %s", tt.duration, result, tt.expected)
		}
	}
}

func TestCreateDefaultSegments(t *testing.T) {
	sg := NewSubtitleGenerator()

	originalText := "This is the first sentence. This is the second sentence. This is the third sentence."
	translatedText := "Esta es la primera frase. Esta es la segunda frase. Esta es la tercera frase."
	videoDuration := 30 * time.Second

	sg.CreateDefaultSegments(originalText, translatedText, videoDuration)

	if len(sg.segments) == 0 {
		t.Error("Expected segments to be created")
	}

	// Check that timing is reasonable
	for i, segment := range sg.segments {
		if segment.StartTime >= segment.EndTime {
			t.Errorf("Segment %d: StartTime (%v) should be less than EndTime (%v)",
				i, segment.StartTime, segment.EndTime)
		}

		if segment.EndTime > videoDuration {
			t.Errorf("Segment %d: EndTime (%v) exceeds video duration (%v)",
				i, segment.EndTime, videoDuration)
		}
	}

	// Check that segments don't overlap
	for i := 0; i < len(sg.segments)-1; i++ {
		if sg.segments[i].EndTime > sg.segments[i+1].StartTime {
			t.Errorf("Segments %d and %d overlap", i, i+1)
		}
	}
}

func TestSplitIntoSentences(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"Hello world. How are you?", 2},
		{"Single sentence", 1},
		{"One. Two. Three!", 3},
		{"Question? Answer! Statement.", 3},
		{"", 0},
	}

	for _, tt := range tests {
		result := splitIntoSentences(tt.input)
		if len(result) != tt.expected {
			t.Errorf("splitIntoSentences(%q) returned %d sentences, expected %d",
				tt.input, len(result), tt.expected)
		}
	}
}

func TestSubtitleGenerator_MultipleSegments(t *testing.T) {
	sg := NewSubtitleGenerator()

	// Add multiple segments
	for i := 0; i < 5; i++ {
		start := time.Duration(i*3) * time.Second
		end := time.Duration(i*3+3) * time.Second
		text := fmt.Sprintf("Segment %d", i+1)
		sg.AddSegment(start, end, text, "")
	}

	srt := sg.GenerateSRT(false, "bottom")

	// Check that all segments are present
	for i := 1; i <= 5; i++ {
		segmentMarker := fmt.Sprintf("%d\n", i)
		if !strings.Contains(srt, segmentMarker) {
			t.Errorf("SRT should contain segment %d", i)
		}

		segmentText := fmt.Sprintf("Segment %d", i)
		if !strings.Contains(srt, segmentText) {
			t.Errorf("SRT should contain text for segment %d", i)
		}
	}
}

// Benchmark tests
func BenchmarkSubtitleGenerator_AddSegment(b *testing.B) {
	sg := NewSubtitleGenerator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sg.AddSegment(0, 3*time.Second, "Test subtitle", "")
	}
}

func BenchmarkSubtitleGenerator_GenerateSRT(b *testing.B) {
	sg := NewSubtitleGenerator()
	for i := 0; i < 100; i++ {
		start := time.Duration(i*3) * time.Second
		end := time.Duration(i*3+3) * time.Second
		sg.AddSegment(start, end, "Test subtitle", "Original text")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sg.GenerateSRT(true, "bottom")
	}
}

func BenchmarkSubtitleGenerator_GenerateVTT(b *testing.B) {
	sg := NewSubtitleGenerator()
	for i := 0; i < 100; i++ {
		start := time.Duration(i*3) * time.Second
		end := time.Duration(i*3+3) * time.Second
		sg.AddSegment(start, end, "Test subtitle", "Original text")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sg.GenerateVTT(true, "bottom")
	}
}
