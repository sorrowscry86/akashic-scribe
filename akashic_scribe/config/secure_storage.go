// Package config provides secure configuration management for Akashic Scribe.
// It uses OS-native secure storage mechanisms to protect sensitive data like API keys.
package config

import (
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	// ServiceName is the identifier used for storing credentials in the OS keyring
	ServiceName = "AkashicScribe"

	// OpenAIKeyName is the key name for the OpenAI API key
	OpenAIKeyName = "OpenAI_API_Key"
)

var (
	// ErrKeyNotFound is returned when a requested key is not found in secure storage
	ErrKeyNotFound = errors.New("key not found in secure storage")
)

// SecureStorage provides methods for securely storing and retrieving sensitive configuration data.
type SecureStorage struct {
	serviceName string
}

// NewSecureStorage creates a new SecureStorage instance.
func NewSecureStorage() *SecureStorage {
	return &SecureStorage{
		serviceName: ServiceName,
	}
}

// SetOpenAIKey stores the OpenAI API key in the OS-native secure storage.
// On macOS: Keychain
// On Windows: Credential Manager
// On Linux: Secret Service (via libsecret)
func (s *SecureStorage) SetOpenAIKey(apiKey string) error {
	if apiKey == "" {
		return errors.New("API key cannot be empty")
	}

	err := keyring.Set(s.serviceName, OpenAIKeyName, apiKey)
	if err != nil {
		return fmt.Errorf("failed to store API key in secure storage: %w", err)
	}

	return nil
}

// GetOpenAIKey retrieves the OpenAI API key from the OS-native secure storage.
// Returns ErrKeyNotFound if the key has not been set.
func (s *SecureStorage) GetOpenAIKey() (string, error) {
	apiKey, err := keyring.Get(s.serviceName, OpenAIKeyName)
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", ErrKeyNotFound
		}
		return "", fmt.Errorf("failed to retrieve API key from secure storage: %w", err)
	}

	if apiKey == "" {
		return "", ErrKeyNotFound
	}

	return apiKey, nil
}

// DeleteOpenAIKey removes the OpenAI API key from the OS-native secure storage.
func (s *SecureStorage) DeleteOpenAIKey() error {
	err := keyring.Delete(s.serviceName, OpenAIKeyName)
	if err != nil {
		if err == keyring.ErrNotFound {
			// Already deleted, not an error
			return nil
		}
		return fmt.Errorf("failed to delete API key from secure storage: %w", err)
	}

	return nil
}

// HasOpenAIKey checks if an OpenAI API key is stored in secure storage.
func (s *SecureStorage) HasOpenAIKey() bool {
	_, err := s.GetOpenAIKey()
	return err == nil
}
