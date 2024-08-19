package main

import (
	"net/http"

	"github.com/erikkvale/yggdrasil"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yourusername/batch-processor-api/routes"
)

func main() {
	pool := yggdrasil.NewWorkerPool(10, 5)
	pool.Start()

	// Create a new Echo instance
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "API is up and running!")
	})
	routes.RegisterRoutes(e, pool)

	e.Logger.Fatal(e.Start(":8080"))
}
