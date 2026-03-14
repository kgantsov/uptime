package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/glebarez/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/kgantsov/uptime/app/auth"
	"github.com/kgantsov/uptime/app/handler"
	"github.com/kgantsov/uptime/app/model"
	"github.com/kgantsov/uptime/app/monitor"
	"github.com/kgantsov/uptime/app/repository"
	"github.com/kgantsov/uptime/app/service"

	_ "github.com/kgantsov/uptime/app/cmd/uptime/docs"
)

type HTTPBox struct {
	*rice.Box
}

func (hb *HTTPBox) Open(name string) (fs.File, error) {
	return hb.Box.Open(name)
}

func InitStaticServer(e *echo.Echo) {
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
	dbPathPtr := flag.String("db-path", "./test.db", "A path to a DB file")
	flag.Parse()

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	// log.SetFormatter(new(handler.StackdriverFormatter))

	db, err := gorm.Open(sqlite.Open(*dbPathPtr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Info("failed to connect database")
		return
	}

	serviceRepo := repository.NewServiceRepository(db)
	heartbeatRepo := repository.NewHeartbeatRepository(db)
	notifRepo := repository.NewNotificationRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	userRepo := repository.NewUserRepository(db)

	dispatcher := monitor.NewDispatcher(serviceRepo, heartbeatRepo, log)
	dispatcher.Start()

	e := echo.New()

	heartbeatSvc := service.NewHeartbeatService(heartbeatRepo)
	serviceSvc := service.NewServiceService(serviceRepo, notifRepo, dispatcher)
	notifSvc := service.NewNotificationService(notifRepo, dispatcher)
	tokenSvc := service.NewTokenService(userRepo, tokenRepo, handler.Key)

	h := handler.NewHandler(log, heartbeatSvc, serviceSvc, notifSvc, tokenSvc)

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

	initUser(userRepo)

	h.ConfigureMiddleware(e)
	h.RegisterRoutes(e)
	InitStaticServer(e)

	go cleanupOldHeartbeats(heartbeatRepo, log)

	go func() {
		e.Logger.Fatal(e.Start(":1323"))
	}()

	e.Logger.Infof("Started uptime monitor")

	<-done

	e.Logger.Info("Stopped monitoring")
}

func initUser(userRepo repository.UserRepository) {
	count, err := userRepo.Count()

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

	createUser(userRepo, firstName, lastName, email, password)
}

func createUser(userRepo repository.UserRepository, firstName, lastName, email, password string) {
	h, _ := auth.HashPassword(password)
	user := &model.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  h,
	}
	userRepo.Create(user)
}

func cleanupOldHeartbeats(heartbeatRepo repository.HeartbeatRepository, logger *logrus.Logger) {
	ticker := time.NewTicker(time.Duration(60) * time.Second)

	for {
		select {
		case <-ticker.C:
			thresholdDate := time.Now().AddDate(0, 0, -30)
			logger.Infof("Deleting heartbeats older than %s", thresholdDate)

			if err := heartbeatRepo.DeleteOlderThan(thresholdDate); err != nil {
				log.Fatal(err)
			}
		}
	}
}
