package routes

import (
	"github.com/erikkvale/yggdrasil"
	"github.com/labstack/echo/v4"
	"github.com/yourusername/batch-processor-api/handlers"
)

func RegisterRoutes(e *echo.Echo, pool *yggdrasil.WorkerPool) {
	e.POST("/api/jobs/start", func(c echo.Context) error {
		return handlers.StartJob(c, pool)
	})
}
