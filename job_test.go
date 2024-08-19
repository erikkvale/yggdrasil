package yggdrasil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type JobTestSuite struct {
	suite.Suite
}

func (suite *JobTestSuite) TestJobStatusInitialization() {
	jobStatus := JobStatus{
		ID:        1,
		Status:    "pending",
		Submitted: time.Now(),
	}

	suite.Equal(1, jobStatus.ID, "Job ID should be 1")
	suite.Equal("pending", jobStatus.Status, "Job status should be 'pending'")
	suite.True(jobStatus.Started.IsZero(), "Job start time should be zero")
	suite.True(jobStatus.Completed.IsZero(), "Job completion time should be zero")
}

// SetupTest and TearDownTest can be added if there’s a need for setup or cleanup before/after each test.

func TestJobTestSuite(t *testing.T) {
	suite.Run(t, new(JobTestSuite))
}
