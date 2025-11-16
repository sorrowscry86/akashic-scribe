// Package validation provides input validation and sanitization for Akashic Scribe.
// It prevents security vulnerabilities like path traversal, command injection, and resource exhaustion.
package validation

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Validation errors
var (
	ErrEmptyInput           = errors.New("input cannot be empty")
	ErrInvalidURL           = errors.New("invalid URL format")
	ErrUnsupportedProtocol  = errors.New("unsupported URL protocol (only http/https allowed)")
	ErrFileNotFound         = errors.New("file does not exist")
	ErrInvalidFileExtension = errors.New("invalid file extension")
	ErrPathTraversal        = errors.New("path traversal detected")
	ErrFileTooLarge         = errors.New("file size exceeds maximum allowed")
	ErrInvalidPath          = errors.New("invalid file path")
	ErrDangerousPath        = errors.New("path points to dangerous system location")
)

const (
	// MaxFileSize is the maximum allowed file size (2 GB)
	MaxFileSize = 2 * 1024 * 1024 * 1024

	// MaxURLLength is the maximum allowed URL length
	MaxURLLength = 2048
)

// AllowedVideoExtensions lists the permitted video file extensions
var AllowedVideoExtensions = []string{".mp4", ".mkv", ".webm", ".avi", ".mov", ".flv", ".ts", ".m4v"}

// AllowedAudioExtensions lists the permitted audio file extensions
var AllowedAudioExtensions = []string{".mp3", ".wav", ".flac", ".m4a", ".aac", ".ogg", ".opus"}

// DangerousPathPrefixes lists path prefixes that should be blocked for output
var DangerousPathPrefixes = []string{
	"/etc/",
	"/sys/",
	"/proc/",
	"/dev/",
	"/boot/",
	"/root/",
	"C:\\Windows\\",
	"C:\\Program Files\\",
	"C:\\Program Files (x86)\\",
}

// InputValidator provides methods for validating and sanitizing user inputs
type InputValidator struct{}

// NewInputValidator creates a new InputValidator instance
func NewInputValidator() *InputValidator {
	return &InputValidator{}
}

// ValidateURL validates that a URL is properly formatted and uses an allowed protocol
func (v *InputValidator) ValidateURL(urlStr string) error {
	if urlStr == "" {
		return ErrEmptyInput
	}

	// Check URL length to prevent DoS
	if len(urlStr) > MaxURLLength {
		return fmt.Errorf("URL length %d exceeds maximum %d: %w", len(urlStr), MaxURLLength, ErrInvalidURL)
	}

	// Parse the URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	// Check for valid scheme
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("%w: got %s", ErrUnsupportedProtocol, scheme)
	}

	// Check for host
	if parsedURL.Host == "" {
		return fmt.Errorf("%w: missing host", ErrInvalidURL)
	}

	return nil
}

// ValidateInputFilePath validates that a file path exists, is readable, and has an allowed extension
func (v *InputValidator) ValidateInputFilePath(filePath string) error {
	if filePath == "" {
		return ErrEmptyInput
	}

	// Check for path traversal attempts
	if strings.Contains(filePath, "..") {
		return ErrPathTraversal
	}

	// Clean and normalize the path
	cleanPath := filepath.Clean(filePath)

	// Check if file exists
	fileInfo, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", ErrFileNotFound, cleanPath)
		}
		return fmt.Errorf("%w: %v", ErrInvalidPath, err)
	}

	// Check if it's a regular file (not a directory or device)
	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("%w: not a regular file", ErrInvalidPath)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(cleanPath))
	if !v.isAllowedVideoExtension(ext) {
		return fmt.Errorf("%w: %s (allowed: %v)", ErrInvalidFileExtension, ext, AllowedVideoExtensions)
	}

	// Check file size
	if fileInfo.Size() > MaxFileSize {
		return fmt.Errorf("%w: %d bytes (max: %d)", ErrFileTooLarge, fileInfo.Size(), MaxFileSize)
	}

	return nil
}

// ValidateOutputDirectory validates that an output directory path is safe to write to
func (v *InputValidator) ValidateOutputDirectory(dirPath string) error {
	if dirPath == "" {
		// Empty is allowed (will use default)
		return nil
	}

	// Check for path traversal attempts
	if strings.Contains(dirPath, "..") {
		return ErrPathTraversal
	}

	// Clean and normalize the path
	cleanPath := filepath.Clean(dirPath)

	// Convert to absolute path for safety checks
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidPath, err)
	}

	// Check against dangerous system paths
	for _, prefix := range DangerousPathPrefixes {
		if strings.HasPrefix(absPath, prefix) || strings.HasPrefix(absPath, strings.ToLower(prefix)) {
			return fmt.Errorf("%w: cannot write to %s", ErrDangerousPath, absPath)
		}
	}

	return nil
}

// ValidateAudioFilePath validates an audio file path (for custom voice samples)
func (v *InputValidator) ValidateAudioFilePath(filePath string) error {
	if filePath == "" {
		return ErrEmptyInput
	}

	// Check for path traversal attempts
	if strings.Contains(filePath, "..") {
		return ErrPathTraversal
	}

	// Clean and normalize the path
	cleanPath := filepath.Clean(filePath)

	// Check if file exists
	fileInfo, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", ErrFileNotFound, cleanPath)
		}
		return fmt.Errorf("%w: %v", ErrInvalidPath, err)
	}

	// Check if it's a regular file
	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("%w: not a regular file", ErrInvalidPath)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(cleanPath))
	if !v.isAllowedAudioExtension(ext) {
		return fmt.Errorf("%w: %s (allowed: %v)", ErrInvalidFileExtension, ext, AllowedAudioExtensions)
	}

	// Check file size (audio files should be smaller, max 100MB)
	maxAudioSize := int64(100 * 1024 * 1024)
	if fileInfo.Size() > maxAudioSize {
		return fmt.Errorf("%w: %d bytes (max: %d)", ErrFileTooLarge, fileInfo.Size(), maxAudioSize)
	}

	return nil
}

// SanitizePath cleans and normalizes a file path
func (v *InputValidator) SanitizePath(path string) string {
	return filepath.Clean(path)
}

// isAllowedVideoExtension checks if a file extension is in the allowed video extensions list
func (v *InputValidator) isAllowedVideoExtension(ext string) bool {
	ext = strings.ToLower(ext)
	for _, allowed := range AllowedVideoExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}

// isAllowedAudioExtension checks if a file extension is in the allowed audio extensions list
func (v *InputValidator) isAllowedAudioExtension(ext string) bool {
	ext = strings.ToLower(ext)
	for _, allowed := range AllowedAudioExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}
