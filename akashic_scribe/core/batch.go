package core

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
)

// BatchJob represents a single job in a batch processing queue.
type BatchJob struct {
	ID          string             // Unique identifier for the job
	Options     ScribeOptions      // Processing options for this job
	Status      JobStatus          // Current status of the job
	Error       error              // Error if job failed
	Result      *ScribeResult      // Result if job succeeded
	StartTime   time.Time          // When the job started processing
	EndTime     time.Time          // When the job completed
	Progress    float64            // Current progress (0.0 to 1.0)
	StatusMsg   string             // Current status message
	jobCancel   context.CancelFunc // Function to cancel this specific job
}

// JobStatus represents the current state of a batch job.
type JobStatus int

const (
	JobPending JobStatus = iota
	JobRunning
	JobCompleted
	JobFailed
	JobCancelled
)

// String returns the string representation of a JobStatus.
func (s JobStatus) String() string {
	switch s {
	case JobPending:
		return "Pending"
	case JobRunning:
		return "Running"
	case JobCompleted:
		return "Completed"
	case JobFailed:
		return "Failed"
	case JobCancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}

// BatchProcessor manages batch processing of multiple video files.
type BatchProcessor struct {
	engine         ScribeEngine
	jobs           map[string]*BatchJob
	jobQueue       chan *BatchJob
	maxConcurrent  int
	workers        sync.WaitGroup
	jobWaitGroup   sync.WaitGroup // Tracks active jobs for efficient waiting
	ctx            context.Context
	cancel         context.CancelFunc
	mu             sync.RWMutex
	progressChan   chan BatchProgress
}

// BatchProgress represents progress for the entire batch operation.
type BatchProgress struct {
	TotalJobs      int
	CompletedJobs  int
	FailedJobs     int
	RunningJobs    int
	PendingJobs    int
	OverallPercent float64
	CurrentJobID   string
	CurrentJobMsg  string
}

// NewBatchProcessor creates a new batch processor.
//
// Parameters:
//   - engine: The ScribeEngine to use for processing
//   - maxConcurrent: Maximum number of jobs to process simultaneously (0 = number of CPUs)
func NewBatchProcessor(engine ScribeEngine, maxConcurrent int) *BatchProcessor {
	if maxConcurrent <= 0 {
		maxConcurrent = runtime.NumCPU() // Default to number of CPUs for performance
	}

	ctx, cancel := context.WithCancel(context.Background())

	bp := &BatchProcessor{
		engine:        engine,
		jobs:          make(map[string]*BatchJob),
		jobQueue:      make(chan *BatchJob, 100), // Buffer for queued jobs
		maxConcurrent: maxConcurrent,
		ctx:           ctx,
		cancel:        cancel,
		progressChan:  make(chan BatchProgress, 10),
	}

	// Start worker goroutines
	for i := 0; i < maxConcurrent; i++ {
		bp.workers.Add(1)
		go bp.worker(i)
	}

	return bp
}

// AddJob adds a new job to the batch queue.
// This method will block if the queue is full to prevent unbounded goroutine creation.
func (bp *BatchProcessor) AddJob(options ScribeOptions) string {
	bp.mu.Lock()

	// Generate unique job ID using UUID to prevent collisions
	jobID := fmt.Sprintf("job_%s", uuid.New().String())

	job := &BatchJob{
		ID:      jobID,
		Options: options,
		Status:  JobPending,
	}

	bp.jobs[jobID] = job
	bp.jobWaitGroup.Add(1) // Track this job for Wait()
	bp.mu.Unlock()

	// Add to queue (blocks if full, preventing goroutine leak)
	log.Printf("Batch: Job %s added to queue", jobID)
	bp.jobQueue <- job

	return jobID
}

// worker processes jobs from the queue.
func (bp *BatchProcessor) worker(workerID int) {
	defer bp.workers.Done()

	log.Printf("Batch: Worker %d started", workerID)

	for {
		select {
		case <-bp.ctx.Done():
			log.Printf("Batch: Worker %d shutting down", workerID)
			return

		case job := <-bp.jobQueue:
			bp.processJob(workerID, job)
		}
	}
}

// processJob processes a single job.
func (bp *BatchProcessor) processJob(workerID int, job *BatchJob) {
	defer bp.jobWaitGroup.Done() // Mark job complete for Wait()

	// Create a cancellable context for this job
	jobCtx, jobCancel := context.WithCancel(bp.ctx)
	defer jobCancel()

	bp.mu.Lock()
	job.Status = JobRunning
	job.StartTime = time.Now()
	job.jobCancel = jobCancel // Store cancel function for CancelJob()
	bp.mu.Unlock()

	log.Printf("Batch: Worker %d processing job %s", workerID, job.ID)
	bp.sendProgress()

	// Create a progress channel for this job
	progressChan := make(chan ProgressUpdate, 10)

	// Listen for progress updates
	done := make(chan bool)
	go func() {
		for {
			select {
			case update := <-progressChan:
				bp.mu.Lock()
				job.Progress = update.Percentage
				job.StatusMsg = update.Message
				bp.mu.Unlock()
				bp.sendProgress()

			case <-done:
				return
			}
		}
	}()

	// Process the job
	result, err := bp.engine.ProcessWithContext(jobCtx, job.Options, progressChan)

	close(done)
	close(progressChan)

	bp.mu.Lock()
	job.EndTime = time.Now()

	if err != nil {
		if jobCtx.Err() == context.Canceled {
			job.Status = JobCancelled
			log.Printf("Batch: Job %s cancelled", job.ID)
		} else {
			job.Status = JobFailed
			job.Error = err
			log.Printf("Batch: Job %s failed: %v", job.ID, err)
		}
	} else {
		job.Status = JobCompleted
		job.Result = result
		job.Progress = 1.0
		log.Printf("Batch: Job %s completed successfully", job.ID)
	}

	bp.mu.Unlock()
	bp.sendProgress()
}

// sendProgress sends batch progress updates to the progress channel.
func (bp *BatchProcessor) sendProgress() {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	progress := BatchProgress{
		TotalJobs: len(bp.jobs),
	}

	var totalProgress float64
	var currentJobID string
	var currentJobMsg string

	for _, job := range bp.jobs {
		switch job.Status {
		case JobCompleted:
			progress.CompletedJobs++
			totalProgress += 1.0
		case JobFailed:
			progress.FailedJobs++
			totalProgress += 1.0
		case JobRunning:
			progress.RunningJobs++
			totalProgress += job.Progress
			if currentJobID == "" {
				currentJobID = job.ID
				currentJobMsg = job.StatusMsg
			}
		case JobPending:
			progress.PendingJobs++
		}
	}

	if progress.TotalJobs > 0 {
		progress.OverallPercent = totalProgress / float64(progress.TotalJobs)
	}

	progress.CurrentJobID = currentJobID
	progress.CurrentJobMsg = currentJobMsg

	// Non-blocking send with logging when dropped
	select {
	case bp.progressChan <- progress:
	default:
		log.Printf("Warning: Progress update dropped (channel full) [OverallPercent=%.2f%%]", progress.OverallPercent*100)
	}
}

// GetProgress returns the progress channel for batch updates.
func (bp *BatchProcessor) GetProgress() <-chan BatchProgress {
	return bp.progressChan
}

// GetJob returns information about a specific job.
func (bp *BatchProcessor) GetJob(jobID string) (*BatchJob, bool) {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	job, exists := bp.jobs[jobID]
	return job, exists
}

// GetAllJobs returns all jobs in the batch.
func (bp *BatchProcessor) GetAllJobs() []*BatchJob {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	jobs := make([]*BatchJob, 0, len(bp.jobs))
	for _, job := range bp.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// CancelJob cancels a specific job if it's pending or running.
func (bp *BatchProcessor) CancelJob(jobID string) error {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	job, exists := bp.jobs[jobID]
	if !exists {
		return fmt.Errorf("job %s not found", jobID)
	}

	if job.Status == JobCompleted || job.Status == JobFailed || job.Status == JobCancelled {
		return fmt.Errorf("job %s already finished with status: %s", jobID, job.Status)
	}

	// Actually cancel the running job by calling its cancel function
	if job.jobCancel != nil {
		job.jobCancel()
		log.Printf("Batch: Cancelling job %s", jobID)
	}

	job.Status = JobCancelled
	return nil
}

// CancelAll cancels all pending and running jobs.
func (bp *BatchProcessor) CancelAll() {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	for _, job := range bp.jobs {
		if job.Status == JobPending || job.Status == JobRunning {
			if job.jobCancel != nil {
				job.jobCancel()
			}
			job.Status = JobCancelled
		}
	}
}

// Wait waits for all jobs to complete and shuts down the processor.
func (bp *BatchProcessor) Wait() {
	// Efficiently wait for all jobs using WaitGroup instead of polling
	bp.jobWaitGroup.Wait()

	// Shutdown workers
	bp.cancel()
	bp.workers.Wait()
	close(bp.progressChan)
}

// Shutdown immediately cancels all jobs and shuts down the processor.
func (bp *BatchProcessor) Shutdown() {
	bp.CancelAll()
	bp.cancel()
	bp.workers.Wait()
	close(bp.progressChan)
}

// GetSummary returns a summary of the batch processing results.
func (bp *BatchProcessor) GetSummary() string {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	var completed, failed, cancelled, pending, running int

	for _, job := range bp.jobs {
		switch job.Status {
		case JobCompleted:
			completed++
		case JobFailed:
			failed++
		case JobCancelled:
			cancelled++
		case JobPending:
			pending++
		case JobRunning:
			running++
		}
	}

	return fmt.Sprintf("Batch Summary: Total=%d, Completed=%d, Failed=%d, Cancelled=%d, Running=%d, Pending=%d",
		len(bp.jobs), completed, failed, cancelled, running, pending)
}
