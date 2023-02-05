package auth

import (
	"net/http"
	"strings"

	"github.com/kgantsov/uptime/app/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func AuthSkipperFunc(c echo.Context) bool {
	// Skip authentication for non API requests and for login API request

	if !strings.HasPrefix(c.Path(), "/API/") {
		return true
	}

	if c.Path() == "/API/v1/tokens" && c.Request().Method == "POST" {
		return true
	}

	return false
}

func CheckTokenMiddleware(db *gorm.DB, logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if AuthSkipperFunc(c) {
				return next(c)
			}

			tokenID, err := GetCurrentTokenID(c)

			if err != nil {
				return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
			}

			token := &model.Token{}
			err = db.Model(&model.Token{}).First(&token, tokenID).Error

			if err != nil {
				logger.WithFields(logrus.Fields{
					"RequestID": c.Get(echo.HeaderXRequestID),
				}).Infof("TOKEN WAS NOT FOUND %s", err)

				return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
			}

			return next(c)
		}
	}
}
