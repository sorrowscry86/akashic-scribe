package validation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateURL(t *testing.T) {
	validator := NewInputValidator()

	tests := []struct {
		name    string
		url     string
		wantErr error
	}{
		{
			name:    "valid HTTP URL",
			url:     "http://example.com/video.mp4",
			wantErr: nil,
		},
		{
			name:    "valid HTTPS URL",
			url:     "https://youtube.com/watch?v=123",
			wantErr: nil,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: ErrEmptyInput,
		},
		{
			name:    "invalid protocol - FTP",
			url:     "ftp://example.com/video.mp4",
			wantErr: ErrUnsupportedProtocol,
		},
		{
			name:    "invalid protocol - file",
			url:     "file:///etc/passwd",
			wantErr: ErrUnsupportedProtocol,
		},
		{
			name:    "malformed URL",
			url:     "not a url at all",
			wantErr: ErrUnsupportedProtocol, // Malformed URLs often result in empty scheme
		},
		{
			name:    "URL without host",
			url:     "http://",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "URL too long",
			url:     "http://example.com/" + strings.Repeat("a", 3000),
			wantErr: ErrInvalidURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateURL(tt.url)
			if tt.wantErr == nil && err != nil {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != nil && !strings.Contains(err.Error(), tt.wantErr.Error()) {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateInputFilePath(t *testing.T) {
	validator := NewInputValidator()

	// Create a temporary test file
	tmpDir := t.TempDir()
	validFile := filepath.Join(tmpDir, "test_video.mp4")
	if err := os.WriteFile(validFile, []byte("test content"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create a file with invalid extension
	invalidExtFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(invalidExtFile, []byte("test"), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr error
	}{
		{
			name:    "valid file path",
			path:    validFile,
			wantErr: nil,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: ErrEmptyInput,
		},
		{
			name:    "non-existent file",
			path:    filepath.Join(tmpDir, "nonexistent.mp4"),
			wantErr: ErrFileNotFound,
		},
		{
			name:    "path traversal attempt",
			path:    "../../../etc/passwd",
			wantErr: ErrPathTraversal,
		},
		{
			name:    "invalid extension",
			path:    invalidExtFile,
			wantErr: ErrInvalidFileExtension,
		},
		{
			name:    "directory instead of file",
			path:    tmpDir,
			wantErr: ErrInvalidPath, // Directories aren't regular files
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateInputFilePath(tt.path)
			if tt.wantErr == nil && err != nil {
				t.Errorf("ValidateInputFilePath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != nil && !strings.Contains(err.Error(), tt.wantErr.Error()) {
				t.Errorf("ValidateInputFilePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateOutputDirectory(t *testing.T) {
	validator := NewInputValidator()

	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		wantErr error
	}{
		{
			name:    "valid directory",
			path:    tmpDir,
			wantErr: nil,
		},
		{
			name:    "empty path (allowed)",
			path:    "",
			wantErr: nil,
		},
		{
			name:    "path traversal attempt",
			path:    tmpDir + "/../../../etc",
			wantErr: ErrPathTraversal,
		},
		{
			name:    "dangerous system path - /etc/",
			path:    "/etc/myapp",
			wantErr: ErrDangerousPath,
		},
		{
			name:    "dangerous system path - /sys/",
			path:    "/sys/devices",
			wantErr: ErrDangerousPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateOutputDirectory(tt.path)
			if tt.wantErr == nil && err != nil {
				t.Errorf("ValidateOutputDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != nil && err != nil && !strings.Contains(err.Error(), tt.wantErr.Error()) {
				t.Errorf("ValidateOutputDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAudioFilePath(t *testing.T) {
	validator := NewInputValidator()

	// Create a temporary test file
	tmpDir := t.TempDir()
	validAudio := filepath.Join(tmpDir, "voice.mp3")
	if err := os.WriteFile(validAudio, []byte("audio content"), 0o644); err != nil {
		t.Fatal(err)
	}

	invalidExtFile := filepath.Join(tmpDir, "voice.txt")
	if err := os.WriteFile(invalidExtFile, []byte("test"), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr error
	}{
		{
			name:    "valid audio file",
			path:    validAudio,
			wantErr: nil,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: ErrEmptyInput,
		},
		{
			name:    "non-existent file",
			path:    filepath.Join(tmpDir, "nonexistent.mp3"),
			wantErr: ErrFileNotFound,
		},
		{
			name:    "path traversal attempt",
			path:    "../../../etc/passwd",
			wantErr: ErrPathTraversal,
		},
		{
			name:    "invalid extension",
			path:    invalidExtFile,
			wantErr: ErrInvalidFileExtension,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateAudioFilePath(tt.path)
			if tt.wantErr == nil && err != nil {
				t.Errorf("ValidateAudioFilePath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != nil && !strings.Contains(err.Error(), tt.wantErr.Error()) {
				t.Errorf("ValidateAudioFilePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizePath(t *testing.T) {
	validator := NewInputValidator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "clean simple path",
			input:    "/home/user/video.mp4",
			expected: "/home/user/video.mp4",
		},
		{
			name:     "path with double slashes",
			input:    "/home//user///video.mp4",
			expected: "/home/user/video.mp4",
		},
		{
			name:     "path with dot segments",
			input:    "/home/./user/./video.mp4",
			expected: "/home/user/video.mp4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.SanitizePath(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizePath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAllowedExtensions(t *testing.T) {
	validator := NewInputValidator()

	videoTests := []struct {
		ext      string
		expected bool
	}{
		{".mp4", true},
		{".MP4", true}, // Case insensitive
		{".mkv", true},
		{".webm", true},
		{".txt", false},
		{".exe", false},
		{".sh", false},
	}

	for _, tt := range videoTests {
		result := validator.isAllowedVideoExtension(tt.ext)
		if result != tt.expected {
			t.Errorf("isAllowedVideoExtension(%s) = %v, want %v", tt.ext, result, tt.expected)
		}
	}

	audioTests := []struct {
		ext      string
		expected bool
	}{
		{".mp3", true},
		{".MP3", true}, // Case insensitive
		{".wav", true},
		{".flac", true},
		{".txt", false},
		{".exe", false},
	}

	for _, tt := range audioTests {
		result := validator.isAllowedAudioExtension(tt.ext)
		if result != tt.expected {
			t.Errorf("isAllowedAudioExtension(%s) = %v, want %v", tt.ext, result, tt.expected)
		}
	}
}
