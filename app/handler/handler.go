package handler

import (
	"github.com/kgantsov/uptime/app/monitor"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
}

func (h *Handler) ConfigureMiddleware(e *echo.Echo) {
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			h.Logger.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")

			return nil
		},
	}))
}
