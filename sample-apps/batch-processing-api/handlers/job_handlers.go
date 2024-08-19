package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/erikkvale/yggdrasil"
	"github.com/labstack/echo/v4"
)

type JobRequest struct {
	JobType    string `json:"job_type"`
	DataSource string `json:"data_source"`
}

type JobStatus struct {
	ID        int       `json:"id"`
	Status    string    `json:"status"`
	Result    string    `json:"result,omitempty"`
	Error     string    `json:"error,omitempty"`
	Submitted time.Time `json:"submitted"`
	Started   time.Time `json:"started,omitempty"`
	Completed time.Time `json:"completed,omitempty"`
}

func StartJob(c echo.Context, pool *yggdrasil.WorkerPool) error {
	req := new(JobRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	job := func() error {
		// Fill in logic later
		fmt.Println("Processing job with type: ", req.JobType, " and data source: ", req.DataSource)
		return nil
	}

	pool.AddJob(job)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Job Started Successfully",
	})
}
