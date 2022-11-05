package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type Health struct {
	Status string `json:"status"`
}

func main() {
	portPtr := flag.String("port", "1313", "A port for the server")
	flag.Parse()

	e := echo.New()
	log := logrus.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")

			return nil
		},
	}))

	rand.Seed(time.Now().UnixNano())

	e.GET("/health", func(c echo.Context) error {
		min := 0
		max := 10
		num := rand.Intn(max-min+1) + min

		if num == 0 {
			h := &Health{
				Status: "Failed",
			}
			log.Info("-----> FAIL")
			return c.JSON(http.StatusInternalServerError, h)
		}

		if num == 1 {
			log.Info("-----> TIMEOUT")
			time.Sleep(time.Second * time.Duration(2))
		} else {
			log.Info("-----> SUCCESS")
		}

		h := &Health{
			Status: "OK",
		}

		return c.JSON(http.StatusOK, h)
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *portPtr)))
}
