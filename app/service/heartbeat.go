package service

import (
	"github.com/kgantsov/uptime/app/model"
	"github.com/kgantsov/uptime/app/repository"
)

type HeartbeatService interface {
	GetLatencies(serviceID string, size int) ([]model.Heartbeat, error)
	GetLastLatencies(size int) ([]model.Heartbeat, error)
	GetStats(days int) ([]model.HeartbeatStatsPoint, error)
	DeleteByServiceID(serviceID uint) error
}

type heartbeatService struct {
	repo repository.HeartbeatRepository
}

func NewHeartbeatService(repo repository.HeartbeatRepository) HeartbeatService {
	return &heartbeatService{repo: repo}
}

func (s *heartbeatService) GetLatencies(serviceID string, size int) ([]model.Heartbeat, error) {
	return s.repo.GetLatencies(serviceID, size)
}

func (s *heartbeatService) GetLastLatencies(size int) ([]model.Heartbeat, error) {
	return s.repo.GetLastLatencies(size)
}

func (s *heartbeatService) GetStats(days int) ([]model.HeartbeatStatsPoint, error) {
	return s.repo.GetStats(days)
}

func (s *heartbeatService) DeleteByServiceID(serviceID uint) error {
	return s.repo.DeleteByServiceID(serviceID)
}
