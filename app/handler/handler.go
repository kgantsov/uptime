package handler

import (
	"time"

	"github.com/kgantsov/uptime/app/auth"
	"github.com/kgantsov/uptime/app/monitor"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	echoSwagger "github.com/swaggo/echo-swagger"
)

type (
	Handler struct {
		DB         *gorm.DB
		Dispatcher *monitor.Dispatcher
		Logger     *logrus.Logger
	}
)

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)

func NewHandler(logger *logrus.Logger, db *gorm.DB, dispatcher *monitor.Dispatcher) *Handler {
	h := &Handler{Logger: logger, DB: db, Dispatcher: dispatcher}

	return h
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	v1 := e.Group("/API/v1")

	v1.POST("/tokens", h.CreateToken)
	v1.DELETE("/tokens", h.DeleteToken)

	v1.GET("/heartbeats/latencies", h.GetHeartbeatsLatencies)
	v1.GET("/heartbeats/latencies/last", h.GetHeartbeatsLastLatencies)
	v1.GET("/heartbeats/stats/:days", h.GetHeartbeatStats)

	v1.GET("/services", h.GetServices)
	v1.POST("/services", h.CreateService)
	v1.GET("/services/:service_id", h.GetService)
	v1.PATCH("/services/:service_id", h.UpdateService)
	v1.DELETE("/services/:service_id", h.DeleteService)
	v1.POST("/services/:service_id/notifications/:notification_name", h.ServiceAddNotification)
	v1.DELETE("/services/:service_id/notifications/:notification_name", h.ServiceDeleteNotification)

	v1.GET("/notifications", h.GetNotifications)
	v1.POST("/notifications", h.CreateNotification)
	v1.GET("/notifications/:notification_name", h.GetNotification)
	v1.PATCH("/notifications/:notification_name", h.UpdateNotification)
	v1.DELETE("/notifications/:notification_name", h.DeleteNotification)

	e.GET("/docs/*", echoSwagger.WrapHandler)
}

func (h *Handler) ConfigureMiddleware(e *echo.Echo) {
	e.HideBanner = true
	e.Logger.SetLevel(log.DEBUG)

	e.Use(middleware.RequestID())

	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		RequestIDHandler: func(c echo.Context, rid string) {
			c.Set(echo.HeaderXRequestID, rid)
		},
	}))

	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(Key),
		ContextKey: "token",

		Skipper: auth.AuthSkipperFunc,
	}))
	e.Use(auth.CheckTokenMiddleware(h.DB, h.Logger))

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogRequestID: true,

		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			h.Logger.WithFields(logrus.Fields{
				"RequestID": values.RequestID,
			}).Infof("%s %s %d", values.Method, values.URI, values.Status)

			return nil
		},
	}))

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		// LogLevel:  log.Lvl,
	}))

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			h.Logger.WithFields(logrus.Fields{
				"RequestID": c.Get(echo.HeaderXRequestID),
			}).Warn("GotTimeout")
		},
		Timeout: 30 * time.Second,
	}))
	e.Use(middleware.SecureWithConfig(middleware.DefaultSecureConfig))

	e.Use(echoprometheus.NewMiddleware("uptime"))  // adds middleware to gather metrics
	e.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics

}
