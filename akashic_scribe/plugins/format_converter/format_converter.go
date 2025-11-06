package format_converter

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"akashic_scribe/core"
)

// FormatConverterPlugin converts between different subtitle formats
type FormatConverterPlugin struct {
	*core.BasePlugin
	inputFormats  []string
	outputFormats []string
}

// NewFormatConverterPlugin creates a new format converter plugin
func NewFormatConverterPlugin() core.Plugin {
	base := core.NewBasePlugin(
		"voidcat.format_converter",
		"Subtitle Format Converter",
		"1.0.0",
		"Convert between SRT, VTT, ASS/SSA, and plain text subtitle formats",
		"VoidCat RDC",
	)

	plugin := &FormatConverterPlugin{
		BasePlugin:    base,
		inputFormats:  []string{"srt", "vtt", "ass", "ssa", "txt"},
		outputFormats: []string{"srt", "vtt", "ass", "txt"},
	}

	return plugin
}

// GetCapabilities returns the plugin's capabilities
func (p *FormatConverterPlugin) GetCapabilities() []core.PluginCapability {
	return []core.PluginCapability{
		core.CapabilityFormatConversion,
		core.CapabilitySubtitleGeneration,
	}
}

// SubtitleEntry represents a single subtitle entry
type SubtitleEntry struct {
	Index     int
	StartTime time.Duration
	EndTime   time.Duration
	Text      string
	Style     string // For ASS/SSA format
}

// ConvertFormat converts a subtitle file from one format to another
func (p *FormatConverterPlugin) ConvertFormat(inputPath, outputPath, targetFormat string, options map[string]interface{}) error {
	ctx := p.GetContext()

	// Validate input file
	if _, err := os.Stat(inputPath); err != nil {
		return fmt.Errorf("input file not found: %w", err)
	}

	// Detect input format
	inputFormat := strings.ToLower(filepath.Ext(inputPath)[1:])
	if inputFormat == "" {
		return fmt.Errorf("could not determine input format")
	}

	ctx.LogInfo(fmt.Sprintf("Converting %s to %s: %s -> %s", inputFormat, targetFormat, inputPath, outputPath))

	// Parse input file
	entries, err := p.parseSubtitles(inputPath, inputFormat)
	if err != nil {
		return fmt.Errorf("failed to parse input file: %w", err)
	}

	ctx.LogInfo(fmt.Sprintf("Parsed %d subtitle entries", len(entries)))

	// Write output file
	if err := p.writeSubtitles(outputPath, targetFormat, entries, options); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	ctx.LogInfo("Format conversion completed successfully")
	return nil
}

// parseSubtitles parses subtitle file based on format
func (p *FormatConverterPlugin) parseSubtitles(filePath, format string) ([]SubtitleEntry, error) {
	switch format {
	case "srt":
		return p.parseSRT(filePath)
	case "vtt":
		return p.parseVTT(filePath)
	case "ass", "ssa":
		return p.parseASS(filePath)
	case "txt":
		return p.parsePlainText(filePath)
	default:
		return nil, fmt.Errorf("unsupported input format: %s", format)
	}
}

// parseSRT parses SRT subtitle format
func (p *FormatConverterPlugin) parseSRT(filePath string) ([]SubtitleEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []SubtitleEntry
	scanner := bufio.NewScanner(file)

	var currentEntry SubtitleEntry
	var textLines []string
	state := 0 // 0: index, 1: timing, 2: text

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			// End of entry
			if len(textLines) > 0 {
				currentEntry.Text = strings.Join(textLines, "\n")
				entries = append(entries, currentEntry)
				textLines = nil
				currentEntry = SubtitleEntry{}
			}
			state = 0
			continue
		}

		switch state {
		case 0: // Index
			if idx, err := strconv.Atoi(line); err == nil {
				currentEntry.Index = idx
				state = 1
			}
		case 1: // Timing
			start, end, err := parseSRTTiming(line)
			if err == nil {
				currentEntry.StartTime = start
				currentEntry.EndTime = end
				state = 2
			}
		case 2: // Text
			textLines = append(textLines, line)
		}
	}

	// Add last entry if exists
	if len(textLines) > 0 {
		currentEntry.Text = strings.Join(textLines, "\n")
		entries = append(entries, currentEntry)
	}

	return entries, scanner.Err()
}

// parseVTT parses WebVTT subtitle format
func (p *FormatConverterPlugin) parseVTT(filePath string) ([]SubtitleEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []SubtitleEntry
	scanner := bufio.NewScanner(file)

	// Skip header
	if scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "WEBVTT") {
			return nil, fmt.Errorf("invalid VTT file: missing WEBVTT header")
		}
	}

	var currentEntry SubtitleEntry
	var textLines []string
	index := 1

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			// End of entry
			if len(textLines) > 0 {
				currentEntry.Index = index
				currentEntry.Text = strings.Join(textLines, "\n")
				entries = append(entries, currentEntry)
				index++
				textLines = nil
				currentEntry = SubtitleEntry{}
			}
			continue
		}

		// Check if it's a timing line
		if strings.Contains(line, "-->") {
			start, end, err := parseVTTTiming(line)
			if err == nil {
				currentEntry.StartTime = start
				currentEntry.EndTime = end
			}
		} else if !strings.Contains(line, "-->") && currentEntry.StartTime > 0 {
			// It's text
			textLines = append(textLines, line)
		}
	}

	// Add last entry if exists
	if len(textLines) > 0 {
		currentEntry.Index = index
		currentEntry.Text = strings.Join(textLines, "\n")
		entries = append(entries, currentEntry)
	}

	return entries, scanner.Err()
}

// parseASS parses ASS/SSA subtitle format (simplified)
func (p *FormatConverterPlugin) parseASS(filePath string) ([]SubtitleEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []SubtitleEntry
	scanner := bufio.NewScanner(file)
	index := 1

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Look for dialogue lines
		if strings.HasPrefix(line, "Dialogue:") {
			parts := strings.Split(line, ",")
			if len(parts) >= 10 {
				// ASS format: Dialogue: Layer,Start,End,Style,Name,MarginL,MarginR,MarginV,Effect,Text
				start := parseASSTime(parts[1])
				end := parseASSTime(parts[2])
				style := parts[3]
				text := strings.Join(parts[9:], ",")

				// Remove ASS formatting tags
				text = regexp.MustCompile(`\{[^}]*\}`).ReplaceAllString(text, "")
				text = strings.ReplaceAll(text, "\\N", "\n")

				entries = append(entries, SubtitleEntry{
					Index:     index,
					StartTime: start,
					EndTime:   end,
					Text:      text,
					Style:     style,
				})
				index++
			}
		}
	}

	return entries, scanner.Err()
}

// parsePlainText parses plain text (creates one subtitle per line)
func (p *FormatConverterPlugin) parsePlainText(filePath string) ([]SubtitleEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []SubtitleEntry
	scanner := bufio.NewScanner(file)
	index := 1
	currentTime := time.Duration(0)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Assume 3 seconds per subtitle
		entry := SubtitleEntry{
			Index:     index,
			StartTime: currentTime,
			EndTime:   currentTime + 3*time.Second,
			Text:      line,
		}

		entries = append(entries, entry)
		currentTime += 3 * time.Second
		index++
	}

	return entries, scanner.Err()
}

// writeSubtitles writes subtitles in the specified format
func (p *FormatConverterPlugin) writeSubtitles(filePath, format string, entries []SubtitleEntry, options map[string]interface{}) error {
	switch format {
	case "srt":
		return p.writeSRT(filePath, entries)
	case "vtt":
		return p.writeVTT(filePath, entries)
	case "ass":
		return p.writeASS(filePath, entries, options)
	case "txt":
		return p.writePlainText(filePath, entries)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// writeSRT writes SRT format
func (p *FormatConverterPlugin) writeSRT(filePath string, entries []SubtitleEntry) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, entry := range entries {
		fmt.Fprintf(file, "%d\n", entry.Index)
		fmt.Fprintf(file, "%s --> %s\n", formatSRTTime(entry.StartTime), formatSRTTime(entry.EndTime))
		fmt.Fprintf(file, "%s\n\n", entry.Text)
	}

	return nil
}

// writeVTT writes WebVTT format
func (p *FormatConverterPlugin) writeVTT(filePath string, entries []SubtitleEntry) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "WEBVTT")
	fmt.Fprintln(file)

	for _, entry := range entries {
		fmt.Fprintf(file, "%s --> %s\n", formatVTTTime(entry.StartTime), formatVTTTime(entry.EndTime))
		fmt.Fprintf(file, "%s\n\n", entry.Text)
	}

	return nil
}

// writeASS writes ASS format (simplified)
func (p *FormatConverterPlugin) writeASS(filePath string, entries []SubtitleEntry, options map[string]interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write ASS header
	fmt.Fprintln(file, "[Script Info]")
	fmt.Fprintln(file, "Title: Akashic Scribe Subtitle")
	fmt.Fprintln(file, "ScriptType: v4.00+")
	fmt.Fprintln(file)
	fmt.Fprintln(file, "[V4+ Styles]")
	fmt.Fprintln(file, "Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding")
	fmt.Fprintln(file, "Style: Default,Arial,20,&H00FFFFFF,&H000000FF,&H00000000,&H00000000,0,0,0,0,100,100,0,0,1,2,2,2,10,10,10,1")
	fmt.Fprintln(file)
	fmt.Fprintln(file, "[Events]")
	fmt.Fprintln(file, "Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text")

	for _, entry := range entries {
		style := entry.Style
		if style == "" {
			style = "Default"
		}
		text := strings.ReplaceAll(entry.Text, "\n", "\\N")
		fmt.Fprintf(file, "Dialogue: 0,%s,%s,%s,,0,0,0,,%s\n",
			formatASSTime(entry.StartTime),
			formatASSTime(entry.EndTime),
			style,
			text)
	}

	return nil
}

// writePlainText writes plain text format
func (p *FormatConverterPlugin) writePlainText(filePath string, entries []SubtitleEntry) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, entry := range entries {
		fmt.Fprintln(file, entry.Text)
	}

	return nil
}

// GetSupportedInputFormats returns supported input formats
func (p *FormatConverterPlugin) GetSupportedInputFormats() []string {
	return p.inputFormats
}

// GetSupportedOutputFormats returns supported output formats
func (p *FormatConverterPlugin) GetSupportedOutputFormats() []string {
	return p.outputFormats
}

// Timing parsing and formatting helpers

func parseSRTTiming(line string) (time.Duration, time.Duration, error) {
	parts := strings.Split(line, " --> ")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid timing line")
	}
	start, err := parseSRTTime(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, err
	}
	end, err := parseSRTTime(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, err
	}
	return start, end, nil
}

func parseSRTTime(timeStr string) (time.Duration, error) {
	// Format: HH:MM:SS,mmm
	timeStr = strings.ReplaceAll(timeStr, ",", ".")
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid time format")
	}

	hours, _ := strconv.Atoi(parts[0])
	minutes, _ := strconv.Atoi(parts[1])
	secondsParts := strings.Split(parts[2], ".")
	seconds, _ := strconv.Atoi(secondsParts[0])
	var milliseconds int
	if len(secondsParts) > 1 {
		milliseconds, _ = strconv.Atoi(secondsParts[1])
	}

	return time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second +
		time.Duration(milliseconds)*time.Millisecond, nil
}

func parseVTTTiming(line string) (time.Duration, time.Duration, error) {
	parts := strings.Split(line, " --> ")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid timing line")
	}
	start, err := parseVTTTime(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, err
	}
	end, err := parseVTTTime(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, err
	}
	return start, end, nil
}

func parseVTTTime(timeStr string) (time.Duration, error) {
	// Format: HH:MM:SS.mmm or MM:SS.mmm
	parts := strings.Split(timeStr, ":")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid time format")
	}

	var hours, minutes, seconds, milliseconds int

	if len(parts) == 3 {
		hours, _ = strconv.Atoi(parts[0])
		minutes, _ = strconv.Atoi(parts[1])
		secondsParts := strings.Split(parts[2], ".")
		seconds, _ = strconv.Atoi(secondsParts[0])
		if len(secondsParts) > 1 {
			milliseconds, _ = strconv.Atoi(secondsParts[1])
		}
	} else {
		minutes, _ = strconv.Atoi(parts[0])
		secondsParts := strings.Split(parts[1], ".")
		seconds, _ = strconv.Atoi(secondsParts[0])
		if len(secondsParts) > 1 {
			milliseconds, _ = strconv.Atoi(secondsParts[1])
		}
	}

	return time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second +
		time.Duration(milliseconds)*time.Millisecond, nil
}

func parseASSTime(timeStr string) time.Duration {
	// Format: H:MM:SS.cc (centiseconds)
	timeStr = strings.TrimSpace(timeStr)
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return 0
	}

	hours, _ := strconv.Atoi(parts[0])
	minutes, _ := strconv.Atoi(parts[1])
	secondsParts := strings.Split(parts[2], ".")
	seconds, _ := strconv.Atoi(secondsParts[0])
	var centiseconds int
	if len(secondsParts) > 1 {
		centiseconds, _ = strconv.Atoi(secondsParts[1])
	}

	return time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second +
		time.Duration(centiseconds)*10*time.Millisecond
}

func formatSRTTime(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000
	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, milliseconds)
}

func formatVTTTime(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}

func formatASSTime(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	centiseconds := int(d.Milliseconds()/10) % 100
	return fmt.Sprintf("%d:%02d:%02d.%02d", hours, minutes, seconds, centiseconds)
}

// Ensure FormatConverterPlugin implements the FormatConverter interface
var _ core.FormatConverter = (*FormatConverterPlugin)(nil)
