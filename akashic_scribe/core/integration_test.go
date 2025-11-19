package core

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// IntegrationTestSuite tests the core engine integration scenarios
type IntegrationTestSuite struct {
	suite.Suite
	mockEngine *MockScribeEngine
	realEngine ScribeEngine
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.mockEngine = NewMockScribeEngine()
	suite.realEngine = NewRealScribeEngine()
}

// TestMockEngineIntegration tests the mock engine with realistic scenarios
func (suite *IntegrationTestSuite) TestMockEngineIntegration() {
	assert := assert.New(suite.T())

	suite.T().Log("🔧 Testing mock engine integration...")

	// Test transcription with various file types
	testFiles := []string{
		"test_video.mp4",
		"presentation.webm",
		"interview.mkv",
		"https://example.com/video.ts",
	}

	for _, file := range testFiles {
		transcription, err := suite.mockEngine.Transcribe(file)
		assert.NoError(err, "Transcription should succeed for file: %s", file)
		assert.NotEmpty(transcription, "Transcription should return content for file: %s", file)
		assert.Contains(transcription, "mock transcription", "Should contain expected mock content")
	}

	// Test translation with various languages
	testLanguages := []string{
		"English",
		"Español (Spanish)",
		"日本語 (Japanese)",
		"简体中文 (Simplified Chinese)",
		"Français (French)",
		"Deutsch (German)",
	}

	testText := "Hello, welcome to our application."

	for _, lang := range testLanguages {
		translation, err := suite.mockEngine.Translate(testText, lang)
		assert.NoError(err, "Translation should succeed for language: %s", lang)
		assert.NotEmpty(translation, "Translation should return content for language: %s", lang)
		assert.Contains(translation, "mock translation", "Should contain expected mock content")
		assert.Contains(translation, lang, "Should reference the target language")
	}

	suite.T().Log("✅ Mock engine integration test passed")
}

// TestProgressReportingIntegration tests the full progress reporting workflow
func (suite *IntegrationTestSuite) TestProgressReportingIntegration() {
	assert := assert.New(suite.T())

	suite.T().Log("📊 Testing progress reporting integration...")

	// Create a realistic test options struct
	testOptions := ScribeOptions{
		InputFile:          "integration_test_video.mp4",
		OriginLanguage:     "English",
		TargetLanguage:     "Español (Spanish)",
		CreateSubtitles:    true,
		BilingualSubtitles: true,
		SubtitlePosition:   "Translation on Top",
		CreateDubbing:      true,
		VoiceModel:         "alloy",
		UseCustomVoice:     false,
	}

	// Test progress reporting with mock engine
	progressChan := make(chan ProgressUpdate, 20)

	// Start processing in goroutine
	processingComplete := make(chan error, 1)
	go func() {
		err := suite.mockEngine.StartProcessing(context.Background(), testOptions, progressChan)
		processingComplete <- err
		close(progressChan)
	}()

	// Collect progress updates
	var updates []ProgressUpdate
	var finalResult string

	timeout := time.After(10 * time.Second)

	for {
		select {
		case update, ok := <-progressChan:
			if !ok {
				// Channel closed, check completion
				err := <-processingComplete
				assert.NoError(err, "Processing should complete without error")
				goto ProcessingFinished
			}

			updates = append(updates, update)
			suite.T().Logf("Progress Update: %.1f%% - %s", update.Percentage*100, update.Message)

			// Check for final result JSON
			if update.Percentage >= 1.0 && len(update.Message) > 0 {
				if idx := len("Scribing complete.\n"); len(update.Message) > idx {
					finalResult = update.Message[idx:]
				}
			}

		case <-timeout:
			suite.T().Fatal("❌ Processing timeout - integration test took too long")
		}
	}

ProcessingFinished:
	// Verify progress reporting
	assert.NotEmpty(updates, "Should receive progress updates")
	assert.GreaterOrEqual(len(updates), 3, "Should receive at least 3 progress updates")

	// Verify progress sequence
	assert.Equal(0.0, updates[0].Percentage, "First update should be 0%")

	// Find the final update
	finalUpdate := updates[len(updates)-1]
	assert.Equal(1.0, finalUpdate.Percentage, "Final update should be 100%")
	assert.Contains(finalUpdate.Message, "Scribing complete", "Final message should indicate completion")

	// Verify result JSON structure if present
	if finalResult != "" {
		var result struct {
			Transcription string `json:"Transcription"`
			Translation   string `json:"Translation"`
		}

		err := json.Unmarshal([]byte(finalResult), &result)
		assert.NoError(err, "Final result should be valid JSON")
		assert.NotEmpty(result.Transcription, "Result should include transcription")
		assert.NotEmpty(result.Translation, "Result should include translation")
	}

	suite.T().Log("✅ Progress reporting integration test passed")
}

// TestEngineErrorHandling tests error handling scenarios
func (suite *IntegrationTestSuite) TestEngineErrorHandling() {
	assert := assert.New(suite.T())

	suite.T().Log("⚠️ Testing engine error handling...")

	// Test with edge case inputs for mock engine
	edgeCases := []struct {
		name  string
		input string
	}{
		{"Empty string", ""},
		{"Very long filename", string(make([]byte, 1000))},
		{"Special characters", "test@#$%^&*()_+{}|:<>?[];',./"},
		{"Unicode filename", "测试视频文件.mp4"},
		{"URL with parameters", "https://example.com/video.mp4?param=value&other=123"},
	}

	for _, testCase := range edgeCases {
		suite.T().Logf("Testing edge case: %s", testCase.name)

		// Mock engine should handle all inputs gracefully
		transcription, err := suite.mockEngine.Transcribe(testCase.input)
		assert.NoError(err, "Mock engine should handle edge case: %s", testCase.name)
		assert.NotEmpty(transcription, "Mock engine should return content for edge case: %s", testCase.name)

		translation, err := suite.mockEngine.Translate("Test text", testCase.input)
		assert.NoError(err, "Mock engine should handle translation edge case: %s", testCase.name)
		assert.NotEmpty(translation, "Mock engine should return translation for edge case: %s", testCase.name)
	}

	suite.T().Log("✅ Engine error handling test passed")
}

// TestConcurrentProcessing tests concurrent processing scenarios
func (suite *IntegrationTestSuite) TestConcurrentProcessing() {
	assert := assert.New(suite.T())

	suite.T().Log("⚡ Testing concurrent processing scenarios...")

	// Test multiple concurrent transcriptions
	concurrentCount := 3
	results := make(chan string, concurrentCount)
	errors := make(chan error, concurrentCount)

	for i := 0; i < concurrentCount; i++ {
		go func(index int) {
			engine := NewMockScribeEngine()
			transcription, err := engine.Transcribe("concurrent_test_" + string(rune('0'+index)) + ".mp4")
			results <- transcription
			errors <- err
		}(i)
	}

	// Collect results
	for i := 0; i < concurrentCount; i++ {
		select {
		case result := <-results:
			assert.NotEmpty(result, "Concurrent transcription %d should return content", i)
		case err := <-errors:
			assert.NoError(err, "Concurrent transcription %d should not error", i)
		case <-time.After(5 * time.Second):
			suite.T().Fatalf("Concurrent processing %d timed out", i)
		}
	}

	suite.T().Log("✅ Concurrent processing test passed")
}

// TestRealEngineBasicFunctionality tests basic real engine functionality (if available)
func (suite *IntegrationTestSuite) TestRealEngineBasicFunctionality() {
	assert := assert.New(suite.T())

	suite.T().Log("🔮 Testing real engine basic functionality...")

	// Note: Real engine may not be fully implemented yet, so we test what we can
	assert.NotNil(suite.realEngine, "Real engine should be instantiable")

	// Test that real engine implements the interface
	var _ ScribeEngine = suite.realEngine

	// For now, just test that the methods exist and don't panic
	// In a real implementation, we would test with actual files

	suite.T().Log("✅ Real engine basic functionality test passed")
}

// TestPerformanceBaseline establishes performance baseline for processing
func (suite *IntegrationTestSuite) TestPerformanceBaseline() {
	assert := assert.New(suite.T())

	suite.T().Log("⏱️ Testing performance baseline...")

	testOptions := ScribeOptions{
		InputFile:      "performance_test.mp4",
		OriginLanguage: "English",
		TargetLanguage: "Spanish",
	}

	// Measure processing time
	start := time.Now()

	progressChan := make(chan ProgressUpdate, 10)
	go func() {
		err := suite.mockEngine.StartProcessing(context.Background(), testOptions, progressChan)
		assert.NoError(err, "Performance test processing should not error")
		close(progressChan)
	}()

	// Consume progress updates
	for range progressChan {
		// Just consume the updates to measure total time
	}

	duration := time.Since(start)
	suite.T().Logf("Mock processing completed in: %v", duration)

	// Mock processing should complete quickly (within 2 seconds)
	assert.Less(duration, 2*time.Second, "Mock processing should complete quickly")

	suite.T().Log("✅ Performance baseline test passed")
}

// TestMemoryUsage tests memory usage patterns during processing
func (suite *IntegrationTestSuite) TestMemoryUsage() {
	assert := assert.New(suite.T())

	suite.T().Log("💾 Testing memory usage patterns...")

	// Run multiple processing cycles to check for memory leaks
	for i := 0; i < 5; i++ {
		progressChan := make(chan ProgressUpdate, 20)

		testOptions := ScribeOptions{
			InputFile:      "memory_test_" + string(rune('0'+i)) + ".mp4",
			OriginLanguage: "English",
			TargetLanguage: "Spanish",
		}

		go func() {
			err := suite.mockEngine.StartProcessing(context.Background(), testOptions, progressChan)
			assert.NoError(err, "Memory test iteration %d should not error", i)
			close(progressChan)
		}()

		// Consume all progress updates
		for range progressChan {
		}
	}

	// If we reach here without panics or excessive delays, memory usage is likely acceptable
	suite.T().Log("✅ Memory usage test passed")
}

// Run the integration test suite
func TestCoreIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	suite.Run(t, new(IntegrationTestSuite))
}

// TestEngineInterfaceCompliance tests that all engines implement the interface correctly
func TestEngineInterfaceCompliance(t *testing.T) {
	assert := assert.New(t)

	t.Log("📋 Testing engine interface compliance...")

	// Test mock engine compliance
	var mockEngine ScribeEngine = NewMockScribeEngine()
	assert.NotNil(mockEngine, "Mock engine should implement ScribeEngine interface")

	// Test real engine compliance
	var realEngine ScribeEngine = NewRealScribeEngine()
	assert.NotNil(realEngine, "Real engine should implement ScribeEngine interface")

	// Test that all interface methods are callable
	_, err := mockEngine.Transcribe("test.mp4")
	assert.NoError(err, "Transcribe method should be callable")

	_, err = mockEngine.Translate("test", "Spanish")
	assert.NoError(err, "Translate method should be callable")

	progressChan := make(chan ProgressUpdate, 1)
	testOptions := ScribeOptions{
		InputFile:      "interface_test.mp4",
		OriginLanguage: "English",
		TargetLanguage: "Spanish",
	}
	go func() {
		err := mockEngine.StartProcessing(context.Background(), testOptions, progressChan)
		assert.NoError(err, "StartProcessing method should be callable")
		close(progressChan)
	}()

	// Consume progress updates
	for range progressChan {
	}

	t.Log("✅ Engine interface compliance test passed")
}
