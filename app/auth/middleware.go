package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// TokenLookup is the minimal interface the middleware needs to verify that
// a token ID extracted from a JWT actually exists in the backing store.
type TokenLookup interface {
	ValidateToken(id uint) error
}

// AuthSkipperFunc returns true when JWT authentication should be skipped.
// It skips all non-API routes and the token creation endpoint.
func AuthSkipperFunc(c *fiber.Ctx) bool {
	if !strings.HasPrefix(strings.ToLower(c.Path()), "/api/") {
		return true
	}

	if strings.ToLower(c.Path()) == "/api/v1/tokens" && c.Method() == fiber.MethodPost {
		return true
	}

	return false
}

// CheckTokenMiddleware returns a Fiber middleware that validates the token ID
// (extracted from the JWT by the upstream JWT middleware) against the backing
// store, ensuring that tokens that have been explicitly deleted are rejected.
func CheckTokenMiddleware(tokenLookup TokenLookup, logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if AuthSkipperFunc(c) {
			return c.Next()
		}

		tokenID, err := GetCurrentTokenID(c)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		if err := tokenLookup.ValidateToken(tokenID); err != nil {
			logger.WithFields(logrus.Fields{
				"RequestID": c.Locals("requestid"),
			}).Infof("TOKEN WAS NOT FOUND %s", err)

			return fiber.NewError(fiber.StatusUnauthorized, "token not found or expired")
		}

		return c.Next()
	}
}
