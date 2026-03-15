package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/kgantsov/uptime/app/model"
)

// ── Input / Output types ──────────────────────────────────────────────────────

type GetServicesOutput struct {
	Body []model.Service
}

type GetServiceInput struct {
	ServiceID int `path:"service_id" doc:"Service ID"`
}

type GetServiceOutput struct {
	Body *model.Service
}

type CreateServiceInput struct {
	Body model.AddService
}

type CreateServiceOutput struct {
	Body *model.Service
}

type UpdateServiceInput struct {
	ServiceID int `path:"service_id" doc:"Service ID"`
	Body      model.UpdateService
}

type UpdateServiceOutput struct {
	Body *model.UpdateService
}

type DeleteServiceInput struct {
	ServiceID int `path:"service_id" doc:"Service ID"`
}

type ServiceNotificationInput struct {
	ServiceID        int    `path:"service_id"        doc:"Service ID"`
	NotificationName string `path:"notification_name" doc:"Notification name"`
}

type ServiceNotificationOutput struct {
	Body *model.ServiceNotification
}

// ── Handlers ──────────────────────────────────────────────────────────────────

// GetServices returns every service in the system.
func (h *Handler) GetServices(
	ctx context.Context,
	input *struct{},
) (*GetServicesOutput, error) {
	services, err := h.ServiceService.GetServices()
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, "failed to retrieve services", err)
	}

	if services == nil {
		services = []model.Service{}
	}

	return &GetServicesOutput{Body: services}, nil
}

// GetService returns a single service by its ID.
func (h *Handler) GetService(
	ctx context.Context,
	input *GetServiceInput,
) (*GetServiceOutput, error) {
	svc, err := h.ServiceService.GetService(uint(input.ServiceID))
	if err != nil {
		return nil, huma.NewError(http.StatusBadRequest, "service not found", err)
	}

	return &GetServiceOutput{Body: svc}, nil
}

// CreateService creates a new service and starts monitoring it.
func (h *Handler) CreateService(
	ctx context.Context,
	input *CreateServiceInput,
) (*CreateServiceOutput, error) {
	notifications := make([]model.Notification, 0, len(input.Body.Notifications))
	for _, n := range input.Body.Notifications {
		notifications = append(notifications, model.Notification{
			Name:           n.Name,
			CallbackType:   n.CallbackType,
			CallbackChatID: n.CallbackChatID,
			Callback:       n.Callback,
		})
	}

	svc := &model.Service{
		Name:               input.Body.Name,
		URL:                input.Body.URL,
		Enabled:            input.Body.Enabled,
		Timeout:            input.Body.Timeout,
		CheckInterval:      input.Body.CheckInterval,
		Retries:            input.Body.Retries,
		AcceptedStatusCode: input.Body.AcceptedStatusCode,
		Notifications:      notifications,
	}

	created, err := h.ServiceService.CreateService(svc)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "failed to create service", err)
	}

	return &CreateServiceOutput{Body: created}, nil
}

// UpdateService updates an existing service and restarts monitoring for it.
func (h *Handler) UpdateService(
	ctx context.Context,
	input *UpdateServiceInput,
) (*UpdateServiceOutput, error) {
	updated, err := h.ServiceService.UpdateService(uint(input.ServiceID), &input.Body)
	if err != nil {
		return nil, huma.NewError(http.StatusBadRequest, "failed to update service", err)
	}

	return &UpdateServiceOutput{Body: updated}, nil
}

// DeleteService stops monitoring and removes a service.
func (h *Handler) DeleteService(
	ctx context.Context,
	input *DeleteServiceInput,
) (*struct{}, error) {
	if err := h.ServiceService.DeleteService(uint(input.ServiceID)); err != nil {
		return nil, huma.NewError(http.StatusBadRequest, "failed to delete service", err)
	}

	return nil, nil
}

// ServiceAddNotification links a notification channel to a service.
func (h *Handler) ServiceAddNotification(
	ctx context.Context,
	input *ServiceNotificationInput,
) (*ServiceNotificationOutput, error) {
	sn, err := h.ServiceService.AddNotification(uint(input.ServiceID), input.NotificationName)
	if err != nil {
		return nil, huma.NewError(http.StatusBadRequest, "failed to add notification", err)
	}

	return &ServiceNotificationOutput{Body: sn}, nil
}

// ServiceDeleteNotification removes a notification channel from a service.
func (h *Handler) ServiceDeleteNotification(
	ctx context.Context,
	input *ServiceNotificationInput,
) (*struct{}, error) {
	if err := h.ServiceService.DeleteNotification(uint(input.ServiceID), input.NotificationName); err != nil {
		return nil, huma.NewError(http.StatusBadRequest, "failed to remove notification", err)
	}

	return nil, nil
}
