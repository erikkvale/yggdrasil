package yggdrasil

import (
	"sync"
	"time"
)

// Worker represents a worker that processes jobs from the WorkerPool.
type Worker struct {
	Pool *WorkerPool
}

// PickUpJob starts the worker to process jobs from the pool.
// The worker listens on the JobQueue channel and processes jobs as they arrive.
// If the job queue is closed, the worker stops processing and exits.
func (w Worker) PickUpJob() {
	for {
		select {
		case job, ok := <-w.Pool.JobQueue:
			if !ok {
				// If the job queue is closed, the worker should exit
				return
			}
			jobID := job.ID
			w.Pool.updateJobStatus(jobID, "running")
			err := job.Execute()
			if err != nil {
				w.Pool.updateJobStatus(jobID, "failed")
				w.Pool.jobTracker[jobID].Error = err.Error()
			} else {
				w.Pool.updateJobStatus(jobID, "completed")
			}
		}
	}
}

// WorkerPool manages a pool of workers to process jobs concurrently.
// It maintains a queue of jobs, tracks the status of each job, and manages workers that execute the jobs.
type WorkerPool struct {
	JobQueue    chan *TrackedJob   // Channel through which jobs are sent to workers
	QueueSize   int                // Maximum number of jobs that can be queued at one time
	Workers     []Worker           // Slice of workers that process jobs
	WorkersSize int                // Number of workers in the pool
	jobTracker  map[int]*JobStatus // Map to track the status of jobs by their ID
	trackerLock sync.Mutex         // Mutex to protect access to the jobTracker
	jobCounter  int                // Counter to generate unique job IDs
	counterLock sync.Mutex         // Mutex to protect access to the jobCounter
}

// NewWorkerPool creates a new WorkerPool with the specified job queue size and number of workers.
// It initializes the job queue, the workers, and the job tracker.
func NewWorkerPool(queueSize int, workersSize int) *WorkerPool {
	wp := &WorkerPool{
		JobQueue:   make(chan *TrackedJob, queueSize),
		Workers:    make([]Worker, workersSize),
		jobTracker: make(map[int]*JobStatus),
	}

	for i := range wp.Workers {
		wp.Workers[i] = Worker{Pool: wp}
	}

	return wp
}

// Start begins the processing of jobs by the workers in the pool.
// It launches each worker as a goroutine that waits to process jobs from the job queue.
func (wp *WorkerPool) Start() {
	for _, w := range wp.Workers {
		go w.PickUpJob()
	}
}

// AddJob adds a new job to the WorkerPool for processing and returns the job ID.
// The job is added to the job queue, and its status is set to "pending" in the job tracker.
func (wp *WorkerPool) AddJob(job Job) int {
	wp.counterLock.Lock()
	wp.jobCounter++
	jobID := wp.jobCounter
	wp.counterLock.Unlock()

	trackedJob := &TrackedJob{
		ID:      jobID,
		Execute: job,
	}

	wp.trackerLock.Lock()
	wp.jobTracker[jobID] = &JobStatus{
		ID:        jobID,
		Status:    "pending",
		Submitted: time.Now(),
	}
	wp.trackerLock.Unlock()

	wp.JobQueue <- trackedJob

	return jobID
}

// updateJobStatus updates the status of a job in the jobTracker.
// It is used internally by workers to update job status as they process jobs.
func (wp *WorkerPool) updateJobStatus(jobID int, status string) {
	wp.trackerLock.Lock()
	defer wp.trackerLock.Unlock()

	jobStatus := wp.jobTracker[jobID]
	jobStatus.Status = status
	if status == "running" {
		jobStatus.Started = time.Now()
	}
	if status == "completed" || status == "failed" {
		jobStatus.Completed = time.Now()
	}
}

// GetJobStatus retrieves the current status of a job by its ID.
// Returns a pointer to the JobStatus or nil if the job ID is not found.
func (wp *WorkerPool) GetJobStatus(jobID int) *JobStatus {
	wp.trackerLock.Lock()
	defer wp.trackerLock.Unlock()

	return wp.jobTracker[jobID]
}

// Shutdown gracefully shuts down the WorkerPool after all jobs have been processed.
// It closes the job queue, causing all workers to finish their current jobs and exit.
func (wp *WorkerPool) Shutdown() {
	close(wp.JobQueue)
}
