package repository

import (
	"fmt"
	"time"

	"github.com/kgantsov/uptime/app/model"
	"gorm.io/gorm"
)

type HeartbeatRepository interface {
	GetLatencies(serviceID string, size int) ([]model.Heartbeat, error)
	GetLastLatencies(size int) ([]model.Heartbeat, error)
	GetStats(days int) ([]model.HeartbeatStatsPoint, error)
	Create(heartbeat *model.Heartbeat) error
	DeleteOlderThan(threshold time.Time) error
	DeleteByServiceID(serviceID uint) error
}

type heartbeatRepository struct {
	db *gorm.DB
}

func NewHeartbeatRepository(db *gorm.DB) HeartbeatRepository {
	return &heartbeatRepository{db: db}
}

func (r *heartbeatRepository) GetLatencies(serviceID string, size int) ([]model.Heartbeat, error) {
	heartbeats := []model.Heartbeat{}
	var err error

	if serviceID != "" {
		err = r.db.Order("id desc").Where("service_id = ?", serviceID).Limit(size).Find(&heartbeats).Error
	} else {
		err = r.db.Order("id desc").Limit(size).Find(&heartbeats).Error
	}

	return heartbeats, err
}

func (r *heartbeatRepository) GetLastLatencies(size int) ([]model.Heartbeat, error) {
	heartbeats := []model.Heartbeat{}

	err := r.db.Raw(
		`
		SELECT * FROM
		(
			SELECT id, service_id, status, created_at, response_time, status_code,
			ROW_NUMBER() OVER (PARTITION BY service_id Order by created_at DESC) AS size
			FROM heartbeats
		) RNK
		WHERE size <= ?
		`,
		size,
	).Scan(&heartbeats).Error

	return heartbeats, err
}

func (r *heartbeatRepository) GetStats(days int) ([]model.HeartbeatStatsPoint, error) {
	heartbeatStatsPoints := []model.HeartbeatStatsPoint{}

	err := r.db.Model(
		&model.Heartbeat{},
	).Select(
		"service_id, status, count(1) as counter, avg(response_time) as average_response_time",
	).Where(
		"created_at > DATE('now', ?)", fmt.Sprintf("-%d days", days),
	).Group("service_id, status").Find(&heartbeatStatsPoints).Error

	return heartbeatStatsPoints, err
}

func (r *heartbeatRepository) Create(heartbeat *model.Heartbeat) error {
	return r.db.Create(heartbeat).Error
}

func (r *heartbeatRepository) DeleteOlderThan(threshold time.Time) error {
	return r.db.Where("created_at < ?", threshold).Delete(&model.Heartbeat{}).Error
}

func (r *heartbeatRepository) DeleteByServiceID(serviceID uint) error {
	return r.db.Where("service_id = ?", serviceID).Delete(&model.Heartbeat{}).Error
}
