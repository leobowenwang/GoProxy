package util

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests by method and status code.",
		},
		[]string{"method", "status_code"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests)
}

func PrometheusMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()

		// Call the next handler.
		err := next(c)

		// Increment the httpRequests counter.
		status := fmt.Sprintf("%d", res.Status)
		httpRequests.WithLabelValues(req.Method, status).Inc()

		return err
	}
}
