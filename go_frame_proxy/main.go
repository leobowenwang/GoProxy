package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

func handleProxy(c echo.Context, remoteHost string, remotePort int) error {
	targetURL := "http://" + remoteHost + ":" + fmt.Sprint(remotePort) + c.Request().RequestURI

	req, err := http.NewRequest(c.Request().Method, targetURL, c.Request().Body)
	if err != nil {
		return err
	}

	req.Header = c.Request().Header.Clone()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

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
	localPort := 3333
	remoteHost := "localhost"
	remotePort := 8080
	e := echo.New()

	e.Any("/*", func(c echo.Context) error {
		return handleProxy(c, remoteHost, remotePort)
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", localPort)))
}
