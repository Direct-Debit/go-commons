package httpext

import (
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

func RunTestServer(address string) *echo.Echo {
	server := echo.New()

	server.GET("/wait", func(c echo.Context) error {
		fmt.Println("Server: Sleeping for a few seconds")

		time.Sleep(2 * time.Second)
		return nil
	})
	server.GET("/noop", func(c echo.Context) error {
		fmt.Println("Server: Doing Nothing")
		return nil
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := server.Start(address)
		if err != nil {
			log.Fatalf("Server stoppped unexpectedly: %v", err)
		}
		wg.Done()
		fmt.Println("Done Listening")
	}()

	return server
}
