package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Health struct {
	Status string `json:"status"`
}

func main() {
	portPtr := flag.String("port", "1313", "A port for the server")
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger_ := log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(logger.New(logger.Config{
		Format: "${time} ${method} ${path} ${status}\n",
		CustomTags: map[string]logger.LogFunc{
			"time": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				return output.WriteString(time.Now().Format(time.RFC3339))
			},
		},
	}))

	rand.Seed(time.Now().UnixNano()) //nolint:staticcheck

	app.Get("/health", func(c *fiber.Ctx) error {
		num := rand.Intn(11) //nolint:gosec

		if num == 0 {
			logger_.Info().Msg("-----> FAIL")
			return c.Status(http.StatusInternalServerError).JSON(&Health{Status: "Failed"})
		}

		if num == 1 {
			logger_.Info().Msg("-----> TIMEOUT")
			time.Sleep(2 * time.Second)
		} else {
			logger_.Info().Msg("-----> SUCCESS")
		}

		delay := rand.Intn(1001) //nolint:gosec
		time.Sleep(time.Duration(delay) * time.Millisecond)

		return c.JSON(&Health{Status: "OK"})
	})

	addr := fmt.Sprintf(":%s", *portPtr)
	logger_.Info().Msgf("Test server listening on %s", addr)

	if err := app.Listen(addr); err != nil {
		logger_.Fatal().Msgf("listen error: %s", err)
	}
}
