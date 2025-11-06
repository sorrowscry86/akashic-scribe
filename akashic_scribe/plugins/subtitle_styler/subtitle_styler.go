package subtitle_styler

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"akashic_scribe/core"
)

// SubtitleStylerPlugin applies styling and themes to subtitle files
type SubtitleStylerPlugin struct {
	*core.BasePlugin
	supportedFormats []string
	themes           map[string]SubtitleTheme
}

// SubtitleTheme defines visual styling for subtitles
type SubtitleTheme struct {
	Name            string
	FontName        string
	FontSize        int
	PrimaryColor    string  // &HAABBGGRR format for ASS
	SecondaryColor  string
	OutlineColor    string
	BackColor       string
	Bold            bool
	Italic          bool
	BorderStyle     int
	Outline         float64
	Shadow          float64
	Alignment       int // 1-9 numpad style
	MarginL         int
	MarginR         int
	MarginV         int
}

// NewSubtitleStylerPlugin creates a new subtitle styler plugin
func NewSubtitleStylerPlugin() core.Plugin {
	base := core.NewBasePlugin(
		"voidcat.subtitle_styler",
		"Subtitle Styler",
		"1.0.0",
		"Apply beautiful styling and themes to subtitle files with various visual presets",
		"VoidCat RDC",
	)

	plugin := &SubtitleStylerPlugin{
		BasePlugin:       base,
		supportedFormats: []string{"srt", "ass", "ssa"},
		themes:           make(map[string]SubtitleTheme),
	}

	// Initialize built-in themes
	plugin.initializeThemes()

	return plugin
}

// GetCapabilities returns the plugin's capabilities
func (p *SubtitleStylerPlugin) GetCapabilities() []core.PluginCapability {
	return []core.PluginCapability{
		core.CapabilitySubtitleStyling,
		core.CapabilitySubtitleGeneration,
	}
}

// initializeThemes sets up built-in subtitle themes
func (p *SubtitleStylerPlugin) initializeThemes() {
	p.themes["default"] = SubtitleTheme{
		Name:           "Default",
		FontName:       "Arial",
		FontSize:       20,
		PrimaryColor:   "&H00FFFFFF", // White
		SecondaryColor: "&H000000FF", // Red
		OutlineColor:   "&H00000000", // Black
		BackColor:      "&H00000000", // Black
		Bold:           false,
		Italic:         false,
		BorderStyle:    1,
		Outline:        2,
		Shadow:         2,
		Alignment:      2, // Bottom center
		MarginL:        10,
		MarginR:        10,
		MarginV:        10,
	}

	p.themes["cinema"] = SubtitleTheme{
		Name:           "Cinema",
		FontName:       "Trebuchet MS",
		FontSize:       24,
		PrimaryColor:   "&H00FFFFFF", // White
		SecondaryColor: "&H00FFFF00", // Yellow
		OutlineColor:   "&H00000000", // Black
		BackColor:      "&H00000000", // Black
		Bold:           true,
		Italic:         false,
		BorderStyle:    1,
		Outline:        3,
		Shadow:         3,
		Alignment:      2,
		MarginL:        20,
		MarginR:        20,
		MarginV:        20,
	}

	p.themes["modern"] = SubtitleTheme{
		Name:           "Modern",
		FontName:       "Segoe UI",
		FontSize:       22,
		PrimaryColor:   "&H00F0F0F0", // Light gray
		SecondaryColor: "&H00FFAA00", // Orange
		OutlineColor:   "&H00202020", // Dark gray
		BackColor:      "&H80000000", // Semi-transparent black
		Bold:           false,
		Italic:         false,
		BorderStyle:    1,
		Outline:        2.5,
		Shadow:         1,
		Alignment:      2,
		MarginL:        15,
		MarginR:        15,
		MarginV:        15,
	}

	p.themes["elegant"] = SubtitleTheme{
		Name:           "Elegant",
		FontName:       "Georgia",
		FontSize:       20,
		PrimaryColor:   "&H00FFFFCC", // Light yellow
		SecondaryColor: "&H00FFCC00", // Gold
		OutlineColor:   "&H00000000", // Black
		BackColor:      "&H00000000", // Black
		Bold:           false,
		Italic:         true,
		BorderStyle:    1,
		Outline:        2,
		Shadow:         2,
		Alignment:      2,
		MarginL:        10,
		MarginR:        10,
		MarginV:        10,
	}

	p.themes["bold_yellow"] = SubtitleTheme{
		Name:           "Bold Yellow",
		FontName:       "Arial",
		FontSize:       24,
		PrimaryColor:   "&H0000FFFF", // Yellow
		SecondaryColor: "&H00FFFFFF", // White
		OutlineColor:   "&H00000000", // Black
		BackColor:      "&H00000000", // Black
		Bold:           true,
		Italic:         false,
		BorderStyle:    1,
		Outline:        3,
		Shadow:         2,
		Alignment:      2,
		MarginL:        10,
		MarginR:        10,
		MarginV:        15,
	}

	p.themes["anime"] = SubtitleTheme{
		Name:           "Anime",
		FontName:       "Arial",
		FontSize:       22,
		PrimaryColor:   "&H00FFFFFF", // White
		SecondaryColor: "&H00FF00FF", // Magenta
		OutlineColor:   "&H00000000", // Black
		BackColor:      "&H00000000", // Black
		Bold:           true,
		Italic:         false,
		BorderStyle:    1,
		Outline:        2.5,
		Shadow:         1.5,
		Alignment:      2,
		MarginL:        10,
		MarginR:        10,
		MarginV:        12,
	}
}

// ProcessSubtitles applies styling to a subtitle file
// Options:
//   - "theme": string - Theme name to apply (default: "default")
//   - "font_size": int - Override font size
//   - "position": string - "top", "center", or "bottom" (default: "bottom")
//   - "add_background": bool - Add background box to subtitles
func (p *SubtitleStylerPlugin) ProcessSubtitles(inputPath, outputPath string, options map[string]interface{}) error {
	ctx := p.GetContext()

	// Validate input file
	if _, err := os.Stat(inputPath); err != nil {
		return fmt.Errorf("input file not found: %w", err)
	}

	// Get theme
	themeName := "default"
	if name, ok := options["theme"].(string); ok {
		themeName = name
	}

	theme, exists := p.themes[themeName]
	if !exists {
		return fmt.Errorf("theme '%s' not found", themeName)
	}

	// Apply option overrides
	if fontSize, ok := options["font_size"].(int); ok {
		theme.FontSize = fontSize
	}

	if position, ok := options["position"].(string); ok {
		switch position {
		case "top":
			theme.Alignment = 8 // Top center
		case "center":
			theme.Alignment = 5 // Middle center
		case "bottom":
			theme.Alignment = 2 // Bottom center
		}
	}

	if addBg, ok := options["add_background"].(bool); ok && addBg {
		theme.BackColor = "&H80000000" // Semi-transparent black
	}

	ctx.LogInfo(fmt.Sprintf("Applying theme '%s' to: %s", themeName, inputPath))

	// Determine input format
	inputFormat := strings.ToLower(filepath.Ext(inputPath)[1:])

	// Process based on format
	switch inputFormat {
	case "srt":
		return p.convertSRTToStyledASS(inputPath, outputPath, theme)
	case "ass", "ssa":
		return p.applyThemeToASS(inputPath, outputPath, theme)
	default:
		return fmt.Errorf("unsupported format: %s", inputFormat)
	}
}

// convertSRTToStyledASS converts SRT to ASS with styling
func (p *SubtitleStylerPlugin) convertSRTToStyledASS(inputPath, outputPath string, theme SubtitleTheme) error {
	// Read SRT file
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Parse SRT entries
	type Entry struct {
		index     int
		startTime string
		endTime   string
		text      string
	}

	var entries []Entry
	var currentEntry Entry
	var textLines []string
	state := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			if len(textLines) > 0 {
				currentEntry.text = strings.Join(textLines, "\\N")
				entries = append(entries, currentEntry)
				textLines = nil
				currentEntry = Entry{}
			}
			state = 0
			continue
		}

		switch state {
		case 0: // Index
			if idx, err := strconv.Atoi(line); err == nil {
				currentEntry.index = idx
				state = 1
			}
		case 1: // Timing
			if strings.Contains(line, "-->") {
				parts := strings.Split(line, "-->")
				if len(parts) == 2 {
					currentEntry.startTime = p.srtTimeToASS(strings.TrimSpace(parts[0]))
					currentEntry.endTime = p.srtTimeToASS(strings.TrimSpace(parts[1]))
					state = 2
				}
			}
		case 2: // Text
			textLines = append(textLines, line)
		}
	}

	// Add last entry
	if len(textLines) > 0 {
		currentEntry.text = strings.Join(textLines, "\\N")
		entries = append(entries, currentEntry)
	}

	// Write ASS file
	return p.writeStyledASS(outputPath, theme, entries)
}

// applyThemeToASS applies theme to existing ASS file
func (p *SubtitleStylerPlugin) applyThemeToASS(inputPath, outputPath string, theme SubtitleTheme) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(inFile)
	inStyles := false
	stylesWritten := false

	for scanner.Scan() {
		line := scanner.Text()

		// Detect styles section
		if strings.HasPrefix(line, "[V4+ Styles]") || strings.HasPrefix(line, "[V4 Styles]") {
			inStyles = true
			fmt.Fprintln(outFile, line)
			continue
		}

		// Write new style when we hit the format line
		if inStyles && strings.HasPrefix(line, "Format:") {
			fmt.Fprintln(outFile, line)
			fmt.Fprintln(outFile, p.formatASSStyle(theme))
			stylesWritten = true
			continue
		}

		// Skip old style definitions
		if inStyles && strings.HasPrefix(line, "Style:") {
			continue
		}

		// Detect end of styles section
		if inStyles && strings.HasPrefix(line, "[") && !strings.HasPrefix(line, "[V4") {
			inStyles = false
		}

		fmt.Fprintln(outFile, line)
	}

	return scanner.Err()
}

type Entry struct {
	index     int
	startTime string
	endTime   string
	text      string
}

// writeStyledASS writes a complete ASS file with theme
func (p *SubtitleStylerPlugin) writeStyledASS(outputPath string, theme SubtitleTheme, entries []Entry) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header
	fmt.Fprintln(file, "[Script Info]")
	fmt.Fprintln(file, "Title: Styled Subtitles - Akashic Scribe")
	fmt.Fprintln(file, "ScriptType: v4.00+")
	fmt.Fprintln(file, "WrapStyle: 0")
	fmt.Fprintln(file, "ScaledBorderAndShadow: yes")
	fmt.Fprintln(file, "YCbCr Matrix: TV.601")
	fmt.Fprintln(file, "PlayResX: 1920")
	fmt.Fprintln(file, "PlayResY: 1080")
	fmt.Fprintln(file)

	// Write styles
	fmt.Fprintln(file, "[V4+ Styles]")
	fmt.Fprintln(file, "Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding")
	fmt.Fprintln(file, p.formatASSStyle(theme))
	fmt.Fprintln(file)

	// Write events
	fmt.Fprintln(file, "[Events]")
	fmt.Fprintln(file, "Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text")

	for _, entry := range entries {
		fmt.Fprintf(file, "Dialogue: 0,%s,%s,Default,,0,0,0,,%s\n",
			entry.startTime, entry.endTime, entry.text)
	}

	return nil
}

// formatASSStyle formats a theme as an ASS style line
func (p *SubtitleStylerPlugin) formatASSStyle(theme SubtitleTheme) string {
	bold := 0
	if theme.Bold {
		bold = -1
	}
	italic := 0
	if theme.Italic {
		italic = -1
	}

	return fmt.Sprintf("Style: Default,%s,%d,%s,%s,%s,%s,%d,%d,0,0,100,100,0,0,%d,%.1f,%.1f,%d,%d,%d,%d,1",
		theme.FontName,
		theme.FontSize,
		theme.PrimaryColor,
		theme.SecondaryColor,
		theme.OutlineColor,
		theme.BackColor,
		bold,
		italic,
		theme.BorderStyle,
		theme.Outline,
		theme.Shadow,
		theme.Alignment,
		theme.MarginL,
		theme.MarginR,
		theme.MarginV,
	)
}

// srtTimeToASS converts SRT time format to ASS format
func (p *SubtitleStylerPlugin) srtTimeToASS(srtTime string) string {
	// SRT: HH:MM:SS,mmm -> ASS: H:MM:SS.cc
	srtTime = strings.ReplaceAll(srtTime, ",", ".")
	parts := strings.Split(srtTime, ":")
	if len(parts) != 3 {
		return "0:00:00.00"
	}

	hours := parts[0]
	// Remove leading zero from hours
	if strings.HasPrefix(hours, "0") && len(hours) > 1 {
		hours = hours[1:]
	}

	secondsParts := strings.Split(parts[2], ".")
	milliseconds := "00"
	if len(secondsParts) > 1 {
		// Convert milliseconds to centiseconds
		ms := secondsParts[1]
		if len(ms) >= 2 {
			milliseconds = ms[:2]
		}
	}

	return fmt.Sprintf("%s:%s:%s.%s", hours, parts[1], secondsParts[0], milliseconds)
}

// GetSupportedFormats returns supported subtitle formats
func (p *SubtitleStylerPlugin) GetSupportedFormats() []string {
	return p.supportedFormats
}

// GetAvailableThemes returns list of available theme names
func (p *SubtitleStylerPlugin) GetAvailableThemes() []string {
	themes := make([]string, 0, len(p.themes))
	for name := range p.themes {
		themes = append(themes, name)
	}
	return themes
}

// GetTheme returns a specific theme
func (p *SubtitleStylerPlugin) GetTheme(name string) (SubtitleTheme, error) {
	theme, exists := p.themes[name]
	if !exists {
		return SubtitleTheme{}, fmt.Errorf("theme '%s' not found", name)
	}
	return theme, nil
}

// AddCustomTheme allows adding custom themes at runtime
func (p *SubtitleStylerPlugin) AddCustomTheme(name string, theme SubtitleTheme) {
	theme.Name = name
	p.themes[name] = theme
	if ctx := p.GetContext(); ctx != nil {
		ctx.LogInfo(fmt.Sprintf("Added custom theme: %s", name))
	}
}

// ApplyTheme is a convenience method to apply a specific theme
func (p *SubtitleStylerPlugin) ApplyTheme(inputPath, outputPath, themeName string) error {
	return p.ProcessSubtitles(inputPath, outputPath, map[string]interface{}{
		"theme": themeName,
	})
}

// Ensure SubtitleStylerPlugin implements the SubtitleProcessor interface
var _ core.SubtitleProcessor = (*SubtitleStylerPlugin)(nil)
