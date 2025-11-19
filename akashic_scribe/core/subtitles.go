package core

import (
	"fmt"
	"strings"
	"time"
)

// SubtitleSegment represents a single subtitle entry with timing information.
type SubtitleSegment struct {
	Index     int           // Subtitle sequence number (1-indexed)
	StartTime time.Duration // Start timestamp
	EndTime   time.Duration // End timestamp
	Text      string        // Subtitle text
	Original  string        // Original language text (for bilingual subtitles)
}

// SubtitleGenerator handles advanced subtitle generation with proper timing.
type SubtitleGenerator struct {
	segments []SubtitleSegment
}

// NewSubtitleGenerator creates a new subtitle generator.
func NewSubtitleGenerator() *SubtitleGenerator {
	return &SubtitleGenerator{
		segments: make([]SubtitleSegment, 0),
	}
}

// AddSegment adds a subtitle segment to the generator.
func (sg *SubtitleGenerator) AddSegment(start, end time.Duration, text string, original string) {
	segment := SubtitleSegment{
		Index:     len(sg.segments) + 1,
		StartTime: start,
		EndTime:   end,
		Text:      text,
		Original:  original,
	}
	sg.segments = append(sg.segments, segment)
}

// GenerateSRT generates subtitles in SRT format with proper timing.
//
// Parameters:
//   - bilingual: If true, includes both original and translated text
//   - position: "top" or "bottom" - determines translation position in bilingual mode
//
// Returns the complete SRT-formatted subtitle file content.
func (sg *SubtitleGenerator) GenerateSRT(bilingual bool, position string) string {
	var builder strings.Builder

	for _, segment := range sg.segments {
		// Write sequence number
		builder.WriteString(fmt.Sprintf("%d\n", segment.Index))

		// Write timestamp in SRT format: HH:MM:SS,mmm --> HH:MM:SS,mmm
		builder.WriteString(fmt.Sprintf("%s --> %s\n",
			formatSRTTimestamp(segment.StartTime),
			formatSRTTimestamp(segment.EndTime)))

		// Write subtitle text
		if bilingual && segment.Original != "" {
			// Bilingual subtitles with original and translation
			if position == "top" {
				// Translation on top, original on bottom
				builder.WriteString(segment.Text + "\n")
				builder.WriteString(segment.Original + "\n")
			} else {
				// Original on top, translation on bottom (default)
				builder.WriteString(segment.Original + "\n")
				builder.WriteString(segment.Text + "\n")
			}
		} else {
			// Monolingual subtitles (translated text only)
			builder.WriteString(segment.Text + "\n")
		}

		// Add blank line between entries
		builder.WriteString("\n")
	}

	return builder.String()
}

// GenerateVTT generates subtitles in WebVTT format.
func (sg *SubtitleGenerator) GenerateVTT(bilingual bool, position string) string {
	var builder strings.Builder

	// VTT header
	builder.WriteString("WEBVTT\n\n")

	for _, segment := range sg.segments {
		// Write timestamp in VTT format: HH:MM:SS.mmm --> HH:MM:SS.mmm
		builder.WriteString(fmt.Sprintf("%s --> %s\n",
			formatVTTTimestamp(segment.StartTime),
			formatVTTTimestamp(segment.EndTime)))

		// Write subtitle text
		if bilingual && segment.Original != "" {
			if position == "top" {
				builder.WriteString(segment.Text + "\n")
				builder.WriteString(segment.Original + "\n")
			} else {
				builder.WriteString(segment.Original + "\n")
				builder.WriteString(segment.Text + "\n")
			}
		} else {
			builder.WriteString(segment.Text + "\n")
		}

		// Add blank line between entries
		builder.WriteString("\n")
	}

	return builder.String()
}

// formatSRTTimestamp converts a time.Duration to SRT timestamp format (HH:MM:SS,mmm)
func formatSRTTimestamp(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000

	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, milliseconds)
}

// formatVTTTimestamp converts a time.Duration to WebVTT timestamp format (HH:MM:SS.mmm)
func formatVTTTimestamp(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}

// CreateDefaultSegments creates default subtitle segments when timing data is unavailable.
// This splits text into segments of approximately equal length with default timing.
func (sg *SubtitleGenerator) CreateDefaultSegments(originalText, translatedText string, videoDuration time.Duration) {
	// Split text into sentences (basic splitting on period, exclamation, question mark)
	sentences := splitIntoSentences(translatedText)
	originalSentences := splitIntoSentences(originalText)

	if len(sentences) == 0 {
		return
	}

	// Calculate duration per segment
	segmentDuration := videoDuration / time.Duration(len(sentences))
	if segmentDuration < 2*time.Second {
		segmentDuration = 2 * time.Second // Minimum 2 seconds per subtitle
	}
	if segmentDuration > 7*time.Second {
		segmentDuration = 7 * time.Second // Maximum 7 seconds per subtitle
	}

	currentTime := time.Duration(0)

	for i, sentence := range sentences {
		if sentence == "" {
			continue
		}

		var original string
		if i < len(originalSentences) {
			original = originalSentences[i]
		}

		endTime := currentTime + segmentDuration
		if endTime > videoDuration {
			endTime = videoDuration
		}

		sg.AddSegment(currentTime, endTime, sentence, original)
		currentTime = endTime
	}
}

// splitIntoSentences splits text into sentences using basic punctuation.
func splitIntoSentences(text string) []string {
	// Replace common sentence endings with a delimiter
	text = strings.ReplaceAll(text, ". ", ".|")
	text = strings.ReplaceAll(text, "! ", "!|")
	text = strings.ReplaceAll(text, "? ", "?|")

	// Split on delimiter
	sentences := strings.Split(text, "|")

	// Clean up and filter empty sentences
	result := make([]string, 0, len(sentences))
	for _, s := range sentences {
		s = strings.TrimSpace(s)
		if s != "" {
			result = append(result, s)
		}
	}

	return result
}

// SyncSubtitlesToAudio synchronizes subtitle timing based on audio analysis.
// This is a placeholder for more advanced timing synchronization that would
// use speech recognition timestamps or audio analysis.
func (sg *SubtitleGenerator) SyncSubtitlesToAudio(audioPath string) error {
	// TODO: Implement audio-based synchronization using:
	// - Whisper API timestamps for each word/phrase
	// - Audio silence detection to adjust segment boundaries
	// - Speech rate analysis for more natural timing
	//
	// For now, this returns nil as the manual timing works adequately.
	return nil
}
