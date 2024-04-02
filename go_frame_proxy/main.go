package main

import (
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

func handleProxy(c echo.Context) error {
	targetURL := "http://localhost:8080"
	req, err := http.NewRequest(c.Request().Method, targetURL+c.Request().RequestURI, c.Request().Body)
	if err != nil {
		return err
	}

	req.Header = c.Request().Header.Clone()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	for key, values := range resp.Header {
		for _, value := range values {
			c.Response().Header().Set(key, value)
		}
	}
	c.Response().WriteHeader(resp.StatusCode)
	_, err = io.Copy(c.Response().Writer, resp.Body)
	return err
}

func main() {
	e := echo.New()
	e.Any("/*", handleProxy)
	e.Logger.Fatal(e.Start(":3333"))
}
