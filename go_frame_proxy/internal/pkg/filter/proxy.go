package filter

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/leobowenwang/go_frame_proxy/internal/pkg/config"
	"github.com/leobowenwang/go_frame_proxy/internal/pkg/util"
	"go.uber.org/zap"
)

func ConfigureProxy(e *echo.Echo, config config.Config) {

	for _, c := range config.Proxy {
		path := e.Group(c.Path)
		host, _ := url.Parse(c.Host)
		proxy := []*middleware.ProxyTarget{
			{
				URL: host,
			},
		}

		path.Use(proxyMethodHandler(c.Methods, proxy))
		path.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(proxy)))
		util.Logger.Info("Configured proxy config", zap.String("path", c.Path), zap.String("host", c.Host))
	}
}

func proxyMethodHandler(methods []string, targets []*middleware.ProxyTarget) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqMethod := c.Request().Method

			allowed := false

			if len(methods) == 0 {
				// Allow all HTTP methods if methods slice is empty
				allowed = true
			} else {
				// Check if request method is allowed
				for _, m := range methods {
					if reqMethod == m {
						allowed = true
						break
					}
				}
			}

			if !allowed {
				return echo.NewHTTPError(http.StatusMethodNotAllowed, "Method not allowed")
			}

			// Setze das c.URL-Feld für jeden target
			for _, target := range targets {
				c.Request().URL = target.URL

				err := next(c)

				// Handle errors returned by subsequent middleware/handler
				if err != nil {
					var code int
					var message interface{}
					switch e := err.(type) {
					case *echo.HTTPError:
						code = e.Code
						message = e.Message
					}
					return c.JSON(code, message)
				}
			}

			return nil
		}
	}
}
