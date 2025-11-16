package config

import (
	"testing"
)

func TestSecureStorage_SetAndGetOpenAIKey(t *testing.T) {
	storage := NewSecureStorage()

	// Clean up any existing key first
	_ = storage.DeleteOpenAIKey()

	// Test setting a new API key
	testKey := "sk-test-1234567890abcdef"
	err := storage.SetOpenAIKey(testKey)
	if err != nil {
		t.Fatalf("Failed to set API key: %v", err)
	}

	// Test retrieving the API key
	retrievedKey, err := storage.GetOpenAIKey()
	if err != nil {
		t.Fatalf("Failed to get API key: %v", err)
	}

	if retrievedKey != testKey {
		t.Errorf("Retrieved key %q does not match set key %q", retrievedKey, testKey)
	}

	// Test HasOpenAIKey
	if !storage.HasOpenAIKey() {
		t.Error("HasOpenAIKey() returned false when key is set")
	}

	// Clean up
	err = storage.DeleteOpenAIKey()
	if err != nil {
		t.Fatalf("Failed to delete API key: %v", err)
	}

	// Verify deletion
	if storage.HasOpenAIKey() {
		t.Error("HasOpenAIKey() returned true after deletion")
	}
}

func TestSecureStorage_GetNonExistentKey(t *testing.T) {
	storage := NewSecureStorage()

	// Clean up first
	_ = storage.DeleteOpenAIKey()

	// Try to get a non-existent key
	_, err := storage.GetOpenAIKey()
	if err != ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got: %v", err)
	}
}

func TestSecureStorage_SetEmptyKey(t *testing.T) {
	storage := NewSecureStorage()

	// Try to set an empty key
	err := storage.SetOpenAIKey("")
	if err == nil {
		t.Error("Expected error when setting empty API key, got nil")
	}
}

func TestSecureStorage_DeleteNonExistentKey(t *testing.T) {
	storage := NewSecureStorage()

	// Clean up first
	_ = storage.DeleteOpenAIKey()

	// Delete should not error on non-existent key
	err := storage.DeleteOpenAIKey()
	if err != nil {
		t.Errorf("DeleteOpenAIKey() should not error on non-existent key, got: %v", err)
	}
}
