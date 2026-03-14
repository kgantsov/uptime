package repository

import (
	"github.com/kgantsov/uptime/app/model"
	"gorm.io/gorm"
)

type TokenRepository interface {
	GetByID(id uint) (*model.Token, error)
	Create(token *model.Token) error
	Delete(id uint) error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) GetByID(id uint) (*model.Token, error) {
	token := &model.Token{}
	err := r.db.Model(&model.Token{}).First(token, id).Error
	return token, err
}

func (r *tokenRepository) Create(token *model.Token) error {
	return r.db.Create(token).Error
}

func (r *tokenRepository) Delete(id uint) error {
	return r.db.Delete(&model.Token{}, id).Error
}
