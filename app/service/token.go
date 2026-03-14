package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kgantsov/uptime/app/auth"
	"github.com/kgantsov/uptime/app/model"
	"github.com/kgantsov/uptime/app/repository"
)

type TokenService interface {
	CreateToken(email, password string) (*model.Token, error)
	GetToken(id uint) (*model.Token, error)
	ValidateToken(id uint) error
	DeleteToken(id uint) error
}

type tokenService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
	jwtKey    string
}

func NewTokenService(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	jwtKey string,
) TokenService {
	return &tokenService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtKey:    jwtKey,
	}
}

func (s *tokenService) CreateToken(email, password string) (*model.Token, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("email or password is incorrect")
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		return nil, fmt.Errorf("email or password is incorrect")
	}

	token := &model.Token{
		UserID:   user.ID,
		ExpireAt: time.Now().Add(time.Hour * 72),
	}

	if err := s.tokenRepo.Create(token); err != nil {
		return nil, err
	}

	jwtToken := jwt.New(jwt.SigningMethodHS256)

	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["id"] = token.ID
	claims["exp"] = token.ExpireAt.Unix()

	signed, err := jwtToken.SignedString([]byte(s.jwtKey))
	if err != nil {
		return nil, err
	}

	token.Token = signed

	return token, nil
}

func (s *tokenService) GetToken(id uint) (*model.Token, error) {
	return s.tokenRepo.GetByID(id)
}

func (s *tokenService) ValidateToken(id uint) error {
	_, err := s.tokenRepo.GetByID(id)
	return err
}

func (s *tokenService) DeleteToken(id uint) error {
	return s.tokenRepo.Delete(id)
}
