package handler

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kgantsov/uptime/app/auth"
	"github.com/kgantsov/uptime/app/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// CreateService godoc
// @Summary      Create an auth token
// @Description  Create an auth token
// @Tags         tokens
// @Accept       json
// @Produce      json
// @Param        body  body     model.CreateToken  true  "Create an auth token"
// @Success      200  {object}  model.Token
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /API/v1/tokens [post]
func (h *Handler) CreateToken(c echo.Context) (err error) {
	t := new(model.CreateToken)
	user := new(model.User)

	if err := c.Bind(t); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err = h.DB.Model(&model.User{}).Where("email = ?", "admin@uptime.io").First(&user).Error
	if err != nil {
		return &echo.HTTPError{
			Code: http.StatusBadRequest, Message: "Email or password is incorrect",
		}
	}

	if !auth.CheckPasswordHash(t.Password, user.Password) {
		return &echo.HTTPError{
			Code: http.StatusBadRequest, Message: "Email or password is incorrect",
		}
	}

	token := &model.Token{
		UserID:   user.ID,
		ExpireAt: time.Now().Add(time.Hour * 72),
	}
	h.DB.Create(token)

	jwtToken := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["id"] = token.ID
	claims["exp"] = token.ExpireAt.Unix()

	// Generate encoded token and send it as response
	token.Token, err = jwtToken.SignedString([]byte(Key))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, token)
}

// DeleteToken godoc
// @Summary      Delete an auth token
// @Description  Delete an auth token
// @Tags         tokens
// @Accept       json
// @Produce      json
// @Success      204  {object}  model.Token
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Security     HttpBearer
// @Router       /API/v1/tokens [delete]
func (h *Handler) DeleteToken(c echo.Context) (err error) {
	tokenId, _ := auth.GetCurrentTokenID(c)

	h.Logger.WithFields(logrus.Fields{
		"RequestID": c.Get(echo.HeaderXRequestID),
	}).Infof("FOUND USER ID %+v", tokenId)

	err = h.DB.Delete(&model.Token{}, tokenId).Error

	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	return c.JSON(http.StatusNoContent, &model.Token{})
}
