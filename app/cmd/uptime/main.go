package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"

	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/kgantsov/uptime/app/auth"
	"github.com/kgantsov/uptime/app/handler"
	"github.com/kgantsov/uptime/app/model"
	"github.com/kgantsov/uptime/app/monitor"
	"github.com/kgantsov/uptime/app/repository"
	"github.com/kgantsov/uptime/app/service"
)

// riceHTTPBox wraps *rice.Box so that its Open method satisfies http.FileSystem.
// rice.Box.Open returns *rice.File which implements http.File, but the return
// type mismatch prevents the compiler from accepting *rice.Box directly as an
// http.FileSystem value.
type riceHTTPBox struct {
	*rice.Box
}

func (b *riceHTTPBox) Open(name string) (http.File, error) {
	return b.Box.Open(name)
}

func InitStaticServer(app *fiber.App) {
	appStaticBox, err := rice.FindBox("../../../frontend/build/static/")
	if err != nil {
		log.Fatal(err)
	}

	appIndexBox, err := rice.FindBox("../../../frontend/build/")
	if err != nil {
		log.Fatal(err)
	}

	app.Use("/static", filesystem.New(filesystem.Config{
		Root: &riceHTTPBox{appStaticBox},
	}))

	app.Use("/", filesystem.New(filesystem.Config{
		Root:         &riceHTTPBox{appIndexBox},
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))
}

func main() {
	dbPathPtr := flag.String("db-path", "./test.db", "A path to a DB file")
	flag.Parse()

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

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

	heartbeatSvc := service.NewHeartbeatService(heartbeatRepo)
	serviceSvc := service.NewServiceService(serviceRepo, notifRepo, dispatcher)
	notifSvc := service.NewNotificationService(notifRepo, dispatcher)
	tokenSvc := service.NewTokenService(userRepo, tokenRepo, handler.Key)

	h := handler.NewHandler(log, heartbeatSvc, serviceSvc, notifSvc, tokenSvc)

	app, _ := handler.NewFiberApp(h)

	done := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Infof("Got signal: %s", sig)
		log.Info("Stopping monitoring")

		dispatcher.Stop()
		time.Sleep(100 * time.Millisecond)

		if err := app.Shutdown(); err != nil {
			log.Errorf("Error shutting down fiber: %s", err)
		}

		done <- struct{}{}
	}()

	model.MigrateDB(db)

	initUser(userRepo)

	InitStaticServer(app)

	go cleanupOldHeartbeats(heartbeatRepo, log)

	go func() {
		if err := app.Listen(":1323"); err != nil {
			log.Fatalf("Fiber listen error: %s", err)
		}
	}()

	log.Infof("Started uptime monitor")

	<-done

	log.Info("Stopped monitoring")
}

func initUser(userRepo repository.UserRepository) {
	count, err := userRepo.Count()

	if err == nil && count > 0 {
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

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
