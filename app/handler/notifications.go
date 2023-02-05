package handler

import (
	"net/http"

	"github.com/kgantsov/uptime/app/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// GetNotifications godoc
// @Summary      Get notifications
// @Description  Returns all notifications
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Success      200  {object}  []model.Notification
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Security     HttpBearer
// @Router       /API/v1/notifications [get]
func (h *Handler) GetNotifications(c echo.Context) error {
	notifications := []model.Notification{}

	err := h.DB.Model(&model.Notification{}).Order("created_at desc").Find(&notifications).Error

	if err != nil {
		h.Logger.WithFields(logrus.Fields{
			"RequestID": c.Get(echo.HeaderXRequestID),
		}).Infof("Got an error getting notifications %s", err)

		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err}
	}

	return c.JSON(http.StatusOK, notifications)
}

// GetNotification godoc
// @Summary      Get a notification
// @Description  Returns a notification
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Param        notification_name    path     string  true  "Get a notification by notification_name"
// @Success      200  {object}  model.Notification
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Security     HttpBearer
// @Router       /API/v1/notifications/{notification_name} [get]
func (h *Handler) GetNotification(c echo.Context) error {
	notificationName := c.Param("notification_name")

	notification := &model.Notification{}

	err := h.DB.Model(&model.Notification{}).Where("name = ?", notificationName).First(&notification).Error

	if err != nil {
		h.Logger.WithFields(logrus.Fields{
			"RequestID": c.Get(echo.HeaderXRequestID),
		}).Infof("Got an error getting a notification %s %s", notificationName, err)

		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	return c.JSON(http.StatusOK, notification)
}

// CreateNotification godoc
// @Summary      Create a new notification
// @Description  Creates notifications
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Param        body  body      model.AddNotification  true  "Add notification"
// @Success      200  {object}  model.Notification
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Security     HttpBearer
// @Router       /API/v1/notifications [post]
func (h *Handler) CreateNotification(c echo.Context) error {
	notification := new(model.Notification)

	if err := c.Bind(notification); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	h.DB.Create(notification)

	return c.JSON(http.StatusOK, notification)
}

// UpdateNotification godoc
// @Summary      Update a notification
// @Description  Updates a notification
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Param        notification_name    path     string  true  "Updates a notification by notification_name"
// @Param        body  body      model.UpdateNotification  true  "Update notification"
// @Success      200  {object}  model.Notification
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Security     HttpBearer
// @Router       /API/v1/notifications/{notification_name} [patch]
func (h *Handler) UpdateNotification(c echo.Context) error {
	notificationName := c.Param("notification_name")

	notification := &model.Notification{}
	updateNotification := &model.UpdateNotification{}

	if err := c.Bind(updateNotification); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	err := h.DB.Model(&model.Notification{}).Where("name = ?", notificationName).First(&notification).Error

	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	if updateNotification.Callback != nil {
		notification.Callback = *updateNotification.Callback
	}

	if updateNotification.CallbackChatID != nil {
		notification.CallbackChatID = *updateNotification.CallbackChatID
	}

	if updateNotification.CallbackType != nil {
		notification.CallbackType = *updateNotification.CallbackType
	}

	h.DB.Save(notification)

	h.Dispatcher.Stop()
	h.Dispatcher.Start()

	return c.JSON(http.StatusOK, notification)
}

// DeleteNotification godoc
// @Summary      Delete a notification
// @Description  Deletes notifications
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Param        notification_name    path     string  true  "Delete by notification_name"
// @Success      200  {object}  model.Notification
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Security     HttpBearer
// @Router       /API/v1/notifications/{notification_name} [delete]
func (h *Handler) DeleteNotification(c echo.Context) error {
	notificationName := c.Param("notification_name")

	err := h.DB.Where("name = ?", notificationName).Delete(&model.Notification{}).Error

	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	return c.NoContent(http.StatusNoContent)
}
