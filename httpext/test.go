package httpext

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func RunTestServer(address string) *echo.Echo {
	server := echo.New()

	server.GET("/wait", func(c echo.Context) error {
		fmt.Println("Server: Sleeping for a few seconds")

		time.Sleep(2 * time.Second)
		return c.JSON(http.StatusOK, "Here is some bytes in the body")
	})
	server.GET("/noop", func(c echo.Context) error {
		fmt.Println("Server: Doing Nothing")
		return nil
	})

	go func() {
		// The server will run as long as the parent process is running, which is cool.
		_ = server.Start(address)
		fmt.Println("Done Listening")
	}()

	return server
}
