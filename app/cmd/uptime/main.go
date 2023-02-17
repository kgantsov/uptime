package main

import (
	"bufio"
	"fmt"
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

	_ "github.com/kgantsov/uptime/app/cmd/uptime/docs"
)

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

	model.MigrateDB(db)

	initUser(db)

	h.ConfigureMiddleware(e)
	h.RegisterRoutes(e)
	h.InitStaticServer(e)

	go func() {
		e.Logger.Fatal(e.Start(":1323"))
	}()

	e.Logger.Infof("Started uptime monitor")

	<-done

	e.Logger.Info("Stopped monitoring")
}

func initUser(db *gorm.DB) {
	var count int64

	err := db.Model(&model.User{}).Count(&count).Error

	if err == nil && count > 0 {
		return
	}

	scanner := bufio.NewScanner((os.Stdin))

	fmt.Println("Enter your First Name: ")
	scanner.Scan()
	firstName := scanner.Text()

	fmt.Println("Enter your Last Name: ")
	scanner.Scan()
	lastName := scanner.Text()

	fmt.Println("Enter your Email: ")
	scanner.Scan()
	email := scanner.Text()

	fmt.Println("Enter your Password: ")
	scanner.Scan()
	password := scanner.Text()

	createUser(db, firstName, lastName, email, password)
}

func createUser(db *gorm.DB, firstName, lastName, email, password string) {
	h, _ := auth.HashPassword(password)
	user := &model.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  h,
	}
	db.Create(user)
}
