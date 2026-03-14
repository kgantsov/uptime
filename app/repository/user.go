package repository

import (
	"github.com/kgantsov/uptime/app/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByEmail(email string) (*model.User, error)
	Count() (int64, error)
	Create(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := r.db.Model(&model.User{}).Where("email = ?", email).First(user).Error
	return user, err
}

func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Count(&count).Error
	return count, err
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}
