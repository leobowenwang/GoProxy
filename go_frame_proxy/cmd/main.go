package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/project-sesame/sesame-gateway/internal/pkg/config"
	"github.com/project-sesame/sesame-gateway/internal/pkg/database"
	"github.com/project-sesame/sesame-gateway/internal/pkg/filter"
	"github.com/project-sesame/sesame-gateway/internal/pkg/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config, err := config.LoadConfig("internal/pkg/config/config.yml")
	if err != nil {
		panic(err)
	}

	// Connect to the database.
	err = database.Connect(config)
	if err != nil {
		panic(err)
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(util.PrometheusMiddleware)

	// Custom Middleware
	e.Use(filter.BasicToJWT(config.Server.CertFile, config.Server.KeyFile))
	e.Use(filter.AddUUIDHeader)
	filter.ConfigureProxy(e, config)

	// Routes
	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Prometheus
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Start server
	var serverPort = fmt.Sprintf("%d", config.Server.Port)
	e.Logger.Fatal(e.Start(":" + serverPort))
}
