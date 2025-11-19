package core

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// realScribeEngine is the actual implementation of the ScribeEngine interface.
// It uses external tools like yt-dlp and ffmpeg to process video and audio.
type realScribeEngine struct{}

// NewRealScribeEngine creates a new instance of the realScribeEngine.
func NewRealScribeEngine() ScribeEngine {
	return &realScribeEngine{}
}

// checkDependencies verifies that required external tools are available.
func (e *realScribeEngine) checkDependencies() error {
	// Check for yt-dlp
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		return fmt.Errorf("yt-dlp not found: %w. Please install yt-dlp to process video URLs", err)
	}

	// Check for ffmpeg
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found: %w. Please install ffmpeg to process video/audio files", err)
	}

	return nil
}

// parseYtDlpProgress parses yt-dlp output to extract download progress.
// yt-dlp outputs progress like: "[download]  45.2% of 100.00MiB at 1.50MiB/s ETA 00:36"
func parseYtDlpProgress(line string) (float64, bool) {
	// Match patterns like "[download]  45.2%" or "[download] 100%"
	re := regexp.MustCompile(`\[download\]\s+(\d+\.?\d*)%`)
	matches := re.FindStringSubmatch(line)
	if len(matches) >= 2 {
		if percent, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return percent / 100.0, true
		}
	}
	return 0, false
}

// parseFfmpegProgress parses ffmpeg output to extract encoding progress.
// ffmpeg outputs progress like: "time=00:01:23.45 bitrate=1234.5kbits/s speed=2.5x"
func parseFfmpegProgress(line string, duration float64) (float64, bool) {
	// Match time= pattern
	re := regexp.MustCompile(`time=(\d+):(\d+):(\d+\.?\d*)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) >= 4 {
		hours, _ := strconv.ParseFloat(matches[1], 64)
		minutes, _ := strconv.ParseFloat(matches[2], 64)
		seconds, _ := strconv.ParseFloat(matches[3], 64)
		currentTime := hours*3600 + minutes*60 + seconds

		if duration > 0 {
			progress := currentTime / duration
			if progress > 1.0 {
				progress = 1.0
			}
			return progress, true
		}
	}
	return 0, false
}

// runCommandWithProgress runs a command and reports progress via a channel.
// It monitors stdout/stderr and calls the progress parser function.
// The command can be cancelled via the context.
func (e *realScribeEngine) runCommandWithProgress(
	ctx context.Context,
	cmd *exec.Cmd,
	baseProgress float64,
	progressRange float64,
	progressChan chan<- ProgressUpdate,
	statusMsg string,
	parseFunc func(string) (float64, bool),
) error {
	// Create pipes for stderr (where yt-dlp and ffmpeg output progress)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Also capture stdout for some tools
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Channel to signal when command completes
	done := make(chan error, 1)

	// Read progress from both stdout and stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			if parseFunc != nil {
				if progress, ok := parseFunc(line); ok {
					actualProgress := baseProgress + (progress * progressRange)
					select {
					case progressChan <- ProgressUpdate{actualProgress, statusMsg}:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			if parseFunc != nil {
				if progress, ok := parseFunc(line); ok {
					actualProgress := baseProgress + (progress * progressRange)
					select {
					case progressChan <- ProgressUpdate{actualProgress, statusMsg}:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	// Wait for command to complete in goroutine
	go func() {
		done <- cmd.Wait()
	}()

	// Wait for either completion or cancellation
	select {
	case <-ctx.Done():
		// Context cancelled - kill the process
		if cmd.Process != nil {
			if err := cmd.Process.Kill(); err != nil {
				log.Printf("Warning: failed to kill process: %v", err)
			}
		}
		// Wait for the command to actually exit
		<-done
		return fmt.Errorf("operation cancelled: %w", ctx.Err())
	case err := <-done:
		// Command completed
		return err
	}
}

// Transcribe takes a video source (local path or URL) and returns the transcription.
func (e *realScribeEngine) Transcribe(videoSource string) (string, error) {
	// Check dependencies
	if err := e.checkDependencies(); err != nil {
		return "", err
	}

	// Create a temporary directory for processing.
	tempDir, err := os.MkdirTemp("", "akashic_scribe_*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			log.Printf("Warning: failed to clean up temp directory %s: %v", tempDir, err)
		}
	}()

	var videoPath string

	// Check if the videoSource is a URL or a local file.
	if strings.HasPrefix(videoSource, "http") {
		// Download the video from the URL using yt-dlp.
		videoPath = filepath.Join(tempDir, "downloaded_video.%(ext)s")
		cmd := exec.Command("yt-dlp", "-o", videoPath, videoSource)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to download video: %w", err)
		}

		// Find the actual downloaded file
		files, err := filepath.Glob(filepath.Join(tempDir, "downloaded_video.*"))
		if err != nil || len(files) == 0 {
			return "", errors.New("downloaded file not found")
		}
		videoPath = files[0]
	} else {
		// Verify local file exists
		if _, err := os.Stat(videoSource); err != nil {
			return "", fmt.Errorf("input file not found: %w", err)
		}
		videoPath = videoSource
	}

	// Extract the audio from the video file using ffmpeg.
	audioPath := filepath.Join(tempDir, "extracted_audio.wav")
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vn", "-acodec", "pcm_s16le", "-ar", "16000", "-ac", "1", audioPath)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to extract audio: %w", err)
	}

	// Verify audio file was created
	if _, err := os.Stat(audioPath); err != nil {
		return "", fmt.Errorf("audio extraction failed - no output file: %w", err)
	}

	// --- Placeholder for actual transcription ---
	// In a real implementation, you would call a transcription service (e.g., Whisper) here.
	// For now, we will just return a dummy transcription based on the file.
	dummyTranscription := fmt.Sprintf("This is a dummy transcription of the video from: %s", filepath.Base(videoPath))

	return dummyTranscription, nil
}

// Translate takes a text and a target language, and returns the translation.
func (e *realScribeEngine) Translate(text string, targetLanguage string) (string, error) {
	// --- Placeholder for actual translation ---
	// In a real implementation, you would call a translation service here.
	// For now, we will just return a dummy translation.
	dummyTranslation := fmt.Sprintf("This is a dummy translation of the text to %s.", targetLanguage)

	return dummyTranslation, nil
}

// setDefaultDubbingParams fills in default values for any unset dubbing parameters.
func setDefaultDubbingParams(opts *ScribeOptions) {
	if opts.VoiceSpeed == 0 {
		opts.VoiceSpeed = 1.0
	}
	if opts.VoiceStability == 0 {
		opts.VoiceStability = 0.5
	}
	if opts.AudioFormat == "" {
		opts.AudioFormat = "mp3"
	}
	if opts.AudioQuality == "" {
		opts.AudioQuality = "high"
	}
	if opts.AudioSampleRate == 0 {
		opts.AudioSampleRate = 44100
	}
	if opts.AudioBitRate == 0 {
		opts.AudioBitRate = 192
	}
	if opts.AudioChannels == 0 {
		opts.AudioChannels = 2
	}
}

// validateDubbingParams validates dubbing parameters and returns an error if invalid.
func validateDubbingParams(opts ScribeOptions) error {
	if opts.CreateDubbing {
		// Validate voice model or custom voice
		if !opts.UseCustomVoice && opts.VoiceModel == "" {
			return errors.New("voice model must be specified when dubbing is enabled")
		}
		if opts.UseCustomVoice && opts.CustomVoicePath == "" {
			return errors.New("custom voice path must be specified when using custom voice")
		}
		if opts.UseCustomVoice {
			if _, err := os.Stat(opts.CustomVoicePath); err != nil {
				return fmt.Errorf("custom voice file not found: %w", err)
			}
		}

		// Validate voice speed
		if opts.VoiceSpeed < 0.25 || opts.VoiceSpeed > 4.0 {
			return fmt.Errorf("voice speed must be between 0.25 and 4.0, got %.2f", opts.VoiceSpeed)
		}

		// Validate voice pitch
		if opts.VoicePitch < -20 || opts.VoicePitch > 20 {
			return fmt.Errorf("voice pitch must be between -20 and 20 semitones, got %.1f", opts.VoicePitch)
		}

		// Validate voice stability
		if opts.VoiceStability < 0 || opts.VoiceStability > 1.0 {
			return fmt.Errorf("voice stability must be between 0.0 and 1.0, got %.2f", opts.VoiceStability)
		}

		// Validate audio format
		validFormats := map[string]bool{"mp3": true, "wav": true, "flac": true, "aac": true, "ogg": true}
		if !validFormats[opts.AudioFormat] {
			return fmt.Errorf("invalid audio format: %s (must be mp3, wav, flac, aac, or ogg)", opts.AudioFormat)
		}

		// Validate audio quality
		validQualities := map[string]bool{"low": true, "medium": true, "high": true, "lossless": true}
		if !validQualities[opts.AudioQuality] {
			return fmt.Errorf("invalid audio quality: %s (must be low, medium, high, or lossless)", opts.AudioQuality)
		}

		// Validate sample rate
		validSampleRates := map[int]bool{8000: true, 16000: true, 22050: true, 44100: true, 48000: true}
		if !validSampleRates[opts.AudioSampleRate] {
			return fmt.Errorf("invalid sample rate: %d Hz (must be 8000, 16000, 22050, 44100, or 48000)", opts.AudioSampleRate)
		}

		// Validate bit rate
		if opts.AudioBitRate < 64 || opts.AudioBitRate > 320 {
			return fmt.Errorf("bit rate must be between 64 and 320 kbps, got %d", opts.AudioBitRate)
		}

		// Validate channels
		if opts.AudioChannels != 1 && opts.AudioChannels != 2 {
			return fmt.Errorf("audio channels must be 1 (mono) or 2 (stereo), got %d", opts.AudioChannels)
		}
	}
	return nil
}

// GenerateDubbing generates dubbed audio from translated text using TTS.
// This function supports both OpenAI TTS and custom voice synthesis.
func (e *realScribeEngine) GenerateDubbing(translation string, opts ScribeOptions, outputDir string) (string, error) {
	// Set default parameters
	setDefaultDubbingParams(&opts)

	// Validate parameters
	if err := validateDubbingParams(opts); err != nil {
		return "", fmt.Errorf("invalid dubbing parameters: %w", err)
	}

	// Create temp file for raw TTS output
	tempDir, err := os.MkdirTemp("", "akashic_scribe_tts_*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			log.Printf("Warning: failed to clean up temp directory %s: %v", tempDir, err)
		}
	}()

	rawAudioPath := filepath.Join(tempDir, "raw_tts.mp3")

	// Generate TTS audio
	if opts.UseCustomVoice {
		// For custom voices, we would need a voice cloning service
		// For now, we'll use a simulated approach with pitch/speed modifications
		return "", errors.New("custom voice synthesis not yet implemented - requires voice cloning service integration")
	} else {
		// Use OpenAI TTS API
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return "", errors.New("OPENAI_API_KEY environment variable not set - required for TTS generation")
		}

		if err := e.generateOpenAITTS(translation, opts.VoiceModel, opts.VoiceSpeed, apiKey, rawAudioPath); err != nil {
			return "", fmt.Errorf("failed to generate TTS audio: %w", err)
		}
	}

	// Process audio with ffmpeg to apply additional effects and format conversion
	finalAudioPath := filepath.Join(outputDir, fmt.Sprintf("dubbed_audio.%s", opts.AudioFormat))
	if err := e.processAudio(rawAudioPath, finalAudioPath, opts); err != nil {
		return "", fmt.Errorf("failed to process audio: %w", err)
	}

	return finalAudioPath, nil
}

// generateOpenAITTS calls the OpenAI TTS API to generate speech audio.
func (e *realScribeEngine) generateOpenAITTS(text, model string, speed float64, apiKey, outputPath string) error {
	// Validate model
	validModels := map[string]bool{
		"alloy":   true,
		"echo":    true,
		"fable":   true,
		"onyx":    true,
		"nova":    true,
		"shimmer": true,
	}
	if !validModels[model] {
		return fmt.Errorf("invalid OpenAI TTS model: %s", model)
	}

	// Prepare API request
	requestBody := map[string]interface{}{
		"model": "tts-1-hd", // Use high-quality TTS model
		"input": text,
		"voice": model,
		"speed": speed,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make API request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/speech", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	// Save audio to file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return fmt.Errorf("failed to save audio: %w", err)
	}

	return nil
}

// processAudio applies audio effects and converts to the desired format using ffmpeg.
func (e *realScribeEngine) processAudio(inputPath, outputPath string, opts ScribeOptions) error {
	// Build ffmpeg command with audio filters
	args := []string{"-i", inputPath}

	// Build filter chain
	var filters []string

	// Apply pitch adjustment if specified
	if opts.VoicePitch != 0 {
		// Convert semitones to frequency ratio: ratio = 2^(semitones/12)
		// For ffmpeg rubberband or asetrate/atempo combination
		filters = append(filters, fmt.Sprintf("asetrate=%d*2^(%.2f/12),aresample=%d", opts.AudioSampleRate, opts.VoicePitch, opts.AudioSampleRate))
	}

	// Apply normalization if requested
	if opts.NormalizeAudio {
		filters = append(filters, "loudnorm")
	}

	// Remove silence if requested
	if opts.RemoveSilence {
		filters = append(filters, "silenceremove=start_periods=1:start_duration=0.1:start_threshold=-50dB:detection=peak,aformat=dblp,areverse,silenceremove=start_periods=1:start_duration=0.1:start_threshold=-50dB:detection=peak,aformat=dblp,areverse")
	}

	// Apply filter chain if any filters were added
	if len(filters) > 0 {
		args = append(args, "-af", strings.Join(filters, ","))
	}

	// Set audio codec based on format
	switch opts.AudioFormat {
	case "mp3":
		args = append(args, "-acodec", "libmp3lame")
	case "wav":
		args = append(args, "-acodec", "pcm_s16le")
	case "flac":
		args = append(args, "-acodec", "flac")
	case "aac":
		args = append(args, "-acodec", "aac")
	case "ogg":
		args = append(args, "-acodec", "libvorbis")
	}

	// Set quality/bitrate based on audio quality setting
	if opts.AudioFormat != "wav" && opts.AudioFormat != "flac" {
		// For lossy formats, set bitrate
		qualityBitrate := opts.AudioBitRate
		if opts.AudioQuality == "low" && qualityBitrate > 128 {
			qualityBitrate = 128
		} else if opts.AudioQuality == "medium" && qualityBitrate > 192 {
			qualityBitrate = 192
		}
		args = append(args, "-b:a", fmt.Sprintf("%dk", qualityBitrate))
	}

	// Set sample rate
	args = append(args, "-ar", fmt.Sprintf("%d", opts.AudioSampleRate))

	// Set channels (mono/stereo)
	args = append(args, "-ac", fmt.Sprintf("%d", opts.AudioChannels))

	// Overwrite output file
	args = append(args, "-y", outputPath)

	// Execute ffmpeg
	cmd := exec.Command("ffmpeg", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg processing failed: %w\nStderr: %s", err, stderr.String())
	}

	return nil
}

// StartProcessing runs the full pipeline and reports progress.
// The operation can be cancelled via the provided context.
func (e *realScribeEngine) StartProcessing(ctx context.Context, options ScribeOptions, progress chan<- ProgressUpdate) error {
	// Check for cancellation before starting
	select {
	case <-ctx.Done():
		return fmt.Errorf("operation cancelled before start: %w", ctx.Err())
	default:
	}

	// Check dependencies first
	if err := e.checkDependencies(); err != nil {
		progress <- ProgressUpdate{0.0, fmt.Sprintf("Dependency check failed: %v", err)}
		return err
	}

	opts := options

	progress <- ProgressUpdate{0.01, "Preparing job..."}

	// Step 1: Download or use local file
	var videoPath string
	if opts.InputFile != "" {
		// Verify file exists
		if _, err := os.Stat(opts.InputFile); err != nil {
			progress <- ProgressUpdate{0.0, fmt.Sprintf("Input file not found: %v", err)}
			return fmt.Errorf("input file not found: %w", err)
		}
		videoPath = opts.InputFile
		progress <- ProgressUpdate{0.05, "Using local video file."}
	} else if opts.InputURL != "" {
		// Check for cancellation before download
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation cancelled: %w", ctx.Err())
		default:
		}

		progress <- ProgressUpdate{0.05, "Starting video download..."}

		// Create temporary directory for downloads
		tempDir, err := os.MkdirTemp("", "akashic_scribe_download_*")
		if err != nil {
			progress <- ProgressUpdate{0.0, fmt.Sprintf("Failed to create temp directory: %v", err)}
			return fmt.Errorf("failed to create temp directory: %w", err)
		}
		defer func() {
			if err := os.RemoveAll(tempDir); err != nil {
				log.Printf("Warning: failed to clean up temp directory %s: %v", tempDir, err)
			}
		}()

		videoPath = filepath.Join(tempDir, "downloaded_video.%(ext)s")
		cmd := exec.Command("yt-dlp", "-o", videoPath, "--newline", opts.InputURL)

		// Use progress tracking for download (0.05 to 0.20 = 15% range)
		if err := e.runCommandWithProgress(ctx, cmd, 0.05, 0.15, progress, "Downloading video...", parseYtDlpProgress); err != nil {
			progress <- ProgressUpdate{0.0, fmt.Sprintf("Failed to download video: %v", err)}
			return fmt.Errorf("failed to download video: %w", err)
		}

		// Find the actual downloaded file
		files, err := filepath.Glob(filepath.Join(tempDir, "downloaded_video.*"))
		if err != nil || len(files) == 0 {
			progress <- ProgressUpdate{0.0, "Downloaded file not found"}
			return errors.New("downloaded file not found")
		}
		videoPath = files[0]
		progress <- ProgressUpdate{0.20, "Download complete"}
	} else {
		progress <- ProgressUpdate{0.0, "No input file or URL provided."}
		return errors.New("no input file or URL")
	}

	// Step 2: Extract audio (no need to extract separately, Transcribe handles this)
	progress <- ProgressUpdate{0.25, "Preparing for transcription..."}

	// Check for cancellation before transcription
	select {
	case <-ctx.Done():
		return fmt.Errorf("operation cancelled: %w", ctx.Err())
	default:
	}

	// Step 3: Transcription (20% to 50% = 30% range)
	progress <- ProgressUpdate{0.30, "Transcribing audio..."}
	transcription, err := e.Transcribe(videoPath)
	if err != nil {
		progress <- ProgressUpdate{0.0, fmt.Sprintf("Transcription failed: %v", err)}
		return fmt.Errorf("transcription failed: %w", err)
	}
	progress <- ProgressUpdate{0.50, "Transcription complete"}

	// Check for cancellation before translation
	select {
	case <-ctx.Done():
		return fmt.Errorf("operation cancelled: %w", ctx.Err())
	default:
	}

	// Step 4: Translation (50% to 65% = 15% range)
	progress <- ProgressUpdate{0.50, "Translating text..."}
	translation, err := e.Translate(transcription, opts.TargetLanguage)
	if err != nil {
		progress <- ProgressUpdate{0.0, fmt.Sprintf("Translation failed: %v", err)}
		return fmt.Errorf("translation failed: %w", err)
	}
	progress <- ProgressUpdate{0.65, "Translation complete"}

	// Prepare output directory early (needed for dubbing and final outputs)
	outputDir := opts.OutputDir
	if outputDir == "" {
		if opts.InputFile != "" {
			outputDir = filepath.Dir(opts.InputFile)
		} else {
			outputDir = filepath.Join(os.TempDir(), "akashic_scribe_output")
		}
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		progress <- ProgressUpdate{0.0, fmt.Sprintf("Failed to prepare output directory: %v", err)}
		return fmt.Errorf("failed to prepare output directory: %w", err)
	}

	// Step 5: (Optional) Dubbing (65% to 85% = 20% range)
	var dubbedAudioPath string
	if opts.CreateDubbing {
		// Check for cancellation before dubbing
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation cancelled: %w", ctx.Err())
		default:
		}

		progress <- ProgressUpdate{0.68, "Synthesizing dubbed audio with TTS..."}

		// Set default dubbing parameters
		setDefaultDubbingParams(&opts)

		// Generate dubbed audio
		audioPath, err := e.GenerateDubbing(translation, opts, outputDir)
		if err != nil {
			progress <- ProgressUpdate{0.68, fmt.Sprintf("Warning: Dubbing failed: %v", err)}
			// Don't fail the entire process, just log the warning
			// User can still get transcription and translation
		} else {
			dubbedAudioPath = audioPath
			progress <- ProgressUpdate{0.85, "Dubbed audio generated successfully"}
		}
	}

	// Step 6: (Optional) Subtitles (85% to 90% = 5% range)
	if opts.CreateSubtitles {
		// Check for cancellation before subtitle generation
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation cancelled: %w", ctx.Err())
		default:
		}

		progress <- ProgressUpdate{0.87, "Generating subtitles..."}
	}

	// Step 7: Compose final output (90% to 95% = 5% range)
	progress <- ProgressUpdate{0.92, "Preparing final output..."}

	// Step 8: Complete
	result := struct {
		Transcription string
		Translation   string
		DubbedAudio   string `json:",omitempty"`
		SubtitlesFile string `json:",omitempty"`
	}{
		Transcription: transcription,
		Translation:   translation,
		DubbedAudio:   dubbedAudioPath,
	}

	// Save outputs to disk (outputDir already prepared earlier)
	progress <- ProgressUpdate{0.97, "Saving outputs..."}

	// Write transcription and translation
	if err := os.WriteFile(filepath.Join(outputDir, "transcription.txt"), []byte(transcription), 0o644); err != nil {
		progress <- ProgressUpdate{0.0, fmt.Sprintf("Failed to write transcription: %v", err)}
		return fmt.Errorf("failed to write transcription: %w", err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "translation.txt"), []byte(translation), 0o644); err != nil {
		progress <- ProgressUpdate{0.0, fmt.Sprintf("Failed to write translation: %v", err)}
		return fmt.Errorf("failed to write translation: %w", err)
	}
	if opts.CreateSubtitles {
		// Use enhanced subtitle generator
		subtitleGen := NewSubtitleGenerator()

		// Create default subtitle segments (will be improved with Whisper timestamps in future)
		// Assuming 3-minute video for now - in production, get actual duration from video
		videoDuration := 3 * time.Minute
		subtitleGen.CreateDefaultSegments(transcription, translation, videoDuration)

		// Determine subtitle format
		format := opts.SubtitleFormat
		if format == "" {
			format = "srt" // default to SRT
		}

		// Generate subtitle content
		var subtitleContent string
		var extension string

		if format == "vtt" {
			subtitleContent = subtitleGen.GenerateVTT(opts.BilingualSubtitles, opts.SubtitlePosition)
			extension = ".vtt"
		} else {
			subtitleContent = subtitleGen.GenerateSRT(opts.BilingualSubtitles, opts.SubtitlePosition)
			extension = ".srt"
		}

		subtitlesPath := filepath.Join(outputDir, "subtitles"+extension)
		if err := os.WriteFile(subtitlesPath, []byte(subtitleContent), 0o644); err != nil {
			progress <- ProgressUpdate{0.0, fmt.Sprintf("Failed to write subtitles: %v", err)}
			return fmt.Errorf("failed to write subtitles: %w", err)
		}
		result.SubtitlesFile = subtitlesPath
	}

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	completionMsg := fmt.Sprintf("Scribing complete.\nOutput saved to: %s\n%s", outputDir, string(resultJSON))
	progress <- ProgressUpdate{1.0, completionMsg}

	return nil
}

// ProcessWithContext runs the full pipeline with context and returns a structured result.
// This method is used by the batch processor for better result handling.
func (e *realScribeEngine) ProcessWithContext(ctx context.Context, opts ScribeOptions, progress chan<- ProgressUpdate) (*ScribeResult, error) {
	// Check dependencies first
	if err := e.checkDependencies(); err != nil {
		return nil, err
	}

	// Determine output directory
	outputDir := opts.OutputDir
	if outputDir == "" {
		outputDir = filepath.Join(".", "akashic_output_"+time.Now().Format("20060102_150405"))
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Step 1: Obtain video file (download if URL, or use local file)
	progress <- ProgressUpdate{0.0, "Starting..."}

	var videoPath string

	if opts.InputFile != "" {
		// Verify local file exists
		if _, err := os.Stat(opts.InputFile); err != nil {
			return nil, fmt.Errorf("input file not found: %w", err)
		}
		videoPath = opts.InputFile
		progress <- ProgressUpdate{0.05, "Using local file"}
	} else if opts.InputURL != "" {
		// Create temporary directory for download
		tempDir, err := os.MkdirTemp("", "akashic_scribe_download_*")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp directory: %w", err)
		}
		defer func() {
			if err := os.RemoveAll(tempDir); err != nil {
				log.Printf("Warning: failed to clean up temp directory %s: %v", tempDir, err)
			}
		}()

		progress <- ProgressUpdate{0.05, "Downloading video..."}

		// Download video using yt-dlp
		videoPath = filepath.Join(tempDir, "downloaded_video.%(ext)s")
		cmd := exec.CommandContext(ctx, "yt-dlp", "-o", videoPath, opts.InputURL)
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("failed to download video: %w", err)
		}

		// Find the actual downloaded file
		files, err := filepath.Glob(filepath.Join(tempDir, "downloaded_video.*"))
		if err != nil || len(files) == 0 {
			return nil, errors.New("downloaded file not found")
		}
		videoPath = files[0]
		progress <- ProgressUpdate{0.20, "Download complete"}
	} else {
		return nil, errors.New("no input file or URL provided")
	}

	// Step 2: Transcription
	progress <- ProgressUpdate{0.30, "Transcribing audio..."}
	transcription, err := e.Transcribe(videoPath)
	if err != nil {
		return nil, fmt.Errorf("transcription failed: %w", err)
	}
	progress <- ProgressUpdate{0.50, "Transcription complete"}

	// Check for cancellation
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("operation cancelled: %w", ctx.Err())
	default:
	}

	// Step 3: Translation
	progress <- ProgressUpdate{0.50, "Translating text..."}
	translation, err := e.Translate(transcription, opts.TargetLanguage)
	if err != nil {
		return nil, fmt.Errorf("translation failed: %w", err)
	}
	progress <- ProgressUpdate{0.65, "Translation complete"}

	// Check for cancellation
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("operation cancelled: %w", ctx.Err())
	default:
	}

	// Step 4: Optional dubbing
	var dubbedAudioPath string
	if opts.CreateDubbing {
		progress <- ProgressUpdate{0.70, "Generating dubbed audio..."}
		dubbedAudioPath, err = e.GenerateDubbing(translation, opts, outputDir)
		if err != nil {
			log.Printf("Warning: Dubbing failed but continuing: %v", err)
			progress <- ProgressUpdate{0.80, "Dubbing failed, continuing without audio..."}
		} else {
			progress <- ProgressUpdate{0.80, "Dubbing complete"}
		}
	}

	// Check for cancellation
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("operation cancelled: %w", ctx.Err())
	default:
	}

	// Step 5: Optional subtitles
	var subtitlesPath string
	if opts.CreateSubtitles {
		progress <- ProgressUpdate{0.87, "Generating subtitles..."}

		subtitleGen := NewSubtitleGenerator()
		videoDuration := 3 * time.Minute
		subtitleGen.CreateDefaultSegments(transcription, translation, videoDuration)

		format := opts.SubtitleFormat
		if format == "" {
			format = "srt"
		}

		var subtitleContent string
		var extension string

		if format == "vtt" {
			subtitleContent = subtitleGen.GenerateVTT(opts.BilingualSubtitles, opts.SubtitlePosition)
			extension = ".vtt"
		} else {
			subtitleContent = subtitleGen.GenerateSRT(opts.BilingualSubtitles, opts.SubtitlePosition)
			extension = ".srt"
		}

		subtitlesPath = filepath.Join(outputDir, "subtitles"+extension)
		if err := os.WriteFile(subtitlesPath, []byte(subtitleContent), 0o644); err != nil {
			return nil, fmt.Errorf("failed to write subtitles: %w", err)
		}
	}

	// Step 6: Save outputs
	progress <- ProgressUpdate{0.97, "Saving outputs..."}

	if err := os.WriteFile(filepath.Join(outputDir, "transcription.txt"), []byte(transcription), 0o644); err != nil {
		return nil, fmt.Errorf("failed to write transcription: %w", err)
	}

	if err := os.WriteFile(filepath.Join(outputDir, "translation.txt"), []byte(translation), 0o644); err != nil {
		return nil, fmt.Errorf("failed to write translation: %w", err)
	}

	// Create result
	result := &ScribeResult{
		Transcription: transcription,
		Translation:   translation,
		DubbedAudio:   dubbedAudioPath,
		SubtitlesFile: subtitlesPath,
		OutputDir:     outputDir,
	}

	progress <- ProgressUpdate{1.0, "Processing complete"}

	return result, nil
}
