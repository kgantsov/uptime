package handler

import (
	"net/http"

	"github.com/kgantsov/uptime/app/auth"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// CreateToken godoc
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
	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	token, err := h.TokenService.CreateToken(req.Email, req.Password)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Email or password is incorrect",
		}
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
	tokenID, _ := auth.GetCurrentTokenID(c)

	h.Logger.WithFields(logrus.Fields{
		"RequestID": c.Get(echo.HeaderXRequestID),
	}).Infof("FOUND USER ID %+v", tokenID)

	if err := h.TokenService.DeleteToken(tokenID); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	return c.JSON(http.StatusNoContent, struct{}{})
}
