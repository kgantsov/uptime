package service

import (
	"github.com/kgantsov/uptime/app/model"
	"github.com/kgantsov/uptime/app/repository"
)

type NotificationService interface {
	GetNotifications() ([]model.Notification, error)
	GetNotification(name string) (*model.Notification, error)
	CreateNotification(notification *model.Notification) (*model.Notification, error)
	UpdateNotification(name string, update *model.UpdateNotification) (*model.Notification, error)
	DeleteNotification(name string) error
}

type notificationService struct {
	repo       repository.NotificationRepository
	dispatcher DispatcherInterface
}

func NewNotificationService(
	repo repository.NotificationRepository,
	dispatcher DispatcherInterface,
) NotificationService {
	return &notificationService{
		repo:       repo,
		dispatcher: dispatcher,
	}
}

func (s *notificationService) GetNotifications() ([]model.Notification, error) {
	return s.repo.GetAll()
}

func (s *notificationService) GetNotification(name string) (*model.Notification, error) {
	return s.repo.GetByName(name)
}

func (s *notificationService) CreateNotification(notification *model.Notification) (*model.Notification, error) {
	if err := s.repo.Create(notification); err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *notificationService) UpdateNotification(name string, update *model.UpdateNotification) (*model.Notification, error) {
	notification, err := s.repo.GetByName(name)
	if err != nil {
		return nil, err
	}

	if update.Callback != nil {
		notification.Callback = *update.Callback
	}
	if update.CallbackChatID != nil {
		notification.CallbackChatID = *update.CallbackChatID
	}
	if update.CallbackType != nil {
		notification.CallbackType = *update.CallbackType
	}

	if err := s.repo.Save(notification); err != nil {
		return nil, err
	}

	if s.dispatcher != nil {
		s.dispatcher.Stop()
		s.dispatcher.Start()
	}

	return notification, nil
}

func (s *notificationService) DeleteNotification(name string) error {
	return s.repo.Delete(name)
}
