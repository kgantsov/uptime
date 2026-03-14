package repository

import (
	"github.com/kgantsov/uptime/app/model"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	GetAll() ([]model.Service, error)
	GetByID(id uint) (*model.Service, error)
	Create(service *model.Service) error
	Save(service *model.Service) error
	Delete(id uint) error
	DeleteServiceNotifications(serviceID int) error
	CreateServiceNotification(sn *model.ServiceNotification) error
	DeleteServiceNotification(serviceID int, notificationName string) error
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) GetAll() ([]model.Service, error) {
	services := []model.Service{}
	err := r.db.Model(&model.Service{}).Preload("Notifications").Order("id desc").Find(&services).Error
	return services, err
}

func (r *serviceRepository) GetByID(id uint) (*model.Service, error) {
	service := &model.Service{}
	err := r.db.Model(&model.Service{}).Preload("Notifications").First(service, id).Error
	return service, err
}

func (r *serviceRepository) Create(service *model.Service) error {
	return r.db.Create(service).Error
}

func (r *serviceRepository) Save(service *model.Service) error {
	return r.db.Save(service).Error
}

func (r *serviceRepository) Delete(id uint) error {
	return r.db.Delete(&model.Service{}, id).Error
}

func (r *serviceRepository) DeleteServiceNotifications(serviceID int) error {
	return r.db.Where("service_id = ?", serviceID).Delete(&model.ServiceNotification{}).Error
}

func (r *serviceRepository) CreateServiceNotification(sn *model.ServiceNotification) error {
	return r.db.Create(sn).Error
}

func (r *serviceRepository) DeleteServiceNotification(serviceID int, notificationName string) error {
	return r.db.Where(
		"service_id = ? AND notification_name = ?", serviceID, notificationName,
	).Delete(&model.ServiceNotification{}).Error
}
