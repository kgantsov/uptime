package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
)

type Health struct {
	Status string `json:"status"`
}

func main() {
	portPtr := flag.String("port", "1313", "A port for the server")
	flag.Parse()

	log := logrus.New()

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
		Output: log.Writer(),
	}))

	rand.Seed(time.Now().UnixNano()) //nolint:staticcheck

	app.Get("/health", func(c *fiber.Ctx) error {
		num := rand.Intn(11) //nolint:gosec

		if num == 0 {
			log.Info("-----> FAIL")
			return c.Status(http.StatusInternalServerError).JSON(&Health{Status: "Failed"})
		}

		if num == 1 {
			log.Info("-----> TIMEOUT")
			time.Sleep(2 * time.Second)
		} else {
			log.Info("-----> SUCCESS")
		}

		delay := rand.Intn(1001) //nolint:gosec
		time.Sleep(time.Duration(delay) * time.Millisecond)

		return c.JSON(&Health{Status: "OK"})
	})

	addr := fmt.Sprintf(":%s", *portPtr)
	log.Infof("Test server listening on %s", addr)

	if err := app.Listen(addr); err != nil {
		log.Fatalf("listen error: %s", err)
	}
}
