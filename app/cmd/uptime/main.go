package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/kgantsov/uptime/app/auth"
	"github.com/kgantsov/uptime/app/handler"
	"github.com/kgantsov/uptime/app/model"
	"github.com/kgantsov/uptime/app/monitor"
	"github.com/kgantsov/uptime/app/repository"
	"github.com/kgantsov/uptime/app/service"
	"github.com/rs/zerolog/log"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Embed a single file
//
//go:embed build/index.html
var indexHtmlFS embed.FS

// Embed a directory
//
//go:embed build/static/*
var frontendFS embed.FS

func InitStaticServer(app *fiber.App) {
	// Serve static files from the embedded filesystem
	app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(frontendFS),
		PathPrefix: "build/static",
		Browse:     false,
	}))

	// Serve index.html from the embedded filesystem
	app.Get("/*", filesystem.New(filesystem.Config{
		Root:         http.FS(indexHtmlFS),
		Index:        "build/index.html",
		NotFoundFile: "build/index.html",
	}))
}

func main() {
	dbPathPtr := flag.String("db-path", "./test.db", "A path to a DB file")
	logModePtr := flag.String("log-mode", "", "Logging mode: STACKDRIVER or console (default)")
	logLevelPtr := flag.String("log-level", "debug", "Logging level: debug, info, warn, error, fatal, panic")
	flag.Parse()

	handler.ConfigureLogger(*logModePtr, *logLevelPtr)

	db, err := gorm.Open(sqlite.Open(*dbPathPtr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Info().Msg("failed to connect database")
		return
	}

	serviceRepo := repository.NewServiceRepository(db)
	heartbeatRepo := repository.NewHeartbeatRepository(db)
	notifRepo := repository.NewNotificationRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	userRepo := repository.NewUserRepository(db)

	dispatcher := monitor.NewDispatcher(serviceRepo, heartbeatRepo)
	dispatcher.Start()

	heartbeatSvc := service.NewHeartbeatService(heartbeatRepo)
	serviceSvc := service.NewServiceService(serviceRepo, notifRepo, dispatcher)
	notifSvc := service.NewNotificationService(notifRepo, dispatcher)
	tokenSvc := service.NewTokenService(userRepo, tokenRepo, handler.Key)

	h := handler.NewHandler(heartbeatSvc, serviceSvc, notifSvc, tokenSvc)

	app, _ := handler.NewFiberApp(h)

	done := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Info().Msgf("Got signal: %s", sig)
		log.Info().Msg("Stopping monitoring")

		dispatcher.Stop()
		time.Sleep(100 * time.Millisecond)

		if err := app.Shutdown(); err != nil {
			log.Error().Msgf("Error shutting down fiber: %s", err)
		}

		done <- struct{}{}
	}()

	model.MigrateDB(db)

	initUser(userRepo)

	InitStaticServer(app)

	go cleanupOldHeartbeats(heartbeatRepo)

	go func() {
		if err := app.Listen(":1323"); err != nil {
			log.Fatal().Msgf("Fiber listen error: %s", err)
		}
	}()

	log.Info().Msg("Started uptime monitor")

	<-done

	log.Info().Msg("Stopped monitoring")
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

func cleanupOldHeartbeats(heartbeatRepo repository.HeartbeatRepository) {
	ticker := time.NewTicker(time.Duration(60) * time.Second)

	for {
		select {
		case <-ticker.C:
			thresholdDate := time.Now().AddDate(0, 0, -30)
			log.Info().Msgf("Deleting heartbeats older than %s", thresholdDate)

			if err := heartbeatRepo.DeleteOlderThan(thresholdDate); err != nil {
				log.Fatal().Msgf("Failed to delete old heartbeats: %s", err)
			}
		}
	}
}
