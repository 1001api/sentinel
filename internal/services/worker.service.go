package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/internal/constants"
)

// WorkerJob represents a job to be processed by the worker pool.
// It includes metadata (UserID, ProjectID, Timestamp) and a Callback function
// that contains the actual work to be performed.
type WorkerJob struct {
	Timestamp time.Time // When the job was created
	UserID    uuid.UUID // Identifier
	ProjectID uuid.UUID // Identifier
	Callback  func()    // Callback to be executed when the job is called
}

// WorkerPool manages a pool of workers to process jobs asynchronously.
// It includes a buffered channel for incoming jobs and a configurable pool size.
type WorkerPool struct {
	jobChan   chan WorkerJob // Channel to receive jobs
	poolCount int            // Number of worker goroutines
}

func InitWorkerPool(poolCount int, bufferSize int) *WorkerPool {
	return &WorkerPool{
		jobChan:   make(chan WorkerJob, bufferSize),
		poolCount: poolCount,
	}
}

// StartWorker starts the worker pool and spawns worker goroutines.
// Each worker processes jobs with a throttling mechanism to ensure jobs for the same
// UserID and ProjectID are not processed more frequently than every 2 seconds.
// ctx: Context for cancellation and graceful shutdown.
func (wp *WorkerPool) StartWorker(ctx context.Context) {
	// Track the last processed time for each user-project pair.
	lastProcessed := make(map[string]time.Time)

	// Store pending jobs that are waiting for the throttle window to expire.
	pendingJobs := make(map[string]WorkerJob)

	mu := sync.RWMutex{}

	for i := 0; i < wp.poolCount; i++ {
		go func() {
			// Ticker for throttling (process jobs every 5 seconds by default)
			ticker := time.NewTicker(constants.WORKER_TIME_TICKER)
			defer ticker.Stop()

			for {
				select {
				// Handle context cancellation for graceful shutdown.
				case <-ctx.Done():
					return

				// Handle incoming jobs (from another service)
				case job := <-wp.jobChan:
					key := fmt.Sprintf("%s:%s", job.UserID, job.ProjectID)

					mu.Lock()
					// add incoming jobs into pending map.
					pendingJobs[key] = job
					mu.Unlock()

				// Handle ticker events (throttle window)
				case <-ticker.C:
					mu.Lock()

					// Process all pending jobs that meet throttle criteria
					for key, job := range pendingJobs {
						lastTime, exists := lastProcessed[key]

						// Process if first time or enough time has passed
						if !exists || time.Since(lastTime) >= constants.WORKER_TIME_TICKER {
							lastProcessed[key] = time.Now()

							// remove job from pending map
							delete(pendingJobs, key)

							mu.Unlock()

							// Execute the callback function to perform the actual work.
							// This decouples the worker pool from another service (business logic).
							if job.Callback != nil {
								job.Callback()
							}

							mu.Lock()
						}
					}
					mu.Unlock()
				}
			}
		}()
	}
}
