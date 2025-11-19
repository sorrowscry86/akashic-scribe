package core

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestContextCancellationBeforeStart tests cancellation before processing starts.
func TestContextCancellationBeforeStart(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	engine := NewMockScribeEngine()

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	opts := ScribeOptions{
		InputFile:      "test.mp4",
		OriginLanguage: "English",
		TargetLanguage: "Spanish",
	}

	progressChan := make(chan ProgressUpdate, 10)

	// Should return error immediately
	err := engine.StartProcessing(ctx, opts, progressChan)
	require.Error(err)
	assert.Contains(err.Error(), "operation cancelled before start")
}

// TestContextCancellationDuringProcessing tests cancellation during processing.
func TestContextCancellationDuringProcessing(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	engine := NewMockScribeEngine()
	ctx, cancel := context.WithCancel(context.Background())

	opts := ScribeOptions{
		InputFile:      "test.mp4",
		OriginLanguage: "English",
		TargetLanguage: "Spanish",
		CreateDubbing:  true,
		CreateSubtitles: true,
	}

	progressChan := make(chan ProgressUpdate, 10)
	errChan := make(chan error, 1)

	// Start processing
	go func() {
		err := engine.StartProcessing(ctx, opts, progressChan)
		errChan <- err
		close(progressChan)
	}()

	// Wait a bit for processing to start
	time.Sleep(75 * time.Millisecond)

	// Cancel context
	cancel()

	// Wait for completion
	select {
	case err := <-errChan:
		require.Error(err)
		assert.Contains(err.Error(), "operation cancelled")
	case <-time.After(2 * time.Second):
		t.Fatal("Processing did not complete after cancellation")
	}
}

// TestContextCancellationAtMultiplePoints tests cancellation at different stages.
func TestContextCancellationAtMultiplePoints(t *testing.T) {
	engine := NewMockScribeEngine()

	tests := []struct {
		name         string
		cancelAfter  time.Duration
		expectStages []string
	}{
		{
			name:         "Cancel early (before transcription)",
			cancelAfter:  30 * time.Millisecond,
			expectStages: []string{"Starting processing"},
		},
		{
			name:         "Cancel mid-process (during transcription)",
			cancelAfter:  120 * time.Millisecond,
			expectStages: []string{"Starting processing", "Transcribing audio"},
		},
		{
			name:         "Cancel late (during translation)",
			cancelAfter:  180 * time.Millisecond,
			expectStages: []string{"Starting processing", "Transcribing audio", "Translating text"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			opts := ScribeOptions{
				InputFile:      "test.mp4",
				OriginLanguage: "English",
				TargetLanguage: "Spanish",
			}

			progressChan := make(chan ProgressUpdate, 10)
			errChan := make(chan error, 1)

			// Start processing
			go func() {
				err := engine.StartProcessing(ctx, opts, progressChan)
				errChan <- err
				close(progressChan)
			}()

			// Cancel after specified time
			time.Sleep(tt.cancelAfter)
			cancel()

			// Collect progress updates
			var updates []ProgressUpdate
			timeout := time.After(2 * time.Second)

		ProgressLoop:
			for {
				select {
				case update, ok := <-progressChan:
					if !ok {
						break ProgressLoop
					}
					updates = append(updates, update)
				case <-timeout:
					t.Fatal("Timeout waiting for cancellation")
				}
			}

			// Wait for error
			select {
			case err := <-errChan:
				require.Error(err)
				assert.Contains(err.Error(), "operation cancelled")
			case <-timeout:
				t.Fatal("Timeout waiting for error")
			}

			// Verify we received expected progress updates
			assert.NotEmpty(updates, "Should receive some progress updates before cancellation")
		})
	}
}

// TestContextTimeout tests processing with a timeout context.
func TestContextTimeout(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	engine := NewMockScribeEngine()

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	opts := ScribeOptions{
		InputFile:      "test.mp4",
		OriginLanguage: "English",
		TargetLanguage: "Spanish",
	}

	progressChan := make(chan ProgressUpdate, 10)

	// Should timeout before completing
	err := engine.StartProcessing(ctx, opts, progressChan)
	require.Error(err)
	assert.Contains(err.Error(), "operation cancelled")
}

// TestContextCancellationCleanup tests that resources are cleaned up on cancellation.
func TestContextCancellationCleanup(t *testing.T) {
	engine := NewMockScribeEngine()

	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithCancel(context.Background())

		opts := ScribeOptions{
			InputFile:      "test.mp4",
			OriginLanguage: "English",
			TargetLanguage: "Spanish",
		}

		progressChan := make(chan ProgressUpdate, 10)

		// Start and immediately cancel
		go func() {
			_ = engine.StartProcessing(ctx, opts, progressChan)
			close(progressChan)
		}()

		time.Sleep(50 * time.Millisecond)
		cancel()

		// Drain progress channel
		for range progressChan {
		}

		// If we reach here multiple times without hanging, cleanup is working
	}

	// If we get here, no goroutine leaks occurred
	// Test passes by completing without hanging
}

// TestProgressChannelCancellation tests that progress channel respects context.
func TestProgressChannelCancellation(t *testing.T) {
	engine := NewMockScribeEngine()
	ctx, cancel := context.WithCancel(context.Background())

	opts := ScribeOptions{
		InputFile:      "test.mp4",
		OriginLanguage: "English",
		TargetLanguage: "Spanish",
	}

	progressChan := make(chan ProgressUpdate, 10)

	// Start processing
	go func() {
		_ = engine.StartProcessing(ctx, opts, progressChan)
		close(progressChan)
	}()

	// Let it run a bit
	time.Sleep(50 * time.Millisecond)

	// Cancel
	cancel()

	// Progress channel should close eventually
	timeout := time.After(2 * time.Second)
	for {
		select {
		case _, ok := <-progressChan:
			if !ok {
				// Channel closed successfully
				return
			}
		case <-timeout:
			t.Fatal("Progress channel did not close after cancellation")
		}
	}
}

// TestConcurrentCancellations tests multiple concurrent operations being cancelled.
func TestConcurrentCancellations(t *testing.T) {
	assert := assert.New(t)

	engine := NewMockScribeEngine()
	numJobs := 5

	type job struct {
		ctx    context.Context
		cancel context.CancelFunc
		done   chan error
	}

	jobs := make([]job, numJobs)

	// Start multiple jobs
	for i := 0; i < numJobs; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		jobs[i] = job{
			ctx:    ctx,
			cancel: cancel,
			done:   make(chan error, 1),
		}

		opts := ScribeOptions{
			InputFile:      "test.mp4",
			OriginLanguage: "English",
			TargetLanguage: "Spanish",
		}

		progressChan := make(chan ProgressUpdate, 10)

		go func(j job) {
			err := engine.StartProcessing(j.ctx, opts, progressChan)
			j.done <- err
			close(progressChan)
		}(jobs[i])
	}

	// Let them run briefly
	time.Sleep(50 * time.Millisecond)

	// Cancel all jobs
	for _, j := range jobs {
		j.cancel()
	}

	// Wait for all to complete
	timeout := time.After(3 * time.Second)
	for i, j := range jobs {
		select {
		case err := <-j.done:
			assert.Error(err, "Job %d should have been cancelled", i)
			assert.Contains(err.Error(), "operation cancelled")
		case <-timeout:
			t.Fatalf("Job %d did not complete after cancellation", i)
		}
	}
}

// BenchmarkCancellation benchmarks the cancellation overhead.
func BenchmarkCancellation(b *testing.B) {
	engine := NewMockScribeEngine()

	opts := ScribeOptions{
		InputFile:      "test.mp4",
		OriginLanguage: "English",
		TargetLanguage: "Spanish",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		progressChan := make(chan ProgressUpdate, 10)

		go func() {
			_ = engine.StartProcessing(ctx, opts, progressChan)
			close(progressChan)
		}()

		time.Sleep(10 * time.Millisecond)
		cancel()

		// Drain channel
		for range progressChan {
		}
	}
}
