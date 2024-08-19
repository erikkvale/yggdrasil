package yggdrasil

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type WorkerPoolTestSuite struct {
	suite.Suite
	workerPool *WorkerPool
}

func (suite *WorkerPoolTestSuite) SetupTest() {
	suite.workerPool = NewWorkerPool(10, 2)
	suite.workerPool.Start()
}

func (suite *WorkerPoolTestSuite) TearDownTest() {
	suite.workerPool.Shutdown()
}

func (suite *WorkerPoolTestSuite) TestAddJob() {
	job := func() error {
		time.Sleep(1 * time.Second)
		return nil
	}

	jobID := suite.workerPool.AddJob(job)
	status := suite.workerPool.GetJobStatus(jobID)

	suite.NotNil(status, "Job status should not be nil")
	suite.Equal("pending", status.Status, "Job status should be 'pending'")
}

func (suite *WorkerPoolTestSuite) TestJobExecution() {
	job := func() error {
		return nil
	}

	jobID := suite.workerPool.AddJob(job)
	time.Sleep(2 * time.Second) // Give the worker some time to process the job

	status := suite.workerPool.GetJobStatus(jobID)

	suite.NotNil(status, "Job status should not be nil")
	suite.Equal("completed", status.Status, "Job status should be 'completed'")
}

func (suite *WorkerPoolTestSuite) TestFailedJob() {
	job := func() error {
		return errors.New("job failed")
	}

	jobID := suite.workerPool.AddJob(job)
	time.Sleep(2 * time.Second) // Give the worker some time to process the job

	status := suite.workerPool.GetJobStatus(jobID)

	suite.NotNil(status, "Job status should not be nil")
	suite.Equal("failed", status.Status, "Job status should be 'failed'")
	suite.Equal("job failed", status.Error, "Job error should be 'job failed'")
}

func TestWorkerPoolTestSuite(t *testing.T) {
	suite.Run(t, new(WorkerPoolTestSuite))
}
