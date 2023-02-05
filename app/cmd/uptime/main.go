package main

import (
	"io/fs"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/kgantsov/uptime/app/auth"
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
// @securityDefinitions.apikey  HttpBearer
// @in                          header
// @name                        Authorization
// @description                 Description for what is this security definition being used
func main() {
	log := logrus.New()
	log.SetFormatter(new(handler.StackdriverFormatter))

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Info("failed to connect database")
		return
	}

	dispatcher := monitor.NewDispatcher(db)
	dispatcher.Start()

	e := echo.New()

	h := handler.NewHandler(log, db, dispatcher)

	done := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		e.Logger.Infof("Got signal: %s", sig)

		h.Logger.Info("Stopping monitoring")

		dispatcher.Stop()

		time.Sleep(100 * time.Millisecond)

		done <- struct{}{}
	}()

	// Migrate the schema
	db.AutoMigrate(
		&model.User{},
		&model.Token{},
		&model.Heartbeat{},
		&model.Service{},
		&model.Notification{},
		&model.ServiceNotification{},
	)
	db.SetupJoinTable(&model.Service{}, "Notifications", &model.ServiceNotification{})

	user := &model.User{}

	err = db.Model(&model.User{}).Where("email = ?", "admin@uptime.io").First(&user).Error

	if err != nil {
		h, _ := auth.HashPassword("12345qwert")
		user = &model.User{
			FirstName: "!",
			LastName:  "2",
			Email:     "admin@uptime.io",
			Password:  h,
		}
		db.Create(user)
	}

	h.ConfigureMiddleware(e)
	h.RegisterRoutes(e)

	e.GET("/docs/*", echoSwagger.WrapHandler)

	appStaticBox, err := rice.FindBox("../../../frontend/build/static/")
	if err != nil {
		e.Logger.Fatal(err)
	}

	appIndexBox, err := rice.FindBox("../../../frontend/build/")
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.StaticFS("/static/", &HTTPBox{appStaticBox})
	e.GET("/*", echo.StaticFileHandler("index.html", &HTTPBox{appIndexBox}))

	go func() {
		e.Logger.Fatal(e.Start(":1323"))
	}()

	e.Logger.Infof("Started uptime monitor")

	<-done

	e.Logger.Info("Stopped monitoring")
}
