package service

import (
	"github.com/kgantsov/uptime/app/model"
	"github.com/kgantsov/uptime/app/repository"
)

type ServiceService interface {
	GetServices() ([]model.Service, error)
	GetService(id uint) (*model.Service, error)
	CreateService(service *model.Service) (*model.Service, error)
	UpdateService(id uint, update *model.UpdateService) (*model.UpdateService, error)
	DeleteService(id uint) error
	AddNotification(serviceID uint, notificationName string) (*model.ServiceNotification, error)
	DeleteNotification(serviceID uint, notificationName string) error
}

type serviceService struct {
	serviceRepo   repository.ServiceRepository
	notifRepo     repository.NotificationRepository
	heartbeatRepo repository.HeartbeatRepository
	dispatcher    DispatcherInterface
}

func NewServiceService(
	serviceRepo repository.ServiceRepository,
	notifRepo repository.NotificationRepository,
	heartbeatRepo repository.HeartbeatRepository,
	dispatcher DispatcherInterface,
) ServiceService {
	return &serviceService{
		serviceRepo:   serviceRepo,
		notifRepo:     notifRepo,
		heartbeatRepo: heartbeatRepo,
		dispatcher:    dispatcher,
	}
}

func (s *serviceService) GetServices() ([]model.Service, error) {
	return s.serviceRepo.GetAll()
}

func (s *serviceService) GetService(id uint) (*model.Service, error) {
	return s.serviceRepo.GetByID(id)
}

func (s *serviceService) CreateService(service *model.Service) (*model.Service, error) {
	if err := s.serviceRepo.Create(service); err != nil {
		return nil, err
	}

	if s.dispatcher != nil {
		s.dispatcher.AddService(service.ID)
	}

	return service, nil
}

func (s *serviceService) UpdateService(id uint, update *model.UpdateService) (*model.UpdateService, error) {
	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if update.Name != nil {
		service.Name = *update.Name
	}
	if update.URL != nil {
		service.URL = *update.URL
	}
	if update.Enabled != nil {
		service.Enabled = *update.Enabled
	}
	if update.CheckInterval != nil {
		service.CheckInterval = *update.CheckInterval
	}
	if update.Retries != nil {
		service.Retries = *update.Retries
	}
	if update.Timeout != nil {
		service.Timeout = *update.Timeout
	}
	if update.AcceptedStatusCode != nil {
		service.AcceptedStatusCode = *update.AcceptedStatusCode
	}

	if err := s.serviceRepo.DeleteServiceNotifications(int(id)); err != nil {
		return nil, err
	}

	if update.Notifications != nil {
		service.Notifications = *update.Notifications
		for _, notification := range *update.Notifications {
			sn := &model.ServiceNotification{
				ServiceID:        int(service.ID),
				NotificationName: notification.Name,
			}
			if err := s.serviceRepo.CreateServiceNotification(sn); err != nil {
				return nil, err
			}
		}
	}

	if err := s.serviceRepo.Save(service); err != nil {
		return nil, err
	}

	if s.dispatcher != nil {
		s.dispatcher.RestartService(service.ID)
	}

	return update, nil
}

func (s *serviceService) DeleteService(id uint) error {
	if err := s.serviceRepo.DeleteServiceNotifications(int(id)); err != nil {
		return err
	}

	if err := s.heartbeatRepo.DeleteByServiceID(id); err != nil {
		return err
	}

	if err := s.serviceRepo.Delete(id); err != nil {
		return err
	}

	if s.dispatcher != nil {
		s.dispatcher.RemoveService(id)
	}

	return nil
}

func (s *serviceService) AddNotification(serviceID uint, notificationName string) (*model.ServiceNotification, error) {
	if _, err := s.serviceRepo.GetByID(serviceID); err != nil {
		return nil, err
	}

	if _, err := s.notifRepo.GetByName(notificationName); err != nil {
		return nil, err
	}

	sn := &model.ServiceNotification{
		ServiceID:        int(serviceID),
		NotificationName: notificationName,
	}

	if err := s.serviceRepo.CreateServiceNotification(sn); err != nil {
		return nil, err
	}

	return sn, nil
}

func (s *serviceService) DeleteNotification(serviceID uint, notificationName string) error {
	return s.serviceRepo.DeleteServiceNotification(int(serviceID), notificationName)
}
