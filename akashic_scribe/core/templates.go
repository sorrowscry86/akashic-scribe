package core

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ProjectTemplate represents a saved configuration template.
type ProjectTemplate struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Options     ScribeOptions `json:"options"`
	Category    string        `json:"category"` // e.g., "YouTube", "Podcast", "Movie", "Custom"
}

// TemplateManager handles saving, loading, and managing project templates.
type TemplateManager struct {
	templatesDir string
	templates    map[string]*ProjectTemplate
}

// NewTemplateManager creates a new template manager.
func NewTemplateManager(configDir string) (*TemplateManager, error) {
	templatesDir := filepath.Join(configDir, "templates")

	// Create templates directory if it doesn't exist
	if err := os.MkdirAll(templatesDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create templates directory: %w", err)
	}

	tm := &TemplateManager{
		templatesDir: templatesDir,
		templates:    make(map[string]*ProjectTemplate),
	}

	// Load existing templates
	if err := tm.LoadAll(); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	// Create default templates if none exist
	if len(tm.templates) == 0 {
		if err := tm.createDefaultTemplates(); err != nil {
			return nil, fmt.Errorf("failed to create default templates: %w", err)
		}
	}

	return tm, nil
}

// SaveTemplate saves a template to disk.
func (tm *TemplateManager) SaveTemplate(template *ProjectTemplate) error {
	if template.Name == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	// Update timestamps
	now := time.Now()
	if template.CreatedAt.IsZero() {
		template.CreatedAt = now
	}
	template.UpdatedAt = now

	// Save to disk
	filename := sanitizeFilename(template.Name) + ".json"
	filePath := filepath.Join(tm.templatesDir, filename)

	data, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal template: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write template file: %w", err)
	}

	// Add to in-memory cache
	tm.templates[template.Name] = template

	return nil
}

// LoadTemplate loads a template by name.
func (tm *TemplateManager) LoadTemplate(name string) (*ProjectTemplate, error) {
	// Check cache first
	if template, exists := tm.templates[name]; exists {
		return template, nil
	}

	// Try to load from disk
	filename := sanitizeFilename(name) + ".json"
	filePath := filepath.Join(tm.templatesDir, filename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	var template ProjectTemplate
	if err := json.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Add to cache
	tm.templates[template.Name] = &template

	return &template, nil
}

// LoadAll loads all templates from disk.
func (tm *TemplateManager) LoadAll() error {
	entries, err := os.ReadDir(tm.templatesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Directory doesn't exist yet, that's fine
		}
		return fmt.Errorf("failed to read templates directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(tm.templatesDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Warning: Could not read template file %s: %v", filePath, err)
			continue // Skip files we can't read
		}

		var template ProjectTemplate
		if err := json.Unmarshal(data, &template); err != nil {
			log.Printf("Warning: Could not parse template file %s: %v", filePath, err)
			continue // Skip invalid templates
		}

		tm.templates[template.Name] = &template
	}

	return nil
}

// DeleteTemplate deletes a template by name.
func (tm *TemplateManager) DeleteTemplate(name string) error {
	filename := sanitizeFilename(name) + ".json"
	filePath := filepath.Join(tm.templatesDir, filename)

	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	delete(tm.templates, name)
	return nil
}

// ListTemplates returns all available templates.
func (tm *TemplateManager) ListTemplates() []*ProjectTemplate {
	templates := make([]*ProjectTemplate, 0, len(tm.templates))
	for _, template := range tm.templates {
		templates = append(templates, template)
	}
	return templates
}

// ListTemplatesByCategory returns templates filtered by category.
func (tm *TemplateManager) ListTemplatesByCategory(category string) []*ProjectTemplate {
	templates := make([]*ProjectTemplate, 0)
	for _, template := range tm.templates {
		if template.Category == category {
			templates = append(templates, template)
		}
	}
	return templates
}

// GetCategories returns all unique template categories.
func (tm *TemplateManager) GetCategories() []string {
	categoryMap := make(map[string]bool)
	for _, template := range tm.templates {
		if template.Category != "" {
			categoryMap[template.Category] = true
		}
	}

	categories := make([]string, 0, len(categoryMap))
	for category := range categoryMap {
		categories = append(categories, category)
	}
	return categories
}

// createDefaultTemplates creates a set of useful default templates.
func (tm *TemplateManager) createDefaultTemplates() error {
	defaults := []*ProjectTemplate{
		{
			Name:        "YouTube Video",
			Description: "Standard settings for YouTube video translation with subtitles",
			Category:    "YouTube",
			Options: ScribeOptions{
				OriginLanguage:     "en-US",
				TargetLanguage:     "es-ES",
				CreateSubtitles:    true,
				BilingualSubtitles: true,
				SubtitlePosition:   "bottom",
				CreateDubbing:      false,
			},
		},
		{
			Name:        "Podcast Dubbing",
			Description: "High-quality audio dubbing for podcasts",
			Category:    "Podcast",
			Options: ScribeOptions{
				OriginLanguage:  "en-US",
				TargetLanguage:  "ja-JP",
				CreateSubtitles: false,
				CreateDubbing:   true,
				VoiceModel:      "alloy",
				VoiceSpeed:      1.0,
				AudioFormat:     "mp3",
				AudioQuality:    "high",
				AudioSampleRate: 44100,
				AudioBitRate:    192,
				AudioChannels:   2,
				NormalizeAudio:  true,
				RemoveSilence:   true,
			},
		},
		{
			Name:        "Movie Subtitles",
			Description: "Professional bilingual subtitles for movies",
			Category:    "Movie",
			Options: ScribeOptions{
				OriginLanguage:     "en-US",
				TargetLanguage:     "fr-FR",
				CreateSubtitles:    true,
				BilingualSubtitles: true,
				SubtitlePosition:   "bottom",
				CreateDubbing:      false,
			},
		},
		{
			Name:        "Full Production",
			Description: "Complete dubbing and subtitles for professional projects",
			Category:    "Professional",
			Options: ScribeOptions{
				OriginLanguage:     "en-US",
				TargetLanguage:     "de-DE",
				CreateSubtitles:    true,
				BilingualSubtitles: false,
				CreateDubbing:      true,
				VoiceModel:         "nova",
				VoiceSpeed:         1.0,
				AudioFormat:        "flac",
				AudioQuality:       "lossless",
				AudioSampleRate:    48000,
				AudioBitRate:       320,
				AudioChannels:      2,
				NormalizeAudio:     true,
				RemoveSilence:      false,
			},
		},
		{
			Name:        "Quick Translation",
			Description: "Fast translation with basic subtitles",
			Category:    "Quick",
			Options: ScribeOptions{
				OriginLanguage:     "en-US",
				TargetLanguage:     "es-ES",
				CreateSubtitles:    true,
				BilingualSubtitles: false,
				CreateDubbing:      false,
			},
		},
	}

	for _, template := range defaults {
		if err := tm.SaveTemplate(template); err != nil {
			return err
		}
	}

	return nil
}

// ApplyTemplate applies a template to ScribeOptions, preserving input/output paths.
func (tm *TemplateManager) ApplyTemplate(templateName string, currentOptions *ScribeOptions) error {
	template, err := tm.LoadTemplate(templateName)
	if err != nil {
		return err
	}

	// Save current input/output settings
	savedInput := currentOptions.InputFile
	savedURL := currentOptions.InputURL
	savedOutput := currentOptions.OutputDir

	// Apply template options
	*currentOptions = template.Options

	// Restore input/output settings
	currentOptions.InputFile = savedInput
	currentOptions.InputURL = savedURL
	currentOptions.OutputDir = savedOutput

	return nil
}

// CreateTemplateFromOptions creates a new template from current options.
func (tm *TemplateManager) CreateTemplateFromOptions(name, description, category string, options ScribeOptions) error {
	// Clear input/output paths from template (these are job-specific)
	options.InputFile = ""
	options.InputURL = ""
	options.OutputDir = ""

	template := &ProjectTemplate{
		Name:        name,
		Description: description,
		Category:    category,
		Options:     options,
	}

	return tm.SaveTemplate(template)
}

// sanitizeFilename removes or replaces characters that are invalid in filenames.
// Uses strings.Builder for performance and adds a hash suffix to prevent collisions.
func sanitizeFilename(name string) string {
	var builder strings.Builder
	builder.Grow(len(name) + 9) // Pre-allocate for name + "_" + 8-char hash

	// Allow alphanumeric, spaces, dashes, underscores, dots, and parentheses
	// Replace other characters with underscores
	for _, ch := range name {
		switch {
		case ch >= 'a' && ch <= 'z':
			builder.WriteRune(ch)
		case ch >= 'A' && ch <= 'Z':
			builder.WriteRune(ch)
		case ch >= '0' && ch <= '9':
			builder.WriteRune(ch)
		case ch == ' ' || ch == '-' || ch == '_':
			builder.WriteRune('_')
		case ch == '.' || ch == '(' || ch == ')':
			builder.WriteRune(ch) // Allow these safe characters
		default:
			builder.WriteRune('_') // Replace invalid chars with underscore
		}
	}

	// Add hash suffix to prevent collisions (e.g., "My Template!" and "My Template?" both become different)
	hash := md5.Sum([]byte(name))
	builder.WriteString(fmt.Sprintf("_%x", hash[:4])) // Use first 4 bytes (8 hex chars)

	return builder.String()
}
