package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func GetCurrentTokenID(c echo.Context) (uint, error) {
	t, ok := c.Get("token").(*jwt.Token)
	if !ok {
		return 0, fmt.Errorf("Token not found")
	}

	claims, ok := t.Claims.(jwt.MapClaims)

	if !ok {
		return 0, fmt.Errorf("Token is invalid")
	}

	return uint(claims["id"].(float64)), nil
}
