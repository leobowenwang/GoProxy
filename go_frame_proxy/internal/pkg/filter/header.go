package filter

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/leobowenwang/go_frame_proxy/internal/pkg/util"
	"go.uber.org/zap"
)

// AddUUIDHeader adds custom UUID header as Echo middleware
func AddUUIDHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := uuid.New().String()
		c.Response().Before(func() {
			c.Response().Header().Set("X-Sesame-Request", id)
			util.Logger.Debug("added header", zap.String("uuid", id))
		})
		return next(c)
	}
}
