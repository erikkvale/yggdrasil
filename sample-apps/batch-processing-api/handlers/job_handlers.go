package handlers

import (
	"fmt"
	"net/http"

	"../../../yggdrasil"
	"github.com/labstack/echo"
)

type JobRequest struct {
	JobType    string `json:"job_type"`
	DataSource string `json:"data_source"`
}

func StartJob(c echo.Context) error {
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

	pool := yggdrasil.NewWorkerPool(10, 5)
	pool.Start()
	pool.Add(job)
}
