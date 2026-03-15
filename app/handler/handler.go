package handler

import (
	"time"

	fiberprometheus "github.com/ansrivas/fiberprometheus/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/kgantsov/uptime/app/auth"
	"github.com/kgantsov/uptime/app/service"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	HeartbeatService    service.HeartbeatService
	ServiceService      service.ServiceService
	NotificationService service.NotificationService
	TokenService        service.TokenService
	Logger              *logrus.Logger
}

const (
	// Key is the HMAC signing secret for JWTs.
	// In production this should come from configuration / environment.
	Key = "secret"
)

func NewHandler(
	logger *logrus.Logger,
	heartbeatService service.HeartbeatService,
	serviceService service.ServiceService,
	notificationService service.NotificationService,
	tokenService service.TokenService,
) *Handler {
	return &Handler{
		Logger:              logger,
		HeartbeatService:    heartbeatService,
		ServiceService:      serviceService,
		NotificationService: notificationService,
		TokenService:        tokenService,
	}
}

// NewFiberApp creates and returns a configured Fiber application together with
// the Huma API instance mounted on it.  Callers are responsible for adding
// static-file routes and calling app.Listen.
func NewFiberApp(h *Handler) (*fiber.App, huma.API) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ReadTimeout:           30 * time.Second,
		WriteTimeout:          30 * time.Second,
	})

	// ── Prometheus ────────────────────────────────────────────────────────────
	prometheus := fiberprometheus.New("uptime")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	// ── Basic middleware ───────────────────────────────────────────────────────
	app.Use(recover.New())
	app.Use(requestid.New())

	// ── Request logger ────────────────────────────────────────────────────────
	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		h.Logger.WithFields(logrus.Fields{
			"RequestID": c.Locals("requestid"),
		}).Infof("%s %s %d", c.Method(), c.Path(), c.Response().StatusCode())
		return err
	})

	// ── JWT authentication ────────────────────────────────────────────────────
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(Key)},
		Filter:     auth.AuthSkipperFunc,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or missing token")
		},
		ContextKey: "user",
	}))

	// ── Token-in-DB validation ────────────────────────────────────────────────
	app.Use(auth.CheckTokenMiddleware(h.TokenService, h.Logger))

	// ── Huma API (OpenAPI 3.1, auto-generated docs at /docs) ─────────────────
	config := huma.DefaultConfig("Uptime API", "1.0.0")
	config.Info.Description = "Uptime monitoring REST API"
	config.Components = &huma.Components{
		SecuritySchemes: map[string]*huma.SecurityScheme{
			"HttpBearer": {
				Type:   "http",
				Scheme: "bearer",
			},
		},
	}

	api := humafiber.New(app, config)

	h.RegisterRoutes(api)

	return app, api
}

// RegisterRoutes registers all API operations with the given Huma API instance.
func (h *Handler) RegisterRoutes(api huma.API) {
	secBearer := []map[string][]string{{"HttpBearer": {}}}

	// ── Tokens ────────────────────────────────────────────────────────────────
	huma.Register(api, huma.Operation{
		Method:      "POST",
		Path:        "/API/v1/tokens",
		Summary:     "Create an auth token",
		Description: "Authenticate with email/password and receive a signed JWT.",
		Tags:        []string{"tokens"},
	}, h.CreateToken)

	huma.Register(api, huma.Operation{
		Method:        "DELETE",
		Path:          "/API/v1/tokens",
		Summary:       "Delete an auth token",
		Description:   "Invalidate the currently authenticated token.",
		Tags:          []string{"tokens"},
		Security:      secBearer,
		DefaultStatus: 204,
	}, h.DeleteToken)

	// ── Heartbeats ────────────────────────────────────────────────────────────
	huma.Register(api, huma.Operation{
		Method:   "GET",
		Path:     "/API/v1/heartbeats/latencies",
		Summary:  "Get heartbeat latencies",
		Tags:     []string{"heartbeats"},
		Security: secBearer,
	}, h.GetHeartbeatsLatencies)

	huma.Register(api, huma.Operation{
		Method:   "GET",
		Path:     "/API/v1/heartbeats/latencies/last",
		Summary:  "Get last heartbeat latencies",
		Tags:     []string{"heartbeats"},
		Security: secBearer,
	}, h.GetHeartbeatsLastLatencies)

	huma.Register(api, huma.Operation{
		Method:   "GET",
		Path:     "/API/v1/heartbeats/stats/{days}",
		Summary:  "Get heartbeat stats",
		Tags:     []string{"heartbeats"},
		Security: secBearer,
	}, h.GetHeartbeatStats)

	// ── Services ──────────────────────────────────────────────────────────────
	huma.Register(api, huma.Operation{
		Method:   "GET",
		Path:     "/API/v1/services",
		Summary:  "Get all services",
		Tags:     []string{"services"},
		Security: secBearer,
	}, h.GetServices)

	huma.Register(api, huma.Operation{
		Method:   "POST",
		Path:     "/API/v1/services",
		Summary:  "Create a service",
		Tags:     []string{"services"},
		Security: secBearer,
	}, h.CreateService)

	huma.Register(api, huma.Operation{
		Method:   "GET",
		Path:     "/API/v1/services/{service_id}",
		Summary:  "Get a service",
		Tags:     []string{"services"},
		Security: secBearer,
	}, h.GetService)

	huma.Register(api, huma.Operation{
		Method:   "PATCH",
		Path:     "/API/v1/services/{service_id}",
		Summary:  "Update a service",
		Tags:     []string{"services"},
		Security: secBearer,
	}, h.UpdateService)

	huma.Register(api, huma.Operation{
		Method:        "DELETE",
		Path:          "/API/v1/services/{service_id}",
		Summary:       "Delete a service",
		Tags:          []string{"services"},
		Security:      secBearer,
		DefaultStatus: 204,
	}, h.DeleteService)

	huma.Register(api, huma.Operation{
		Method:   "POST",
		Path:     "/API/v1/services/{service_id}/notifications/{notification_name}",
		Summary:  "Add notification to a service",
		Tags:     []string{"services"},
		Security: secBearer,
	}, h.ServiceAddNotification)

	huma.Register(api, huma.Operation{
		Method:        "DELETE",
		Path:          "/API/v1/services/{service_id}/notifications/{notification_name}",
		Summary:       "Remove notification from a service",
		Tags:          []string{"services"},
		Security:      secBearer,
		DefaultStatus: 204,
	}, h.ServiceDeleteNotification)

	// ── Notifications ─────────────────────────────────────────────────────────
	huma.Register(api, huma.Operation{
		Method:   "GET",
		Path:     "/API/v1/notifications",
		Summary:  "Get all notifications",
		Tags:     []string{"notifications"},
		Security: secBearer,
	}, h.GetNotifications)

	huma.Register(api, huma.Operation{
		Method:   "POST",
		Path:     "/API/v1/notifications",
		Summary:  "Create a notification",
		Tags:     []string{"notifications"},
		Security: secBearer,
	}, h.CreateNotification)

	huma.Register(api, huma.Operation{
		Method:   "GET",
		Path:     "/API/v1/notifications/{notification_name}",
		Summary:  "Get a notification",
		Tags:     []string{"notifications"},
		Security: secBearer,
	}, h.GetNotification)

	huma.Register(api, huma.Operation{
		Method:   "PATCH",
		Path:     "/API/v1/notifications/{notification_name}",
		Summary:  "Update a notification",
		Tags:     []string{"notifications"},
		Security: secBearer,
	}, h.UpdateNotification)

	huma.Register(api, huma.Operation{
		Method:        "DELETE",
		Path:          "/API/v1/notifications/{notification_name}",
		Summary:       "Delete a notification",
		Tags:          []string{"notifications"},
		Security:      secBearer,
		DefaultStatus: 204,
	}, h.DeleteNotification)
}
