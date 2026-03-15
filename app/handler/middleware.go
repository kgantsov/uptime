package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ctxKey is a custom type for context keys to avoid collisions.
type ctxKey string

// RequestIDKey is the context key used to store and retrieve the request ID.
// You should use this same key in your middleware when setting the context.
const RequestIDKey ctxKey = "request_id"

func RequestIDMiddleware(logMode string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := c.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
			c.Request().Header.Set("X-Request-ID", reqID)
		}

		var reqLogger zerolog.Logger
		if logMode == "STACKDRIVER" {
			reqDict := zerolog.Dict()

			if reqID := c.Get("X-Request-Id"); reqID != "" {
				reqDict.Str("id", reqID)
			}
			if userAgent := c.Get("User-Agent"); userAgent != "" {
				reqDict.Str("user_agent", userAgent)
			}

			// Attach the dictionary to a new request-scoped logger
			// This results in: {"level":"info", "request": {"id": "123", ...}, "message": "..."}
			reqLogger = log.Logger.With().Dict("request", reqDict).Logger()
		} else {
			reqLogger = log.With().Str("request_id", reqID).Logger()
		}

		// 2. Put the logger into the context
		ctx := reqLogger.WithContext(c.UserContext())

		// 3. Put the raw string into the context for the HTTP client headers
		ctx = context.WithValue(ctx, RequestIDKey, reqID)

		// Save the context back to Fiber
		c.SetUserContext(ctx)
		c.Set("X-Request-ID", reqID)

		return c.Next()
	}
}
