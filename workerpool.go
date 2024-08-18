package yggdrasil

import (
	"fmt"
	"sync"
)

// Job represents a unit of work that can be processed by a worker.
type Job func() error

// Worker represents a worker that processes jobs from the WorkerPool.
type Worker struct {
	Pool *WorkerPool
}

// PickUpJob starts the worker to process jobs from the pool.
func (w Worker) PickUpJob() {
	for {
		select {
		case job, ok := <-w.Pool.JobQueue:
			if !ok {
				return
			}
			err := job()
			if err != nil {
				fmt.Println("Error processing job:", err)
			}
		}
	}
}

// WorkerPool manages a pool of workers to process jobs concurrently.
type WorkerPool struct {
	JobQueue    chan Job
	QueueSize   int
	Workers     []Worker
	WorkersSize int
	once        sync.Once
}

// NewWorkerPool creates a new WorkerPool with the specified job queue size and number of workers.
func NewWorkerPool(queueSize int, workersSize int) *WorkerPool {
	jobQueue := make(chan Job, queueSize)
	wp := &WorkerPool{
		JobQueue:    jobQueue,
		QueueSize:   queueSize,
		Workers:     make([]Worker, workersSize),
		WorkersSize: workersSize,
	}
	// Initialize each worker
	for i := range wp.Workers {
		wp.Workers[i] = Worker{Pool: wp}
	}
	return wp
}

// Start begins the processing of jobs by the workers in the pool.
func (wp *WorkerPool) Start() {
	for _, w := range wp.Workers {
		go w.PickUpJob()
	}
}

// AddJob adds a new job to the WorkerPool for processing.
func (wp *WorkerPool) AddJob(job Job) {
	wp.JobQueue <- job
}

// Shutdown gracefully shuts down the WorkerPool after all jobs have been processed.
func (wp *WorkerPool) Shutdown() {
	wp.once.Do(func() {
		close(wp.JobQueue)
	})
}
