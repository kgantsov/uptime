package handler

import (
	"net/http"
	"strconv"

	"github.com/kgantsov/uptime/app/model"
	"github.com/labstack/echo/v4"
)

// GetServices godoc
// @Summary      Get all services
// @Description  Returns all services
// @Tags         services
// @Accept       json
// @Produce      json
// @Success      200  {object}  []model.Service
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /API/v1/services [get]
func (h *Handler) GetServices(c echo.Context) error {
	services := []model.Service{}

	err := h.DB.Model(&model.Service{}).Preload("Notifications").Order("id desc").Find(&services).Error

	if err != nil {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: err}
	}

	return c.JSON(http.StatusOK, services)
}

// GetService godoc
// @Summary      Get a service
// @Description  Gets a service by its ID
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        service_id    path     string  true  "Gets service by service_id"
// @Success      200  {object}  model.Service
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /API/v1/services/{service_id} [get]
func (h *Handler) GetService(c echo.Context) error {
	id := c.Param("service_id")
	serviceID, err := strconv.Atoi(id)

	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	service := &model.Service{}

	err = h.DB.Model(&model.Service{}).Preload("Notifications").First(&service, serviceID).Error

	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	return c.JSON(http.StatusOK, service)
}

// CreateService godoc
// @Summary      Create a new service
// @Description  Creates a new service and starts monitor it
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        body  body      model.AddService  true  "Add service"
// @Success      200  {object}  model.Service
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /API/v1/services [post]
func (h *Handler) CreateService(c echo.Context) error {
	service := new(model.Service)

	if err := c.Bind(service); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	h.DB.Create(service)

	h.Dispatcher.AddService(service.ID)

	return c.JSON(http.StatusOK, service)
}

// UpdateService godoc
// @Summary      Update a service
// @Description  Updates an existing service and restarts monitoring for it
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        service_id    path     string  true  "Updates by service_id"
// @Param        body  body      model.UpdateService  true  "Update service"
// @Success      200  {object}  model.Service
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /API/v1/services/{service_id} [patch]
func (h *Handler) UpdateService(c echo.Context) error {
	id := c.Param("service_id")
	serviceID, err := strconv.Atoi(id)

	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	service := &model.Service{}
	updateService := &model.UpdateService{}

	if err = c.Bind(updateService); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err = h.DB.Model(&model.Service{}).Preload("Notifications").First(&service, serviceID).Error

	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	if updateService.Name != nil {
		service.Name = *updateService.Name
	}

	if updateService.URL != nil {
		service.URL = *updateService.URL
	}

	if updateService.Enabled != nil {
		service.Enabled = *updateService.Enabled
	}

	if updateService.CheckInterval != nil {
		service.CheckInterval = *updateService.CheckInterval
	}

	if updateService.Notifications != nil {
		service.Notifications = *updateService.Notifications
	}

	if updateService.Timeout != nil {
		service.Timeout = *updateService.Timeout
	}

	if updateService.AcceptedStatusCode != nil {
		service.AcceptedStatusCode = *updateService.AcceptedStatusCode
	}

	err = h.DB.Where("service_id = ?", serviceID).Delete(&model.ServiceNotification{}).Error
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	if updateService.Notifications != nil {
		for _, notification := range *updateService.Notifications {
			serviceNotification := &model.ServiceNotification{
				ServiceID:        int(service.ID),
				NotificationName: notification.Name,
			}

			h.DB.Create(serviceNotification)
		}
	}

	h.DB.Save(service)

	h.Dispatcher.RestartService(service.ID)

	return c.JSON(http.StatusOK, updateService)
}

// DeleteService godoc
// @Summary      Delete a service
// @Description  Stops a service monitoring and deletes it
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        service_id    path     string  true  "Delete by service_id"
// @Success      204  {object}  model.Service
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /API/v1/services/{service_id} [delete]
func (h *Handler) DeleteService(c echo.Context) error {
	serviceIDStr := c.Param("service_id")
	serviceID, err := strconv.Atoi(serviceIDStr)

	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	err = h.DB.Delete(&model.Service{}, serviceID).Error

	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	h.Dispatcher.RemoveService(uint(serviceID))

	return c.NoContent(http.StatusNoContent)
}

// ServiceAddNotification godoc
// @Summary      Add a notification to a service
// @Description  Adds a notification to a service
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        service_id    path     string  true  "service_id"
// @Param        notification_name    path     string  true  "notification_name"
// @Success      200  {object}  model.ServiceNotification
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /API/v1/services/{service_id}/notifications/{notification_name} [post]
func (h *Handler) ServiceAddNotification(c echo.Context) error {
	id := c.Param("service_id")
	notificationName := c.Param("notification_name")
	serviceID, err := strconv.Atoi(id)

	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	service := &model.Service{}

	err = h.DB.Model(&model.Service{}).Preload("Notifications").First(&service, serviceID).Error

	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	notification := &model.Notification{}

	err = h.DB.Model(&model.Notification{}).Where("name = ?", notificationName).First(&notification).Error

	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	serviceNotification := &model.ServiceNotification{
		ServiceID:        int(service.ID),
		NotificationName: notification.Name,
	}

	h.DB.Create(serviceNotification)

	return c.JSON(http.StatusOK, serviceNotification)
}

// ServiceDeleteNotification godoc
// @Summary      Delete a notification to a service
// @Description  Deletes a notification to a service
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        service_id    path     string  true  "service_id"
// @Param        notification_name    path     string  true  "notification_name"
// @Success      204  {object}  model.ServiceNotification
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /API/v1/services/{service_id}/notifications/{notification_name} [delete]
func (h *Handler) ServiceDeleteNotification(c echo.Context) error {
	id := c.Param("service_id")
	notificationName := c.Param("notification_name")
	serviceID, err := strconv.Atoi(id)

	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	err = h.DB.Where(
		"service_id = ? AND notification_name = ?", serviceID, notificationName,
	).Delete(&model.ServiceNotification{}).Error

	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	return c.JSON(http.StatusNoContent, &model.ServiceNotification{})
}
