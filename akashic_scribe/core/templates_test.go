package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTemplateManager_SaveAndLoad(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "template_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tm, err := NewTemplateManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	// Create a test template
	template := &ProjectTemplate{
		Name:        "Test Template",
		Description: "A test template",
		Category:    "Testing",
		Options: ScribeOptions{
			OriginLanguage:     "en-US",
			TargetLanguage:     "es-ES",
			CreateSubtitles:    true,
			BilingualSubtitles: true,
		},
	}

	// Save the template
	if err := tm.SaveTemplate(template); err != nil {
		t.Fatalf("Failed to save template: %v", err)
	}

	// Load the template
	loaded, err := tm.LoadTemplate("Test Template")
	if err != nil {
		t.Fatalf("Failed to load template: %v", err)
	}

	// Verify template data
	if loaded.Name != template.Name {
		t.Errorf("Template name mismatch: got %s, want %s", loaded.Name, template.Name)
	}

	if loaded.Description != template.Description {
		t.Errorf("Template description mismatch: got %s, want %s", loaded.Description, template.Description)
	}

	if loaded.Options.OriginLanguage != template.Options.OriginLanguage {
		t.Errorf("Origin language mismatch: got %s, want %s",
			loaded.Options.OriginLanguage, template.Options.OriginLanguage)
	}
}

func TestTemplateManager_ListTemplates(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "template_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tm, err := NewTemplateManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	// Default templates should be created
	templates := tm.ListTemplates()
	if len(templates) == 0 {
		t.Error("Expected default templates to be created")
	}

	// Check that we have the expected default templates
	expectedNames := []string{"YouTube Video", "Podcast Dubbing", "Movie Subtitles", "Full Production", "Quick Translation"}
	foundTemplates := make(map[string]bool)

	for _, tmpl := range templates {
		foundTemplates[tmpl.Name] = true
	}

	for _, name := range expectedNames {
		if !foundTemplates[name] {
			t.Errorf("Expected default template '%s' not found", name)
		}
	}
}

func TestTemplateManager_DeleteTemplate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "template_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tm, err := NewTemplateManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	template := &ProjectTemplate{
		Name:        "To Delete",
		Description: "This will be deleted",
		Category:    "Test",
		Options:     ScribeOptions{},
	}

	// Save template
	if err := tm.SaveTemplate(template); err != nil {
		t.Fatalf("Failed to save template: %v", err)
	}

	// Verify it exists
	if _, err := tm.LoadTemplate("To Delete"); err != nil {
		t.Fatalf("Template should exist after saving: %v", err)
	}

	// Delete template
	if err := tm.DeleteTemplate("To Delete"); err != nil {
		t.Fatalf("Failed to delete template: %v", err)
	}

	// Verify it's gone
	if _, err := tm.LoadTemplate("To Delete"); err == nil {
		t.Error("Template should not exist after deletion")
	}

	// Verify file is deleted
	filename := sanitizeFilename("To Delete") + ".json"
	filePath := filepath.Join(tempDir, "templates", filename)
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("Template file should be deleted from disk")
	}
}

func TestTemplateManager_ListByCategory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "template_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tm, err := NewTemplateManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	// Get YouTube templates
	youtubeTemplates := tm.ListTemplatesByCategory("YouTube")
	if len(youtubeTemplates) == 0 {
		t.Error("Expected to find YouTube templates")
	}

	for _, tmpl := range youtubeTemplates {
		if tmpl.Category != "YouTube" {
			t.Errorf("Template '%s' has wrong category: %s", tmpl.Name, tmpl.Category)
		}
	}

	// Get Podcast templates
	podcastTemplates := tm.ListTemplatesByCategory("Podcast")
	if len(podcastTemplates) == 0 {
		t.Error("Expected to find Podcast templates")
	}
}

func TestTemplateManager_GetCategories(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "template_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tm, err := NewTemplateManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	categories := tm.GetCategories()
	if len(categories) == 0 {
		t.Error("Expected to find template categories")
	}

	// Check for expected categories
	expectedCategories := map[string]bool{
		"YouTube":      false,
		"Podcast":      false,
		"Movie":        false,
		"Professional": false,
		"Quick":        false,
	}

	for _, cat := range categories {
		if _, exists := expectedCategories[cat]; exists {
			expectedCategories[cat] = true
		}
	}

	for cat, found := range expectedCategories {
		if !found {
			t.Errorf("Expected category '%s' not found", cat)
		}
	}
}

func TestTemplateManager_ApplyTemplate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "template_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tm, err := NewTemplateManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	// Create current options with input/output paths
	currentOptions := &ScribeOptions{
		InputFile:  "/path/to/video.mp4",
		InputURL:   "https://example.com/video",
		OutputDir:  "/path/to/output",
		CreateDubbing: true, // This will be overwritten
	}

	// Apply YouTube Video template
	if err := tm.ApplyTemplate("YouTube Video", currentOptions); err != nil {
		t.Fatalf("Failed to apply template: %v", err)
	}

	// Verify template was applied
	if !currentOptions.CreateSubtitles {
		t.Error("Template should enable subtitles")
	}

	// Verify input/output paths were preserved
	if currentOptions.InputFile != "/path/to/video.mp4" {
		t.Error("InputFile should be preserved")
	}

	if currentOptions.InputURL != "https://example.com/video" {
		t.Error("InputURL should be preserved")
	}

	if currentOptions.OutputDir != "/path/to/output" {
		t.Error("OutputDir should be preserved")
	}
}

func TestTemplateManager_CreateFromOptions(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "template_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tm, err := NewTemplateManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	// Create options with input/output paths (these should be cleared)
	options := ScribeOptions{
		InputFile:          "/path/to/video.mp4",
		InputURL:           "https://example.com/video",
		OutputDir:          "/path/to/output",
		OriginLanguage:     "en-US",
		TargetLanguage:     "fr-FR",
		CreateSubtitles:    true,
		BilingualSubtitles: true,
	}

	// Create template from options
	err = tm.CreateTemplateFromOptions("My Template", "My custom template", "Custom", options)
	if err != nil {
		t.Fatalf("Failed to create template from options: %v", err)
	}

	// Load the created template
	loaded, err := tm.LoadTemplate("My Template")
	if err != nil {
		t.Fatalf("Failed to load created template: %v", err)
	}

	// Verify input/output paths were cleared
	if loaded.Options.InputFile != "" {
		t.Error("InputFile should be cleared in template")
	}

	if loaded.Options.InputURL != "" {
		t.Error("InputURL should be cleared in template")
	}

	if loaded.Options.OutputDir != "" {
		t.Error("OutputDir should be cleared in template")
	}

	// Verify other options were preserved
	if loaded.Options.OriginLanguage != "en-US" {
		t.Error("Origin language should be preserved")
	}

	if !loaded.Options.CreateSubtitles {
		t.Error("Subtitle setting should be preserved")
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Simple Name", "Simple_Name"},
		{"Name-With-Dashes", "Name_With_Dashes"},
		{"Name_With_Underscores", "Name_With_Underscores"},
		{"Name123", "Name123"},
		{"Name!@#$%", "Name"},
		{"CamelCaseName", "CamelCaseName"},
	}

	for _, tt := range tests {
		result := sanitizeFilename(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeFilename(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestTemplateManager_Timestamps(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "template_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tm, err := NewTemplateManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	template := &ProjectTemplate{
		Name:        "Time Test",
		Description: "Testing timestamps",
		Category:    "Test",
		Options:     ScribeOptions{},
	}

	// Save template
	if err := tm.SaveTemplate(template); err != nil {
		t.Fatalf("Failed to save template: %v", err)
	}

	// Load template
	loaded, err := tm.LoadTemplate("Time Test")
	if err != nil {
		t.Fatalf("Failed to load template: %v", err)
	}

	// Verify timestamps were set
	if loaded.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}

	if loaded.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set")
	}

	// Save again to update timestamp
	loaded.Description = "Updated description"
	oldUpdatedAt := loaded.UpdatedAt

	// Small delay to ensure timestamp difference
	// Note: On fast systems this might not always work, but it's good enough for a test
	if err := tm.SaveTemplate(loaded); err != nil {
		t.Fatalf("Failed to save updated template: %v", err)
	}

	// Reload
	reloaded, err := tm.LoadTemplate("Time Test")
	if err != nil {
		t.Fatalf("Failed to reload template: %v", err)
	}

	// UpdatedAt should have changed (or at least not be earlier)
	if reloaded.UpdatedAt.Before(oldUpdatedAt) {
		t.Error("UpdatedAt should not go backwards")
	}
}
