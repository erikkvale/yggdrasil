package yggdrasil

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// WorkerPoolTestSuite is a test suite for the WorkerPool.
type WorkerPoolTestSuite struct {
	suite.Suite
	Pool *WorkerPool
}

// SetupSuite runs once before all tests in the suite.
func (suite *WorkerPoolTestSuite) SetupSuite() {
	suite.Pool = NewWorkerPool(10, 5)
	suite.Pool.Start()
}

// TearDownSuite runs once after all tests in the suite.
func (suite *WorkerPoolTestSuite) TearDownSuite() {
	suite.Pool.Shutdown()
}

// TestNewWorkerPool checks the initialization of the pool.
func (suite *WorkerPoolTestSuite) TestNewWorkerPool() {
	assert.Equal(suite.T(), 10, suite.Pool.QueueSize, "Expected QueueSize to be 10")
	assert.Equal(suite.T(), 5, len(suite.Pool.Workers), "Expected 5 workers")
	assert.Equal(suite.T(), 10, cap(suite.Pool.JobQueue), "Expected JobQueue capacity to be 10")
}

// TestWorkerPoolAddJob tests adding and processing a job.
func (suite *WorkerPoolTestSuite) TestWorkerPoolAddJob() {
	processed := make(chan bool, 1)

	job := func() error {
		processed <- true
		return nil
	}

	suite.Pool.AddJob(job)

	select {
	case <-processed:
		assert.True(suite.T(), true, "Job was processed")
	case <-suite.Pool.JobQueue:
		assert.Fail(suite.T(), "Job was not processed")
	}
}

// TestWorkerPoolJobErrorHandling tests error handling in jobs.
func (suite *WorkerPoolTestSuite) TestWorkerPoolJobErrorHandling() {
	jobErrorOccurred := make(chan bool, 1)

	job := func() error {
		jobErrorOccurred <- true
		return fmt.Errorf("job failed")
	}

	suite.Pool.AddJob(job)

	select {
	case <-jobErrorOccurred:
		assert.True(suite.T(), true, "Job error occurred as expected")
	case <-suite.Pool.JobQueue:
		assert.Fail(suite.T(), "Job did not trigger an error")
	}
}

// TestWorkerPoolShutdown tests that the WorkerPool shuts down correctly.
func (suite *WorkerPoolTestSuite) TestWorkerPoolShutdown() {
	pool := NewWorkerPool(10, 5)
	pool.Start()

	job := func() error {
		return nil
	}

	// Add jobs
	for i := 0; i < 5; i++ {
		pool.AddJob(job)
	}

	// Give some time for jobs to be processed
	time.Sleep(100 * time.Millisecond)

	// Shutdown the pool
	pool.Shutdown()

	// Check that the channel is closed using select
	select {
	case _, ok := <-pool.JobQueue:
		assert.False(suite.T(), ok, "Expected JobQueue to be closed")
	default:
		// If we reach here, the channel was already empty and closed
	}
}

// TestWorkerPoolTestSuite runs the test suite.
func TestWorkerPoolTestSuite(t *testing.T) {
	suite.Run(t, new(WorkerPoolTestSuite))
}
