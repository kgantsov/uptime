package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/kgantsov/uptime/app/model"
)

// ── Input / Output types ──────────────────────────────────────────────────────

type GetNotificationsOutput struct {
	Body []model.Notification
}

type GetNotificationInput struct {
	NotificationName string `path:"notification_name" doc:"Notification name"`
}

type GetNotificationOutput struct {
	Body *model.Notification
}

type CreateNotificationInput struct {
	Body model.AddNotification
}

type CreateNotificationOutput struct {
	Body *model.Notification
}

type UpdateNotificationInput struct {
	NotificationName string `path:"notification_name" doc:"Notification name"`
	Body             model.UpdateNotification
}

type UpdateNotificationOutput struct {
	Body *model.Notification
}

type DeleteNotificationInput struct {
	NotificationName string `path:"notification_name" doc:"Notification name"`
}

// ── Handlers ──────────────────────────────────────────────────────────────────

// GetNotifications returns all notification channels.
func (h *Handler) GetNotifications(
	ctx context.Context,
	input *struct{},
) (*GetNotificationsOutput, error) {
	notifications, err := h.NotificationService.GetNotifications()
	if err != nil {
		h.Logger.Infof("error getting notifications: %s", err)
		return nil, huma.NewError(http.StatusInternalServerError, "failed to retrieve notifications", err)
	}

	if notifications == nil {
		notifications = []model.Notification{}
	}

	return &GetNotificationsOutput{Body: notifications}, nil
}

// GetNotification returns a single notification channel by name.
func (h *Handler) GetNotification(
	ctx context.Context,
	input *GetNotificationInput,
) (*GetNotificationOutput, error) {
	notification, err := h.NotificationService.GetNotification(input.NotificationName)
	if err != nil {
		h.Logger.Infof("error getting notification %q: %s", input.NotificationName, err)
		return nil, huma.NewError(http.StatusBadRequest, "notification not found", err)
	}

	return &GetNotificationOutput{Body: notification}, nil
}

// CreateNotification creates a new notification channel.
func (h *Handler) CreateNotification(
	ctx context.Context,
	input *CreateNotificationInput,
) (*CreateNotificationOutput, error) {
	notification := &model.Notification{
		Name:           input.Body.Name,
		CallbackType:   input.Body.CallbackType,
		CallbackChatID: input.Body.CallbackChatID,
		Callback:       input.Body.Callback,
	}

	created, err := h.NotificationService.CreateNotification(notification)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "failed to create notification", err)
	}

	return &CreateNotificationOutput{Body: created}, nil
}

// UpdateNotification updates an existing notification channel.
func (h *Handler) UpdateNotification(
	ctx context.Context,
	input *UpdateNotificationInput,
) (*UpdateNotificationOutput, error) {
	notification, err := h.NotificationService.UpdateNotification(input.NotificationName, &input.Body)
	if err != nil {
		return nil, huma.NewError(http.StatusBadRequest, "failed to update notification", err)
	}

	return &UpdateNotificationOutput{Body: notification}, nil
}

// DeleteNotification removes a notification channel.
func (h *Handler) DeleteNotification(
	ctx context.Context,
	input *DeleteNotificationInput,
) (*struct{}, error) {
	if err := h.NotificationService.DeleteNotification(input.NotificationName); err != nil {
		return nil, huma.NewError(http.StatusBadRequest, "failed to delete notification", err)
	}

	return nil, nil
}
