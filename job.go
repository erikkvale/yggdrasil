package yggdrasil

import (
	"time"
)

// Job represents a unit of work that can be processed by a worker.
// It is defined as a function that returns an error if the job fails.
type Job func() error

// TrackedJob represents a job with an ID that can be tracked by the WorkerPool.
// The ID is used to track the job's status in the job tracker.
type TrackedJob struct {
	ID      int // Unique identifier for the job
	Execute Job // Function to execute the job
}

// JobStatus represents the status of a job being processed by the worker pool.
// It includes details such as the job's ID, current status, result, any error encountered, and timestamps for submission, start, and completion.
type JobStatus struct {
	ID        int       `json:"id"`                  // Unique identifier for the job
	Status    string    `json:"status"`              // Current status of the job (e.g., "pending", "running", "completed", "failed")
	Result    string    `json:"result,omitempty"`    // Result of the job, if any (optional)
	Error     string    `json:"error,omitempty"`     // Error message if the job failed (optional)
	Submitted time.Time `json:"submitted"`           // Timestamp when the job was submitted
	Started   time.Time `json:"started,omitempty"`   // Timestamp when the job started (optional)
	Completed time.Time `json:"completed,omitempty"` // Timestamp when the job completed (optional)
}
