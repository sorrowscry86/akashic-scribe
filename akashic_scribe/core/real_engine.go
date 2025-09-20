package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	defer os.RemoveAll(tempDir)

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

// StartProcessing runs the full pipeline and reports progress.
func (e *realScribeEngine) StartProcessing(options ScribeOptions, progress chan<- ProgressUpdate) error {
	// Check dependencies first
	if err := e.checkDependencies(); err != nil {
		progress <- ProgressUpdate{0.0, fmt.Sprintf("Dependency check failed: %v", err)}
		return err
	}

	opts := options

	progress <- ProgressUpdate{0.01, "Preparing job..."}
	time.Sleep(200 * time.Millisecond)

	// Step 1: Download or use local file
	var videoPath string
	if opts.InputFile != "" {
		// Verify file exists
		if _, err := os.Stat(opts.InputFile); err != nil {
			progress <- ProgressUpdate{0.0, fmt.Sprintf("Input file not found: %v", err)}
			return fmt.Errorf("input file not found: %w", err)
		}
		videoPath = opts.InputFile
		progress <- ProgressUpdate{0.10, "Using local video file."}
	} else if opts.InputURL != "" {
		progress <- ProgressUpdate{0.10, "Downloading video from URL..."}

		// Create temporary directory for downloads
		tempDir, err := os.MkdirTemp("", "akashic_scribe_download_*")
		if err != nil {
			progress <- ProgressUpdate{0.0, fmt.Sprintf("Failed to create temp directory: %v", err)}
			return fmt.Errorf("failed to create temp directory: %w", err)
		}
		defer os.RemoveAll(tempDir)

		videoPath = filepath.Join(tempDir, "downloaded_video.%(ext)s")
		cmd := exec.Command("yt-dlp", "-o", videoPath, opts.InputURL)
		if err := cmd.Run(); err != nil {
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
	} else {
		progress <- ProgressUpdate{0.0, "No input file or URL provided."}
		return errors.New("no input file or URL")
	}

	// Step 2: Extract audio
	progress <- ProgressUpdate{0.25, "Extracting audio..."}
	time.Sleep(400 * time.Millisecond)

	// Step 3: Transcription
	progress <- ProgressUpdate{0.40, "Transcribing audio..."}
	transcription, err := e.Transcribe(videoPath)
	if err != nil {
		progress <- ProgressUpdate{0.0, fmt.Sprintf("Transcription failed: %v", err)}
		return fmt.Errorf("transcription failed: %w", err)
	}
	time.Sleep(400 * time.Millisecond)

	// Step 4: Translation
	progress <- ProgressUpdate{0.60, "Translating text..."}
	translation, err := e.Translate(transcription, opts.TargetLanguage)
	if err != nil {
		progress <- ProgressUpdate{0.0, fmt.Sprintf("Translation failed: %v", err)}
		return fmt.Errorf("translation failed: %w", err)
	}
	time.Sleep(400 * time.Millisecond)

	// Step 5: (Optional) Dubbing
	if opts.CreateDubbing {
		progress <- ProgressUpdate{0.75, "Synthesizing dubbed audio..."}
		time.Sleep(400 * time.Millisecond)
	}

	// Step 6: (Optional) Subtitles
	if opts.CreateSubtitles {
		progress <- ProgressUpdate{0.85, "Generating subtitles..."}
		time.Sleep(400 * time.Millisecond)
	}

	// Step 7: Compose final output
	progress <- ProgressUpdate{0.95, "Composing final output..."}
	time.Sleep(400 * time.Millisecond)

	// Step 8: Complete
	result := struct {
		Transcription string
		Translation   string
	}{transcription, translation}

	// Save outputs to disk
	// Work out output directory
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
		srt := "1\n00:00:00,000 --> 00:00:03,000\n" + translation + "\n\n"
		if err := os.WriteFile(filepath.Join(outputDir, "subtitles.srt"), []byte(srt), 0o644); err != nil {
			progress <- ProgressUpdate{0.0, fmt.Sprintf("Failed to write subtitles: %v", err)}
			return fmt.Errorf("failed to write subtitles: %w", err)
		}
	}

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	completionMsg := fmt.Sprintf("Scribing complete.\nOutput saved to: %s\n%s", outputDir, string(resultJSON))
	progress <- ProgressUpdate{1.0, completionMsg}

	return nil
}
