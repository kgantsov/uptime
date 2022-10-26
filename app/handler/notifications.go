package handler

import (
	"net/http"

	"github.com/kgantsov/uptime/app/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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
// @Router       /API/v1/notifications [get]
func (h *Handler) GetNotifications(c echo.Context) error {
	var notifications []model.Notification

	err := h.DB.Model(&model.Notification{}).Order("created_at desc").Find(&notifications).Error

	if err != nil {
		log.Errorf("Got an error %s", err)
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
// @Router       /API/v1/notifications/{notification_name} [get]
func (h *Handler) GetNotification(c echo.Context) error {
	notificationName := c.Param("notification_name")

	notification := &model.Notification{}

	err := h.DB.Model(&model.Notification{}).Where("name = ?", notificationName).First(&notification).Error

	if err != nil {
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
// @Param        account  body      model.AddNotification  true  "Add notification"
// @Success      200  {object}  model.Notification
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
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
// @Param        account  body      model.UpdateNotification  true  "Update notification"
// @Success      200  {object}  model.Notification
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
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

	notification.Callback = updateNotification.Callback
	notification.CallbackChatID = updateNotification.CallbackChatID
	notification.CallbackType = updateNotification.CallbackType

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
// @Router       /API/v1/notifications/{notification_name} [delete]
func (h *Handler) DeleteNotification(c echo.Context) error {
	notificationName := c.Param("notification_name")

	err := h.DB.Delete(&model.Service{}, notificationName).Error

	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	return c.NoContent(http.StatusNoContent)
}
