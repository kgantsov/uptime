package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/kgantsov/uptime/app/handler"
	"github.com/kgantsov/uptime/app/model"
	"github.com/kgantsov/uptime/app/monitor"

	rice "github.com/GeertJohan/go.rice"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/kgantsov/uptime/app/cmd/uptime/docs"
)

type HTTPBox struct {
	*rice.Box
}

func (hb *HTTPBox) Open(name string) (fs.File, error) {
	return hb.Box.Open(name)
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func main() {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		fmt.Print("failed to connect database\n")
		return
	}

	dispatcher := monitor.NewDispatcher(db)
	dispatcher.Start()

	done := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Printf("Got signal: %s\n", sig)

		fmt.Print("Stopping monitoring\n")

		dispatcher.Stop()

		time.Sleep(100 * time.Millisecond)

		done <- struct{}{}
	}()

	// Migrate the schema
	db.AutoMigrate(
		&model.Heartbeat{}, &model.Service{}, &model.Notification{}, &model.ServiceNotification{},
	)
	db.SetupJoinTable(&model.Service{}, "Notifications", &model.ServiceNotification{})

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

	h := handler.NewHandler(db, dispatcher)
	h.RegisterRoutes(e)

	e.GET("/docs/*", echoSwagger.WrapHandler)

	appStaticBox, err := rice.FindBox("../../../frontend/build/static/")
	if err != nil {
		log.Fatal(err)
	}

	appIndexBox, err := rice.FindBox("../../../frontend/build/")
	if err != nil {
		log.Fatal(err)
	}

	e.StaticFS("/static/", &HTTPBox{appStaticBox})
	e.GET("/*", echo.StaticFileHandler("index.html", &HTTPBox{appIndexBox}))

	go func() {
		e.Logger.Fatal(e.Start(":1323"))
	}()

	log.Infof("Started uptime monitor\n")

	<-done

	fmt.Print("Stopped monitoring\n")
}
