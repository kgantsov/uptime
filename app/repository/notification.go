package repository

import (
	"github.com/kgantsov/uptime/app/model"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	GetAll() ([]model.Notification, error)
	GetByName(name string) (*model.Notification, error)
	Create(notification *model.Notification) error
	Save(notification *model.Notification) error
	Delete(name string) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) GetAll() ([]model.Notification, error) {
	notifications := []model.Notification{}
	err := r.db.Model(&model.Notification{}).Order("created_at desc").Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) GetByName(name string) (*model.Notification, error) {
	notification := &model.Notification{}
	err := r.db.Model(&model.Notification{}).Where("name = ?", name).First(notification).Error
	return notification, err
}

func (r *notificationRepository) Create(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) Save(notification *model.Notification) error {
	return r.db.Save(notification).Error
}

func (r *notificationRepository) Delete(name string) error {
	return r.db.Where("name = ?", name).Delete(&model.Notification{}).Error
}
