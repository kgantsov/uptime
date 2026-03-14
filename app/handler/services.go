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
// @Security     HttpBearer
// @Router       /API/v1/services [get]
func (h *Handler) GetServices(c echo.Context) error {
	services, err := h.ServiceService.GetServices()
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
// @Security     HttpBearer
// @Router       /API/v1/services/{service_id} [get]
func (h *Handler) GetService(c echo.Context) error {
	id := c.Param("service_id")
	serviceID, err := strconv.Atoi(id)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	service, err := h.ServiceService.GetService(uint(serviceID))
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
// @Security     HttpBearer
// @Router       /API/v1/services [post]
func (h *Handler) CreateService(c echo.Context) error {
	service := new(model.Service)

	if err := c.Bind(service); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	created, err := h.ServiceService.CreateService(service)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err}
	}

	return c.JSON(http.StatusOK, created)
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
// @Security     HttpBearer
// @Router       /API/v1/services/{service_id} [patch]
func (h *Handler) UpdateService(c echo.Context) error {
	id := c.Param("service_id")
	serviceID, err := strconv.Atoi(id)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	updateService := &model.UpdateService{}
	if err = c.Bind(updateService); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	updated, err := h.ServiceService.UpdateService(uint(serviceID), updateService)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	return c.JSON(http.StatusOK, updated)
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
// @Security     HttpBearer
// @Router       /API/v1/services/{service_id} [delete]
func (h *Handler) DeleteService(c echo.Context) error {
	serviceIDStr := c.Param("service_id")
	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	if err := h.ServiceService.DeleteService(uint(serviceID)); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

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
// @Security     HttpBearer
// @Router       /API/v1/services/{service_id}/notifications/{notification_name} [post]
func (h *Handler) ServiceAddNotification(c echo.Context) error {
	id := c.Param("service_id")
	notificationName := c.Param("notification_name")

	serviceID, err := strconv.Atoi(id)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	sn, err := h.ServiceService.AddNotification(uint(serviceID), notificationName)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	return c.JSON(http.StatusOK, sn)
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
// @Security     HttpBearer
// @Router       /API/v1/services/{service_id}/notifications/{notification_name} [delete]
func (h *Handler) ServiceDeleteNotification(c echo.Context) error {
	id := c.Param("service_id")
	notificationName := c.Param("notification_name")

	serviceID, err := strconv.Atoi(id)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	if err := h.ServiceService.DeleteNotification(uint(serviceID), notificationName); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	return c.JSON(http.StatusNoContent, &model.ServiceNotification{})
}
