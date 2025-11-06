package elevenlabs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"akashic_scribe/core"
)

// ElevenLabsPlugin provides ElevenLabs AI voice synthesis
type ElevenLabsPlugin struct {
	*core.BasePlugin
	apiKey     string
	apiBaseURL string
	client     *http.Client
}

// VoiceInfo represents an ElevenLabs voice
type VoiceInfo struct {
	VoiceID       string            `json:"voice_id"`
	Name          string            `json:"name"`
	Category      string            `json:"category"`
	Description   string            `json:"description"`
	PreviewURL    string            `json:"preview_url"`
	Labels        map[string]string `json:"labels"`
	AvailableFor  []string          `json:"available_for_tiers"`
}

// TTSRequest represents a text-to-speech request
type TTSRequest struct {
	Text                 string  `json:"text"`
	ModelID              string  `json:"model_id"`
	VoiceSettings        *VoiceSettings `json:"voice_settings,omitempty"`
}

// VoiceSettings for fine-tuning voice output
type VoiceSettings struct {
	Stability       float64 `json:"stability"`        // 0.0 to 1.0
	SimilarityBoost float64 `json:"similarity_boost"` // 0.0 to 1.0
	Style           float64 `json:"style,omitempty"`  // 0.0 to 1.0 (v2 models)
	UseSpeakerBoost bool    `json:"use_speaker_boost,omitempty"`
}

// NewElevenLabsPlugin creates a new ElevenLabs plugin
func NewElevenLabsPlugin() core.Plugin {
	base := core.NewBasePlugin(
		"voidcat.elevenlabs",
		"ElevenLabs Voice Synthesis",
		"1.0.0",
		"Generate natural-sounding speech using ElevenLabs AI voice models with advanced customization",
		"VoidCat RDC",
	)

	plugin := &ElevenLabsPlugin{
		BasePlugin: base,
		apiBaseURL: "https://api.elevenlabs.io/v1",
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}

	return plugin
}

// GetCapabilities returns the plugin's capabilities
func (p *ElevenLabsPlugin) GetCapabilities() []core.PluginCapability {
	return []core.PluginCapability{
		core.CapabilityVoiceSynthesis,
		core.CapabilityAPIIntegration,
	}
}

// Initialize sets up the plugin with API credentials
func (p *ElevenLabsPlugin) Initialize(context core.PluginContext) error {
	if err := p.BasePlugin.Initialize(context); err != nil {
		return err
	}

	// Get API key from environment or config
	p.apiKey = os.Getenv("ELEVENLABS_API_KEY")
	if p.apiKey == "" {
		// Try to get from config
		config := context.GetConfig()
		if key, ok := config["api_key"].(string); ok {
			p.apiKey = key
		}
	}

	if p.apiKey == "" {
		context.LogWarning("ElevenLabs API key not set. Set ELEVENLABS_API_KEY environment variable or configure via context.")
		// Don't fail initialization - plugin can still be used to query voices etc
	} else {
		context.LogInfo("ElevenLabs API key configured successfully")
	}

	return nil
}

// HealthCheck verifies plugin is functional
func (p *ElevenLabsPlugin) HealthCheck() error {
	if err := p.BasePlugin.HealthCheck(); err != nil {
		return err
	}

	if p.apiKey == "" {
		return fmt.Errorf("API key not configured")
	}

	return nil
}

// GenerateSpeech generates speech from text using ElevenLabs TTS
// Options:
//   - "voice_id": string - Voice ID (required)
//   - "model_id": string - Model ID (default: "eleven_monolingual_v1")
//   - "stability": float64 - Voice stability 0.0-1.0 (default: 0.5)
//   - "similarity_boost": float64 - Similarity boost 0.0-1.0 (default: 0.75)
//   - "style": float64 - Style intensity 0.0-1.0 (default: 0.0)
//   - "speaker_boost": bool - Enable speaker boost (default: true)
func (p *ElevenLabsPlugin) GenerateSpeech(text string, outputPath string, options map[string]interface{}) error {
	ctx := p.GetContext()

	if p.apiKey == "" {
		return fmt.Errorf("API key not configured")
	}

	// Extract voice ID (required)
	voiceID, ok := options["voice_id"].(string)
	if !ok || voiceID == "" {
		return fmt.Errorf("voice_id is required")
	}

	// Extract model ID
	modelID := "eleven_monolingual_v1"
	if model, ok := options["model_id"].(string); ok {
		modelID = model
	}

	// Extract voice settings
	settings := &VoiceSettings{
		Stability:       0.5,
		SimilarityBoost: 0.75,
		UseSpeakerBoost: true,
	}

	if stability, ok := options["stability"].(float64); ok {
		settings.Stability = stability
	}
	if similarity, ok := options["similarity_boost"].(float64); ok {
		settings.SimilarityBoost = similarity
	}
	if style, ok := options["style"].(float64); ok {
		settings.Style = style
	}
	if boost, ok := options["speaker_boost"].(bool); ok {
		settings.UseSpeakerBoost = boost
	}

	ctx.LogInfo(fmt.Sprintf("Generating speech with voice ID: %s, model: %s", voiceID, modelID))

	// Create request
	request := TTSRequest{
		Text:          text,
		ModelID:       modelID,
		VoiceSettings: settings,
	}

	// Make API call
	audioData, err := p.callTTSAPI(voiceID, request)
	if err != nil {
		return fmt.Errorf("TTS API call failed: %w", err)
	}

	// Save to file
	if err := os.WriteFile(outputPath, audioData, 0644); err != nil {
		return fmt.Errorf("failed to save audio file: %w", err)
	}

	ctx.LogInfo(fmt.Sprintf("Speech generated successfully: %s", outputPath))
	return nil
}

// callTTSAPI makes the API request to ElevenLabs
func (p *ElevenLabsPlugin) callTTSAPI(voiceID string, request TTSRequest) ([]byte, error) {
	// Create request body
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/text-to-speech/%s", p.apiBaseURL, voiceID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "audio/mpeg")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", p.apiKey)

	// Make request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	// Read audio data
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return audioData, nil
}

// GetVoices retrieves available voices from ElevenLabs
func (p *ElevenLabsPlugin) GetVoices() ([]VoiceInfo, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("API key not configured")
	}

	ctx := p.GetContext()
	ctx.LogInfo("Fetching available voices from ElevenLabs")

	// Create request
	url := fmt.Sprintf("%s/voices", p.apiBaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("xi-api-key", p.apiKey)

	// Make request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var response struct {
		Voices []VoiceInfo `json:"voices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	ctx.LogInfo(fmt.Sprintf("Retrieved %d voices", len(response.Voices)))
	return response.Voices, nil
}

// GetVoiceByName finds a voice by name
func (p *ElevenLabsPlugin) GetVoiceByName(name string) (*VoiceInfo, error) {
	voices, err := p.GetVoices()
	if err != nil {
		return nil, err
	}

	for _, voice := range voices {
		if voice.Name == name {
			return &voice, nil
		}
	}

	return nil, fmt.Errorf("voice '%s' not found", name)
}

// CloneVoice creates a new voice from audio samples
func (p *ElevenLabsPlugin) CloneVoice(name, description string, audioFiles []string, labels map[string]string) (*VoiceInfo, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("API key not configured")
	}

	ctx := p.GetContext()
	ctx.LogInfo(fmt.Sprintf("Cloning voice: %s", name))

	// This is a simplified version - actual implementation would handle multipart form data
	// For now, return an error indicating this needs additional implementation
	return nil, fmt.Errorf("voice cloning requires multipart form data upload - implementation pending")
}

// GetModels returns available TTS models
func (p *ElevenLabsPlugin) GetModels() []string {
	return []string{
		"eleven_monolingual_v1",      // English only, fast
		"eleven_multilingual_v1",     // Multiple languages
		"eleven_multilingual_v2",     // Improved multilingual with style
		"eleven_turbo_v2",           // Fastest, low latency
	}
}

// GenerateWithPreset generates speech using a preset configuration
func (p *ElevenLabsPlugin) GenerateWithPreset(text, voiceID, preset, outputPath string) error {
	var options map[string]interface{}

	switch preset {
	case "natural":
		options = map[string]interface{}{
			"voice_id":         voiceID,
			"model_id":         "eleven_multilingual_v2",
			"stability":        0.5,
			"similarity_boost": 0.75,
			"style":            0.0,
			"speaker_boost":    true,
		}
	case "expressive":
		options = map[string]interface{}{
			"voice_id":         voiceID,
			"model_id":         "eleven_multilingual_v2",
			"stability":        0.3,
			"similarity_boost": 0.8,
			"style":            0.6,
			"speaker_boost":    true,
		}
	case "stable":
		options = map[string]interface{}{
			"voice_id":         voiceID,
			"model_id":         "eleven_monolingual_v1",
			"stability":        0.8,
			"similarity_boost": 0.5,
			"speaker_boost":    false,
		}
	case "fast":
		options = map[string]interface{}{
			"voice_id":         voiceID,
			"model_id":         "eleven_turbo_v2",
			"stability":        0.5,
			"similarity_boost": 0.75,
			"speaker_boost":    true,
		}
	default:
		return fmt.Errorf("unknown preset: %s", preset)
	}

	return p.GenerateSpeech(text, outputPath, options)
}

// ProcessSubtitles generates audio for subtitle file
func (p *ElevenLabsPlugin) ProcessSubtitles(inputPath, outputDir string, options map[string]interface{}) error {
	ctx := p.GetContext()
	ctx.LogInfo(fmt.Sprintf("Processing subtitle file: %s", inputPath))

	// This would parse the subtitle file and generate audio for each entry
	// For now, return a helpful error
	return fmt.Errorf("subtitle processing requires subtitle parser integration - use with Format Converter plugin")
}

// GetPresets returns available generation presets
func (p *ElevenLabsPlugin) GetPresets() []string {
	return []string{
		"natural",     // Balanced, natural speech
		"expressive",  // More emotion and variation
		"stable",      // Consistent, predictable
		"fast",        // Fastest generation
	}
}

// SetAPIKey updates the API key at runtime
func (p *ElevenLabsPlugin) SetAPIKey(apiKey string) {
	p.apiKey = apiKey
	if ctx := p.GetContext(); ctx != nil {
		ctx.SetConfig("api_key", apiKey)
		ctx.LogInfo("API key updated")
	}
}

// GetUsage retrieves API usage statistics
func (p *ElevenLabsPlugin) GetUsage() (map[string]interface{}, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("API key not configured")
	}

	url := fmt.Sprintf("%s/user/subscription", p.apiBaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("xi-api-key", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var usage map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&usage); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return usage, nil
}

// SaveVoiceToCache saves frequently used voice info to cache
func (p *ElevenLabsPlugin) SaveVoiceToCache(voice VoiceInfo) error {
	ctx := p.GetContext()
	cacheDir := ctx.GetPluginCacheDir()
	cachePath := filepath.Join(cacheDir, fmt.Sprintf("voice_%s.json", voice.VoiceID))

	data, err := json.Marshal(voice)
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}

// LoadVoiceFromCache loads voice info from cache
func (p *ElevenLabsPlugin) LoadVoiceFromCache(voiceID string) (*VoiceInfo, error) {
	ctx := p.GetContext()
	cacheDir := ctx.GetPluginCacheDir()
	cachePath := filepath.Join(cacheDir, fmt.Sprintf("voice_%s.json", voiceID))

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var voice VoiceInfo
	if err := json.Unmarshal(data, &voice); err != nil {
		return nil, err
	}

	return &voice, nil
}
