package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const TokenIDContextKey contextKey = "tokenID"

// GetCurrentTokenID extracts the JWT token ID from a Fiber context.
// The gofiber/contrib/jwt middleware stores the parsed token in c.Locals("user").
func GetCurrentTokenID(c *fiber.Ctx) (uint, error) {
	t, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return 0, fmt.Errorf("token not found")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("token is invalid")
	}

	id, ok := claims["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("token id claim is missing or invalid")
	}

	return uint(id), nil
}

// GetCurrentTokenIDFromCtx extracts the JWT token ID injected into a
// context.Context by the token-injection middleware.
func GetCurrentTokenIDFromCtx(ctx context.Context) (uint, error) {
	id, ok := ctx.Value(TokenIDContextKey).(uint)
	if !ok || id == 0 {
		return 0, fmt.Errorf("token not found in context")
	}
	return id, nil
}

// ParseTokenIDFromHeader parses the JWT from a raw Authorization header value
// (e.g. "Bearer eyJhbGc...") and returns the token ID claim.
// This is used in handlers that need the token ID but only have access to a
// context.Context (not the underlying Fiber context).
func ParseTokenIDFromHeader(authHeader string, signingKey string) (uint, error) {
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == "" {
		return 0, fmt.Errorf("authorization header is empty")
	}

	t, err := jwt.Parse(tokenStr, func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tok.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("token claims are invalid")
	}

	id, ok := claims["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("token id claim is missing or invalid")
	}

	return uint(id), nil
}
